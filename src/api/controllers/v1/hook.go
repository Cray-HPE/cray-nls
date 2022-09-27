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

	v1 "github.com/Cray-HPE/cray-nls/src/api/models/v1"
	"github.com/Cray-HPE/cray-nls/src/api/services"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/gin-gonic/gin"
)

// HookController data type
type HookController struct {
	workflowService services.WorkflowService
	logger          utils.Logger
}

// NewHookController creates new Ncn controller
func NewHookController(workflowService services.WorkflowService, logger utils.Logger) HookController {
	return HookController{
		workflowService: workflowService,
		logger:          logger,
	}
}

// AddHooks
func (u HookController) AddHooks(c *gin.Context) {
	var requestBody v1.SyncRequest
	var response v1.SyncResponse
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(400, errResponse)
		return
	}
	u.logger.Infof(
		"[%s] Hook changed, phase: %s, observed generation: %d",
		requestBody.Parent.Name,
		requestBody.Parent.Status.Phase,
		requestBody.Parent.Status.ObservedGeneration,
	)
	response = v1.SyncResponse{
		Status:             v1.HookStatus{Phase: "created"},
		ResyncAfterSeconds: 0,
	}
	u.logger.Infof("[%s] Hook created, namespace: %s, resourceVersion: %s", requestBody.Parent.Name, requestBody.Parent.Namespace, requestBody.Parent.ResourceVersion)
	c.JSON(200, response)
}
