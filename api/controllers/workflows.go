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
	"github.com/Cray-HPE/cray-nls/services"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/gin-gonic/gin"
)

// NcnController data type
type WorkflowsController struct {
	service services.NcnService
	logger  utils.Logger
}

// NewNcnController creates new Ncn controller
func NewWorkflowsController(NcnService services.NcnService, logger utils.Logger) NcnController {
	return NcnController{
		service: NcnService,
		logger:  logger,
	}
}

// NcnGetWorkflows
// @Summary   Get status of a ncn workflow
// @Param     workflow_ids  query  []string  true  "workflow ids"  collectionFormat(csv)
// @Tags      Workflow
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v1/workflows [get]
// @Security  OAuth2Application[admin,read]
func (u NcnController) NcnGetWorkflows(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnDeleteWorkflow
// @Summary   Delete a ncn workflow
// @Param     name  path  string  true  "name of workflow"
// @Tags      Workflow
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v1/workflows/{name} [delete]
// @Security  OAuth2Application[admin,read]
func (u NcnController) NcnDeleteWorkflow(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnRetryWorkflow
// @Summary   Retry a failed ncn workflow, skip passed steps
// @Param     name  path  string  true  "name of workflow"
// @Tags      Workflow
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v1/workflows/{name}/retry [put]
// @Security  OAuth2Application[admin,read]
func (u NcnController) NcnRetryWorkflow(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnRerunWorkflow
// @Summary   Rerun a workflow, all steps will run
// @Param     name  path  string  true  "name of workflow"
// @Tags      Workflow
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v1/workflows/{name}/rerun [put]
// @Security  OAuth2Application[admin,read]
func (u NcnController) NcnRerunWorkflow(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}
