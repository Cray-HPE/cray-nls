/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022,2025 Hewlett Packard Enterprise Development LP
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
	"sort"
	"time"
	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s iufService) CreateActivity(req iuf.CreateActivityRequest) (iuf.Activity, error) {
	// construct activity object from create req
	reqBytes, _ := json.Marshal(req)
	var activity iuf.Activity
	err := json.Unmarshal(reqBytes, &activity)
	if err != nil {
		s.logger.Error(err)
		return iuf.Activity{}, err
	}

	if activity.Name == "" {
		err := fmt.Errorf("activity name is not set")
		s.logger.Error(err)
		return iuf.Activity{}, err
	}

	// store activity
	configmap, err := s.iufObjectToConfigMapData(activity, activity.Name, LABEL_ACTIVITY)
	if err != nil {
		s.logger.Error(err)
		return iuf.Activity{}, err
	}
	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Create(
			context.TODO(),
			&configmap,
			v1.CreateOptions{},
		)
	if err != nil {
		s.logger.Error(err)
		return iuf.Activity{}, err
	}

	// store history
	err = s.CreateHistoryEntry(activity.Name, iuf.ActivityStateWaitForAdmin, "Activity created")
	if err != nil {
		return iuf.Activity{}, err
	}

	return activity, nil
}

func (s iufService) CreateHistoryEntry(activityName string, activityState iuf.ActivityState, comment string) error {
	name := utils.GenerateName(activityName)
	iufHistory := iuf.History{
		ActivityState: activityState,
		StartTime:     int32(time.Now().UnixMilli()),
		Name:          name,
		Comment:       comment,
	}
	configmap, err := s.iufObjectToConfigMapData(iufHistory, name, LABEL_HISTORY)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	configmap.Labels[LABEL_ACTIVITY_REF] = activityName
	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Create(
			context.TODO(),
			&configmap,
			v1.CreateOptions{},
		)
	if err != nil {
		s.logger.Errorf("CreateHistoryEntry: error when saving history entry in config maps for activity %s and data %#v: %v", activityName, iufHistory, err)
		return err
	}

	return nil
}

func (s iufService) GetActivity(name string) (iuf.Activity, error) {
	rawConfigMapData, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Get(
			context.TODO(),
			name,
			v1.GetOptions{},
		)
	if err != nil {
		s.logger.Errorf("GetActivity.1: An error occurred while trying to get activity %s, %v", name, err)
		return iuf.Activity{}, err
	}

	res, err := s.configMapDataToActivity(rawConfigMapData.Data[LABEL_ACTIVITY])
	if err != nil {
		return res, err
	}
	return res, err
}

