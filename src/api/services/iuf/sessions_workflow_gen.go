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
package services_iuf

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/oliveagle/jsonpath"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
	"reflect"
	"sigs.k8s.io/yaml"
	"sort"
	"strings"
)

const ARGO_TASKS_SIZE_LIMIT = 1
const LABEL_PRODUCT_PREFIX = "product_"
const LABEL_PARTIAL_WORKFLOW = "partial_workflow"

func (s iufService) workflowGen(session *iuf.Session) (workflow v1alpha1.Workflow, err error, skipStage bool) {
	stageName := session.CurrentStage
	if stageName == "" {
		noStageError := utils.GenericError{Message: "No current stage to run."}
		s.logger.Error(noStageError)
		return v1alpha1.Workflow{}, noStageError, false
	}

	stagesMetadata, err := s.GetStages()
	var stageMetadata iuf.Stage
	for _, stage := range stagesMetadata.Stages {
		if stage.Name == stageName {
			stageMetadata = stage
			break
		}
	}
	if stageMetadata.Name == "" {
		err := fmt.Errorf("stage: %s is invalid", stageName)
		s.logger.Error(err)
		return v1alpha1.Workflow{}, err, false
	}
	res := v1alpha1.Workflow{}

	// note that we don't have to care about the length of the prefix here abiding by the 63 character limit because
	//  K8S already trims the prefix accordingly. See
	//  https://github.com/kubernetes/kubernetes/blob/b0b7a323cc5a4a2019b2e9520c21c7830b7f708e/staging/src/k8s.io/apiserver/pkg/storage/names/generate.go#L50
	res.GenerateName = session.Name + "-" + stageName + "-"

	labels := map[string]string{
		"session":    session.Name,
		"activity":   session.ActivityRef,
		"stage":      stageMetadata.Name,
		"stage_type": stageMetadata.Type,
		"iuf":        "true",
	}

	res.Spec.PodMetadata = &v1alpha1.Metadata{
		Labels: map[string]string{
			"iuf":      "true",
			"session":  session.Name,
			"activity": session.ActivityRef,
		},
		Annotations: map[string]string{"sidecar.istio.io/inject": "false"},
	}
	hostPathDir := corev1.HostPathDirectory
	res.Spec.Volumes = []corev1.Volume{
		{
			Name:         "iuf",
			VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: s.env.MediaDirBase, Type: &hostPathDir}},
		},
		{
			Name:         "ssh",
			VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/root/.ssh", Type: &hostPathDir}},
		},
		{
			Name:         "host-usr-bin",
			VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/usr/bin", Type: &hostPathDir}},
		},
		{
			Name:         "ca-bundle",
			VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/lib/ca-certificates", Type: &hostPathDir}},
		},
	}
	res.Spec.PodPriorityClassName = "system-node-critical"

	var concurrency int64 = 10 // default concurrency is 10
	if session.InputParameters.Concurrency > 0 {
		concurrency = session.InputParameters.Concurrency
	}

	res.Spec.Parallelism = &concurrency

	// TODO: commenting this out because devs are finding it confusing why tasks are being retried automatically.
	//retryLimit := intstr.FromInt(3)
	//retryBackoffFactor := intstr.FromInt(2)
	//
	//res.Spec.RetryStrategy = &v1alpha1.RetryStrategy{
	//	Limit:       &retryLimit,
	//	RetryPolicy: v1alpha1.RetryPolicyAlways,
	//	Backoff: &v1alpha1.Backoff{
	//		Duration:    "1m",
	//		Factor:      &retryBackoffFactor,
	//		MaxDuration: "10m",
	//	},
	//}

	res.Spec.PodGC = &v1alpha1.PodGC{Strategy: v1alpha1.PodGCOnPodCompletion}

	// TODO: commenting this out because adding this seems to make it harder to debug
	//var secondsAfterSuccess int32 = 60
	//res.Spec.TTLStrategy = &v1alpha1.TTLStrategy{
	//	SecondsAfterSuccess: &secondsAfterSuccess,
	//}

	res.Spec.Tolerations = []corev1.Toleration{
		{
			Key:      "node-role.kubernetes.io/master",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		},
	}
	res.Spec.Affinity = &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
				{
					Weight: 50,
					Preference: corev1.NodeSelectorTerm{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "node-role.kubernetes.io/master",
								Operator: corev1.NodeSelectorOpExists,
							},
						},
					},
				},
			},
		},
	}
	if !stageMetadata.NoHooks {
		// if we have hooks, then we have to run on ncn-m001. This is a limitation we have for now, because we can only
		// reference hook scripts on ncn-m001 since the rbd mount only exists on ncn-m001.
		// Note that administrator can supply a different media host other than ncn-m001
		if session.InputParameters.MediaHost == "" {
			session.InputParameters.MediaHost = "ncn-m001"
		}

		res.Spec.NodeSelector = map[string]string{"kubernetes.io/hostname": session.InputParameters.MediaHost}
	} else {
		// if we don't have hooks, run this on ncn-m002
		// TODO: we need to find a better way to do this. Perhaps allow specifying the node on which the NoHooks stage
		// 	will run? Not sure.
		res.Spec.NodeSelector = map[string]string{"kubernetes.io/hostname": "ncn-m002"}
	}
	res.Spec.Entrypoint = "main"

	// global stages have product-less global parameters.
	globalParams := s.getGlobalParams(*session, iuf.Product{}, stagesMetadata)
	globalParamsContent, err := json.Marshal(globalParams)
	if err != nil {
		marshalErr := utils.GenericError{Message: fmt.Sprintf("Could not marshal globalParams %v %v", globalParams, err)}
		s.logger.Error(marshalErr)
		return v1alpha1.Workflow{}, marshalErr, false
	}

	globalParamsStr := string(globalParamsContent)
	const globalParamsName = "global_params"

	// Generate global_params for all products in advance
	globalParamsPerProduct := map[string]string{}
	globalParamsNamesPerProduct := map[string]string{}
	for _, product := range session.Products {
		productGlobalParams := s.getGlobalParams(*session, product, stagesMetadata)
		b, err := json.Marshal(productGlobalParams)
		if err != nil {
			marshalErr := utils.GenericError{Message: fmt.Sprintf("Could not marshal globalParams %v %v", productGlobalParams, err)}
			s.logger.Error(marshalErr)
			continue
		}
		productKey := s.getProductVersionKey(product)
		globalParamsPerProduct[productKey] = string(b)
		globalParamsNamesPerProduct[productKey] = productKey
	}

	// generate auth token in advance
	authToken, err := s.keycloakService.NewKeycloakAccessToken()
	if err != nil {
		marshalErr := utils.GenericError{Message: fmt.Sprintf("Could not generate authToken %v", err)}
		s.logger.Error(marshalErr)
		return v1alpha1.Workflow{}, marshalErr, false
	}
	const authTokenName = "auth_token"

	dagTasks, products, err := s.getDAGTasks(session, stageMetadata, stagesMetadata, globalParamsNamesPerProduct, globalParamsName, authTokenName)
	if err != nil {
		s.logger.Error(err)
		return v1alpha1.Workflow{}, err, false
	} else if len(dagTasks) == 0 {
		s.logger.Infof("No DAG tasks for stage %s in session %s, skipping this stage.", stageName, session.Name)
		return v1alpha1.Workflow{}, nil, true
	}

	// Encode products into the labels. This is helpful in tracking what products were composed for each stage.
	for _, product := range products {
		key := LABEL_PRODUCT_PREFIX + product.Name
		value := product.Version
		labels[key] = value
	}

	if len(products) != len(session.Products) {
		labels[LABEL_PARTIAL_WORKFLOW] = "true"

		// update the set of products in the session. Note that the session is a pointer, so this eventually gets saved.
		if session.ProcessedProductsByStage == nil {
			session.ProcessedProductsByStage = make(map[string]map[string]bool)
		}

		processedProducts := session.ProcessedProductsByStage[session.CurrentStage]
		if processedProducts == nil {
			processedProducts = make(map[string]bool)
		}

		for _, product := range products {
			processedProducts[s.getProductVersionKey(product)] = true
		}
		session.ProcessedProductsByStage[session.CurrentStage] = processedProducts
	}

	res.ObjectMeta.Labels = labels

	failFast := false

	res.Spec.Templates = []v1alpha1.Template{
		{
			Name: "main",
			DAG: &v1alpha1.DAGTemplate{
				Tasks:    dagTasks,
				FailFast: &failFast,
			},
		},
	}

	exitHandlers := s.getOnExitHandlers(session, stageMetadata, stagesMetadata.Hooks, globalParamsNamesPerProduct, authTokenName)

	// only run add this field IF this is the last workflow in a set of partial workflows that we need to execute.
	if len(exitHandlers) > 0 && (labels[LABEL_PARTIAL_WORKFLOW] == "" || len(session.ProcessedProductsByStage[session.CurrentStage]) == len(session.Products)) {
		res.Spec.OnExit = "onExitHandlers"
		// list of all the tasks that we picked from the products onExit field
		// only add this template IF this is the last workflow in a set of partial workflows that we need to execute.
		onexit := v1alpha1.Template{
			Name: "onExitHandlers",
			Steps: []v1alpha1.ParallelSteps{
				// note: we do not want to run the exit handlers in parallel, so a sequential list of exit handlers.
				//  This is because certain products like CSM will be upgrading k8s and need to be run in isolation.
				{
					Steps: exitHandlers,
				},
			},
		}
		res.Spec.Templates = append(res.Spec.Templates, onexit)
	}

	var specArgumentsParameters []v1alpha1.Parameter
	for productKey, globalParams := range globalParamsPerProduct {
		param := v1alpha1.Parameter{
			Name:  globalParamsNamesPerProduct[productKey],
			Value: v1alpha1.AnyStringPtr(globalParams),
		}
		specArgumentsParameters = append(specArgumentsParameters, param)
	}

	specArgumentsParameters = append(specArgumentsParameters, v1alpha1.Parameter{
		Name:  authTokenName,
		Value: v1alpha1.AnyStringPtr(authToken),
	})

	specArgumentsParameters = append(specArgumentsParameters, v1alpha1.Parameter{
		Name:  globalParamsName,
		Value: v1alpha1.AnyStringPtr(globalParamsStr),
	})

	res.Spec.Arguments = v1alpha1.Arguments{
		Parameters: specArgumentsParameters,
	}

	return res, nil, false
}

