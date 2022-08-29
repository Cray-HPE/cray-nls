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

// NcnsCreateRebuildWorkflow
// @Summary  End to end rolling rebuild ncns
// @Param    include  body  models.CreateRebuildWorkflowRequest  true  "hostnames to include"
// @Tags     NCN Lifecycle Events
// @Accept   json
// @Produce  json
// @Success  200  {object}  models.CreateRebuildWorkflowResponse
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /v1/ncns/rebuild [post]
func (u NcnController) NcnsCreateRebuildWorkflow(c *gin.Context) {
	var requestBody models.CreateRebuildWorkflowRequest
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(400, errResponse)
		return
	}
	u.createRebuildWorkflow(requestBody, c)
}

// NcnsCreateRebootWorkflow
// @Summary  End to end rolling reboot ncns
// @Param    include  body  models.CreateRebootWorkflowRequest  true  "hostnames to include"
// @Tags     NCN Lifecycle Events
// @Accept   json
// @Produce  json
// @Success  200  {object}  models.CreateRebootWorkflowResponse
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /v1/ncns/reboot [post]
func (u NcnController) NcnsCreateRebootWorkflow(c *gin.Context) {
	var requestBody models.CreateRebootWorkflowRequest
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(400, errResponse)
		return
	}
	//u.createRebootWorkflow(requestBody, c)
}

// NcnsGetHooks
// @Summary  Get ncn lifecycle hooks
// @Tags     NCN Lifecycle Hooks
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v1/ncns/hooks [get]
func (u NcnController) NcnsGetHooks(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// NcnsAddHooks
// @Summary  Get ncn lifecycle hooks
// @Tags     NCN Lifecycle Hooks
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v1/ncns/hooks [post]
func (u NcnController) NcnsAddHooks(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// NcnsRemoveHook
// @Summary  Get ncn lifecycle hooks
// @Param    hook_id  path  string  true  "id of a hook"
// @Tags     NCN Lifecycle Hooks
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v1/ncns/hooks/{hook_id} [delete]
func (u NcnController) NcnsRemoveHook(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// NcnsUpdateHook
// @Summary  Update a ncn lifecycle hook
// @Param    hook_id  path  string  true  "id of a hook"
// @Tags     NCN Lifecycle Hooks
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v1/ncns/hooks/{hook_id} [put]
func (u NcnController) NcnsUpdateHook(c *gin.Context) {
	c.JSON(501, "not implemented")
}

func (u NcnController) createRebuildWorkflow(req models.CreateRebuildWorkflowRequest, c *gin.Context) {
	req.Hosts = removeDuplicateHostnames(req.Hosts)

	err := u.validator.ValidateHostnames(req.Hosts)
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(400, errResponse)
		return
	}
	u.logger.Infof("Hostnames: %v, dryRun: %v", req.Hosts, req.DryRun)

	workflow, err := u.workflowService.CreateRebuildWorkflow(req)

	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(500, errResponse)
		return
	} else {
		myWorkflow := models.CreateRebuildWorkflowResponse{
			Name:       workflow.Name,
			TargetNcns: req.Hosts,
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
