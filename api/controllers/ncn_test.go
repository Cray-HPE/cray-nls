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
package controllers

// import (
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	workflowServiceMock "github.com/Cray-HPE/cray-nls/api/mocks/services"
// 	"github.com/alecthomas/assert"
// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// )

// func TestNcnCreateRebuildWorkflow(t *testing.T) {
// 	var (
// 		invitationId   = "3da465a6-be13-405e-a653-c68adf59f2be"
// 		firstName      = "Tom"
// 		lastName       = "Sudchai"
// 		roleId         = uint(1)
// 		operatorCode   = "velo"
// 		email          = "toms@gmail.com"
// 		password       = "P@ssw0rd"
// 		hashedPassword = "$2y$12$S0Gbs0Qm5rJGibfFBTARa.6ap9OBuXYbYJ.deCzsOo4uQNJR1KbJO"
// 	)

// 	gin.SetMode(gin.TestMode)
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	workflowSvcMock := workflowServiceMock.NewMockWorkflowServiceInterface(ctrl)

// 	executeWithContext := func(
// 		createWorkerRebuildWorkflow *workflowServiceMock.MockWorkflowServiceInterface,
// 		jsonRequestBody []byte,
// 		operatorCode string,
// 	) *httptest.ResponseRecorder {
// 		response := httptest.NewRecorder()
// 		context, ginEngine := gin.CreateTestContext(response)

// 		requestUrl := "/v1/operators/staffs"
// 		httpRequest, _ := http.NewRequest("POST", requestUrl, strings.NewReader(string(jsonRequestBody)))

// 		NewEndpointHTTPHandler(ginEngine, workflowSvcMock)
// 		ginEngine.ServeHTTP(response, httpRequest)
// 		return response
// 	}

// 	createdStaffEntity := entities.OperatorStaff{
// 		ID:        roleId,
// 		FirstName: firstName,
// 		LastName:  lastName,
// 		Email:     email,
// 		Password:  hashedPassword,
// 		Operators: []entities.StaffOperator{{
// 			OperatorCode: operatorCode, RoleID: roleId,
// 		}},
// 	}

// 	t.Run("Happy", func(t *testing.T) {
// 		jsonRequestBody, _ := json.Marshal(createStaffFromInviteRequestJSON{
// 			InvitationId:    invitationId,
// 			FirstName:       firstName,
// 			LastName:        lastName,
// 			Password:        password,
// 			ConfirmPassword: password,
// 		})

// 		staffMock.EXPECT().CreateStaff(gomock.Any(), gomock.Any()).Return(&createdStaffEntity, nil)

// 		res := executeWithContext(staffMock, jsonRequestBody, operatorCode)
// 		assert.Equal(t, http.StatusOK, res.Code)
// 	})
// }
