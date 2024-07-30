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
	core_v1 "k8s.io/api/core/v1"
	"sort"
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
		s.logger.Errorf("ListSessions.1: An error occurred while retrieving list of sessions for activity %s: %v", activityName, err)
		return []iuf.Session{}, err
	}

	sort.Slice(rawConfigMapList.Items, func(i, j int) bool {
		return rawConfigMapList.Items[i].CreationTimestamp.Before(&rawConfigMapList.Items[j].CreationTimestamp)
	})

	var res []iuf.Session
	for _, rawConfigMap := range rawConfigMapList.Items {
		tmp, err := s.ConfigMapDataToSession(rawConfigMap.Data[LABEL_SESSION])
		if err != nil {
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
		s.logger.Errorf("ConfigMapDataToSession: An error occurred while parsing JSON: %s: %v", data, err)
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

func (s iufService) UpdateSessionAndActivity(session iuf.Session, comment string) error {
	err := s.UpdateSession(session)
	if err != nil {
		return err
	}

	// if the session update was successful, we also want to update the activity
	s.logger.Infof("UpdateSessionAndActivity.1: update activity activity %s from session %s with comment %s: %#v", session.ActivityRef, session.Name, comment, session)
	err = s.UpdateActivityStateFromSessionState(session, comment)
	if err != nil {
		return err
	}

	return nil
}

// IsSessionLocked is the session locked by another worker? See LockSession
func (s iufService) IsSessionLocked(session iuf.Session) bool {
	lockName := session.Name + "-lock"
	_, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Get(context.TODO(), lockName, v1.GetOptions{})

	return err == nil
}

// LockSession this is a poor-man's distributed lock. Unfortunately, in the absence of proper distributed caching or some
//
//	transactional database, we are going to have to make do with locking using configmaps.
//	But note that the configmap is stored in etcd, which is eventually consistent :\
func (s iufService) LockSession(session iuf.Session) bool {
	if s.IsSessionLocked(session) {
		return false
	}

	configmap := core_v1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name: session.Name + "-lock",
			Labels: map[string]string{
				"type": LABEL_SESSION_LOCK,
			},
		},
		Data: map[string]string{LABEL_SESSION_LOCK: "true"},
	}

	_, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Create(context.TODO(), &configmap, v1.CreateOptions{})

	if err != nil {
		s.logger.Errorf("LockSession.1: error while creating a lock configmap resource %s %s in activity %s: %v", configmap.Name, session.Name, session.ActivityRef, err)
		return false
	} else {
		return true
	}
}

// UnlockSession unlocks the session. See LockSession
func (s iufService) UnlockSession(session iuf.Session) {
	lockName := session.Name + "-lock"
	err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Delete(context.TODO(), lockName, v1.DeleteOptions{})

	if err != nil {
		s.logger.Errorf("UnlockSession.1: error while deleting a lock configmap resource %s for session %s in activity %s: %v", lockName, session.Name, session.ActivityRef, err)
	} else {
		s.logger.Debugf("UnlockSession.2: Successfully deleted the lock configmap resource %s for session %s in activity %s: %v", lockName, session.Name, session.ActivityRef, err)
	}
}

func (s iufService) UpdateSession(session iuf.Session) error {
	configmap, err := s.iufObjectToConfigMapData(session, session.Name, LABEL_SESSION)
	if err != nil {
		s.logger.Errorf("UpdateSession.1: error while update session %s in activity %s with contents %#v: %v", session.Name, session.ActivityRef, session, err)
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
		_, err2 := s.k8sRestClientSet.
			CoreV1().
			ConfigMaps(DEFAULT_NAMESPACE).
			Get(context.TODO(), configmap.Name, v1.GetOptions{})
		if err2 != nil {
			_, err3 := s.k8sRestClientSet.
				CoreV1().
				ConfigMaps(DEFAULT_NAMESPACE).
				Create(context.TODO(), &configmap, v1.CreateOptions{})
			if err3 != nil {
				s.logger.Errorf("UpdateSession.2: error while creating a new session %s in activity %s with contents %#v: %v", session.Name, session.ActivityRef, session, err3)
				return err
			}
		} else {
			s.logger.Errorf("UpdateSession.3: error while update session %s in activity %s with contents %#v: %v", session.Name, session.ActivityRef, session, err)
			return err
		}
	}

	return nil
}

