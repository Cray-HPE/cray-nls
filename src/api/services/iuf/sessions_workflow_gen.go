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
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
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
	res.GenerateName = stageName + "-" + session.Name + "-"
	res.ObjectMeta.Labels = map[string]string{
		"session":    session.Name,
		"stage":      stageMetadata.Name,
		"stage_type": stageMetadata.Type,
	}
	res.Spec.PodMetadata = &v1alpha1.Metadata{Annotations: map[string]string{"sidecar.istio.io/inject": "false"}}
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
	res.Spec.PodGC = &v1alpha1.PodGC{Strategy: v1alpha1.PodGCOnPodCompletion}
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
		res.Spec.NodeSelector = map[string]string{"kubernetes.io/hostname": "ncn-m001"}
	} else {
		// if we don't have hooks, run this on ncn-m002
		// TODO: we need to find a better way to do this. Perhaps allow specifying the node on which the NoHooks stage
		// 	will run? Not sure.
		res.Spec.NodeSelector = map[string]string{"kubernetes.io/hostname": "ncn-m002"}
	}
	res.Spec.Entrypoint = "main"

	dagTasks, err := s.getDAGTasks(session, stageMetadata, stagesMetadata)
	if err != nil {
		s.logger.Error(err)
		return v1alpha1.Workflow{}, err, false
	} else if len(dagTasks) == 0 {
		s.logger.Infof("No DAG tasks for stage %s in session %s, skipping this stage.", stageName, session.Name)
		return v1alpha1.Workflow{}, nil, true
	}

	res.Spec.Templates = []v1alpha1.Template{
		{
			Name: "main",
			DAG: &v1alpha1.DAGTemplate{
				Tasks: dagTasks,
			},
		},
	}
	return res, nil, false
}

// Gets DAG tasks for the given session and stage
func (s iufService) getDAGTasks(session iuf.Session, stageInfo iuf.Stage, stages iuf.Stages) ([]v1alpha1.DAGTask, error) {
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

	// Generate global_params for all products in advance
	globalParamsPerProduct := map[string][]byte{}
	for _, product := range session.Products {
		globalParams := s.getGlobalParams(session, product, stages)
		b, err := json.Marshal(globalParams)
		if err != nil {
			s.logger.Error(err)
			continue
		}
		globalParamsPerProduct[product.Name] = b
	}

	// generate auth token in advance
	authToken, err := s.keycloakService.NewKeycloakAccessToken()
	if err != nil {
		return []v1alpha1.DAGTask{}, err
	}

	preSteps, postSteps := s.getProductHookTasks(session, stageInfo, stages, existingArgoUploadedTemplateMap, globalParamsPerProduct, authToken)

	if stageInfo.Type == "product" {
		res = s.getDAGTasksForProductStage(session, stageInfo, existingArgoUploadedTemplateMap, preSteps, postSteps, globalParamsPerProduct, authToken, res)
	} else {
		res = s.getDAGTasksForGlobalStage(session, stageInfo, stages, existingArgoUploadedTemplateMap, preSteps, postSteps, authToken, res)
	}

	return res, nil
}

// Gets the DAG tasks for a product stage
func (s iufService) getDAGTasksForProductStage(session iuf.Session, stageInfo iuf.Stage, templateMap map[string]bool,
	preSteps map[string]v1alpha1.DAGTask, postSteps map[string]v1alpha1.DAGTask,
	globalParamsPerProduct map[string][]byte, authToken string,
	res []v1alpha1.DAGTask) []v1alpha1.DAGTask {

	for _, product := range session.Products {

		// the initial dependency is the name of the hook script for that product, if any.
		preStageHook, exists := preSteps[product.Name]
		var lastOpDependency string
		if exists {
			lastOpDependency = preStageHook.Name
			res = append(res, preStageHook)
		}

		for _, operation := range stageInfo.Operations {
			if !templateMap[operation.Name] {
				s.logger.Warnf("The template %v cannot be found in Argo. Make sure you have run upload-rebuild-templates.sh from docs-csm", operation.Name)
				continue
			}

			opName := product.Name + "-" + operation.Name

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
						Value: v1alpha1.AnyStringPtr(authToken),
					},
					{
						Name:  "global_params",
						Value: v1alpha1.AnyStringPtr(string(globalParamsPerProduct[product.Name])),
					},
				},
			}
			task.TemplateRef = &v1alpha1.TemplateRef{
				Name:     operation.Name,
				Template: "main",
			}
			res = append(res, task)
		}

		// add the post-stage hook for this product
		postStageHook, exists := postSteps[product.Name]
		if exists {
			postStageHook.Dependencies = []string{
				lastOpDependency,
			}
			res = append(res, postStageHook)
		}
	}
	return res
}

// Gets the DAG tasks for a global stage
func (s iufService) getDAGTasksForGlobalStage(session iuf.Session, stageInfo iuf.Stage, stages iuf.Stages,
	existingArgoUploadedTemplateMap map[string]bool,
	preSteps map[string]v1alpha1.DAGTask, postSteps map[string]v1alpha1.DAGTask,
	authToken string, res []v1alpha1.DAGTask) []v1alpha1.DAGTask {

	var lastOpDependencies []string

	for _, product := range session.Products {
		// the initial dependency is the name of the hook script for that product, if any.
		preStageHook, exists := preSteps[product.Name]
		if exists {
			lastOpDependencies = append(lastOpDependencies, preStageHook.Name)
			res = append(res, preStageHook)
		}
	}

	// global stages have product-less global parameters.
	globalParams := s.getGlobalParams(session, iuf.Product{}, stages)
	b, _ := json.Marshal(globalParams)

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
					Value: v1alpha1.AnyStringPtr(authToken),
				},
				{
					Name:  "global_params",
					Value: v1alpha1.AnyStringPtr(string(b)),
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
		postStageHook, exists := postSteps[product.Name]
		if exists {
			postStageHook.Dependencies = lastOpDependencies
			res = append(res, postStageHook)
		}
	}
	return res
}