func (s iufService) getOnExitHandlers(session *iuf.Session, stage iuf.Stage,
	hookTemplateMap map[string]string,
	workflowParamNamesGlobalParamsPerProduct map[string]string, workflowParamNameAuthToken string) []v1alpha1.WorkflowStep {

	var workflowSteps []v1alpha1.WorkflowStep

	listTemplates := workflowtemplate.WorkflowTemplateListRequest{
		Namespace: DEFAULT_NAMESPACE,
	}
	templates, err := s.workflowTemplateClient.ListWorkflowTemplates(context.TODO(), &listTemplates)
	if err != nil {
		return []v1alpha1.WorkflowStep{}
	}
	var existingArgoUploadedTemplateMap = map[string]bool{}

	for _, t := range templates.Items {
		existingArgoUploadedTemplateMap[t.Name] = true
	}

	for _, product := range session.Products {
		s.logger.Infof("Processing exit handler for %v - %v", product.Name, product.Version)
		manifest, err := s.getProductManifestAsInterface(product)
		if err != nil {
			s.logger.Errorf("Manifest not found or error while parsing manifest for %v - %v, %v", product.Name, product.Version, err)
			continue
		}

		stageName := strings.Replace(stage.Name, "-", "_", -1)
		jsonPath := fmt.Sprintf("$.onExit.%s.script_path", stageName)
		scriptPathInterface, err := jsonpath.JsonPathLookup(manifest, jsonPath)
		if err != nil || scriptPathInterface == nil {
			s.logger.Debugf("No exit handler found for %v - %v", product.Name, product.Version)
			continue
		}

		executionHost := "master_host"
		jsonExecutionHostPath := fmt.Sprintf("$.onExit.%s.execution_context", stage.Name)
		executionHostInterface, err := jsonpath.JsonPathLookup(manifest, jsonExecutionHostPath)
		if err == nil && executionHostInterface != nil {
			executionHost = fmt.Sprintf("%s", executionHostInterface)
		}

		script := iuf.ManifestHookScript{
			ScriptPath:       fmt.Sprintf("%s", scriptPathInterface),
			ExecutionContext: executionHost,
		}

		// now we assemble the Argo task
		productKey := s.getProductVersionKey(product)

		// find the original location
		var originalLocation string
		for _, product := range session.Products {
			if s.getProductVersionKey(product) == productKey {
				originalLocation = product.OriginalLocation
				break
			}
		}

		originalLocation = filepath.Clean(originalLocation)
		filePath := filepath.Clean(filepath.Join(originalLocation, script.ScriptPath))
		if strings.Index(filePath, originalLocation) != 0 {
			// possible hack attempt ... reading a parent directory through relative paths.
			s.logger.Warnf("Bad hook script path %v found for product %s in stage %s ... reading a parent directory through relative paths.", filePath, productKey, stage.Name)
			continue
		}

		name := utils.GenerateName("onExitHandler" + productKey)

		hookTemplateName := hookTemplateMap[executionHost]

		templateExists := existingArgoUploadedTemplateMap[hookTemplateName]
		if hookTemplateName == "" || !templateExists {
			// this is a backend error so we don't use a template to inform the user here.
			s.logger.Warnf("The template %s is not available in Argo.", hookTemplateName)
			continue
		}

		step := v1alpha1.WorkflowStep{
			Name: name,
			Arguments: v1alpha1.Arguments{
				Parameters: []v1alpha1.Parameter{
					{
						Name:  "auth_token",
						Value: v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.%s}}", workflowParamNameAuthToken)),
					},
					{
						Name:  "global_params",
						Value: v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.%s}}", workflowParamNamesGlobalParamsPerProduct[productKey])),
					},
					{
						Name:  "script_path",
						Value: v1alpha1.AnyStringPtr(filePath),
					},
				},
			},
			TemplateRef: &v1alpha1.TemplateRef{
				Name:     hookTemplateName,
				Template: "main",
			},
		}

		workflowSteps = append(workflowSteps, step)
	}

	return workflowSteps
}