func (s iufService) UpdateActivityStateFromSessionState(session iuf.Session, comment string) error {
	var activityState iuf.ActivityState
	if session.CurrentState == iuf.SessionStateCompleted || session.CurrentState == iuf.SessionStateAborted {
		activityState = iuf.ActivityStateWaitForAdmin
	} else {
		activityState = iuf.ActivityState(session.CurrentState)
	}
	activity, err := s.GetActivity(session.ActivityRef)
	if err != nil {
		return err
	}

	activity.ActivityState = activityState
	configmap, err := s.iufObjectToConfigMapData(activity, activity.Name, LABEL_ACTIVITY)
	if err != nil {
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
		s.logger.Errorf("UpdateActivityStateFromSessionState.1: An error occurred while trying to save activity %s with contents %#v: %v", activity.Name, activity, err)
		return err
	}

	// store history
	name := utils.GenerateName(activity.Name)
	iufHistory := iuf.History{
		ActivityState: activityState,
		StartTime:     int32(time.Now().UnixMilli()),
		Name:          name,
		SessionName:   session.Name,
		Comment:       comment,
	}
	configmap, err = s.iufObjectToConfigMapData(iufHistory, name, LABEL_HISTORY)
	if err != nil {
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

	if err != nil {
		s.logger.Errorf("UpdateActivityStateFromSessionState.2: An error occurred while trying to save activity %s with contents %#v: %v", activity.Name, activity, err)
	}

	return err
}

func (s iufService) CreateIufWorkflow(session *iuf.Session) (retWorkflow *v1alpha1.Workflow, err error, skipStage bool) {
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

// RunNextPartialWorkflow Runs another workflow for the same stage with the remaining set of products.
func (s iufService) RunNextPartialWorkflow(session *iuf.Session) (response iuf.SyncResponse, err error, sessionCompleted bool) {
	s.logger.Infof("RunNextPartialWorkflow.1: About to run the next partial workflow for session %s in activity %s", session.Name, session.ActivityRef)

	// need to figure out what products are remaining
	remainingProducts := s.getRemainingProducts(session)

	if len(remainingProducts) == 0 {
		s.logger.Infof("RunNextPartialWorkflow.2: Could not find any remaining products for session %s in activity %s, hence will try to proceed to next stage", session.Name, session.ActivityRef)

		// if we don't have any products that are remaining, we have a decision to make. Either proceed to the next stage,
		//  or put the session into DEBUG state. We will only proceed to the next stage if all partial workflows have been
		//  successful. And we will put the session into DEBUG state if any of the workflows had failed or had errors.

		workflows := s.FindAllPartialWorkflowForCurrentStage(session)
		allWorkflowsSuccessful := true
		for _, workflow := range workflows {
			if workflow.Status.Phase != v1alpha1.WorkflowSucceeded {
				allWorkflowsSuccessful = false
				break
			}
		}

		if allWorkflowsSuccessful {
			s.logger.Infof("RunNextPartialWorkflow.3: After no remaining products were found, since all partial workflows for this session %s in activity %s have completed successfully, going to next stage if possible.", session.Name, session.ActivityRef)
			return s.RunNextStage(session)
		}

		s.logger.Infof("RunNextPartialWorkflow.4: After no remaining products were found, since some partial workflow(s) for this session %s in activity %s were not completed successfully, putting session into DEBUG.", session.Name, session.ActivityRef)

		// other workflow(s) have been unsuccessful, so we'll have to mark this as being DEBUG state
		session.CurrentState = iuf.SessionStateDebug
		err = s.UpdateSessionAndActivity(*session, fmt.Sprintf("At least one partial workflow failed %s", workflows[0].Name))
		if err != nil {
			response = iuf.SyncResponse{
				ResyncAfterSeconds: 30,
			}
		} else {
			response = iuf.SyncResponse{}
		}

		return response, nil, true
	}

	s.logger.Infof("RunNextPartialWorkflow.5: Found %#v remaining products for session %s in activity %s, hence will try to run the next partial workflow", len(remainingProducts), session.Name, session.ActivityRef)

	// the run stage will automatically pick up the remaining products.
	return s.RunStage(session, session.CurrentStage)
}

// getRemainingProducts gets the products that have not been processed yet for the current stage.
func (s iufService) getRemainingProducts(session *iuf.Session) []iuf.Product {
	var remainingProducts []iuf.Product
	processedProducts := session.ProcessedProductsByStage[session.CurrentStage]
	for _, product := range session.Products {
		if !processedProducts[s.getProductVersionKey(product)] {
			remainingProducts = append(remainingProducts, product)
		}
	}
	return remainingProducts
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

	err := s.UpdateSessionAndActivity(*session, fmt.Sprintf("Completed %s", session.CurrentStage))
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

	workflow, err, skipStage := s.CreateIufWorkflow(session)
	if err != nil {
		s.logger.Error(err)

		session.CurrentState = iuf.SessionStateDebug
		s.logger.Infof("Update session: %v", session)
		err = s.UpdateSessionAndActivity(*session, fmt.Sprintf("Error in creating workflow %s", err))

		return iuf.SyncResponse{}, err, skipStage
	} else if !skipStage {
		s.logger.Infof("workflow: %s has been created", workflow.Name)
		session.Workflows = append(session.Workflows, iuf.SessionWorkflow{Id: workflow.Name})
	}

	s.logger.Infof("Update session: %v", session)
	err = s.UpdateSessionAndActivity(*session, fmt.Sprintf("Running %s", stageToRun))
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
	switch workflow.ObjectMeta.Labels["stage_type"] {
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
				nodeStatus.Outputs != nil &&
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
		if workflow.ObjectMeta.Labels["stage"] == "process-media" {
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
					nodeStatus.Outputs != nil &&
					len(nodeStatus.Outputs.Parameters) > 0 {
					operationName := nodeStatus.TemplateScope[len("namespaced/"):len(nodeStatus.TemplateScope)]
					stepName := nodeStatus.DisplayName
					s.logger.Infof("output parameter value is %s",nodeStatus.Outputs.Parameters[0].Value)
					if nodeStatus.Outputs.Parameters[0].Value.String() == "" {
						s.logger.Infof("Inside if")
						continue
					}
					s.logger.Infof("process output for Activity %s, Operation %s, step %s with value %v, for stage %s", activity.Name, operationName, stepName, nodeStatus.Outputs,workflow.ObjectMeta.Labels["stage"])
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
		return fmt.Errorf("stage_type: %s is not supported", workflow.ObjectMeta.Labels["stage_type"])
	}

}

func (s iufService) processOutputOfProcessMedia(activity *iuf.Activity, workflow *v1alpha1.Workflow) error {
	nodesWithOutputs := workflow.Status.Nodes.Filter(func(nodeStatus v1alpha1.NodeStatus) bool {
		return nodeStatus.Outputs.HasOutputs() && len(nodeStatus.Outputs.Parameters) == 2
	})
	if len(nodesWithOutputs) == 0 {
		return nil
	}

	if activity.OperationOutputs == nil {
		activity.OperationOutputs = make(map[string]interface{})
	}

	if activity.OperationOutputs["stage_params"] == nil {
		activity.OperationOutputs["stage_params"] = make(map[string]interface{})
	}
	stageParams := activity.OperationOutputs["stage_params"].(map[string]interface{})

	if stageParams["process-media"] == nil {
		stageParams["process-media"] = make(map[string]interface{})
	}
	outputStage := stageParams["process-media"].(map[string]interface{})

	outputStage["products"] = map[string]interface{}{}
	productsMap := outputStage["products"].(map[string]interface{})

	activity.Products = []iuf.Product{}
	for _, nodeStatus := range nodesWithOutputs {
		var manifest map[string]interface{}
		if nodeStatus.Outputs == nil || len(nodeStatus.Outputs.Parameters) == 0 || nodeStatus.Outputs.Parameters[0].Value == nil {
			continue
		}

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

			productsMap[fmt.Sprintf("%v", productKey)] = make(map[string]interface{})

			productsMap[fmt.Sprintf("%v", productKey)].(map[string]interface{})["parent_directory"] = nodeStatus.Outputs.Parameters[1].Value.String()
		}
	}

	outputStage["products"] = productsMap
	stageParams["process-media"] = outputStage
	activity.OperationOutputs["stage_params"] = stageParams

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
		// we skip all output parameters that are marked as "skipped"
		if param.Value != nil && *(param.Value) != "skipped" {
			outputStep[param.Name] = param.Value
			changed = true
		}
	}

	if !changed {
		// fail fast
		return false, nil
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

func (s iufService) PauseSession(session *iuf.Session, comment string) error {
	// first, set session and activity to paused state
	session.CurrentState = iuf.SessionStatePaused

	err := s.UpdateSessionAndActivity(*session, comment)
	if err != nil {
		s.logger.Errorf("PauseSession: An error(s) occurred while setting session %s to Paused: %v", session.Name, err)
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

func (s iufService) ResumeSession(session *iuf.Session, comment string) error {
	var err error
	lastWorkflowIndex := len(session.Workflows) - 1
	if lastWorkflowIndex < 0 {
		// how can there be no workflows for the current stage? Only explanation here
		//  is that workflows were deleted either manually or through forced abort. In this case,
		//  let's go rerun the current stage
		err = utils.GenericError{Message: fmt.Sprintf("ResumeSession.1: There are no workflows to resume for session %s in activity %s: %v", session.Name, session.ActivityRef, err)}
		s.logger.Warn(err)

		return s.RestartCurrentStage(session, comment)
	}

	workflows := s.FindAllPartialWorkflowForCurrentStage(session)
	if len(workflows) == 0 {
		lastWorkflow := s.FindLastWorkflowForCurrentStage(session)
		if lastWorkflow != nil {
			workflows = append(workflows, lastWorkflow)
		}
	}

	if len(workflows) == 0 {
		err = utils.GenericError{Message: fmt.Sprintf("ResumeSession.2: There are no workflows to resume for session %s in activity %s: %v", session.Name, session.ActivityRef, err)}
		s.logger.Warn(err)

		return s.RestartCurrentStage(session, comment)
	}

	allSuccessful := true
	restartCurrentStage := false

	// if there are some workflows, then let's attempt the last failed or error workflow that belongs to the current stage
	for _, lastWorkflow := range workflows {
		if lastWorkflow.Status.Phase == v1alpha1.WorkflowFailed ||
			lastWorkflow.Status.Phase == v1alpha1.WorkflowError {
			allSuccessful = false

			// not successful? retry that error workflow
			_, err = s.workflowClient.RetryWorkflow(context.TODO(), &workflow.WorkflowRetryRequest{
				Name:              lastWorkflow.Name,
				Namespace:         "argo",
				RestartSuccessful: false,
			})

			// if there was an error with retrying, then let's resubmit
			if err != nil {
				s.logger.Errorf("ResumeSession.2: An error occurred while retrying workflow %s in session %s in activity %s: %v. Going to try resubmit instead", lastWorkflow.Name, session.Name, session.ActivityRef, err)
				restartCurrentStage = true
				break

				//  Note: we used to have more code that would try to resubmit the workflow. This would create a new workflow.
				//  However, with the introduction of partial workflows to split up large number of products,
				//  this resubmitting workflows will mess up the determination if the stage has been successful or not,
				//  since that relies upon aggregate status of all the workflows belonging to a stage.
			}
		} else if lastWorkflow.Status.Phase == v1alpha1.WorkflowRunning {
			// try resuming the workflow...if there is an error, that's ok, let it complete on its own
			s.workflowClient.ResumeWorkflow(context.TODO(), &workflow.WorkflowResumeRequest{
				Name:      lastWorkflow.Name,
				Namespace: "argo",
			})
		} // other states are Pending, both for which we do nothing except update session and activity as below
	}

	if allSuccessful {
		// umm...the last workflows were actually successful. We need to go to the next stage instead.
		return s.GotoNextStage(session, comment)
	}

	if restartCurrentStage {
		return s.RestartCurrentStage(session, comment)
	}

	// set session and activity to in progress state
	session.CurrentState = iuf.SessionStateInProgress

	err = s.UpdateSessionAndActivity(*session, comment)
	if err != nil {
		return err
	}

	return nil
}

func (s iufService) FindLastWorkflowForCurrentStage(session *iuf.Session) *v1alpha1.Workflow {
	if len(session.Workflows) == 0 {
		return nil
	}

	var lastWorkflow *v1alpha1.Workflow
	for i := len(session.Workflows) - 1; i >= 0; i-- {
		w := session.Workflows[i]

		lastWorkflowObj, err := s.workflowClient.GetWorkflow(context.TODO(), &workflow.WorkflowGetRequest{
			Name:      w.Id,
			Namespace: "argo",
		})

		if err == nil && lastWorkflowObj != nil &&
			lastWorkflowObj.ObjectMeta.Labels != nil && lastWorkflowObj.ObjectMeta.Labels["stage"] == session.CurrentStage {
			lastWorkflow = lastWorkflowObj
			break
		}
	}
	return lastWorkflow
}

func (s iufService) FindAllPartialWorkflowForCurrentStage(session *iuf.Session) []*v1alpha1.Workflow {
	if len(session.Workflows) == 0 {
		return nil
	}

	var workflows []*v1alpha1.Workflow
	for i := len(session.Workflows) - 1; i >= 0; i-- {
		w := session.Workflows[i]

		lastWorkflowObj, err := s.workflowClient.GetWorkflow(context.TODO(), &workflow.WorkflowGetRequest{
			Name:      w.Id,
			Namespace: "argo",
		})

		if err == nil && lastWorkflowObj != nil &&
			lastWorkflowObj.ObjectMeta.Labels != nil &&
			lastWorkflowObj.ObjectMeta.Labels["stage"] == session.CurrentStage &&
			lastWorkflowObj.ObjectMeta.Labels[LABEL_PARTIAL_WORKFLOW] == "true" {
			workflows = append(workflows, lastWorkflowObj)
		}
	}
	return workflows
}

func (s iufService) GotoNextStage(session *iuf.Session, comment string) error {
	session.CurrentState = ""
	err := s.UpdateSessionAndActivity(*session, comment)
	if err != nil {
		return err
	}

	return nil
}

func (s iufService) RestartCurrentStage(session *iuf.Session, comment string) error {
	if len(session.InputParameters.Stages) == 0 {
		// should have never happened. This is just a bad request.
		err := utils.GenericError{Message: fmt.Sprintf("RestartCurrentStage.1: There are no stages to resume for session %s in activity %s", session.Name, session.ActivityRef)}
		s.logger.Error(err)
		return err
	}

	if session.CurrentStage == "" || session.CurrentStage == session.InputParameters.Stages[0] {
		// if we are still on the current stage, then restart
		session.CurrentState = ""
		session.CurrentStage = ""
		err := s.UpdateSessionAndActivity(*session, comment)
		if err != nil {
			return err
		}

		return nil
	}

	// set the current state to empty so that it can get picked up by the Sync call to RunNextStage.
	session.CurrentState = ""

	// find the previous stage
	session.CurrentStage = session.InputParameters.Stages[0]

	for i := 1; i < len(session.InputParameters.Stages); i++ {
		if session.CurrentStage == session.InputParameters.Stages[i] {
			break
		} else {
			session.CurrentStage = session.InputParameters.Stages[i]
		}
	}

	err := s.UpdateSessionAndActivity(*session, comment)
	if err != nil {
		return err
	}

	return nil
}

func (s iufService) AbortSession(session *iuf.Session, comment string, force bool, workflowList *v1alpha1.WorkflowList) error {
	// first, set session and activity to aborted state
	session.CurrentState = iuf.SessionStateAborted

	err := s.UpdateSessionAndActivity(*session, comment)
	if err != nil {
		s.logger.Errorf("AbortSession: An error(s) occurred while setting session %s to aborted: %v", session.Name, err)
		return err
	}

	// only do the next part if force=true. If force=false, we want the stage to finish whatever it was doing first.
	if !force {
		return nil
	}

	// now terminate the workflows, so any callbacks right after is correctly ignored because of session aborted state
	var errors []error
	var workflowIDsToCheck []string
	for _, workflowObj := range workflowList.Items {
		_, err := s.workflowClient.TerminateWorkflow(context.TODO(), &workflow.WorkflowTerminateRequest{
			Name:      workflowObj.Name,
			Namespace: "argo",
		})

		if err != nil {
			// delete the workflow right away.
			_, err := s.workflowClient.DeleteWorkflow(context.TODO(), &workflow.WorkflowDeleteRequest{
				Name:      workflowObj.Name,
				Namespace: "argo",
			})

			if err != nil {
				errors = append(errors, err)
			}
		} else {
			workflowIDsToCheck = append(workflowIDsToCheck, workflowObj.Name)
		}
	}

	// do a check again before going for a more aggressive delete workflow option.
	terminatedAll := true
	for _, workflowToCheckId := range workflowIDsToCheck {
		workflowToCheck, err := s.workflowClient.GetWorkflow(context.TODO(), &workflow.WorkflowGetRequest{
			Name:      workflowToCheckId,
			Namespace: "argo",
			Fields:    "status.phase",
		})

		if err == nil && (workflowToCheck.Status.Phase == v1alpha1.WorkflowPending || workflowToCheck.Status.Phase == v1alpha1.WorkflowRunning) {
			terminatedAll = false
		}
	}

	if terminatedAll {
		return nil
	}

	// wait 30 seconds before checking that all workflows have in fact been terminated.
	time.Sleep(30 * time.Second)

	for _, workflowToCheckId := range workflowIDsToCheck {
		workflowToCheck, err := s.workflowClient.GetWorkflow(context.TODO(), &workflow.WorkflowGetRequest{
			Name:      workflowToCheckId,
			Namespace: "argo",
			Fields:    "status.phase",
		})

		if err != nil || workflowToCheck.Status.Phase == v1alpha1.WorkflowPending || workflowToCheck.Status.Phase == v1alpha1.WorkflowRunning {
			// good candidate to nuke the workflow.
			_, err := s.workflowClient.DeleteWorkflow(context.TODO(), &workflow.WorkflowDeleteRequest{
				Name:      workflowToCheckId,
				Namespace: "argo",
			})

			if err != nil {
				errors = append(errors, err)
			}
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

func (s iufService) SyncWorkflowsToSession(session *iuf.Session) error {
	workflows, err := s.workflowClient.ListWorkflows(context.TODO(), &workflow.WorkflowListRequest{
		Namespace: "argo",
		ListOptions: &v1.ListOptions{
			LabelSelector: fmt.Sprintf("session=%s,iuf=true", session.Name),
		},
		Fields: "-items.spec,-items.status",
	})
	if err != nil {
		s.logger.Errorf("SyncWorkflowsToSession.1: An error occurred while retrieving list of workflows for session %s in activity %s: %v", session.Name, session.ActivityRef, err)
		return err
	}

	sort.Slice(workflows.Items, func(i, j int) bool {
		return workflows.Items[i].CreationTimestamp.Before(&workflows.Items[j].CreationTimestamp)
	})

	// now we make sure that the stored workflows in the session are in sync with the workflows
	needsSync := len(session.Workflows) != len(workflows.Items)
	if len(workflows.Items) > 0 {
		for i, workflowObj := range workflows.Items {
			if i >= len(session.Workflows) || workflowObj.Name != session.Workflows[i].Id {
				needsSync = true
				break
			}
		}
	}

	if needsSync {
		session.Workflows = []iuf.SessionWorkflow{}
		for _, workflowObj := range workflows.Items {
			session.Workflows = append(session.Workflows, iuf.SessionWorkflow{Id: workflowObj.Name})
		}

		// try to update the session and ignore errors because this is meant to be eventually persistent
		s.UpdateSession(*session)
	}

	return nil
}
