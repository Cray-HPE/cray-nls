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
	"os"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	services_shared "github.com/Cray-HPE/cray-nls/src/api/services/shared"
	"github.com/Cray-HPE/cray-nls/src/utils"
	core_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	GetSession(activityName string, sessionName string) (iuf.Session, error)
}

// IufService service layer
type iufService struct {
	logger           utils.Logger
	k8sRestClientSet *kubernetes.Clientset
	env              utils.Env
}

// NewIufService creates a new Iufservice
func NewIufService(logger utils.Logger, argoService services_shared.ArgoService, env utils.Env) IufService {

	var config *rest.Config
	config, err := rest.InClusterConfig()
	if err != nil {
		// use k3d kubeconfig in development mode
		home, _ := os.UserHomeDir()
		config, err = clientcmd.BuildConfigFromFlags("", home+"/.k3d/kubeconfig-mycluster.yaml")
		if err != nil {
			panic(err.Error())
		}
	}
	k8sRestClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	iufSvc := iufService{
		logger:           logger,
		k8sRestClientSet: k8sRestClientSet,
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
