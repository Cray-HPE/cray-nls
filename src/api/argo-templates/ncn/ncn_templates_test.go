//go:build templates
// +build templates

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
package ncn_templates

import (
	"context"
	mocks "github.com/Cray-HPE/cray-nls/src/api/mocks/services"
	services_shared "github.com/Cray-HPE/cray-nls/src/api/services/shared"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"
	"regexp"
	"testing"
)

// NOTE: This test is not run as part of the normal test suite. It is run as part of the integration test suite.
// it is mainly used for local development and debugging.
// it reads from .env so you can use this file for debugging argo templates
func TestNcnWorkflowTemplates(t *testing.T) {
	argoSvc, k8sSvc := ncnTemplatesTestSetup(t)

	t.Run("add-label: it can add labels", func(t *testing.T) {
		wf, err := argoSvc.Client.NewWorkflowServiceClient().SubmitWorkflow(
			context.Background(),
			&workflow.WorkflowSubmitRequest{
				Namespace:    "argo",
				ResourceKind: "WorkflowTemplate",
				ResourceName: "add-labels",
				SubmitOptions: &v1alpha1.SubmitOpts{
					GenerateName: uuid.NewString(),
					Entrypoint:   "",
					Parameters:   []string{"targetNcn=ncn-w001"},
					DryRun:       false,
				},
			},
		)
		assert.Nil(t, err)
		assert.NotNil(t, wf)
		node, _ := k8sSvc.Client.CoreV1().Nodes().Get(context.Background(), "ncn-w001", metav1.GetOptions{})
		assert.Equal(t, "ncn-w001", node.Labels["cray.nls"])
	})

}

func ncnTemplatesTestSetup(t *testing.T) (services_shared.ArgoService, services_shared.K8sService) {
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

	mockTokenValue := "{{workflow.parameters.auth_token}}"
	keycloakServiceMock := mocks.NewMockKeycloakService(ctrl)
	keycloakServiceMock.EXPECT().NewKeycloakAccessToken().Return(mockTokenValue, nil).AnyTimes()

	logger := utils.GetLogger()
	argoSvc := services_shared.NewArgoService(env)
	k8sSvc := services_shared.NewK8sService()

	services_shared.NewWorkflowService(logger, argoSvc, k8sSvc, env)

	return argoSvc, k8sSvc
}
