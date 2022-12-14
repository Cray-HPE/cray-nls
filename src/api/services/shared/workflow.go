/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package services_shared

//go:generate mockgen -destination=../mocks/services/workflow.go -package=mocks -source=workflow.go

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/argoproj/pkg/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	argo_templates "github.com/Cray-HPE/cray-nls/src/api/argo-templates"
	models_nls "github.com/Cray-HPE/cray-nls/src/api/models/nls"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
)

type WorkflowService interface {
	GetWorkflows(ctx *gin.Context) (*v1alpha1.WorkflowList, error)
	GetWorkflowByName(name string, ctx *gin.Context) (*v1alpha1.Workflow, error)
	DeleteWorkflow(ctx *gin.Context) error
	RerunWorkflow(ctx *gin.Context) error
	RetryWorkflow(ctx *gin.Context) error
	CreateRebuildWorkflow(req models_nls.CreateRebuildWorkflowRequest) (*v1alpha1.Workflow, error)
	InitializeWorkflowTemplate(template []byte) error
}

// WorkflowService service layer
type workflowService struct {
	logger                 utils.Logger
	ctx                    context.Context
	workflowClient         workflow.WorkflowServiceClient
	workflowTemplateClient workflowtemplate.WorkflowTemplateServiceClient
	k8sRestClientSet       *kubernetes.Clientset
	env                    utils.Env
}

// NewWorkflowService creates a new Workflowservice
func NewWorkflowService(logger utils.Logger, argoService ArgoService, k8sSvc K8sService, env utils.Env) WorkflowService {

	workflowTemplateClient, _ := argoService.Client.NewWorkflowTemplateServiceClient()

	workflowSvc := workflowService{
		logger:                 logger,
		ctx:                    argoService.Context,
		workflowClient:         argoService.Client.NewWorkflowServiceClient(),
		workflowTemplateClient: workflowTemplateClient,
		k8sRestClientSet:       k8sSvc.Client,
		env:                    env,
	}
	workflowTemplates, _ := argo_templates.GetWorkflowTemplate()
	for _, workflowTemplate := range workflowTemplates {
		err := workflowSvc.InitializeWorkflowTemplate(workflowTemplate)
		if err != nil {
			return nil
		}
	}
	return workflowSvc
}

func (s workflowService) DeleteWorkflow(ctx *gin.Context) error {
	wfName := ctx.Param("name")
	workflowToDelete, err := s.workflowClient.GetWorkflow(
		s.ctx,
		&workflow.WorkflowGetRequest{
			Namespace: "argo",
			Name:      wfName,
		},
	)
	if err != nil {
		s.logger.Error(err)
		return fmt.Errorf("failed to find workflow with name: %s", wfName)
	}
	// only delete rebuild workflow
	if workflowToDelete.Labels["type"] != "rebuild" {
		err := fmt.Errorf("workflow type is wrong: %s", workflowToDelete.Labels["type"])
		s.logger.Error(err)
		return err
	}

	_, err = s.workflowClient.DeleteWorkflow(
		s.ctx,
		&workflow.WorkflowDeleteRequest{
			Namespace: "argo",
			Name:      wfName,
		},
	)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}

func (s workflowService) RerunWorkflow(ctx *gin.Context) error {
	wfName := ctx.Param("name")
	wf, err := s.workflowClient.GetWorkflow(ctx, &workflow.WorkflowGetRequest{
		Namespace: "argo",
		Name:      wfName,
	})
	if err != nil || wf == nil {
		if err == nil {
			err = fmt.Errorf("failed to get workflow: %s", wfName)
		}
		s.logger.Error(err)
		return err
	}
	workflows, err := s.checkRunningOrFailedWorkflows(models_nls.RebuildWorkflowType(wf.Labels["node-type"]))
	if err != nil {
		s.logger.Error(err)
		return err
	}

	if workflows.Len() == 1 && workflows[0].Name != wfName {
		err := fmt.Errorf("another ncn rebuild workflow is still running: %s", workflows[0].Name)
		s.logger.Error(err)
		return err
	}

	_, err = s.workflowClient.ResubmitWorkflow(
		s.ctx,
		&workflow.WorkflowResubmitRequest{
			Namespace: "argo",
			Name:      wfName,
		},
	)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s workflowService) RetryWorkflow(ctx *gin.Context) error {
	wfName := ctx.Param("name")
	wf, err := s.workflowClient.GetWorkflow(ctx, &workflow.WorkflowGetRequest{
		Namespace: "argo",
		Name:      wfName,
	})
	if err != nil || wf == nil {
		if err == nil {
			err = fmt.Errorf("failed to get workflow: %s", wfName)
		}
		s.logger.Error(err)
		return err
	}
	workflows, err := s.checkRunningOrFailedWorkflows(models_nls.RebuildWorkflowType(wf.Labels["node-type"]))
	if err != nil {
		s.logger.Error(err)
		return err
	}

	if workflows.Len() == 1 && workflows[0].Name != wfName {
		err := fmt.Errorf("another ncn rebuild workflow is still running: %s", workflows[0].Name)
		s.logger.Error(err)
		return err
	}

	var requestBody models_nls.RetryWorkflowRequestBody
	if err := ctx.BindJSON(&requestBody); err != nil {
		s.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		ctx.JSON(400, errResponse)
		return err
	}

	_, err = s.workflowClient.RetryWorkflow(
		s.ctx,
		&workflow.WorkflowRetryRequest{
			Namespace:         "argo",
			Name:              wfName,
			RestartSuccessful: requestBody.RestartSuccessful,
			NodeFieldSelector: fmt.Sprintf("name=%s.%s", wfName, requestBody.StepName),
		},
	)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}

