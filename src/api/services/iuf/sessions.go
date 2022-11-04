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
	"github.com/google/uuid"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
