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
	"path"
	"time"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	"github.com/imdario/mergo"
	"golang.org/x/exp/slices"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func (s iufService) GetSession(sessionName string) (iuf.Session, error) {
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
		return iuf.Session{}, err
	}

	res, err := s.ConfigMapDataToSession(rawConfigMapData.Data[LABEL_SESSION])
	if err != nil {
		s.logger.Error(err)
		return res, err
	}
	return res, err
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

func (s iufService) UpdateSession(session iuf.Session) error {
	configmap, err := s.iufObjectToConfigMapData(session, session.Name, LABEL_SESSION)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	configmap.Labels[LABEL_ACTIVITY_REF] = session.ActivityRef
	// set completed label so metacontroller won't sync it again
	if session.CurrentState == iuf.SessionStateCompleted {
		configmap.Labels["completed"] = "true"
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
	return nil
}

func (s iufService) UpdateActivityStateFromSessionState(session iuf.Session) error {
	var activityState iuf.ActivityState
	if session.CurrentState == iuf.SessionStateCompleted {
		activityState = iuf.ActivityStateWaitForAdmin
	} else {
		activityState = iuf.ActivityState(session.CurrentState)
	}
	activity, err := s.GetActivity(session.ActivityRef)
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
	myWorkflow, err := s.workflowGen(session)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

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

func (s iufService) workflowGen(session iuf.Session) (v1alpha1.Workflow, error) {
	stages, err := s.GetStages()
	if err != nil {
		s.logger.Error(err)
		return v1alpha1.Workflow{}, err
	}
	stageName := session.InputParameters.Stages[len(session.Workflows)]
	var stageInfo iuf.Stage
	for _, stage := range stages.Stages {
		if stage.Name == stageName {
			stageInfo = stage
			break
		}
	}
	if stageInfo.Name == "" {
		err := fmt.Errorf("stage: %s is invalid", stageName)
		s.logger.Error(err)
		return v1alpha1.Workflow{}, err
	}
	res := v1alpha1.Workflow{}
	res.GenerateName = session.Name + "-"
	res.ObjectMeta.Labels = map[string]string{
		"session":    session.Name,
		"stage":      stageInfo.Name,
		"stage_type": stageInfo.Type,
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
	res.Spec.NodeSelector = map[string]string{"kubernetes.io/hostname": "ncn-m001"}
	res.Spec.Entrypoint = "main"

	dagTasks, err := s.getDagTasks(session, stageInfo)
	if err != nil {
		s.logger.Error(err)
		return v1alpha1.Workflow{}, err
	}

	res.Spec.Templates = []v1alpha1.Template{
		{
			Name: "main",
			DAG: &v1alpha1.DAGTemplate{
				Tasks: dagTasks,
			},
		},
	}
	return res, nil
}

func (s iufService) RunNextStage(session *iuf.Session) (iuf.SyncResponse, error) {
	// get list of stages
	stages := session.InputParameters.Stages
	workflow, err := s.CreateIufWorkflow(*session)
	if err != nil {
		s.logger.Error(err)
		return iuf.SyncResponse{}, err
	}
	s.logger.Infof("workflow: %s has been created", workflow.Name)

	session.Workflows = append(session.Workflows, iuf.SessionWorkflow{Id: workflow.Name})
	session.CurrentStage = stages[len(session.Workflows)-1]
	session.CurrentState = iuf.SessionStateInProgress
	s.logger.Infof("Update activity state, session state: %s", session.CurrentState)
	err = s.UpdateActivityStateFromSessionState(*session)
	if err != nil {
		s.logger.Error(err)
		return iuf.SyncResponse{}, err
	}
	s.logger.Infof("Update session: %v", session)
	err = s.UpdateSession(*session)
	if err != nil {
		s.logger.Error(err)
		return iuf.SyncResponse{}, err
	}
	response := iuf.SyncResponse{
		ResyncAfterSeconds: 5,
	}
	return response, nil
}

func (s iufService) ProcessOutput(session *iuf.Session, workflow *v1alpha1.Workflow) error {
	// get activity
	activity, err := s.GetActivity(session.ActivityRef)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	// get tasks we care about (top level dag)
	tasks := workflow.Spec.Templates[0].DAG.Tasks
	switch workflow.Labels["stage_type"] {
	case "product":
		for _, task := range tasks {
			operationName := task.TemplateRef.Name
			nodeStatus := workflow.Status.Nodes.FindByDisplayName(task.Name)
			var productName string
			for _, param := range nodeStatus.Inputs.Parameters {
				if param.Name == "global_params" {
					var valueJson map[string]interface{}
					json.Unmarshal([]byte(param.Value.String()), &valueJson)
					productManifest := valueJson["product_manifest"].(map[string]interface{})
					currentProduct := productManifest["current_product"].(map[string]interface{})
					manifest := currentProduct["manifest"].(map[string]interface{})
					productName = manifest["name"].(string)
					break
				}
			}
			s.logger.Infof("process output of: %s, product: %s, %v", operationName, productName, nodeStatus.Outputs)
			s.updateActivityOperationOutputFromWorkflow(activity, *session, nodeStatus, operationName, productName)
		}
		return nil
	case "global":
		// special handling of process media
		if workflow.Labels["stage"] == "process-media" {
			err := s.processOutputOfProcessMedia(&activity, workflow)
			if err != nil {
				s.logger.Error(err)
				return err
			}
			session.Products = activity.Products
			// update activity
			_, err = s.updateActivity(activity)
			if err != nil {
				s.logger.Error(err)
				return err
			}
			return nil
		} else {
			for _, task := range tasks {
				operationName := task.TemplateRef.Name
				nodeStatus := workflow.Status.Nodes.FindByDisplayName(task.Name)
				s.logger.Infof("process output of: %s, %v", operationName, nodeStatus.Outputs)
				s.updateActivityOperationOutputFromWorkflow(activity, *session, nodeStatus, operationName, "")
			}
			return nil
		}
	default:
		return fmt.Errorf("stage_type: %s is not supported", workflow.Labels["stage_type"])
	}

}

func (s iufService) processOutputOfProcessMedia(activity *iuf.Activity, workflow *v1alpha1.Workflow) error {
	nodesWithOutputs := workflow.Status.Nodes.Filter(func(nodeStatus v1alpha1.NodeStatus) bool {
		return nodeStatus.Outputs.HasOutputs() && len(nodeStatus.Outputs.Parameters) == 2
	})
	if len(nodesWithOutputs) == 0 {
		return nil
	}
	activity.OperationOutputs = map[string]interface{}{
		"stage_params": map[string]interface{}{
			"process-media": map[string]interface{}{
				"products": map[string]interface{}{},
			},
		},
	}
	activity.Products = []iuf.Product{}
	for _, nodeStatus := range nodesWithOutputs {
		var manifest map[string]interface{}
		err := yaml.Unmarshal([]byte(nodeStatus.Outputs.Parameters[0].Value.String()), &manifest)
		if err != nil {
			s.logger.Error(err)
			return err
		}
		// validate iuf product manifest
		data, _ := yaml.Marshal(manifest)
		validated := true
		err = iuf.Validate(data)
		if err != nil {
			s.logger.Error(err)
			validated = false
		}
		jsonManifest, _ := json.Marshal(manifest)
		if manifest["name"] != nil {
			s.logger.Infof("manifest: %s - %s", manifest["name"], manifest["version"])
			// add product to activity object
			activity.Products = append(activity.Products, iuf.Product{
				Name:             fmt.Sprintf("%v", manifest["name"]),
				Version:          fmt.Sprintf("%v", manifest["version"]),
				Validated:        validated,
				Manifest:         string(jsonManifest),
				OriginalLocation: nodeStatus.Outputs.Parameters[1].Value.String(),
			})
			activity.OperationOutputs["stage_params"].(map[string]interface{})["process-media"].(map[string]interface{})["products"].(map[string]interface{})[fmt.Sprintf("%v", manifest["name"])] = make(map[string]interface{})
			activity.OperationOutputs["stage_params"].(map[string]interface{})["process-media"].(map[string]interface{})["products"].(map[string]interface{})[fmt.Sprintf("%v", manifest["name"])].(map[string]interface{})["parent_directory"] = nodeStatus.Outputs.Parameters[1].Value.String()
		}
	}
	return nil
}

func (s iufService) updateActivityOperationOutputFromWorkflow(
	activity iuf.Activity,
	session iuf.Session,
	nodeStatus *v1alpha1.NodeStatus,
	operationName string,
	productName string,
) error {
	// no-op if there is no outputs
	if nodeStatus.Outputs == nil {
		return nil
	}
	if activity.OperationOutputs == nil {
		activity.OperationOutputs = make(map[string]interface{})
	}
	if activity.OperationOutputs[session.CurrentStage] == nil {
		activity.OperationOutputs[session.CurrentStage] = make(map[string]interface{})
	}
	outputStage := activity.OperationOutputs[session.CurrentStage].(map[string]interface{})
	if outputStage[operationName] == nil {
		outputStage[operationName] = make(map[string]interface{})
	}
	outputOperation := outputStage[operationName].(map[string]interface{})
	if productName != "" {
		if outputOperation[productName] == nil {
			outputOperation[productName] = make(map[string]interface{})
		}
		operationOutputOfProduct := outputOperation[productName].(map[string]interface{})
		for _, param := range nodeStatus.Outputs.Parameters {
			operationOutputOfProduct[param.Name] = param.Value
		}

	} else {
		for _, param := range nodeStatus.Outputs.Parameters {
			outputOperation[param.Name] = param.Value
		}
	}
	activity.OperationOutputs[session.CurrentStage] = outputStage
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
	return nil
}

func (s iufService) getDagTasks(session iuf.Session, stageInfo iuf.Stage) ([]v1alpha1.DAGTask, error) {
	res := []v1alpha1.DAGTask{}
	stage := stageInfo.Name
	s.logger.Infof("create DAG for stage: %s", stage)

	authToken, err := s.keycloakService.NewKeycloakAccessToken()
	if err != nil {
		return []v1alpha1.DAGTask{}, err
	}

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
				globaParams := s.getGlobalParams(session, product)
				b, _ := json.Marshal(globaParams)
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
		}
	} else {
		for index, operation := range stageInfo.Operations {
			task := v1alpha1.DAGTask{
				Name: operation.Name,
			}
			if index != 0 {
				task.Dependencies = []string{
					stageInfo.Operations[index-1].Name,
				}
			}
			globaParams := s.getGlobalParams(session, iuf.Product{})
			b, _ := json.Marshal(globaParams)
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
	}

	return res, nil
}

func (s iufService) getGlobalParams(session iuf.Session, in_product iuf.Product) map[string]interface{} {
	return map[string]interface{}{
		"product_manifest": s.getGlobalParamsProductManifest(session, in_product),
		"input_params":     s.getGlobalParamsInputParams(session, in_product),
		"site_params":      s.getGlobalParamsSiteParams(session, in_product),
		"stage_params":     s.getGlobalParamsStageParams(session, in_product),
	}
}

func (s iufService) getGlobalParamsProductManifest(session iuf.Session, in_product iuf.Product) map[string]interface{} {
	resProducts := make(map[string]interface{})
	var currentProductManifest map[string]interface{}
	for _, product := range session.Products {
		manifestBytes := []byte(product.Manifest)
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
	}
	return map[string]interface{}{
		"products": resProducts,
		"current_product": map[string]interface{}{
			"name":              in_product.Name,
			"manifest":          currentProductManifest,
			"original_location": in_product.OriginalLocation,
		},
	}
}

func (s iufService) getGlobalParamsInputParams(session iuf.Session, in_product iuf.Product) map[string]interface{} {
	var productsArray []string
	for _, product := range session.Products {
		productsArray = append(productsArray, product.Name)
	}
	return map[string]interface{}{
		"products":  productsArray,
		"media_dir": path.Join(s.env.MediaDirBase, session.InputParameters.MediaDir),
		//todo: bootprep_config_managed
		//todo: bootprep_config_management
		"limit_nodes": session.InputParameters.LimitNodes,
	}
}

func (s iufService) getGlobalParamsStageParams(session iuf.Session, in_product iuf.Product) map[string]interface{} {
	res := make(map[string]interface{})
	activity, _ := s.GetActivity(session.ActivityRef)
	if activity.OperationOutputs == nil {
		return map[string]interface{}{}
	}
	stages, _ := s.GetStages()
	stageParams := activity.OperationOutputs["stage_params"].(map[string]interface{})
	// loop through each stage's output
	for stageName, v := range stageParams {
		idx := slices.IndexFunc(stages.Stages, func(stage iuf.Stage) bool { return stage.Name == stageName })
		stageType := stages.Stages[idx].Type
		outputValue := v.(map[string]interface{})
		res[stageName] = make(map[string]interface{})
		s.logger.Debugf("stage: %s, type: %s, outputs: %v", stageName, stageType, v)
		if stageType == "product" || stageName == "process-media" {
			var currentProduct map[string]interface{}
			var products map[string]interface{}
			for _, value := range outputValue {
				mergo.Merge(&products, value.(map[string]interface{}))
				mergo.Merge(&currentProduct, value.(map[string]interface{})[in_product.Name])
			}
			res[stageName].(map[string]interface{})["products"] = products
			res[stageName].(map[string]interface{})["current_product"] = currentProduct
		} else {
			res[stageName].(map[string]interface{})["global"] = outputValue
		}

	}
	return res
}

func (s iufService) getGlobalParamsSiteParams(session iuf.Session, in_product iuf.Product) map[string]interface{} {
	//todo: site_parameters
	return map[string]interface{}{}
}
