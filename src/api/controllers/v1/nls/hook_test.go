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
package controllers_v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mocks "github.com/Cray-HPE/cray-nls/src/api/mocks/services"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestAddHook(t *testing.T) {

	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executeWithContext := func(
		workflowService *mocks.MockWorkflowService,
		requestBody string,
	) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		context, ginEngine := gin.CreateTestContext(response)

		requestUrl := "/v1/ncns/hooks"

		context.Request, _ = http.NewRequest("POST", requestUrl, strings.NewReader(requestBody))

		ginEngine.POST(requestUrl, NewHookController(workflowService, *utils.GetLogger().GetGinLogger().Logger).AddHooks)
		ginEngine.ServeHTTP(response, context.Request)
		return response
	}

	var tests = []struct {
		name        string
		requestBody string
		statusCode  int
	}{
		{
			"return immediately for already created hook",
			`{
				"controller":{
				"kind":"CompositeController",
				"apiVersion":"metacontroller.k8s.io/v1alpha1",
				"metadata":{
					"name":"nls-hooks",
					"uid":"c9adf1d6-b187-42af-8333-21b18e5930da",
					"resourceVersion":"3888",
					"generation":2,
					"creationTimestamp":"2022-09-27T00:13:33Z",
					"labels":{
						"app.kubernetes.io/managed-by":"Helm"
					},
					"annotations":{
						"meta.helm.sh/release-name":"argo-only",
						"meta.helm.sh/release-namespace":"argo"
					}
				},
				"spec":{
					"parentResource":{
						"apiVersion":"cray-nls.hpe.com/v1",
						"resource":"hooks"
					},
					"hooks":{
						"sync":{
							"webhook":{
							"url":"http://host.k3d.internal:3000/apis/nls/v1/ncns/hooks"
							}
						}
					},
					"generateSelector":true
				},
				"status":{

				}
				},
				"parent":{
				"apiVersion":"cray-nls.hpe.com/v1",
				"kind":"Hook",
				"metadata":{
					"annotations":{
						"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"cray-nls.hpe.com/v1\",\"kind\":\"Hook\",\"metadata\":{\"annotations\":{},\"name\":\"install-csi\",\"namespace\":\"default\"},\"spec\":{\"hookName\":\"test before all hook\",\"template\":\"this is a template\\nthat i am testing\\n\"}}\n"
					},
					"creationTimestamp":"2022-09-27T00:54:28Z",
					"generation":24,
					"name":"install-csi",
					"namespace":"default",
					"resourceVersion":"57328",
					"uid":"6bcdbf1f-29c5-4713-9dec-e61b8b52f1a6"
				},
				"spec":{
					"hookName":"test before all hook",
					"template":"this is a template2\nthat i am testing22\n"
				},
				"status":{
					"observedGeneration":23,
					"phase":"created"
				}
				},
				"children":{

				},
				"related":{

				},
				"finalizing":false
		 	}`,
			200,
		},
		{
			"return 400 for bad request",
			`{asdf}`,
			400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workflowServiceMock := mocks.NewMockWorkflowService(ctrl)
			res := executeWithContext(
				workflowServiceMock,
				tt.requestBody,
			)
			assert.Equal(t, tt.statusCode, res.Code)
		})
	}
}