// Gets DAG tasks for the given session and stage
func (s iufService) getDAGTasks(session *iuf.Session, stageInfo iuf.Stage, stages iuf.Stages,
	workflowParamNamesGlobalParamsPerProduct map[string]string, workflowParamNameGlobalParamsForGlobalStage string,
	workflowParamNameAuthToken string) ([]v1alpha1.DAGTask, []iuf.Product, error) {
	var res []v1alpha1.DAGTask
	stage := stageInfo.Name
	s.logger.Infof("getDAGTasks: create workflow DAG for stage %s in session %s in activity %s", stage, session.Name, session.ActivityRef)

	// first find out what templates are available in the system.
	listTemplates := workflowtemplate.WorkflowTemplateListRequest{
		Namespace: DEFAULT_NAMESPACE,
	}
	templates, err := s.workflowTemplateClient.ListWorkflowTemplates(context.TODO(), &listTemplates)
	if err != nil {
		return res, nil, err
	}
	var existingArgoUploadedTemplateMap = map[string]bool{}

	for _, t := range templates.Items {
		existingArgoUploadedTemplateMap[t.Name] = true
	}

	prevStepsSuccessful := map[string]map[string]string{}
	prevStepsAlreadyProcessed := map[string]map[string]bool{}
	for _, product := range session.Products {
		productKey := s.getProductVersionKey(product)
		prevStepsSuccessful[productKey] = make(map[string]string)
		prevStepsAlreadyProcessed[productKey] = make(map[string]bool)
	}

	// we only skip existing operations if force=false AND stage type is product. Note that it is dangerous to skip
	//  stages that are global, because product content may have changed.
	if !session.InputParameters.Force && stageInfo.Type == "product" {
		// go through all the previous sessions of the activity, and see if we can pick up something that is already completed.
		workflows, err := s.workflowClient.ListWorkflows(context.TODO(), &workflow.WorkflowListRequest{
			Namespace: "argo",
			ListOptions: &v1.ListOptions{
				LabelSelector: fmt.Sprintf("activity=%s,stage=%s,iuf=true", session.ActivityRef, stage),
			},
			Fields: "-items.status.nodes,-items.spec",
		})

		if err == nil {
			sort.Slice(workflows.Items, func(i, j int) bool {
				// Note: this is reverse-sort (latest item first)
				return !workflows.Items[i].CreationTimestamp.Before(&workflows.Items[j].CreationTimestamp)
			})

			for _, workflowObjWithName := range workflows.Items {
				workflowObj, err := s.workflowClient.GetWorkflow(context.TODO(), &workflow.WorkflowGetRequest{
					Name:      workflowObjWithName.Name,
					Namespace: "argo",
				})

				if err != nil {
					s.logger.Infof("getDAGTasks: For session %s in activity %s, when generating a DAG for stage %s, an error occurred while checking the previous workflow %s: %v", session.Name, session.ActivityRef, stage, workflowObj.Name, err)
					continue
				}

				//s.logger.Infof("getDAGTasks: For session %s in activity %s, when generating a DAG for stage %s, about to check if previous workflow %s has any successful operations, because force=%v and stage-type=%s...", session.Name, session.ActivityRef, stage, workflowObj.Name, session.InputParameters.Force, stageInfo.Type)

				// for this workflow only, construct a map of previously failed steps so that we can check if grouped
				//  steps have failed
				prevStepsFailedInWorkflow := map[string]map[string]bool{}
				prevStepsSuccessfulInWorkflow := map[string]map[string]string{}
				for _, product := range session.Products {
					productKey := s.getProductVersionKey(product)
					prevStepsFailedInWorkflow[productKey] = make(map[string]bool)
					prevStepsSuccessfulInWorkflow[productKey] = make(map[string]string)
				}

				//s.logger.Infof("getDAGTasks: For session %s in activity %s, when generating a DAG for stage %s, about to check if any of the %v nodes in previous workflow %s have any successful operations, because force=%v and stage-type=%s...", session.Name, session.ActivityRef, stage, len(workflowObj.Status.Nodes), workflowObj.Name, session.InputParameters.Force, stageInfo.Type)

				for _, nodeStatus := range workflowObj.Status.Nodes {
					//s.logger.Infof("getDAGTasks: For session %s in activity %s, when generating a DAG for stage %s, about to check if step %s of operation type %s in previous workflow %s has any successful operations, because force=%v and stage-type=%s...", session.Name, session.ActivityRef, stage, nodeStatus.Name, nodeStatus.TemplateScope, workflowObj.Name, session.InputParameters.Force, stageInfo.Type)
					if strings.HasPrefix(nodeStatus.TemplateScope, "namespaced/") {
						var operationName string

						if strings.Contains(nodeStatus.Name, "-pre-hook-") {
							operationName = "-pre-hook-" + stage
						} else if strings.Contains(nodeStatus.Name, "-post-hook-") {
							operationName = "-post-hook-" + stage
						} else {
							operationName = nodeStatus.TemplateScope[len("namespaced/"):len(nodeStatus.TemplateScope)]
						}

						// go through the products and see which product this belongs to
						for productKey := range prevStepsSuccessfulInWorkflow {

							if strings.Contains(nodeStatus.Name, productKey) {

								if nodeStatus.Phase == v1alpha1.NodeSucceeded {
									// do not join the two ifs in one block -- see note below in else.
									if !prevStepsFailedInWorkflow[productKey][operationName] {
										// if we have determined that previously at least one node in the subgraph of
										//  productKey-operationName has failed, then do not mark this as succeeded.
										prevStepsSuccessfulInWorkflow[productKey][operationName] = workflowObj.Name
									}
								} else {
									// anything other than succeeded needs to be marked as necessary to run.
									// Note that because we are traversing through a DAG, there maybe child steps that have
									//  errors but not the parent steps and vice versa. As such, what we are saying here is that
									//  if any node in the subgraph of a particular productKey-operationName has not succeeded,
									//  then the entire subgraph (i.e. the operation itself) must be retried.
									prevStepsSuccessfulInWorkflow[productKey][operationName] = ""
									prevStepsFailedInWorkflow[productKey][operationName] = true
								}
								break
							}
						}
					}
				}

				// now go through all the failed and successful workflows and mark them as successful or already processed
				for productKey, opMap := range prevStepsSuccessfulInWorkflow {
					for opKey, success := range opMap {
						if success != "" && !prevStepsAlreadyProcessed[productKey][opKey] {
							prevStepsAlreadyProcessed[productKey][opKey] = true
							prevStepsSuccessful[productKey][opKey] = success
							s.logger.Infof("getDAGTasks: For session %s in activity %s, when generating a DAG for stage %s, skipping previously successful operation %s for product %s because force=%v and stage-type=%s", session.Name, session.ActivityRef, stage, opKey, productKey, session.InputParameters.Force, stageInfo.Type)
						}
					}
				}
				for productKey, opMap := range prevStepsFailedInWorkflow {
					for opKey, failed := range opMap {
						if failed && !prevStepsAlreadyProcessed[productKey][opKey] {
							prevStepsAlreadyProcessed[productKey][opKey] = true
							prevStepsSuccessful[productKey][opKey] = ""
							s.logger.Infof("getDAGTasks: For session %s in activity %s, when generating a DAG for stage %s, not going to skip previously unsuccessful operation %s for product %s because force=%v and stage-type=%s", session.Name, session.ActivityRef, stage, opKey, productKey, session.InputParameters.Force, stageInfo.Type)
						}
					}
				}
			}

			for _, product := range session.Products {
				for _, op := range stageInfo.Operations {
					productKey := s.getProductVersionKey(product)
					opKey := op.Name
					if !prevStepsAlreadyProcessed[productKey][opKey] {
						s.logger.Infof("getDAGTasks: For session %s in activity %s, when generating a DAG for stage %s, couldn't determine whether or not to skip operation %s for product %s because force=%v and stage-type=%s", session.Name, session.ActivityRef, stage, opKey, productKey, session.InputParameters.Force, stageInfo.Type)
					}
				}
			}
		} else {
			s.logger.Errorf("getDAGTasks: Got an error while trying to List all workflows for session %s in activity %s, when generating a DAG for stage %s, not attempting to skip previously successful operations because force=%v and stage-type=%s: %v", session.Name, session.ActivityRef, stage, session.InputParameters.Force, stageInfo.Type, err)
		}
	} else {
		s.logger.Infof("getDAGTasks: For session %s in activity %s, when generating a DAG for stage %s, not attempting to skip previously successful operations because force=%v and stage-type=%s", session.Name, session.ActivityRef, stage, session.InputParameters.Force, stageInfo.Type)
	}

	preSteps, postSteps := s.getProductHookTasks(*session, stageInfo, stages, prevStepsSuccessful, existingArgoUploadedTemplateMap, workflowParamNamesGlobalParamsPerProduct, workflowParamNameAuthToken)

	if stageInfo.Type == "product" {
		return s.getDAGTasksForProductStage(*session, s.getRemainingProducts(session), stageInfo, prevStepsSuccessful, existingArgoUploadedTemplateMap, preSteps, postSteps, workflowParamNamesGlobalParamsPerProduct, workflowParamNameAuthToken)
	} else {
		res, err = s.getDAGTasksForGlobalStage(*session, stageInfo, stages, existingArgoUploadedTemplateMap, preSteps, postSteps, workflowParamNameGlobalParamsForGlobalStage, workflowParamNameAuthToken)
		return res, session.Products, err
	}
}