func (s iufService) PatchActivity(activity iuf.Activity, patchParams iuf.PatchActivityRequest) (iuf.Activity, error) {
	s.logger.Infof("Called: PatchActivity(activity: %v, patchParams: %v)", activity, patchParams)

	if patchParams.InputParameters.MediaDir != nil {
		activity.InputParameters.MediaDir = *(patchParams.InputParameters.MediaDir)
	}
	if patchParams.InputParameters.SiteParameters != nil {
		activity.InputParameters.SiteParameters = *(patchParams.InputParameters.SiteParameters)
	}
	if patchParams.InputParameters.LimitManagementNodes != nil {
		activity.InputParameters.LimitManagementNodes = *(patchParams.InputParameters.LimitManagementNodes)
	}
	if patchParams.InputParameters.LimitManagedNodes != nil {
		activity.InputParameters.LimitManagedNodes = *(patchParams.InputParameters.LimitManagedNodes)
	}
	if patchParams.InputParameters.ManagementRolloutStrategy != nil {
		activity.InputParameters.ManagementRolloutStrategy = *(patchParams.InputParameters.ManagementRolloutStrategy)
	}
	if patchParams.InputParameters.ManagedRolloutStrategy != nil {
		activity.InputParameters.ManagedRolloutStrategy = *(patchParams.InputParameters.ManagedRolloutStrategy)
	}
	if patchParams.InputParameters.ConcurrentManagementRolloutPercentage != nil {
		activity.InputParameters.ConcurrentManagementRolloutPercentage = *(patchParams.InputParameters.ConcurrentManagementRolloutPercentage)
	}
	if patchParams.InputParameters.MediaHost != nil {
		activity.InputParameters.MediaHost = *(patchParams.InputParameters.MediaHost)
	}
	if patchParams.InputParameters.Concurrency != nil {
		activity.InputParameters.Concurrency = *(patchParams.InputParameters.Concurrency)
	}
	if patchParams.InputParameters.BootprepConfigManaged != nil {
		activity.InputParameters.BootprepConfigManaged = *(patchParams.InputParameters.BootprepConfigManaged)
	}
	if patchParams.InputParameters.BootprepConfigManagement != nil {
		activity.InputParameters.BootprepConfigManagement = *(patchParams.InputParameters.BootprepConfigManagement)
	}
	if patchParams.InputParameters.Stages != nil {
		activity.InputParameters.Stages = *(patchParams.InputParameters.Stages)
	}
	if patchParams.InputParameters.Force != nil {
		activity.InputParameters.Force = *(patchParams.InputParameters.Force)
	}

	// Add new parameters
	if patchParams.InputParameters.CfsConfigurationManagement != nil {
		activity.InputParameters.CfsConfigurationManagement = *(patchParams.InputParameters.CfsConfigurationManagement)
	}
	if patchParams.InputParameters.BootImageManagement != nil {
		activity.InputParameters.BootImageManagement = *(patchParams.InputParameters.BootImageManagement)
	}

	// patch site parameters...all or nothing for Products and Global attributes
	if len(patchParams.SiteParameters.Products) > 0 {
		activity.SiteParameters.Products = patchParams.SiteParameters.Products
	}
	if len(patchParams.SiteParameters.Global) > 0 {
		activity.SiteParameters.Global = patchParams.SiteParameters.Global
	}

	// only allow patching activity state in a limited way.
	switch patchParams.ActivityState {
	case iuf.ActivityStateBlocked:
		switch activity.ActivityState {
		// allow from anything except "in_progress"
		case iuf.ActivityStateInProgress:
			return iuf.Activity{}, utils.GenericError{
				Message: fmt.Sprintf("Illegal activity state transition from %s to %s",
					activity.ActivityState, patchParams.ActivityState)}
		default:
			activity.ActivityState = patchParams.ActivityState
		}
	case iuf.ActivityStatePaused:
		switch activity.ActivityState {
		// allow only from in_progress
		case iuf.ActivityStateInProgress:
			activity.ActivityState = patchParams.ActivityState
		default:
			return iuf.Activity{}, utils.GenericError{
				Message: fmt.Sprintf("Illegal activity state transition from %s to %s",
					activity.ActivityState, patchParams.ActivityState)}
		}
	case "":
		break
	default:
		return iuf.Activity{}, utils.GenericError{
			Message: fmt.Sprintf("Illegal activity state transition from %s to %s",
				activity.ActivityState, patchParams.ActivityState)}
	}

	// when you update site or input parameters of an activity, you also have to update all the Sessions that have not
	// already completed. This is so that the next time a workflow for a stage is created, that workflow can pick up
	// the input and site parameters from the session
	sessions, _ := s.ListSessions(activity.Name)
	for _, session := range sessions {
		if session.CurrentState != iuf.SessionStateCompleted {
			session.InputParameters = activity.InputParameters
			session.SiteParameters = activity.SiteParameters
			err := s.UpdateSession(session)
			if err != nil {
				return iuf.Activity{}, err
			}
		}
	}

	return s.updateActivity(activity)
}

func (s iufService) ListActivities() ([]iuf.Activity, error) {
	rawConfigMapList, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		List(
			context.TODO(),
			v1.ListOptions{
				LabelSelector: fmt.Sprintf("type=%s", LABEL_ACTIVITY),
			},
		)
	if err != nil {
		s.logger.Error(err)
		return []iuf.Activity{}, err
	}

	sort.Slice(rawConfigMapList.Items, func(i, j int) bool {
		return rawConfigMapList.Items[i].CreationTimestamp.Before(&rawConfigMapList.Items[j].CreationTimestamp)
	})

	var res []iuf.Activity
	for _, rawConfigMap := range rawConfigMapList.Items {
		tmp, err := s.configMapDataToActivity(rawConfigMap.Data[LABEL_ACTIVITY])
		if err != nil {
			s.logger.Error(err)
			return []iuf.Activity{}, err
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (s iufService) configMapDataToActivity(data string) (iuf.Activity, error) {
	var res iuf.Activity
	err := json.Unmarshal([]byte(data), &res)
	if err != nil {
		s.logger.Errorf("configMapDataToActivity.1: An error occurred while trying to parse activity configmap %s, %v", data, err)
		return res, err
	}
	return res, err
}

func (s iufService) updateActivity(activity iuf.Activity) (iuf.Activity, error) {
	configmap, err := s.iufObjectToConfigMapData(activity, activity.Name, LABEL_ACTIVITY)
	if err != nil {
		return iuf.Activity{}, err
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
		s.logger.Errorf("updateActivity: error while saving activity %s with %#v: %v", activity.Name, activity, err)
		return iuf.Activity{}, err
	}
	return activity, err
}
