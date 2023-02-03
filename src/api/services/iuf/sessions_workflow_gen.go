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
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sigs.k8s.io/yaml"
	"strings"
)

func (s iufService) workflowGen(session iuf.Session) (workflow v1alpha1.Workflow, err error, skipStage bool) {
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

	res.ObjectMeta.Labels = map[string]string{
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
	globalParams := s.getGlobalParams(session, iuf.Product{}, stagesMetadata)
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
		productGlobalParams := s.getGlobalParams(session, product, stagesMetadata)
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

	dagTasks, err := s.getDAGTasks(session, stageMetadata, stagesMetadata, globalParamsNamesPerProduct, globalParamsName, authTokenName)
	if err != nil {
		s.logger.Error(err)
		return v1alpha1.Workflow{}, err, false
	} else if len(dagTasks) == 0 {
		s.logger.Infof("No DAG tasks for stage %s in session %s, skipping this stage.", stageName, session.Name)
		return v1alpha1.Workflow{}, nil, true
	}

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

// Gets DAG tasks for the given session and stage
func (s iufService) getDAGTasks(session iuf.Session, stageInfo iuf.Stage, stages iuf.Stages,
	workflowParamNamesGlobalParamsPerProduct map[string]string, workflowParamNameGlobalParamsForGlobalStage string,
	workflowParamNameAuthToken string) ([]v1alpha1.DAGTask, error) {
	var res []v1alpha1.DAGTask
	stage := stageInfo.Name
	s.logger.Infof("create DAG for stage: %s", stage)

	// first find out what templates are available in the system.
	listTemplates := workflowtemplate.WorkflowTemplateListRequest{
		Namespace: DEFAULT_NAMESPACE,
	}
	templates, err := s.workflowTemplateClient.ListWorkflowTemplates(context.TODO(), &listTemplates)
	if err != nil {
		return res, err
	}
	var existingArgoUploadedTemplateMap = map[string]bool{}

	for _, t := range templates.Items {
		existingArgoUploadedTemplateMap[t.Name] = true
	}

	prevStepsCompleted := map[string]map[string]bool{}
	for _, product := range session.Products {
		productKey := s.getProductVersionKey(product)
		prevStepsCompleted[productKey] = make(map[string]bool)
	}

	if !session.InputParameters.Force {
		// go through all the previous workflows of this session and see what steps were completed previously
		for _, workflowId := range session.Workflows {
			workflowObj, err := s.workflowClient.GetWorkflow(context.TODO(), &workflow.WorkflowGetRequest{
				Name:      workflowId.Id,
				Namespace: "argo",
			})
			if err != nil {
				// we don't really care about errors because it should not block us from creating a new DAG.
				continue
			}

			if workflowObj.ObjectMeta.Labels["stage"] != session.CurrentStage {
				// not the current stage? Don't use this then.
				continue
			}

			for _, nodeStatus := range workflowObj.Status.Nodes {
				if nodeStatus.Type == v1alpha1.NodeTypePod &&
					strings.HasPrefix(nodeStatus.TemplateScope, "namespaced/") &&
					nodeStatus.Phase == v1alpha1.NodeSucceeded {
					var operationName string

					if strings.Contains(nodeStatus.Name, "-pre-hook-") {
						operationName = "-pre-hook-" + session.CurrentStage
					} else if strings.Contains(nodeStatus.Name, "-post-hook-") {
						operationName = "-post-hook-" + session.CurrentStage
					} else {
						operationName = nodeStatus.TemplateScope[len("namespaced/"):len(nodeStatus.TemplateScope)]
					}

					// go through the products and see which product this belongs to
					for productKey, opMap := range prevStepsCompleted {

						if strings.Contains(nodeStatus.Name, productKey) {
							opMap[operationName] = true
							prevStepsCompleted[productKey] = opMap
							break
						}
					}
				}
			}
		}
	}

	preSteps, postSteps := s.getProductHookTasks(session, stageInfo, stages, prevStepsCompleted, existingArgoUploadedTemplateMap, workflowParamNamesGlobalParamsPerProduct, workflowParamNameAuthToken)

	if stageInfo.Type == "product" {
		res = s.getDAGTasksForProductStage(session, stageInfo, prevStepsCompleted, existingArgoUploadedTemplateMap, preSteps, postSteps, workflowParamNamesGlobalParamsPerProduct, workflowParamNameAuthToken, res)
	} else {
		res = s.getDAGTasksForGlobalStage(session, stageInfo, stages, existingArgoUploadedTemplateMap, preSteps, postSteps, workflowParamNameGlobalParamsForGlobalStage, workflowParamNameAuthToken, res)
	}

	return res, nil
}

// Gets the DAG tasks for a product stage
func (s iufService) getDAGTasksForProductStage(session iuf.Session, stageInfo iuf.Stage,
	prevStepsCompleted map[string]map[string]bool,
	templateMap map[string]bool,
	preSteps map[string]v1alpha1.DAGTask, postSteps map[string]v1alpha1.DAGTask,
	workflowParamNamesGlobalParamsPerProduct map[string]string, workflowParamNameAuthToken string,
	res []v1alpha1.DAGTask) []v1alpha1.DAGTask {

	var resPtrs []*v1alpha1.DAGTask

	for _, product := range session.Products {
		// the initial dependency is the name of the hook script for that product, if any.
		productKey := s.getProductVersionKey(product)
		preStageHook, exists := preSteps[productKey]
		var lastOpDependency string
		if exists {
			lastOpDependency = preStageHook.Name
			resPtrs = append(resPtrs, &preStageHook)
		}

		for _, operation := range stageInfo.Operations {
			if prevStepsCompleted[productKey][operation.Name] {
				s.logger.Warnf("getDAGTasksForProductStage: Operation %s is being skipped for product %s because it was previously completed successfully in session %s in activity %s", operation.Name, productKey, session.Name, session.ActivityRef)
				continue
			}

			if !templateMap[operation.Name] {
				s.logger.Infof("getDAGTasksForProductStage: The template %v cannot be found in Argo. Make sure you have run upload-rebuild-templates.sh from docs-csm", operation.Name)
				continue
			}

			if operation.RequiredManifestAttributes != nil && len(operation.RequiredManifestAttributes) > 0 {
				// check if the operation's required manifest attributes are satisfied in the product's manifest
				manifestBytes := []byte(product.Manifest)
				manifestJsonBytes, err := yaml.YAMLToJSON(manifestBytes)
				if err != nil {
					s.logger.Warnf("getDAGTasksForProductStage: Cannot convert JSON to YAML for product %s while creating a task for operation %s during session %s in activity %s. YAML Manifest: %s. Error: %v", s.getProductVersionKey(product), operation.Name, session.Name, session.ActivityRef, product.Manifest, err)
					continue
				}
				var manifestJson map[string]interface{}
				err = json.Unmarshal(manifestJsonBytes, &manifestJson)
				if err != nil {
					s.logger.Warnf("getDAGTasksForProductStage: Cannot parse manifest for product %s while creating a task for operation %s during session %s in activity %s. YAML Manifest: %s. Error: %v", s.getProductVersionKey(product), operation.Name, session.Name, session.ActivityRef, product.Manifest, err)
					continue
				}

				found := true

				for _, requiredAttributes := range operation.RequiredManifestAttributes {
					attributeHierarchy := strings.Split(requiredAttributes, ".")
					var jsonStruct map[string]interface{}
					jsonStruct = manifestJson
					for _, key := range attributeHierarchy {
						if jsonStruct == nil || jsonStruct[key] == nil {
							s.logger.Warnf("getDAGTasksForProductStage: Skipping operation %s for product %s in the session %s in activity %s", operation.Name, s.getProductVersionKey(product), session.Name, session.ActivityRef)
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
					continue
				}
			}

			opName := utils.GenerateName(productKey + "-" + operation.Name)

			task := v1alpha1.DAGTask{
				Name: opName,
			}
			// dep with a stage
			if lastOpDependency != "" {
				task.Dependencies = []string{
					lastOpDependency,
				}
			}

			lastOpDependency = opName

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

			resPtrs = append(resPtrs, &postStageHook)
		}
	}

	for _, step := range resPtrs {
		res = append(res, *step)
	}

	return res
}

// Gets the DAG tasks for a global stage
func (s iufService) getDAGTasksForGlobalStage(session iuf.Session, stageInfo iuf.Stage, stages iuf.Stages,
	existingArgoUploadedTemplateMap map[string]bool,
	preSteps map[string]v1alpha1.DAGTask, postSteps map[string]v1alpha1.DAGTask,
	workflowParamNameGlobalParamsForGlobalStage string, workflowParamNameAuthToken string, res []v1alpha1.DAGTask) []v1alpha1.DAGTask {

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
		if !existingArgoUploadedTemplateMap[operation.Name] {
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
		task.TemplateRef = &v1alpha1.TemplateRef{
			Name:     operation.Name,
			Template: "main",
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
	return res
}
