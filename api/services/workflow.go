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
package services

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

	argo_templates "github.com/Cray-HPE/cray-nls/api/argo-templates"
	"github.com/Cray-HPE/cray-nls/api/models"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/gin-gonic/gin"
)

type WorkflowService interface {
	GetWorkflows(ctx *gin.Context) (*v1alpha1.WorkflowList, error)
	DeleteWorkflow(ctx *gin.Context) error
	RerunWorkflow(ctx *gin.Context) error
	RetryWorkflow(ctx *gin.Context) error
	CreateRebuildWorkflow(req models.CreateRebuildWorkflowRequest) (*v1alpha1.Workflow, error)
	InitializeWorkflowTemplate(template []byte) error
}

// WorkflowService service layer
type workflowService struct {
	logger                utils.Logger
	ctx                   context.Context
	workflowCient         workflow.WorkflowServiceClient
	workflowTemplateCient workflowtemplate.WorkflowTemplateServiceClient
	env                   utils.Env
}

// NewWorkflowService creates a new Workflowservice
func NewWorkflowService(logger utils.Logger, argoService ArgoService, env utils.Env) WorkflowService {

	workflowTemplateCient, _ := argoService.Client.NewWorkflowTemplateServiceClient()

	workflowSvc := workflowService{
		logger:                logger,
		ctx:                   argoService.Context,
		workflowCient:         argoService.Client.NewWorkflowServiceClient(),
		workflowTemplateCient: workflowTemplateCient,
		env:                   env,
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
	workflowToDelete, err := s.workflowCient.GetWorkflow(
		s.ctx,
		&workflow.WorkflowGetRequest{
			Namespace: "argo",
			Name:      wfName,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to find workflow with name: %s", wfName)
	}
	// only delete rebuild workflow
	if workflowToDelete.Labels["type"] != "rebuild" {
		return fmt.Errorf("workflow type is wrong: %s", workflowToDelete.Labels["type"])
	}

	_, err = s.workflowCient.DeleteWorkflow(
		s.ctx,
		&workflow.WorkflowDeleteRequest{
			Namespace: "argo",
			Name:      wfName,
		},
	)
	return err
}

func (s workflowService) RerunWorkflow(ctx *gin.Context) error {
	wfName := ctx.Param("name")
	workflows, err := s.checkRunningWorkflows()
	if err != nil {
		return err
	}

	if workflows.Len() == 1 && workflows[0].Name != wfName {
		return fmt.Errorf("another ncn rebuild workflow is still running: %s", workflows[0].Name)
	}

	_, err = s.workflowCient.ResubmitWorkflow(
		s.ctx,
		&workflow.WorkflowResubmitRequest{
			Namespace: "argo",
			Name:      wfName,
		},
	)
	return err
}

func (s workflowService) RetryWorkflow(ctx *gin.Context) error {
	wfName := ctx.Param("name")
	workflows, err := s.checkRunningWorkflows()
	if err != nil {
		return err
	}

	if workflows.Len() == 1 && workflows[0].Name != wfName {
		return fmt.Errorf("another ncn rebuild workflow is still running: %s", workflows[0].Name)
	}

	var requestBody models.RetryWorkflowRequestBody
	if err := ctx.BindJSON(&requestBody); err != nil {
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		ctx.JSON(400, errResponse)
		return err
	}

	_, err = s.workflowCient.RetryWorkflow(
		s.ctx,
		&workflow.WorkflowRetryRequest{
			Namespace:         "argo",
			Name:              wfName,
			RestartSuccessful: requestBody.RestartSuccessful,
			NodeFieldSelector: fmt.Sprintf("name=%s.%s", wfName, requestBody.StepName),
		},
	)
	return err
}

func (s workflowService) GetWorkflows(ctx *gin.Context) (*v1alpha1.WorkflowList, error) {
	labelSelector := ctx.Query("labelSelector")
	return s.workflowCient.ListWorkflows(
		s.ctx,
		&workflow.WorkflowListRequest{
			Namespace: "argo",
			ListOptions: &v1.ListOptions{
				LabelSelector: labelSelector,
			},
		},
	)
}

func (s workflowService) CreateRebuildWorkflow(req models.CreateRebuildWorkflowRequest) (*v1alpha1.Workflow, error) {
	for _, hostname := range req.Hosts {
		// only support worker rebuild for now
		isWorker, err := regexp.Match(`^ncn-w[0-9]*$`, []byte(hostname))
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		if !isWorker {
			err := fmt.Errorf("only worker nodes rebuild is supported")
			s.logger.Error(err)
			return nil, err
		}
	}

	workflows, err := s.checkRunningWorkflows()
	if err != nil {
		return nil, err
	}

	if workflows.Len() > 0 {
		return nil, fmt.Errorf("another ncn rebuild workflow is still running: %s", workflows[0].Name)
	}

	s.logger.Infof("Creating workflow for: %v", req.Hosts)
	workerRebuildWorkflowFS := os.DirFS(s.env.WorkerRebuildWorkflowFiles)
	workerRebuildWorkflow, err := argo_templates.GetWorkerRebuildWorkflow(workerRebuildWorkflowFS, req)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	jsonTmp, err := yaml.YAMLToJSONStrict(workerRebuildWorkflow)
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

	res, err := s.workflowCient.CreateWorkflow(s.ctx, &workflow.WorkflowCreateRequest{
		Namespace: "argo",
		Workflow:  &myWorkflow,
	})
	if err != nil {
		s.logger.Infof("Creating workflow for: %v FAILED", req.Hosts)
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
	}
	s.logger.Infof("Initializing workflow template: %s", myWorkflowTemplate.Name)
	for {
		workflowTemplateList, err := s.workflowTemplateCient.ListWorkflowTemplates(s.ctx, &workflowtemplate.WorkflowTemplateListRequest{Namespace: "argo"})
		if err != nil {
			s.logger.Errorf("Failded to get a list of workflow templates: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, workflowTemplate := range workflowTemplateList.Items {
			if workflowTemplate.Name == myWorkflowTemplate.Name && myWorkflowTemplate.ObjectMeta.Labels["version"] != workflowTemplate.ObjectMeta.Labels["version"] {
				s.logger.Info("workflow template has already been initialized")
				s.workflowTemplateCient.DeleteWorkflowTemplate(s.ctx, &workflowtemplate.WorkflowTemplateDeleteRequest{
					Namespace: "argo",
					Name:      workflowTemplate.Name,
				})
				break
			}
		}

		_, err = s.workflowTemplateCient.CreateWorkflowTemplate(
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

func (s workflowService) checkRunningWorkflows() (v1alpha1.Workflows, error) {
	workflows, err := s.workflowCient.ListWorkflows(s.ctx, &workflow.WorkflowListRequest{
		Namespace: "argo",
		ListOptions: &v1.ListOptions{
			LabelSelector: "workflows.argoproj.io/phase!=Succeeded,workflows.argoproj.io/complated!=true,type=rebuild",
		},
	})
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if workflows.Items.Len() > 1 {
		err := fmt.Errorf("another ncn rebuild workflow is still running")
		s.logger.Error(err)
		return nil, err
	}

	return workflows.Items, nil
}
