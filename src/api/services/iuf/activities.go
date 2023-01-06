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
	"time"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/google/uuid"
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
	name := activity.Name + "-" + uuid.NewString()
	iufHistory := iuf.History{
		ActivityState: iuf.ActivityStateWaitForAdmin,
		StartTime:     int32(time.Now().UnixMilli()),
		Name:          name,
	}
	configmap, err = s.iufObjectToConfigMapData(iufHistory, name, LABEL_HISTORY)
	if err != nil {
		s.logger.Error(err)
		return iuf.Activity{}, err
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

	return activity, err
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
		s.logger.Error(err)
		return iuf.Activity{}, err
	}

	res, err := s.configMapDataToActivity(rawConfigMapData.Data[LABEL_ACTIVITY])
	if err != nil {
		s.logger.Error(err)
		return res, err
	}
	return res, err
}

func (s iufService) PatchActivity(activity iuf.Activity, patchParams iuf.PatchActivityRequest) (iuf.Activity, error) {
	if patchParams.InputParameters.MediaDir != "" {
		// input parameters exists
		activity.InputParameters = patchParams.InputParameters
	}

	// patch site parameters
	if len(patchParams.SiteParameters.Products) > 0 || len(patchParams.SiteParameters.Global) > 0 {
		activity.SiteParameters = patchParams.SiteParameters
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
		s.logger.Error(err)
		return res, err
	}
	return res, err
}

func (s iufService) updateActivity(activity iuf.Activity) (iuf.Activity, error) {
	configmap, err := s.iufObjectToConfigMapData(activity, activity.Name, LABEL_ACTIVITY)
	if err != nil {
		s.logger.Error(err)
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
		s.logger.Error(err)
		return iuf.Activity{}, err
	}
	return activity, err
}
