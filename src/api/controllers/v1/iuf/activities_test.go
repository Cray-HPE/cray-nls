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
package iuf

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

func TestCreateActivity(t *testing.T) {

	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executeWithContext := func(
		workflowService *mocks.MockWorkflowService,
		iufServices *mocks.MockIufService,
		requestBody string,
	) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		context, ginEngine := gin.CreateTestContext(response)

		requestUrl := "/iuf/v1/activities"

		context.Request, _ = http.NewRequest("POST", requestUrl, strings.NewReader(string(requestBody)))

		ginEngine.POST("/iuf/v1/activities", NewIufController(workflowService, iufServices, *utils.GetLogger().GetGinLogger().Logger).CreateActivity)
		ginEngine.ServeHTTP(response, context.Request)
		return response
	}

	t.Run("Input validate", func(t *testing.T) {

		workflowServiceMock := mocks.NewMockWorkflowService(ctrl)
		iufServiceMock := mocks.NewMockIufService(ctrl)
		iufServiceMock.EXPECT().CreateActivity(gomock.Any()).Return(nil).AnyTimes()
		var tests = []struct {
			name        string
			requestBody string
			statusCode  int
		}{
			{name: "no name", requestBody: "{}", statusCode: http.StatusBadRequest},
			{name: "no input", requestBody: `{"name":"asdf"}`, statusCode: http.StatusCreated},
			{name: "with input", requestBody: `{"name":"asdf", "input_parameters":{}}`, statusCode: http.StatusCreated},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := executeWithContext(workflowServiceMock, iufServiceMock, tt.requestBody)
				assert.Equal(t, tt.statusCode, res.Code)
			})
		}

	})

}
