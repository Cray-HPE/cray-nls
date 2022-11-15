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
	"testing"

	mocks "github.com/Cray-HPE/cray-nls/src/api/mocks/services"
	"github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestGetHistory(t *testing.T) {

	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executeWithContext := func(
		workflowService *mocks.MockWorkflowService,
		iufServices *mocks.MockIufService,
		ginRequestPath string,
		requestUrl string,
	) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		context, ginEngine := gin.CreateTestContext(response)

		context.Request, _ = http.NewRequest("GET", requestUrl, nil)

		ginEngine.GET(ginRequestPath, NewIufController(workflowService, iufServices, *utils.GetLogger().GetGinLogger().Logger).GetHistory)
		ginEngine.ServeHTTP(response, context.Request)
		return response
	}
	t.Run("404: not found", func(t *testing.T) {
		workflowServiceMock := mocks.NewMockWorkflowService(ctrl)
		iufServiceMock := mocks.NewMockIufService(ctrl)
		iufServiceMock.EXPECT().GetActivityHistory(gomock.Any(), gomock.Any()).Return(iuf.History{}, nil).AnyTimes()
		res := executeWithContext(workflowServiceMock, iufServiceMock, "/iuf/v1/activities/:activity_name/history/:start_time", "/iuf/v1/activities/test/history/123")
		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("400: wrong start time", func(t *testing.T) {
		workflowServiceMock := mocks.NewMockWorkflowService(ctrl)
		iufServiceMock := mocks.NewMockIufService(ctrl)
		iufServiceMock.EXPECT().GetActivityHistory(gomock.Any(), gomock.Any()).Return(iuf.History{}, nil).AnyTimes()
		res := executeWithContext(workflowServiceMock, iufServiceMock, "/iuf/v1/activities/:activity_name/history/:start_time", "/iuf/v1/activities/test/history/asdf")
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

}