func (s workflowService) GetWorkflows(ctx *gin.Context) (*v1alpha1.WorkflowList, error) {
	labelSelector := ctx.Query("labelSelector")
	return s.workflowClient.ListWorkflows(
		s.ctx,
		&workflow.WorkflowListRequest{
			Namespace: "argo",
			ListOptions: &v1.ListOptions{
				LabelSelector: labelSelector,
			},
		},
	)
}

func (s workflowService) GetWorkflowByName(name string, ctx *gin.Context) (*v1alpha1.Workflow, error) {
	return s.workflowClient.GetWorkflow(
		ctx,
		&workflow.WorkflowGetRequest{
			Name:      name,
			Namespace: "argo",
		},
	)
}

func (s workflowService) CreateRebuildWorkflow(req models_nls.CreateRebuildWorkflowRequest) (*v1alpha1.Workflow, error) {
	// support worker rebuild and storage rebuild for now
	workerNodeSet, storageNodeSet := false, false
	var rebuildType models_nls.RebuildWorkflowType
	workerRegEx, err := regexp.Compile(`^ncn-w[0-9]*$`)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	storageRegEx, err := regexp.Compile(`^ncn-s[0-9]*$`)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	for _, hostname := range req.Hosts {
		isWorker := workerRegEx.Match([]byte(hostname))
		if isWorker {
			workerNodeSet = true
			rebuildType = models_nls.WORKER
		}
		isStorage := storageRegEx.Match([]byte(hostname))
		if isStorage {
			storageNodeSet = true
			rebuildType = models_nls.STORAGE
		}
		if !isWorker && !isStorage {
			err = fmt.Errorf("invalid worker or storage node hostname: %s", hostname)
			s.logger.Error(err)
			return nil, err
		}
		// check that hostnames do not contain both worker and storage nodes
		if workerNodeSet && storageNodeSet {
			err = fmt.Errorf("hostnames cannot contain both worker and storage nodes. Only one node type is supported at a time")
			s.logger.Error(err)
			return nil, err
		}
	}

	_, err = s.checkRunningOrFailedWorkflows(rebuildType)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	s.logger.Infof("Creating workflow for: %v", req.Hosts)
	var rebuildWorkflow []byte
	var getWorkflowErr error
	if workerNodeSet {
		rebuildHooks, err := s.getRebuildHooks()
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		// rebuild worker nodes
		workerRebuildWorkflowFS := os.DirFS(s.env.WorkerRebuildWorkflowFiles)
		rebuildWorkflow, getWorkflowErr = argo_templates.GetWorkerRebuildWorkflow(workerRebuildWorkflowFS, req, rebuildHooks)
	} else {
		// rebuild storage nodes
		storageRebuildWorkflowFS := os.DirFS(s.env.StorageRebuildWorkflowFiles)
		rebuildWorkflow, getWorkflowErr = argo_templates.GetStorageRebuildWorkflow(storageRebuildWorkflowFS, req)
	}
	if getWorkflowErr != nil {
		s.logger.Error(getWorkflowErr)
		return nil, getWorkflowErr
	}

	jsonTmp, err := yaml.YAMLToJSONStrict(rebuildWorkflow)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	var myWorkflow v1alpha1.Workflow
	err = json.Unmarshal(jsonTmp, &myWorkflow)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	res, err := s.workflowClient.CreateWorkflow(s.ctx, &workflow.WorkflowCreateRequest{
		Namespace: "argo",
		Workflow:  &myWorkflow,
	})
	if err != nil {
		s.logger.Errorf("Creating workflow for: %v FAILED", req.Hosts)
		s.logger.Error(err)
		return nil, err
	}
	return res, nil
}

