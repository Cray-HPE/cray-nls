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
package services

//go:generate mockgen -destination=../mocks/services/iuf.go -package=mocks -source=iuf.go

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	core_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	DEFAULT_NAMESPACE      = "argo"
	LABEL_ACTIVITY         = "iuf-activity"
	CONFIGMAP_KEY_ACTIVITY = "Activity"
)

type IufService interface {
	GetSessionsByActivityName(activityName string) ([]iuf.Session, error)
	CreateActivity(req iuf.CreateActivityRequest) error
	ListActivities() ([]iuf.Activity, error)
	GetActivity(name string) (iuf.Activity, error)
	PatchActivity(name string, req iuf.PatchActivityRequest) (iuf.Activity, error)
}

// IufService service layer
type iufService struct {
	logger           utils.Logger
	k8sRestClientSet *kubernetes.Clientset
	env              utils.Env
}

// NewIufService creates a new Iufservice
func NewIufService(logger utils.Logger, argoService ArgoService, env utils.Env) IufService {

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

func (s iufService) GetSessionsByActivityName(activityName string) ([]iuf.Session, error) {
	var mySessions []iuf.Session
	if s.k8sRestClientSet == nil {
		return mySessions, nil
	}

	sessions, err := s.k8sRestClientSet.
		RESTClient().Get().
		AbsPath("/apis/iuf.hpe.com/v1").
		Resource("sessions").
		Param("labelSelector", fmt.Sprintf("activityName=%s", activityName)).
		DoRaw(context.TODO())
	if err != nil {
		s.logger.Error(err)
		return mySessions, err
	}

	err = json.Unmarshal(sessions, &mySessions)
	if err != nil {
		s.logger.Error(err)
		return mySessions, err
	}
	return mySessions, nil
}

func (s iufService) CreateActivity(req iuf.CreateActivityRequest) error {
	// TODO: validate input parameters
	req.ActivityState = iuf.ActivityStateWaitForAdmin
	configmap, err := s.iufObjectToConfigMapData(req, req.Name)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Create(
			context.TODO(),
			&configmap,
			v1.CreateOptions{},
		)
	// TODO: add activity history
	return err
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

	res, err := s.configMapDataToActivity(rawConfigMapData.Data[CONFIGMAP_KEY_ACTIVITY])
	if err != nil {
		s.logger.Error(err)
		return res, err
	}
	return res, err
}

func (s iufService) PatchActivity(name string, req iuf.PatchActivityRequest) (iuf.Activity, error) {
	tmp, err := s.GetActivity(name)
	if err != nil {
		s.logger.Error(err)
		return iuf.Activity{}, err
	}

	// TODO: block request if activity is in_progress, paused
	// TODO: validate input parameters
	// TODO: support partial update
	tmp.InputParameters = req.InputParameters
	configmap, err := s.iufObjectToConfigMapData(tmp, tmp.Name)
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
	return tmp, err
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
		tmp, err := s.configMapDataToActivity(rawConfigMap.Data[CONFIGMAP_KEY_ACTIVITY])
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

func (s iufService) iufObjectToConfigMapData(activity interface{}, name string) (core_v1.ConfigMap, error) {
	reqBytes, err := json.Marshal(activity)
	if err != nil {
		s.logger.Error(err)
		return core_v1.ConfigMap{}, err
	}
	res := core_v1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"type": LABEL_ACTIVITY,
			},
		},
		Data: map[string]string{CONFIGMAP_KEY_ACTIVITY: string(reqBytes)},
	}
	return res, nil
}
