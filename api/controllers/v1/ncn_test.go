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

func TestNcnsCreateRebuildWorkflow(t *testing.T) {

	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executeWithContext := func(
		workflowService *mocks.MockWorkflowService,
		requestBody string,
	) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		context, ginEngine := gin.CreateTestContext(response)

		requestUrl := "/v1/ncns/rebuild"

		context.Request, _ = http.NewRequest("POST", requestUrl, strings.NewReader(requestBody))

		ginEngine.POST("/v1/ncns/rebuild", NewNcnController(workflowService, *utils.GetLogger().GetGinLogger().Logger).NcnsCreateRebuildWorkflow)
		ginEngine.ServeHTTP(response, context.Request)
		return response
	}

	t.Run("Happy", func(t *testing.T) {

		workflowServiceMock := mocks.NewMockWorkflowService(ctrl)
		workflowServiceMock.EXPECT().CreateRebuildWorkflow(gomock.Any(), gomock.Any(), gomock.Any()).Return(
			&v1alpha1.Workflow{
				ObjectMeta: v1.ObjectMeta{Name: "mocked", Labels: map[string]string{"targetNcn": "mocked-target-ncn"}},
			}, nil)
		res := executeWithContext(
			workflowServiceMock,
			`{
				"hosts": [
					"ncn-w003",
					"ncn-w003",
					"ncn-w003",
					"ncn-w003"
				]
			}`,
		)
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Error", func(t *testing.T) {

		workflowServiceMock := mocks.NewMockWorkflowService(ctrl)
		workflowServiceMock.EXPECT().CreateRebuildWorkflow(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mocked error"))
		res := executeWithContext(
			workflowServiceMock,
			`{
				"hosts": [
					"ncn-w003",
					"ncn-w003",
					"ncn-w003",
					"ncn-w003"
				]
			}`,
		)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("wrong hostname - invalid", func(t *testing.T) {

		workflowServiceMock := mocks.NewMockWorkflowService(ctrl)
		res := executeWithContext(
			workflowServiceMock,
			`{
				"hosts": [
					"ncn-s003",
					"ncn-m003",
					"ncn-w003",
					"ncn-w003"
				]
			}`,
		)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("invalid request", func(t *testing.T) {

		workflowServiceMock := mocks.NewMockWorkflowService(ctrl)
		res := executeWithContext(
			workflowServiceMock,
			`{
				"hosts": ["ncn-s003]
		  	}`,
		)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
}
