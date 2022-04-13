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
	"fmt"

	"github.com/Cray-HPE/cray-nls/api/models"
	"github.com/Cray-HPE/cray-nls/api/services"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/gin-gonic/gin"
)

// NcnController data type
type NcnController struct {
	workflowService services.WorkflowService
	logger          utils.Logger
	validator       utils.Validator
}

// NewNcnController creates new Ncn controller
func NewNcnController(workflowService services.WorkflowService, logger utils.Logger) NcnController {
	return NcnController{
		workflowService: workflowService,
		logger:          logger,
	}
}

// NcnCreateRebuildWorkflow
// @Summary   End to end rebuild of a single ncn (worker only)
// @Param     hostname  path  string  true  "hostname"
// @Tags      NCNs
// @Accept    json
// @Produce   json
// @Success   200  {object}  models.Workflow
// @Failure   400  {object}  utils.ResponseError
// @Failure   404  {object}  utils.ResponseError
// @Failure   500  {object}  utils.ResponseError
// @Router    /v1/ncns/{hostname}/rebuild [post]
// @Security  OAuth2Application[admin]
func (u NcnController) NcnCreateRebuildWorkflow(c *gin.Context) {
	hostname := c.Param("hostname")
	u.createRebuildWorkflow([]string{hostname}, c)
}

// NcnsCreateRebuildWorkflow
// @Summary   End to end rolling rebuild ncns (workers only)
// @Param     include  body  models.CreateRebuildWorkflowRequest  true  "hostnames to include"
// @Tags      NCNs
// @Accept    json
// @Produce   json
// @Success   200  {object}  models.Workflow
// @Failure   400  {object}  utils.ResponseError
// @Failure   404  {object}  utils.ResponseError
// @Failure   500  {object}  utils.ResponseError
// @Router    /v1/ncns/rebuild [post]
// @Security  OAuth2Application[admin]
func (u NcnController) NcnsCreateRebuildWorkflow(c *gin.Context) {
	var requestBody models.CreateRebuildWorkflowRequest
	if err := c.BindJSON(&requestBody); err != nil {
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(400, errResponse)
		return
	}
	u.createRebuildWorkflow(requestBody.Hosts, c)
}

func (u NcnController) createRebuildWorkflow(hostnames []string, c *gin.Context) {
	hostnames = removeDuplicateHostnames(hostnames)

	err := u.validator.ValidateWorkerHostnames(hostnames)
	if err != nil {
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(400, errResponse)
		return
	}
	u.logger.Infof("Hostnames: %v", hostnames)

	workflow, err := u.workflowService.CreateRebuildWorkflow(hostnames)

	if err != nil {
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(500, errResponse)
		return
	} else {
		myWorkflow := models.Workflow{
			Name:       workflow.Name,
			TargetNcns: []string{workflow.Labels["targetNcn"]}, //todo
		}
		c.JSON(200, myWorkflow)
		return
	}
}

func removeDuplicateHostnames(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}