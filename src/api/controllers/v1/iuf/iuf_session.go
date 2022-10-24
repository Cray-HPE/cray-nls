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
	"fmt"

	v1 "github.com/Cray-HPE/cray-nls/src/api/models/iuf/v1"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/gin-gonic/gin"
)

// AddIufSession
func (u IufController) AddIufSession(c *gin.Context) {
	var requestBody v1.IufSyncRequest
	var response v1.IufSyncResponse
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(400, errResponse)
		return
	}
	response.Status = requestBody.Parent.Status

	// switch requestBody.Parent.Status.Phase {
	// case v1.IufSessionUnknown:
	// 	wf, err := u.workflowService.CreateIufWorkflow(requestBody.Parent.Spec)
	// 	if err != nil {
	// 		u.logger.Warn("Failed to create argo workflow")
	// 		u.logger.Warn(err)
	// 		response = v1.IufSyncResponse{
	// 			Status: v1.IufSessionStatus{
	// 				Phase:        v1.IufSessionError,
	// 				Operations:   [][]string{},
	// 				ArgoWorkflow: wf.GetName(),
	// 				Message:      fmt.Sprintf("%v", err),
	// 			},
	// 			ResyncAfterSeconds: 0,
	// 		}
	// 	}
	// 	u.logger.Infof("Argo workflow created: %s", wf.GetName())
	// 	response = v1.IufSyncResponse{
	// 		Status: v1.IufSessionStatus{
	// 			Phase:        v1.IufSessionPending,
	// 			Operations:   [][]string{},
	// 			ArgoWorkflow: wf.GetName(),
	// 			Message:      "",
	// 		},
	// 		ResyncAfterSeconds: 0,
	// 	}
	// case v1.IufSessionPending, v1.IufSessionRunning:
	// 	wf, err := u.workflowService.GetWorkflowByName(requestBody.Parent.Status.ArgoWorkflow, c)
	// 	if err == nil && wf.Status.Phase != "" {
	// 		response.Status.Phase = v1.IufSessionPhase(wf.Status.Phase)
	// 		operations := [][]string{{}}
	// 		for _, node := range wf.Status.Nodes {
	// 			if node.Type == v1alpha1.NodeTypeRetry {
	// 				operations[0] = append(operations[0], node.Name+":"+string(node.Phase))
	// 			}
	// 		}
	// 		sort.Strings(operations[0])
	// 		response.Status.Operations = operations
	// 		response.ResyncAfterSeconds = 10
	// 	} else {
	// 		u.logger.Error("Failed to get argo workflow")
	// 		u.logger.Error(err)
	// 		// resync in 10s
	// 		// this will allow us to check workflow status again
	// 		response.Status.Message = fmt.Sprintf("%v", err)
	// 		response.ResyncAfterSeconds = 10
	// 	}
	// case v1.IufSessionSucceeded, v1.IufSessionError, v1.IufSessionFailed:
	// 	response.ResyncAfterSeconds = 0
	// }

	u.logger.Infof("Sync Response.Status: %v", response.Status)
	c.JSON(200, response)

}
