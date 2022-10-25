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
package iuf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mocks "github.com/Cray-HPE/cray-nls/src/api/mocks/services"
	v1 "github.com/Cray-HPE/cray-nls/src/api/models/iuf/v1"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestIufActivitySync(t *testing.T) {
	type mymocks struct {
		iufSvcMock *mocks.MockIufService
	}

	gin.SetMode(gin.TestMode)

	executeWithContext := func(
		workflowService *mocks.MockWorkflowService,
		iufService *mocks.MockIufService,
		requestBody string,
	) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		context, ginEngine := gin.CreateTestContext(response)

		requestUrl := "/v1/iuf/activities/sync"

		context.Request, _ = http.NewRequest("POST", requestUrl, strings.NewReader(requestBody))

		ginEngine.POST(requestUrl, NewIufController(workflowService, iufService, *utils.GetLogger().GetGinLogger().Logger).IufActivitySync)
		ginEngine.ServeHTTP(response, context.Request)
		return response
	}

	var tests = []struct {
		name        string
		requestBody v1.IufActivitiesSyncRequest
		prepare     func(f *mymocks)
		statusCode  int
	}{
		{
			"return 200 when activity is marked as completed",
			v1.IufActivitiesSyncRequest{
				Parent: v1.IufActivity{
					Spec: v1.IufActivitiesSpec{
						SharedInput:    v1.SharedInput{},
						IsCompleted:    true,
						CurrentComment: "",
					},
					Status: v1.IufActivitiesStatus{
						SharedInput: v1.SharedInput{},
					},
				},
			},
			func(f *mymocks) {},
			200,
		},
		{
			"return 200 when activity is NOT marked as completed",
			v1.IufActivitiesSyncRequest{
				Parent: v1.IufActivity{
					Spec: v1.IufActivitiesSpec{
						SharedInput:    v1.SharedInput{},
						IsCompleted:    false,
						CurrentComment: "",
					},
					Status: v1.IufActivitiesStatus{
						SharedInput: v1.SharedInput{},
					},
				},
			},
			func(f *mymocks) {
				// mock: return error
				f.iufSvcMock.EXPECT().GetSessionsByActivityName(gomock.Any()).Return(nil, fmt.Errorf("mocked error"))
			},
			400,
		},
		{
			"return 200 when associated sessions are empty",
			v1.IufActivitiesSyncRequest{
				Parent: v1.IufActivity{
					Spec: v1.IufActivitiesSpec{
						SharedInput:    v1.SharedInput{},
						IsCompleted:    false,
						CurrentComment: "",
					},
					Status: v1.IufActivitiesStatus{
						SharedInput: v1.SharedInput{},
					},
				},
			},
			func(f *mymocks) {
				// mock: return empty list
				f.iufSvcMock.EXPECT().GetSessionsByActivityName(gomock.Any()).Return([]v1.IufSession{}, nil)
			},
			200,
		},
		{
			"return 200 when one session is in progress",
			v1.IufActivitiesSyncRequest{
				Parent: v1.IufActivity{
					Spec: v1.IufActivitiesSpec{
						SharedInput:    v1.SharedInput{},
						IsCompleted:    false,
						CurrentComment: "",
					},
					Status: v1.IufActivitiesStatus{
						SharedInput: v1.SharedInput{},
					},
				},
			},
			func(f *mymocks) {
				f.iufSvcMock.EXPECT().GetSessionsByActivityName(gomock.Any()).Return([]v1.IufSession{
					// mock: return list with ONE element
					{Status: v1.IufSessionStatus{CurrentState: v1.IufSessionCurrentState{Type: v1.IufSessionStageInProgress}}},
				}, nil)
			},
			200,
		},
		{
			"return 200 when associated sessions are empty and request is updating sharedinput",
			v1.IufActivitiesSyncRequest{
				Parent: v1.IufActivity{
					Spec: v1.IufActivitiesSpec{
						SharedInput:    v1.SharedInput{MediaDir: "new"},
						IsCompleted:    false,
						CurrentComment: "",
					},
					Status: v1.IufActivitiesStatus{
						SharedInput: v1.SharedInput{MediaDir: "old"},
					},
				},
			},
			func(f *mymocks) {
				// mock: return empty list
				f.iufSvcMock.EXPECT().GetSessionsByActivityName(gomock.Any()).Return([]v1.IufSession{}, nil)
			},
			200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			workflowServiceMock := mocks.NewMockWorkflowService(ctrl)
			f := mymocks{
				iufSvcMock: mocks.NewMockIufService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			b, _ := json.Marshal(tt.requestBody)
			res := executeWithContext(
				workflowServiceMock,
				f.iufSvcMock,
				string(b),
			)
			assert.Equal(t, tt.statusCode, res.Code)
		})
	}
}
