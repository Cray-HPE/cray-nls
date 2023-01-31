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

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
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
	activity, err := s.GetActivity(activityName)
	if err != nil {
		s.logger.Error(err)
		return iuf.Session{}, err
	}

	activity, err = s.PatchActivity(activity, iuf.PatchActivityRequest{
		InputParameters: req.InputParameters,
		SiteParameters:  req.SiteParameters,
	})
	if err != nil {
		s.logger.Error(err)
		return iuf.Session{}, err
	}

	// store session
	name := utils.GenerateName(activity.Name)
	session := iuf.Session{
		InputParameters: activity.InputParameters,
		SiteParameters:  s.getSiteParams(activity.InputParameters.SiteParameters, activity.SiteParameters),
		Products:        activity.Products,
		Name:            name,
		ActivityRef:     activityName,
	}
	return s.CreateSession(session, name, activity)
}

func (s iufService) HistoryAbortAction(activityName string, req iuf.HistoryAbortRequest) (iuf.Session, error) {
	// go through the sessions and if there is any session that is not completed or aborted, then mark it as aborted
	// and terminate its workflows.
	sessions, err := s.ListSessions(activityName)
	if err != nil {
		s.logger.Errorf("HistoryAbortAction: An error occurred while listing sessions for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	var errors []error
	for _, session := range sessions {
		if session.CurrentState != iuf.SessionStateCompleted && session.CurrentState != iuf.SessionStateAborted {
			err := s.AbortSession(&session, req.Force)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		s.logger.Errorf("HistoryAbortAction: An error(s) occurred while aborting sessions for activity %s: %v", activityName, errors)
		return iuf.Session{}, err
	}

	// add a history entry for aborted sessions
	comment := req.Comment
	if comment == "" {
		comment = "Aborted"
	}

	err = s.CreateHistoryEntry(activityName, iuf.ActivityStateWaitForAdmin, comment)
	if err != nil {
		s.logger.Errorf("HistoryAbortAction: An error occurred while creating history entry for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	if len(sessions) > 0 {
		return sessions[len(sessions)-1], nil
	} else {
		return iuf.Session{}, nil
	}
}

func (s iufService) HistoryPausedAction(activityName string, req iuf.HistoryActionRequest) (iuf.Session, error) {
	// go through the sessions and if there is any session that is in_progress state, then mark it as paused
	sessions, err := s.ListSessions(activityName)
	if err != nil {
		s.logger.Errorf("HistoryPausedAction: An error occurred while listing sessions for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	var errors []error
	for _, session := range sessions {
		if session.CurrentState == iuf.SessionStateInProgress {
			err := s.PauseSession(&session)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		s.logger.Errorf("HistoryPausedAction: An error(s) occurred while aborting sessions for activity %s: %v", activityName, errors)
		return iuf.Session{}, err
	}

	comment := req.Comment
	if comment == "" {
		comment = "Paused"
	}

	// add a history entry for aborted sessions
	err = s.CreateHistoryEntry(activityName, iuf.ActivityStatePaused, comment)
	if err != nil {
		s.logger.Errorf("HistoryPausedAction: An error occurred while creating history entry for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	if len(sessions) > 0 {
		return sessions[len(sessions)-1], nil
	} else {
		return iuf.Session{}, nil
	}
}

func (s iufService) HistoryResumeAction(activityName string, req iuf.HistoryActionRequest) (iuf.Session, error) {
	// go through the sessions and if there is any session that is in_progress state, then mark it as paused
	sessions, err := s.ListSessions(activityName)
	if err != nil {
		s.logger.Errorf("HistoryResumeAction: An error occurred while listing sessions for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	var errors []error
	for _, session := range sessions {
		if session.CurrentState == iuf.SessionStatePaused {
			err := s.ResumeSession(&session)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		s.logger.Errorf("HistoryResumeAction: An error(s) occurred while aborting sessions for activity %s: %v", activityName, errors)
		return iuf.Session{}, err
	}

	comment := req.Comment
	if comment == "" {
		comment = "Resumed"
	}

	// add a history entry for aborted sessions
	err = s.CreateHistoryEntry(activityName, iuf.ActivityStateInProgress, comment)
	if err != nil {
		s.logger.Errorf("HistoryResumeAction: An error occurred while creating history entry for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	if len(sessions) > 0 {
		return sessions[len(sessions)-1], nil
	} else {
		return iuf.Session{}, nil
	}
}

func (s iufService) HistoryRestartAction(activityName string, req iuf.HistoryActionRequest) (iuf.Session, error) {
	return iuf.Session{}, nil
}

func (s iufService) HistoryBlockedAction(activityName string, req iuf.HistoryActionRequest) (iuf.Session, error) {
	// this is only allowed when activity is in debug, paused, or wait_for_admin state.
	activity, err := s.GetActivity(activityName)
	if err != nil {
		s.logger.Errorf("HistoryBlockedAction: An error occurred while fetching activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	sessions, err := s.ListSessions(activityName)
	if err != nil {
		s.logger.Errorf("HistoryBlockedAction: An error occurred while listing sessions for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	// there shouldn't be any running sessions
	var lastSession iuf.Session
	for _, session := range sessions {
		if session.CurrentState == iuf.SessionStateInProgress || session.CurrentState == iuf.SessionStatePaused {
			err = utils.GenericError{
				Message: fmt.Sprintf("HistoryBlockedAction: For the activity %s, there is currently an session %s that is in state %s.", activityName, session.Name, session.CurrentStage),
			}
			s.logger.Error(err)
			return iuf.Session{}, err
		}

		lastSession = session
	}

	switch activity.ActivityState {
	case iuf.ActivityStateWaitForAdmin, iuf.ActivityStateDebug:
		activity.ActivityState = iuf.ActivityStateBlocked
		_, err := s.updateActivity(activity)
		if err != nil {
			s.logger.Errorf("HistoryBlockedAction: An error occured while updating activity %s to be in blocked state.", activityName)
			return iuf.Session{}, err
		}
	case iuf.ActivityStateBlocked:
		// noop
		return lastSession, nil
	default:
		err = utils.GenericError{
			Message: fmt.Sprintf("HistoryBlockedAction: The activity %s must be in debug or wait_for_admin state for it to be marked as blocked. Currently, it is in %s: %v", activityName, activity.ActivityState, activity.ActivityState),
		}
		s.logger.Error(err)
		return iuf.Session{}, err
	}

	comment := req.Comment
	if comment == "" {
		comment = "Blocked"
	}

	// add a history entry for blocked activity
	err = s.CreateHistoryEntry(activityName, iuf.ActivityStateBlocked, comment)
	if err != nil {
		s.logger.Errorf("HistoryAbortAction: An error occurred while creating history entry for activity %s: %v", activityName, err)
		return iuf.Session{}, err
	}

	return lastSession, nil
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
