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
	"time"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	yaml_v2 "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func (s iufService) GetSession(sessionName string) (iuf.Session, string, error) {
	rawConfigMapData, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Get(
			context.TODO(),
			sessionName,
			v1.GetOptions{},
		)
	if err != nil {
		s.logger.Error(err)
		return iuf.Session{}, "", err
	}

	res, err := s.ConfigMapDataToSession(rawConfigMapData.Data[LABEL_SESSION])
	if err != nil {
		s.logger.Error(err)
		return res, "", err
	}
	return res, rawConfigMapData.Labels[LABEL_ACTIVITY_REF], err
}

func (s iufService) ListSessions(activityName string) ([]iuf.Session, error) {
	rawConfigMapList, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		List(
			context.TODO(),
			v1.ListOptions{
				LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_SESSION, LABEL_ACTIVITY_REF, activityName),
			},
		)
	if err != nil {
		s.logger.Error(err)
		return []iuf.Session{}, err
	}
	var res []iuf.Session
	for _, rawConfigMap := range rawConfigMapList.Items {
		tmp, err := s.ConfigMapDataToSession(rawConfigMap.Data[LABEL_SESSION])
		if err != nil {
			s.logger.Error(err)
			return []iuf.Session{}, err
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (s iufService) ConfigMapDataToSession(data string) (iuf.Session, error) {
	var res iuf.Session
	err := json.Unmarshal([]byte(data), &res)
	if err != nil {
		s.logger.Error(err)
		return res, err
	}
	return res, err
}

func (s iufService) UpdateSession(session iuf.Session, activityRef string) error {
	configmap, err := s.iufObjectToConfigMapData(session, session.Name, LABEL_SESSION)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	configmap.Labels[LABEL_ACTIVITY_REF] = activityRef
	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Update(
			context.TODO(),
			&configmap,
			v1.UpdateOptions{},
		)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s iufService) UpdateActivityStateFromSessionState(session iuf.Session, activityRef string) error {
	var activityState iuf.ActivityState
	if session.CurrentState == iuf.SessionStateCompleted {
		activityState = iuf.ActivityStateWaitForAdmin
	} else {
		activityState = iuf.ActivityState(session.CurrentState)
	}
	activity, err := s.GetActivity(activityRef)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	activity.ActivityState = activityState
	configmap, err := s.iufObjectToConfigMapData(activity, activity.Name, LABEL_ACTIVITY)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Update(
			context.TODO(),
			&configmap,
			v1.UpdateOptions{},
		)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	// store history
	name := activity.Name + "-" + uuid.NewString()
	iufHistory := iuf.History{
		ActivityState: activityState,
		StartTime:     int32(time.Now().UnixMilli()),
		Name:          name,
		SessionName:   session.Name,
	}
	configmap, err = s.iufObjectToConfigMapData(iufHistory, name, LABEL_HISTORY)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	configmap.Labels[LABEL_ACTIVITY_REF] = activity.Name
	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Create(
			context.TODO(),
			&configmap,
			v1.CreateOptions{},
		)

	return err
}

func (s iufService) CreateIufWorkflow(session iuf.Session) (*v1alpha1.Workflow, error) {
	myWorkflow := s.workflowGen(session)

	res, err := s.workflowCient.CreateWorkflow(context.TODO(), &workflow.WorkflowCreateRequest{
		Namespace: "argo",
		Workflow:  &myWorkflow,
	})
	if err != nil {
		s.logger.Errorf("Creating workflow for: %v FAILED", session)
		s.logger.Error(err)
		return nil, err
	}
	return res, nil
}

func (s iufService) workflowGen(session iuf.Session) v1alpha1.Workflow {
	res := v1alpha1.Workflow{}
	res.GenerateName = session.Name + "-"
	res.ObjectMeta.Labels = map[string]string{"session": session.Name}
	res.Spec.PodMetadata = &v1alpha1.Metadata{Annotations: map[string]string{"sidecar.istio.io/inject": "false"}}
	hostPathDir := corev1.HostPathDirectory
	res.Spec.Volumes = []corev1.Volume{
		{
			Name:         "iuf",
			VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/etc/iuf", Type: &hostPathDir}},
		},
		{
			Name:         "ssh",
			VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/root/.ssh", Type: &hostPathDir}},
		},
		{
			Name:         "host-usr-bin",
			VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/usr/bin", Type: &hostPathDir}},
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
	res.Spec.NodeSelector = map[string]string{"kubernetes.io/hostname": "ncn-m001"}
	res.Spec.Entrypoint = "main"
	//todo: find stage info from stages.yaml
	stageInfo := iuf.Stage{
		Name: "process-media",
		Type: "product",
		Operations: []struct {
			Name      string "json:\"name\""
			LocalPath string "json:\"local_path\""
		}{
			{
				Name:      "extract-release-distributions",
				LocalPath: "operations/extract-release-distributions.yaml",
			},
		},
	}
	res.Spec.Templates = []v1alpha1.Template{
		{
			Name: "main",
			DAG: &v1alpha1.DAGTemplate{
				Tasks: s.getDagTasks(session, stageInfo),
			},
		},
	}
	return res
}

func (s iufService) getDagTasks(session iuf.Session, stageInfo iuf.Stage) []v1alpha1.DAGTask {
	res := []v1alpha1.DAGTask{}
	stage := session.InputParameters.Stages[len(session.Workflows)]
	s.logger.Infof("create DAG for stage: %s", stage)
	if stageInfo.Type == "product" {
		for _, product := range session.Products {
			for index, operation := range stageInfo.Operations {
				task := v1alpha1.DAGTask{
					Name: product.Name + "-" + operation.Name,
				}
				// dep with a stage
				if index != 0 {
					task.Dependencies = []string{
						product.Name + "-" + stageInfo.Operations[index-1].Name,
					}
				}
				stageInputs := s.getStageInputs(session, product)
				b, _ := json.Marshal(stageInputs)
				task.Arguments = v1alpha1.Arguments{
					Parameters: []v1alpha1.Parameter{
						{
							Name:  "auth_token",
							Value: v1alpha1.AnyStringPtr("todo"), // todo token
						},
						{
							Name:  "stage_inputs",
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
		}

	} else {
		s.logger.Infof("TODO: support global stage: %s", stageInfo.Type)
	}
	return res
}

func (s iufService) getStageInputs(session iuf.Session, in_product iuf.Product) map[string]interface{} {
	res := map[string]interface{}{
		"products":     make(map[string]interface{}),
		"input_params": make(map[string]interface{}),
		//"site_params":  make(map[string]interface{}),
	}

	// products
	var productsArray []string
	resProducts := make(map[string]interface{})
	var currentProductManifest map[string]interface{}
	for _, product := range session.Products {
		manifest, _ := s.extractManifestFromTarballFile(product.OriginalLocation)
		manifestBytes, _ := yaml_v2.Marshal(manifest)
		manifestJsonBytes, _ := yaml.YAMLToJSON(manifestBytes)
		var manifestJson map[string]interface{}
		json.Unmarshal(manifestJsonBytes, &manifestJson)
		if product.Name == in_product.Name {
			currentProductManifest = manifestJson
		}
		resProducts[product.Name] = map[string]interface{}{
			"manifest":          manifestJson,
			"original_location": product.OriginalLocation,
		}
		productsArray = append(productsArray, product.Name)
	}
	resProducts["current_product"] = map[string]interface{}{
		"name":              in_product.Name,
		"manifest":          currentProductManifest,
		"original_location": in_product.OriginalLocation,
	}
	res["products"] = resProducts
	res["input_params"] = map[string]interface{}{
		"products":  productsArray,
		"media_dir": session.InputParameters.MediaDir,
		//todo: site_parameters
		//todo: bootprep_config_managed
		//todo: bootprep_config_management
		"limit_nodes": session.InputParameters.LimitNodes,
	}
	//todo: stage params
	return res
}