// Gets the DAG tasks for a product stage
func (s iufService) getDAGTasksForProductStage(session iuf.Session,
	products []iuf.Product,
	stageInfo iuf.Stage,
	prevStepsCompleted map[string]map[string]string,
	templateMap map[string]bool,
	preSteps map[string]v1alpha1.DAGTask, postSteps map[string]v1alpha1.DAGTask,
	workflowParamNamesGlobalParamsPerProduct map[string]string, workflowParamNameAuthToken string) (res []v1alpha1.DAGTask, retProducts []iuf.Product, err error) {

	var resPtrs []*v1alpha1.DAGTask

	// this map is to deal with the stageInfo.ProcessProductVariantsSequentially (see docs for that attribute)
	lastOpNamePerProductName := map[string]string{}

	// Here we handle the etcd size limit on resources. Recall that an Argo template is really an etcd resource, and
	//  is constrained to ~1mb. So we have to figure out if the tasks that we are about to submit will go beyond that.
	//  Unfortunately, this is not as straight-forward as `len(json_serialize(tasks)) > 1mb`, because when these tasks
	//  are instantiated into a Workflow, they will carry some metadata (e.g. inputs, outputs, schedule metadata). As
	//  such, we will start with empirically what does work safely (15 products * 10 operations)
	maxProducts := len(products)
	if maxProducts*len(stageInfo.Operations) > ARGO_TASKS_SIZE_LIMIT {
		maxProducts = int(ARGO_TASKS_SIZE_LIMIT / len(stageInfo.Operations))
		s.logger.Infof("Received %s products, but limiting to %s products", len(products), maxProducts)
	}

	for i := 0; i < maxProducts; i++ {
		product := products[i]
		retProducts = append(retProducts, product)
		// the initial dependency is the name of the hook script for that product, if any.
		productKey := s.getProductVersionKey(product)
		preStageHook, exists := preSteps[productKey]
		var lastOpDependency string
		if exists {
			lastOpDependency = preStageHook.Name
			resPtrs = append(resPtrs, &preStageHook)
		}

		isFirstOp := true

		for _, operation := range stageInfo.Operations {

			opName := utils.GenerateName(productKey + "-" + operation.Name)

			task := v1alpha1.DAGTask{
				Name: opName,
			}

			hasEchoTemplate := false

			// do some validations before we are sure to run the operation.
			if prevStepsCompleted[productKey][operation.Name] != "" {
				s.setEchoTemplate(false, &task, fmt.Sprintf("No action needed in %s for product %s because it was previously completed successfully in workflow %s, skipping", operation.Name, productKey, prevStepsCompleted[productKey][operation.Name]))
				hasEchoTemplate = true
			} else if !templateMap[operation.Name] {
				// this is a backend error so we don't use a template to inform the user here.
				s.logger.Errorf("getDAGTasksForProductStage: The template %v cannot be found in Argo. Make sure you have run upload-rebuild-templates.sh from docs-csm", operation.Name)
				continue
			} else {
				manifestBytes := []byte(product.Manifest)
				manifestJsonBytes, err := yaml.YAMLToJSON(manifestBytes)
				if err != nil {
					s.setEchoTemplate(true, &task, fmt.Sprintf("Cannot convert JSON to YAML for product %s while creating a task for operation %s during session %s in activity %s. YAML Manifest: %s. Error: %v", s.getProductVersionKey(product), operation.Name, session.Name, session.ActivityRef, product.Manifest, err))
					hasEchoTemplate = true
				} else {
					var manifestJson map[string]interface{}
					err = json.Unmarshal(manifestJsonBytes, &manifestJson)
					if err != nil {
						s.setEchoTemplate(true, &task, fmt.Sprintf("Cannot parse manifest for product %s while creating a task for operation %s during session %s in activity %s. YAML Manifest: %s. Error: %v", s.getProductVersionKey(product), operation.Name, session.Name, session.ActivityRef, product.Manifest, err))
						hasEchoTemplate = true
					} else if operation.RequiredManifestAttributes != nil && len(operation.RequiredManifestAttributes) > 0 {
						// check if the operation's required manifest attributes are satisfied in the product's manifest
						found := true

						for _, requiredAttributes := range operation.RequiredManifestAttributes {
							attributeHierarchy := strings.Split(requiredAttributes, ".")
							var jsonStruct map[string]interface{}
							jsonStruct = manifestJson
							for _, key := range attributeHierarchy {
								if jsonStruct == nil || jsonStruct[key] == nil {
									found = false
									break
								} else if reflect.TypeOf(jsonStruct[key]).String() == "map[string]interface {}" {
									jsonStruct = jsonStruct[key].(map[string]interface{})
								} else {
									jsonStruct = nil
								}
							}
						}

						if !found {
							s.setEchoTemplate(false, &task, fmt.Sprintf("No action needed in operation %s for product %s, skipping", operation.Name, s.getProductVersionKey(product)))
							hasEchoTemplate = true
						}
					}
				}
			}

			if !hasEchoTemplate {
				task.Arguments = v1alpha1.Arguments{
					Parameters: []v1alpha1.Parameter{
						{
							Name:  "auth_token",
							Value: v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.%s}}", workflowParamNameAuthToken)),
						},
						{
							Name:  "global_params",
							Value: v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.%s}}", workflowParamNamesGlobalParamsPerProduct[productKey])),
						},
					},
				}
				task.TemplateRef = &v1alpha1.TemplateRef{
					Name:     operation.Name,
					Template: "main",
				}
			}

			task.Dependencies = []string{}

			// dep with a stage
			if lastOpDependency != "" {
				task.Dependencies = append(task.Dependencies, lastOpDependency)
			}

			if isFirstOp && stageInfo.ProcessProductVariantsSequentially && lastOpNamePerProductName[product.Name] != "" {
				task.Dependencies = append(task.Dependencies, lastOpNamePerProductName[product.Name])
			}

			isFirstOp = false

			lastOpDependency = opName

			resPtrs = append(resPtrs, &task)
		}

		// add the post-stage hook for this product
		postStageHook, exists := postSteps[productKey]
		if exists {
			if lastOpDependency != "" {
				postStageHook.Dependencies = []string{
					lastOpDependency,
				}
			}

			lastOpDependency = postStageHook.Name
			resPtrs = append(resPtrs, &postStageHook)
		}

		if lastOpDependency != "" {
			lastOpNamePerProductName[product.Name] = lastOpDependency
		}
	}

	for _, step := range resPtrs {
		res = append(res, *step)
	}

	return res, retProducts, nil
}