func (s workflowService) InitializeWorkflowTemplate(template []byte) error {
	var myWorkflowTemplate v1alpha1.WorkflowTemplate
	tmpBytes, _ := yaml.YAMLToJSON(template)
	err := json.Unmarshal(tmpBytes, &myWorkflowTemplate)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	s.logger.Infof("Initializing workflow template: %s", myWorkflowTemplate.Name)
	for {
		workflowTemplateList, err := s.workflowTemplateClient.ListWorkflowTemplates(s.ctx, &workflowtemplate.WorkflowTemplateListRequest{Namespace: "argo"})
		if err != nil {
			s.logger.Errorf("Failded to get a list of workflow templates: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, workflowTemplate := range workflowTemplateList.Items {
			if workflowTemplate.Name == myWorkflowTemplate.Name && myWorkflowTemplate.ObjectMeta.Labels["version"] != workflowTemplate.ObjectMeta.Labels["version"] {
				s.logger.Info("workflow template has already been initialized")
				s.workflowTemplateClient.DeleteWorkflowTemplate(s.ctx, &workflowtemplate.WorkflowTemplateDeleteRequest{
					Namespace: "argo",
					Name:      workflowTemplate.Name,
				})
				break
			}
		}

		_, err = s.workflowTemplateClient.CreateWorkflowTemplate(
			s.ctx,
			&workflowtemplate.WorkflowTemplateCreateRequest{
				Namespace: "argo",
				Template:  &myWorkflowTemplate,
			})
		if err != nil {
			st := status.Convert(err)
			if st != nil && st.Code() == codes.AlreadyExists {
				err = nil
				break
			}
			// retry
			s.logger.Warnf("Failded to initialize workflow templates: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	s.logger.Infof("Workflow template initialized: %s", myWorkflowTemplate.Name)
	return nil
}

func (s workflowService) checkRunningOrFailedWorkflows(rebuildType models_nls.RebuildWorkflowType) (v1alpha1.Workflows, error) {
	workflows, err := s.workflowClient.ListWorkflows(s.ctx, &workflow.WorkflowListRequest{
		Namespace: "argo",
		ListOptions: &v1.ListOptions{
			LabelSelector: fmt.Sprintf("workflows.argoproj.io/phase!=Succeeded,workflows.argoproj.io/complated!=true,type=rebuild,node-type=%s", rebuildType),
		},
	})
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if workflows.Items.Len() > 1 {
		err := fmt.Errorf("another ncn rebuild workflow (type: %s) is running/failed", rebuildType)
		s.logger.Error(err)
		return nil, err
	}

	return workflows.Items, nil
}

func (s workflowService) getRebuildHooks() (models_nls.RebuildHooks, error) {
	var result models_nls.RebuildHooks
	// get all hooks
	var beforeAllHooks unstructured.UnstructuredList
	beforeAllHooks, err := s.getHooksByLabel("before-all=true")
	if err != nil {
		s.logger.Error(err)
		return result, err
	}
	s.logger.Infof("Before All Hooks: %d", len(beforeAllHooks.Items))
	result.BeforeAll = beforeAllHooks.Items

	var beforeEachHooks unstructured.UnstructuredList
	beforeEachHooks, err = s.getHooksByLabel("before-each=true")
	if err != nil {
		s.logger.Error(err)
		return result, err
	}
	s.logger.Infof("Before Each Hooks: %d", len(beforeEachHooks.Items))
	result.BeforeEach = beforeEachHooks.Items

	var afterEachHooks unstructured.UnstructuredList
	afterEachHooks, err = s.getHooksByLabel("after-each=true")
	if err != nil {
		s.logger.Error(err)
		return result, err
	}
	s.logger.Infof("After Each Hooks: %d", len(afterEachHooks.Items))
	result.AfterEach = afterEachHooks.Items

	var afterAllHooks unstructured.UnstructuredList
	afterAllHooks, err = s.getHooksByLabel("after-all=true")
	if err != nil {
		s.logger.Error(err)
		return result, err
	}
	s.logger.Infof("After All Hooks: %d", len(afterAllHooks.Items))
	result.AfterAll = afterAllHooks.Items

	return result, nil
}

func (s workflowService) getHooksByLabel(label string) (unstructured.UnstructuredList, error) {
	var myHooks unstructured.UnstructuredList
	if s.k8sRestClientSet == nil {
		return myHooks, nil
	}

	beforeAllHooks, err := s.k8sRestClientSet.
		RESTClient().Get().
		AbsPath("/apis/cray-nls.hpe.com/v1").
		Resource("hooks").
		Param("labelSelector", label).
		DoRaw(context.TODO())
	if err != nil {
		s.logger.Error(err)
		return myHooks, err
	}

	err = json.Unmarshal(beforeAllHooks, &myHooks)
	if err != nil {
		s.logger.Error(err)
		return myHooks, err
	}
	return myHooks, nil
}
