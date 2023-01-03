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
	models_nls "github.com/Cray-HPE/cray-nls/src/api/models/nls"
	services_shared "github.com/Cray-HPE/cray-nls/src/api/services/shared"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/gin-gonic/gin"
)

// Controller data type
type WorkflowController struct {
	service services_shared.WorkflowService
	logger  utils.Logger
}

// NewController creates new  controller
func NewWorkflowController(Service services_shared.WorkflowService, logger utils.Logger) WorkflowController {
	return WorkflowController{
		service: Service,
		logger:  logger,
	}
}

// GetWorkflows
// @Summary  Get status of a ncn workflow
// @Param    labelSelector  query  string  false  "Label Selector"
// @Tags     Workflow Management
// @Accept   json
// @Produce  json
// @Success  200  {object}  []models.GetWorkflowResponse
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /nls/v1/workflows [get]
func (u WorkflowController) GetWorkflows(c *gin.Context) {
	workflowList, err := u.service.GetWorkflows(c)
	if err != nil {
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(500, errResponse)
		return
	}
	var workflows []models_nls.GetWorkflowResponse
	for _, workflow := range workflowList.Items {
		workflow.Status.StoredTemplates = nil
		workflow.Status.ArtifactRepositoryRef = nil
		tmp := models_nls.GetWorkflowResponse{
			Name:   workflow.Name,
			Labels: workflow.Labels,
			Status: workflow.Status,
		}
		workflows = append(workflows, tmp)
	}
	c.JSON(200, workflows)
}

// DeleteWorkflow
// @Summary  Delete a ncn workflow
// @Param    name  path  string  true  "name of workflow"
// @Tags     Workflow Management
// @Accept   json
// @Produce  json
// @Success  200  {object}  utils.ResponseOk
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /nls/v1/workflows/{name} [delete]
func (u WorkflowController) DeleteWorkflow(c *gin.Context) {
	err := u.service.DeleteWorkflow(c)
	if err != nil {
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(500, errResponse)
		return
	}
	c.JSON(200, gin.H{"data": " deleted"})
}

// RetryWorkflows
// @Summary  Retry a failed ncn workflow, skip passed steps
// @Param    name          path  string                           true  "name of workflow"
// @Param    retryOptions  body  models.RetryWorkflowRequestBody  true  "retry options"
// @Tags     Workflow Management
// @Accept   json
// @Produce  json
// @Success  200  {object}  utils.ResponseOk
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /nls/v1/workflows/{name}/retry [put]
func (u WorkflowController) RetryWorkflow(c *gin.Context) {
	err := u.service.RetryWorkflow(c)
	if err != nil {
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(500, errResponse)
		return
	}
	c.Status(200)
}

// RerunWorkflows
// @Summary  Rerun a workflow, all steps will run
// @Param    name  path  string  true  "name of workflow"
// @Tags     Workflow Management
// @Accept   json
// @Produce  json
// @Success  200  {object}  utils.ResponseOk
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /nls/v1/workflows/{name}/rerun [put]
func (u WorkflowController) RerunWorkflow(c *gin.Context) {
	err := u.service.RerunWorkflow(c)
	if err != nil {
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(500, errResponse)
		return
	}
	c.Status(200)
}
