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

import (
	"github.com/Cray-HPE/cray-nls/api/services"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/gin-gonic/gin"
)

// Controller data type
type WorkflowController struct {
	service services.WorkflowService
	logger  utils.Logger
}

// NewController creates new  controller
func NewWorkflowController(Service services.WorkflowService, logger utils.Logger) WorkflowController {
	return WorkflowController{
		service: Service,
		logger:  logger,
	}
}

// GetWorkflows
// @Summary   Get status of a ncn workflow
// @Tags      Workflow
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v1/workflows [get]
// @Security  OAuth2Application[admin,read]
func (u WorkflowController) GetWorkflows(c *gin.Context) {
	workflowList, err := u.service.GetWorkflows(c)
	if err != nil {
		u.logger.Error(err)
	}
	var workflows []interface{}
	for _, workflow := range workflowList.Items {
		tmp := map[string]interface{}{
			"name":  workflow.Name,
			"phase": workflow.Labels["workflows.argoproj.io/phase"],
		}
		workflows = append(workflows, tmp)
	}
	c.JSON(200, gin.H{"data": workflows})
}

// DeleteWorkflow
// @Summary   Delete a ncn workflow
// @Param     name  path  string  true  "name of workflow"
// @Tags      Workflow
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v1/workflows/{name} [delete]
// @Security  OAuth2Application[admin]
func (u WorkflowController) DeleteWorkflow(c *gin.Context) {
	c.JSON(200, gin.H{"data": " updated"})
}

// RetryWorkflows
// @Summary   Retry a failed ncn workflow, skip passed steps
// @Param     name  path  string  true  "name of workflow"
// @Tags      Workflow
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v1/workflows/{name}/retry [put]
// @Security  OAuth2Application[admin]
func (u WorkflowController) RetryWorkflow(c *gin.Context) {
	c.JSON(200, gin.H{"data": " updated"})
}

// RerunWorkflows
// @Summary   Rerun a workflow, all steps will run
// @Param     name  path  string  true  "name of workflow"
// @Tags      Workflow
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v1/workflows/{name}/rerun [put]
// @Security  OAuth2Application[admin]
func (u WorkflowController) RerunWorkflow(c *gin.Context) {
	c.JSON(200, gin.H{"data": " updated"})
}
