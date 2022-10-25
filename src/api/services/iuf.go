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
	"fmt"
	"os"

	"github.com/argoproj/pkg/json"

	iuf_v1 "github.com/Cray-HPE/cray-nls/src/api/models/iuf/v1"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type IufService interface {
	GetSessionsByActivityName(activityName string) ([]iuf_v1.IufSession, error)
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

func (s iufService) GetSessionsByActivityName(activityName string) ([]iuf_v1.IufSession, error) {
	var mySessions []iuf_v1.IufSession
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
