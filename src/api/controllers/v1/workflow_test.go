//
//  MIT License
//
//  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
//
//  Permission is hereby granted, free of charge, to any person obtaining a
//  copy of this software and associated documentation files (the "Software"),
//  to deal in the Software without restriction, including without limitation
//  the rights to use, copy, modify, merge, publish, distribute, sublicense,
//  and/or sell copies of the Software, and to permit persons to whom the
//  Software is furnished to do so, subject to the following conditions:
//
//  The above copyright notice and this permission notice shall be included
//  in all copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
//  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
//  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
//  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
//  OTHER DEALINGS IN THE SOFTWARE.
//
package controllers_v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mocks "github.com/Cray-HPE/cray-nls/api/mocks/services"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/alecthomas/assert"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetWorkflows(t *testing.T) {

	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executeWithContext := func(
		workflowService *mocks.MockWorkflowService,
		url ...string,
	) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		context, ginEngine := gin.CreateTestContext(response)

		requestUrl := "/v1/workflows"
		if url != nil {
			requestUrl = url[0]
		}

		context.Request, _ = http.NewRequest("GET", requestUrl, strings.NewReader(string("")))

		ginEngine.GET("/v1/workflows", NewWorkflowController(workflowService, *utils.GetLogger().GetGinLogger().Logger).GetWorkflows)
		ginEngine.ServeHTTP(response, context.Request)
		return response
	}

	t.Run("Happy", func(t *testing.T) {

		workflowServiceMock := mocks.NewMockWorkflowService(ctrl)
		workflowServiceMock.EXPECT().GetWorkflows(gomock.Any()).Return(
			&v1alpha1.WorkflowList{
				Items: v1alpha1.Workflows{
					v1alpha1.Workflow{
						ObjectMeta: v1.ObjectMeta{Name: "mocked", Labels: map[string]string{"targetNcn": "mocked-target-ncn"}},
					},
				},
			}, nil)
		res := executeWithContext(workflowServiceMock)
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Error", func(t *testing.T) {

		workflowServiceMock := mocks.NewMockWorkflowService(ctrl)
		workflowServiceMock.EXPECT().GetWorkflows(gomock.Any()).Return(nil, fmt.Errorf("mocked error"))
		res := executeWithContext(workflowServiceMock)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

}
