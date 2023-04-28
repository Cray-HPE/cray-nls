//go:build integration
// +build integration

/*
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
	"encoding/json"
	mocks "github.com/Cray-HPE/cray-nls/src/api/mocks/services"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"log"
	"os"
	"regexp"
	"testing"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	workflowmocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow/mocks"
	workflowtemplatemocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate/mocks"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fake "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/yaml"
)

// NOTE: This test is not run as part of the normal test suite. It is run as part of the integration test suite.
// it is mainly used for local development and debugging.
// it reads from .env so you can use this file for debugging argo templates
func TestNcnRolloutWorkflowGen(t *testing.T) {
	activityName, _, iufSvc := ncnRolloutTestSetup(t)

	t.Run("it can generate management rollout workflow", func(t *testing.T) {
		session := iuf.Session{
			Products: []iuf.Product{},
			Workflows: []iuf.SessionWorkflow{
				{
					Id:  "1",
					Url: "1",
				},
			},
			CurrentStage: "management-nodes-rollout",
			CurrentState: iuf.SessionStateInProgress,
			InputParameters: iuf.InputParameters{
				Stages: []string{"management-nodes-rollout"},
			},
			ActivityRef: activityName,
		}

		workflow, err, _ := iufSvc.workflowGen(session)
		assert.NoError(t, err)
		assert.Equal(t, "ncn-m001", workflow.Spec.NodeSelector["kubernetes.io/hostname"])
		data, _ := yaml.Marshal(&workflow)
		t.Logf("%s", string(data))
	})

}

func ncnRolloutTestSetup(t *testing.T) (string, string, iufService) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	re := regexp.MustCompile(`^(.*` + "cray-nls" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	viper.SetConfigFile(string(rootPath) + "/.env")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	var env utils.Env
	viper.Unmarshal(&env)

	wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
	wfServiceClientMock.On(
		"GetWorkflow",
		mock.Anything,
		mock.Anything,
	).Return(new(v1alpha1.Workflow), nil)
	wfServiceClientMock.On(
		"ListWorkflows",
		mock.Anything,
		mock.Anything,
	).Return(new(v1alpha1.WorkflowList), nil)
	wfTemplateServiceClientMock := &workflowtemplatemocks.WorkflowTemplateServiceClient{}
	availableOps := []string{
		"info-to-read", "verify-images-and-configuration",
		"management-storage-nodes-rollout", "management-two master-nodes-rollout", "management-worker-nodes-rollout",
	}
	var availableTemplates v1alpha1.WorkflowTemplates

	for _, op := range availableOps {
		wftBytes, _ := os.ReadFile(env.IufInstallWorkflowFiles + "/operations/management-nodes-rollout/" + op + ".yaml")
		var wft v1alpha1.WorkflowTemplate
		err := yaml.Unmarshal(wftBytes, &wft)
		if err != nil {
			utils.GetLogger().Error(err)
		}
		availableTemplates = append(availableTemplates, wft)
	}
	templateList := v1alpha1.WorkflowTemplateList{
		Items: availableTemplates,
	}
	wfTemplateServiceClientMock.On(
		"ListWorkflowTemplates",
		mock.Anything,
		mock.Anything,
	).Return(&templateList, nil)
	name := uuid.NewString()
	activity := iuf.Activity{
		Name: name,
	}
	reqBytes, _ := json.Marshal(activity)
	configmap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: DEFAULT_NAMESPACE,
			Labels: map[string]string{
				"type": LABEL_ACTIVITY,
			},
		},
		Data: map[string]string{LABEL_ACTIVITY: string(reqBytes)},
	}
	fakeClient := fake.NewSimpleClientset(&configmap)

	mockTokenValue := "{{workflow.parameters.auth_token}}"
	keycloakServiceMock := mocks.NewMockKeycloakService(ctrl)
	keycloakServiceMock.EXPECT().NewKeycloakAccessToken().Return(mockTokenValue, nil).AnyTimes()

	iufSvc := iufService{
		logger:                 utils.GetLogger(),
		workflowClient:         wfServiceClientMock,
		workflowTemplateClient: wfTemplateServiceClientMock,
		k8sRestClientSet:       fakeClient,
		keycloakService:        keycloakServiceMock,
		env:                    env,
	}
	return name, mockTokenValue, iufSvc
}