func (s iufService) setEchoTemplate(isError bool, task *v1alpha1.DAGTask, message string) {
	errorVal := "false"
	if isError {
		errorVal = "true"
	}

	task.Arguments = v1alpha1.Arguments{
		Parameters: []v1alpha1.Parameter{
			{
				Name:  "message",
				Value: v1alpha1.AnyStringPtr(message),
			},
			{
				Name:  "isError",
				Value: v1alpha1.AnyStringPtr(errorVal),
			},
		},
	}
	task.TemplateRef = &v1alpha1.TemplateRef{
		Name:     "echo-template",
		Template: "echo-message",
	}
}

// Gets the DAG tasks for a global stage
func (s iufService) getDAGTasksForGlobalStage(session iuf.Session, stageInfo iuf.Stage, stages iuf.Stages,
	existingArgoUploadedTemplateMap map[string]bool,
	preSteps map[string]v1alpha1.DAGTask, postSteps map[string]v1alpha1.DAGTask,
	workflowParamNameGlobalParamsForGlobalStage string, workflowParamNameAuthToken string) (res []v1alpha1.DAGTask, err error) {

	var lastOpDependencies []string

	for _, product := range session.Products {
		// the initial dependency is the name of the hook script for that product, if any.
		preStageHook, exists := preSteps[s.getProductVersionKey(product)]
		if exists {
			lastOpDependencies = append(lastOpDependencies, preStageHook.Name)
			res = append(res, preStageHook)
		}
	}

	for _, operation := range stageInfo.Operations {

		if !existingArgoUploadedTemplateMap[operation.Name] && operation.Name != "management-nodes-rollout" {
			s.logger.Warnf("The template %v cannot be found in Argo. Make sure you have run upload-rebuild-templates.sh from docs-csm", operation.Name)
			continue
		}

		task := v1alpha1.DAGTask{
			Name:         operation.Name,
			Dependencies: lastOpDependencies,
		}

		lastOpDependencies = []string{operation.Name}

		task.Arguments = v1alpha1.Arguments{
			Parameters: []v1alpha1.Parameter{
				{
					Name:  "auth_token",
					Value: v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.%s}}", workflowParamNameAuthToken)),
				},
				{
					Name:  "global_params",
					Value: v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.%s}}", workflowParamNameGlobalParamsForGlobalStage)),
				},
			},
		}
		if operation.Name == "management-nodes-rollout" {
			managementRolloutSubOperation, err := s.getManagementNodesRolloutSubOperation(session.InputParameters.LimitManagementNodes)
			if err != nil {
				s.setEchoTemplate(true, &task, fmt.Sprintf("Management-nodes-rollout can not be run: %s", err))
			} else {
				if !existingArgoUploadedTemplateMap[managementRolloutSubOperation] {
					s.logger.Warnf("The template %v cannot be found in Argo. Make sure you have run upload-rebuild-templates.sh from docs-csm", managementRolloutSubOperation)
					break
				}
				task.TemplateRef = &v1alpha1.TemplateRef{
					Name:     managementRolloutSubOperation,
					Template: "main",
				}
			}
		} else {
			task.TemplateRef = &v1alpha1.TemplateRef{
				Name:     operation.Name,
				Template: "main",
			}
		}
		res = append(res, task)
	}

	// now let's add all the post-stage hooks
	for _, product := range session.Products {
		// the initial dependency is the name of the hook script for that product, if any.
		postStageHook, exists := postSteps[s.getProductVersionKey(product)]
		if exists {
			postStageHook.Dependencies = lastOpDependencies
			res = append(res, postStageHook)
		}
	}
	return res, nil
}

// Get the master, worker, or storage workflow for management nodes rollout operation
func (s iufService) getManagementNodesRolloutSubOperation(limitManagementNodes []string) (string, error) {
	validator := utils.NewValidator()
	var workflowType string
	workflowType, err := validator.ValidateLimitManagementNodesInput(limitManagementNodes)
	if err != nil {
		return "", err
	}
	workflowNames := map[string]string{
		"worker":      "management-worker-nodes-rollout",
		"storage":     "management-storage-nodes-rollout",
		"master1":     "management-m001-rollout",
		"masterOther": "management-two-master-nodes-rollout",
	}
	if workflowType == "master" {
		if limitManagementNodes[0] == "ncn-m001" {
			workflowType = "master1"
		} else {
			workflowType = "masterOther"
		}
	}
	return workflowNames[workflowType], nil
}

