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

//go:generate mockgen -destination=../mocks/services/iuf.go -package=mocks -source=iuf.go

import (
	_ "embed"
	"encoding/json"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	services_shared "github.com/Cray-HPE/cray-nls/src/api/services/shared"
	"github.com/Cray-HPE/cray-nls/src/utils"
	core_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	DEFAULT_NAMESPACE  = "argo"
	LABEL_ACTIVITY     = "iuf_activity"
	LABEL_HISTORY      = "iuf_history"
	LABEL_SESSION      = "iuf_session"
	LABEL_ACTIVITY_REF = "iuf_activity_ref"
)

type IufService interface {
	CreateActivity(req iuf.CreateActivityRequest) error
	ListActivities() ([]iuf.Activity, error)
	GetActivity(name string) (iuf.Activity, error)
	PatchActivity(name string, req iuf.PatchActivityRequest) (iuf.Activity, error)
	// history
	ListActivityHistory(activityName string) ([]iuf.History, error)
	HistoryRunAction(activityName string, req iuf.HistoryRunActionRequest) error
	// session
	ListSessions(activityName string) ([]iuf.Session, error)
	GetSession(sessionName string) (iuf.Session, string, error)
	// session operator
	ConfigMapDataToSession(data string) (iuf.Session, error)
	UpdateActivityStateFromSessionState(session iuf.Session, activityRef string) error
	UpdateSession(session iuf.Session, activityRef string) error
}

// IufService service layer
type iufService struct {
	logger           utils.Logger
	k8sRestClientSet *kubernetes.Clientset
	env              utils.Env
}

// NewIufService creates a new Iufservice
func NewIufService(logger utils.Logger, k8sSvc services_shared.K8sService, env utils.Env) IufService {

	iufSvc := iufService{
		logger:           logger,
		k8sRestClientSet: k8sSvc.Client,
		env:              env,
	}
	return iufSvc
}

func (s iufService) iufObjectToConfigMapData(activity interface{}, name string, iufType string) (core_v1.ConfigMap, error) {
	reqBytes, err := json.Marshal(activity)
	if err != nil {
		s.logger.Error(err)
		return core_v1.ConfigMap{}, err
	}
	res := core_v1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"type": iufType,
			},
		},
		Data: map[string]string{iufType: string(reqBytes)},
	}
	return res, nil
}
