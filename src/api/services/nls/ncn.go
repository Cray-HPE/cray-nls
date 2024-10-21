// MIT License
//
// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
package services_nls

//go:generate mockgen -destination=../mocks/services/ncn.go -package=mocks -source=ncn.go

import (
	"context"
	"embed"
	"os"
	"time"

	"github.com/Cray-HPE/cray-nls/src/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

//go:embed cray-nls.hpe.com_hooks.yaml
var nlsHooksFS embed.FS

type NcnService interface{}

// NcnService service layer
type ncnService struct {
	k8sRestClientSet *kubernetes.Clientset
	logger           utils.Logger
}

// NewNcnService creates a new Ncnservice
func NewNcnService(logger utils.Logger) NcnService {
	var config *rest.Config
	config, err := rest.InClusterConfig()
	if err != nil {
		// use k3d kubeconfig in development mode
		home, _ := os.UserHomeDir()
		config, err = clientcmd.BuildConfigFromFlags("", home+"/.k3d/kubeconfig-mycluster.yaml")
		if err != nil {
			config, err = clientcmd.BuildConfigFromFlags("", "/etc/kubernetes/admin.conf")
			if err != nil {
				panic(err.Error())
			}
		}
	}
	k8sRestClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// initialize ncn hooks crd
	_, err = k8sRestClientSet.
		RESTClient().
		Get().
		AbsPath("/apis/apiextensions.k8s.io/v1/customresourcedefinitions/hooks.cray-nls.hpe.com").
		DoRaw(context.TODO())
	if err == nil {
		// delete existing crd before upgrade
		_, err = k8sRestClientSet.
			RESTClient().
			Delete().
			AbsPath("/apis/apiextensions.k8s.io/v1/customresourcedefinitions/hooks.cray-nls.hpe.com").
			DoRaw(context.TODO())
		if err != nil {
			logger.Panic(err)
		}
		time.Sleep(5000 * time.Millisecond)
	}
	// create crd
	hooksCrdBytes, _ := nlsHooksFS.ReadFile("cray-nls.hpe.com_hooks.yaml")
	body, _ := yaml.YAMLToJSON(hooksCrdBytes)
	_, err = k8sRestClientSet.
		RESTClient().
		Post().
		AbsPath("/apis/apiextensions.k8s.io/v1/customresourcedefinitions").
		Body(body).DoRaw(context.TODO())
	if err != nil {
		logger.Info("got panic. PRINTING IF IsCONFLICT")
		logger.Info(errors.IsConflict(err))
		logger.Panic(err)
		// Handle the conflict error here
	}

	ncSvc := ncnService{
		logger:           logger,
		k8sRestClientSet: k8sRestClientSet,
	}
	return ncSvc
}
