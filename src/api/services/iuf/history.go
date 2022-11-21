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

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/google/uuid"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s iufService) ListActivityHistory(activityName string) ([]iuf.History, error) {
	rawConfigMapList, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		List(
			context.TODO(),
			v1.ListOptions{
				LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_HISTORY, LABEL_ACTIVITY_REF, activityName),
			},
		)
	if err != nil {
		s.logger.Error(err)
		return []iuf.History{}, err
	}
	var res []iuf.History
	for _, rawConfigMap := range rawConfigMapList.Items {
		tmp, err := s.configMapDataToHistory(rawConfigMap.Data[LABEL_HISTORY])
		if err != nil {
			s.logger.Error(err)
			return []iuf.History{}, err
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (s iufService) GetActivityHistory(activityName string, startTime int32) (iuf.History, error) {
	rawConfigMapList, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		List(
			context.TODO(),
			v1.ListOptions{
				LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_HISTORY, LABEL_ACTIVITY_REF, activityName),
			},
		)
	if err != nil {
		s.logger.Error(err)
		return iuf.History{}, err
	}
	var res iuf.History
	for _, rawConfigMap := range rawConfigMapList.Items {
		tmp, err := s.configMapDataToHistory(rawConfigMap.Data[LABEL_HISTORY])
		if err != nil {
			s.logger.Error(err)
			return iuf.History{}, err
		}
		if tmp.StartTime == startTime {
			res = tmp
			break
		}
	}
	return res, nil
}

func (s iufService) ReplaceHistoryComment(activityName string, startTime int32, req iuf.ReplaceHistoryCommentRequest) (iuf.History, error) {
	history, err := s.GetActivityHistory(activityName, startTime)
	if err != nil {
		s.logger.Error(err)
		return iuf.History{}, err
	}
	history.Comment = req.Comment

	// update history
	configmap, err := s.iufObjectToConfigMapData(history, history.Name, LABEL_HISTORY)
	if err != nil {
		s.logger.Error(err)
		return iuf.History{}, err
	}
	configmap.Labels[LABEL_ACTIVITY_REF] = activityName
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
		return iuf.History{}, err
	}
	return history, nil
}

func (s iufService) HistoryRunAction(activityName string, req iuf.HistoryRunActionRequest) (iuf.Session, error) {
	patchReq := iuf.PatchActivityRequest{InputParameters: req.InputParameters}
	activity, err := s.PatchActivity(activityName, patchReq)
	if err != nil {
		s.logger.Error(err)
		return iuf.Session{}, err
	}

	// store session
	name := activity.Name + "-" + uuid.NewString()
	session := iuf.Session{
		InputParameters: activity.InputParameters,
		Products:        activity.Products,
		Name:            name,
		ActivityRef:     activityName,
	}
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

func (s iufService) configMapDataToHistory(data string) (iuf.History, error) {
	var res iuf.History
	err := json.Unmarshal([]byte(data), &res)
	if err != nil {
		s.logger.Error(err)
		return res, err
	}
	return res, err
}
