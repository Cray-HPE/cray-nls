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
package services_shared

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sService struct {
	Client *kubernetes.Clientset
}

func NewK8sService() K8sService {
	var config *rest.Config

	// first try local development mode
	home, _ := os.UserHomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", home+"/.k3d/kubeconfig-mycluster.yaml")
	if err != nil {
		// first try: running within a pod
		config, err = rest.InClusterConfig()
		if err != nil {
			// second try: remote development mode ... running outside of a pod (development mode on an ncn)
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
	return K8sService{
		Client: k8sRestClientSet,
	}
}
