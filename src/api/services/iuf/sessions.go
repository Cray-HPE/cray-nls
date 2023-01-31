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
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"strings"
	"time"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
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

func (s iufService) CreateSession(session iuf.Session, name string, activity iuf.Activity) (iuf.Session, error) {
	configmap, err := s.iufObjectToConfigMapData(session, name, LABEL_SESSION)
	if err != nil {
		s.logger.Error(err)
		return iuf.Session{}, err
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
	return session, err
}

func (s iufService) UpdateSessionAndActivity(session iuf.Session) error {
	err := s.UpdateSession(session)
	if err != nil {
		return err
	}

	// if the session update was successful, we also want to update the activity
	s.logger.Infof("Update activity state, session state: %s", session.CurrentState)
	err = s.UpdateActivityStateFromSessionState(session)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}

func (s iufService) UpdateSession(session iuf.Session) error {
	configmap, err := s.iufObjectToConfigMapData(session, session.Name, LABEL_SESSION)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	configmap.Labels[LABEL_ACTIVITY_REF] = session.ActivityRef
	// set completed label so metacontroller won't sync it again
	if session.CurrentState == iuf.SessionStateCompleted || session.CurrentState == iuf.SessionStateAborted {
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
		// does it even exist? If it doesn't, let's create it instead
		_, err := s.k8sRestClientSet.
			CoreV1().
			ConfigMaps(DEFAULT_NAMESPACE).
			Get(context.TODO(), configmap.Name, v1.GetOptions{})
		if err != nil {
			_, err := s.k8sRestClientSet.
				CoreV1().
				ConfigMaps(DEFAULT_NAMESPACE).
				Create(context.TODO(), &configmap, v1.CreateOptions{})
			if err != nil {
				s.logger.Error(err)
				return err
			}
		} else {
			s.logger.Error(err)
			return err
		}
	}

	return nil
}

func (s iufService) UpdateActivityStateFromSessionState(session iuf.Session) error {
	var activityState iuf.ActivityState
	if session.CurrentState == iuf.SessionStateCompleted || session.CurrentState == iuf.SessionStateAborted {
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
	name := utils.GenerateName(activity.Name)
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

func (s iufService) CreateIufWorkflow(session iuf.Session) (retWorkflow *v1alpha1.Workflow, err error, skipStage bool) {
	myWorkflow, err, skipStage := s.workflowGen(session)
	if err != nil {
		s.logger.Error(err)
		return nil, err, false
	} else if skipStage {
		return nil, nil, true
	}

	res, err := s.workflowClient.CreateWorkflow(context.TODO(), &workflow.WorkflowCreateRequest{
		Namespace: "argo",
		Workflow:  &myWorkflow,
	})
	if err != nil {
		s.logger.Errorf("Creating workflow for: %v FAILED", session)
		s.logger.Error(err)
		return nil, err, false
	}
	return res, nil, false
}

// RunNextStage Runs the next stage in the list of stages to execute.
func (s iufService) RunNextStage(session *iuf.Session) (response iuf.SyncResponse, err error, sessionCompleted bool) {
	// find the current stage in the list of stages, and use the next one
	var currentStage string
	found := false
	if session.CurrentStage != "" {
		for _, stage := range session.InputParameters.Stages {
			if !found {
				if stage == session.CurrentStage {
					found = true
				}
			} else {
				currentStage = stage
				break
			}
		}
	}

	if !found {
		if len(session.InputParameters.Stages) > 0 {
			// Someone updated the input parameters, perhaps. Restart from the beginning because we don't know where we are
			//  anymore
			currentStage = session.InputParameters.Stages[0]
		} else {
			// this session is done because we don't have anything to run
			s.logger.Infof("Session completed. No stages to run")
			return s.SetSessionToCompleted(session)
		}
	} else if currentStage == "" { // we found the last stage
		// this session is done
		return s.SetSessionToCompleted(session)
	}

	stage, err, skipStage := s.RunStage(session, currentStage)
	if skipStage {
		return s.RunNextStage(session)
	} else {
		return stage, err, false
	}
}

func (s iufService) SetSessionToCompleted(session *iuf.Session) (iuf.SyncResponse, error, bool) {
	session.CurrentState = iuf.SessionStateCompleted
	s.logger.Infof("Session completed. Last stage was %s", session.CurrentStage)

	err := s.UpdateSessionAndActivity(*session)
	if err != nil {
		s.logger.Errorf("Error while updating the session %v", err)
		return iuf.SyncResponse{}, err, false
	}

	return iuf.SyncResponse{}, nil, true
}

// RunStage Runs a specific stage for the given session. Creates a new Argo workflow behind the scenes for this stage.
func (s iufService) RunStage(session *iuf.Session, stageToRun string) (ret iuf.SyncResponse, err error, skipStage bool) {
	if stageToRun == "" {
		// this session is done
		s.logger.Infof("No stage specified to run. Last stage was %s and list of all stages are %v",
			session.CurrentStage, session.InputParameters.Stages)
		return iuf.SyncResponse{}, nil, false
	}

	session.CurrentStage = stageToRun
	session.CurrentState = iuf.SessionStateInProgress

	workflow, err, skipStage := s.CreateIufWorkflow(*session)
	if err != nil {
		s.logger.Error(err)
		return iuf.SyncResponse{}, err, skipStage
	} else if !skipStage {
		s.logger.Infof("workflow: %s has been created", workflow.Name)
		session.Workflows = append(session.Workflows, iuf.SessionWorkflow{Id: workflow.Name})
	}

	s.logger.Infof("Update session: %v", session)
	err = s.UpdateSessionAndActivity(*session)
	if err != nil {
		s.logger.Error(err)
		return iuf.SyncResponse{}, err, skipStage
	}

	response := iuf.SyncResponse{
		ResyncAfterSeconds: 5,
	}
	return response, nil, skipStage
}

func (s iufService) ProcessOutput(session *iuf.Session, workflow *v1alpha1.Workflow) error {
	// get activity
	activity, err := s.GetActivity(session.ActivityRef)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	switch workflow.Labels["stage_type"] {
	case "product":
		// first generate a map of all productKeys to Products
		productKeyMap := map[string]iuf.Product{}
		for _, product := range session.Products {
			productKeyMap[s.getProductVersionKey(product)] = product
		}

		// now go through all the nodeStatus items
		changed := false
		for _, nodeStatus := range workflow.Status.Nodes {
			if nodeStatus.Type == v1alpha1.NodeTypePod &&
				strings.HasPrefix(nodeStatus.TemplateScope, "namespaced/") &&
				len(nodeStatus.Outputs.Parameters) > 0 {

				// check which product this is for
				for productKey, _ := range productKeyMap {
					if strings.HasPrefix(nodeStatus.DisplayName, productKey) {
						operationName := nodeStatus.TemplateScope[len("namespaced/"):len(nodeStatus.TemplateScope)]
						stepName := nodeStatus.DisplayName
						s.logger.Infof("process output for Activity %s, Operation %s, step %s with value %v", activity.Name, operationName, stepName, nodeStatus.Outputs)
						stepChanged, err := s.updateActivityOperationOutputFromWorkflow(&activity, session, &nodeStatus, operationName, stepName, productKey)
						if err != nil {
							s.logger.Infof("An error occurred while processing output for Activity %s, Operation %s, step %s with value %v: %v", activity.Name, operationName, stepName, nodeStatus.Outputs, err)
						} else if stepChanged {
							changed = true
						}
					}

					break
				}
			}
		}

		if changed {
			_, err := s.updateActivity(activity)
			return err
		} else {
			return nil
		}
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
			changed := false
			for _, nodeStatus := range workflow.Status.Nodes {
				if nodeStatus.Type == v1alpha1.NodeTypePod &&
					strings.HasPrefix(nodeStatus.TemplateScope, "namespaced/") &&
					len(nodeStatus.Outputs.Parameters) > 0 {
					operationName := nodeStatus.TemplateScope[len("namespaced/"):len(nodeStatus.TemplateScope)]
					stepName := nodeStatus.DisplayName
					s.logger.Infof("process output for Activity %s, Operation %s, step %s with value %v", activity.Name, operationName, stepName, nodeStatus.Outputs)
					stepChanged, err := s.updateActivityOperationOutputFromWorkflow(&activity, session, &nodeStatus, operationName, stepName, "")
					if err != nil {
						s.logger.Infof("An error occurred while processing output for Activity %s, Operation %s, step %s with value %v: %v", activity.Name, operationName, stepName, nodeStatus.Outputs, err)
					} else if stepChanged {
						changed = true
					}
				}
			}

			if changed {
				_, err := s.updateActivity(activity)
				return err
			} else {
				return nil
			}
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
		if manifest["name"] != nil && manifest["version"] != nil {
			// normalize the product version so that we force-follow semver format
			productVersion := s.normalizeProductVersion(fmt.Sprintf("%v", manifest["version"]))
			manifest["version"] = productVersion
			s.logger.Infof("manifest: %s - %s", manifest["name"], manifest["version"])
			// add product to activity object
			activity.Products = append(activity.Products, iuf.Product{
				Name:             fmt.Sprintf("%v", manifest["name"]),
				Version:          productVersion,
				Validated:        validated,
				Manifest:         string(jsonManifest),
				OriginalLocation: nodeStatus.Outputs.Parameters[1].Value.String(),
			})
			productKey := s.getProductVersionKeyFromNameAndVersion(manifest["name"].(string), manifest["version"].(string))

			activity.OperationOutputs["stage_params"].(map[string]interface{})["process-media"].(map[string]interface{})["products"].(map[string]interface{})[fmt.Sprintf("%v", productKey)] = make(map[string]interface{})

			activity.OperationOutputs["stage_params"].(map[string]interface{})["process-media"].(map[string]interface{})["products"].(map[string]interface{})[fmt.Sprintf("%v", productKey)].(map[string]interface{})["parent_directory"] = nodeStatus.Outputs.Parameters[1].Value.String()
		}
	}
	return nil
}

func (s iufService) updateActivityOperationOutputFromWorkflow(
	activity *iuf.Activity,
	session *iuf.Session,
	nodeStatus *v1alpha1.NodeStatus,
	operationName string,
	stepName string,
	productKey string,
) (bool, error) {
	// no-op if there is no outputs
	if nodeStatus.Outputs == nil {
		return false, nil
	}

	changed := false
	if activity.OperationOutputs == nil {
		activity.OperationOutputs = make(map[string]interface{})
	}

	if activity.OperationOutputs["stage_params"] == nil {
		activity.OperationOutputs["stage_params"] = make(map[string]interface{})
	}
	stageParams := activity.OperationOutputs["stage_params"].(map[string]interface{})

	if stageParams[session.CurrentStage] == nil {
		stageParams[session.CurrentStage] = make(map[string]interface{})
	}
	outputStage := stageParams[session.CurrentStage].(map[string]interface{})

	var outputGlobalOrProduct map[string]interface{}

	if productKey != "" {
		if outputStage[productKey] == nil {
			outputStage[productKey] = make(map[string]interface{})
		}
		outputGlobalOrProduct = outputStage[productKey].(map[string]interface{})
	} else {
		outputGlobalOrProduct = outputStage
	}

	if outputGlobalOrProduct[operationName] == nil {
		outputGlobalOrProduct[operationName] = make(map[string]interface{})
	}
	outputOperation := outputGlobalOrProduct[operationName].(map[string]interface{})

	if outputOperation[stepName] == nil {
		outputOperation[stepName] = make(map[string]interface{})
	}
	outputStep := outputOperation[stepName].(map[string]interface{})

	for _, param := range nodeStatus.Outputs.Parameters {
		outputStep[param.Name] = param.Value
		changed = true
	}

	outputOperation[stepName] = outputStep
	outputGlobalOrProduct[operationName] = outputOperation
	if productKey != "" {
		outputStage[productKey] = outputGlobalOrProduct
	} else {
		outputStage = outputGlobalOrProduct
	}

	(activity.OperationOutputs["stage_params"].(map[string]interface{}))[session.CurrentStage] = outputStage

	return changed, nil
}

func (s iufService) PauseSession(session *iuf.Session) error {
	// first, set session and activity to aborted state
	session.CurrentState = iuf.SessionStatePaused

	err := s.UpdateSessionAndActivity(*session)
	if err != nil {
		s.logger.Errorf("PauseSession: An error(s) occurred while setting session %s to aborted: %v", session.Name, err)
		return err
	}

	// now pause the workflows
	var errors []error
	for _, workflowRef := range session.Workflows {
		_, err := s.workflowClient.SuspendWorkflow(context.TODO(), &workflow.WorkflowSuspendRequest{
			Name:      workflowRef.Id,
			Namespace: "argo",
		})

		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		s.logger.Errorf("PauseSession: An error(s) occurred while terminating workflows: %v", errors)
		return errors[0]
	} else {
		return nil
	}
}

func (s iufService) ResumeSession(session *iuf.Session) error {
	// set session and activity to aborted state
	session.CurrentState = iuf.SessionStateInProgress

	err := s.UpdateSessionAndActivity(*session)
	if err != nil {
		s.logger.Errorf("ResumeSession: An error(s) occurred while setting session %s to aborted: %v", session.Name, err)
		return err
	}

	// now resume the workflows
	var errors []error
	for _, workflowRef := range session.Workflows {
		_, err := s.workflowClient.ResumeWorkflow(context.TODO(), &workflow.WorkflowResumeRequest{
			Name:      workflowRef.Id,
			Namespace: "argo",
		})

		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		s.logger.Errorf("ResumeSession: An error(s) occurred while terminating workflows: %v", errors)
		return errors[0]
	} else {
		return nil
	}
}

func (s iufService) AbortSession(session *iuf.Session, force bool) error {
	// first, set session and activity to aborted state
	session.CurrentState = iuf.SessionStateAborted

	err := s.UpdateSessionAndActivity(*session)
	if err != nil {
		s.logger.Errorf("AbortSession: An error(s) occurred while setting session %s to aborted: %v", session.Name, err)
		return err
	}

	if !force {
		return nil
	}

	// now terminate the workflows, so any callbacks right after is correctly ignored because of session aborted state
	var errors []error
	for _, workflowRef := range session.Workflows {
		_, err := s.workflowClient.TerminateWorkflow(context.TODO(), &workflow.WorkflowTerminateRequest{
			Name:      workflowRef.Id,
			Namespace: "argo",
		})

		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		s.logger.Errorf("AbortSession: An error(s) occurred while terminating workflows: %v", errors)

		// we don't want to return error when we had issues terminating
		return nil
	} else {
		return nil
	}
}
