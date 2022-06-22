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
package controllers_v2

import (
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

// NcnCreateRebootWorkflow
// @Summary  End to end reboot of a single ncn
// @Param    hostname  path  string  true  "hostname"
// @Tags     V2 NCNs
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v2/ncns/{hostname}/reboot [post]
func (u NcnController) NcnCreateRebootWorkflow(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// NcnsCreateRebootsWorkflow
// @Summary  End to end rolling reboot ncns
// @Param    include  body  []string  false  "hostnames to include"
// @Tags     V2 NCNs
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v2/ncns/reboot [post]
func (u NcnController) NcnsCreateRebootWorkflow(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// NcnsBeforeK8sDrainHook
// @Summary  Add additional steps before k8s drain
// @Tags     V2 NCN Hooks
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v2/ncns/hooks/before-k8s-drain [post]
func (u NcnController) NcnsBeforeK8sDrainHook(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// NcnsBeforeWipeHook
// @Summary  Add additional steps before wipe a ncn
// @Tags     V2 NCN Hooks
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v2/ncns/hooks/before-wipe [post]
func (u NcnController) NcnsBeforeWipeHook(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// NcnsPostBootHook
// @Summary  Add additional steps after a ncn boot(reboot)
// @Tags     V2 NCN Hooks
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v2/ncns/hooks/post-boot [post]
func (u NcnController) NcnsPostBootHook(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// NcnsGetHooks
// @Summary  Get ncn lifecycle hooks
// @Param    filter  query  string  false  "filter"
// @Tags     V2 NCN Hooks
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v2/ncns/hooks [get]
func (u NcnController) NcnsGetHooks(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// NcnsRemoveHook
// @Summary  Remove a ncn lifecycle hook
// @Param    hook_name  path  string  true  "hook_name"
// @Tags     V2 NCN Hooks
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v2/ncns/hooks/{hook_name} [delete]
func (u NcnController) NcnsRemoveHook(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// NcnAdd
// @Summary   Add a ncn
// @Tags     V2 NCNs
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router    /v2/ncn [post]

func (u NcnController) NcnAdd(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// NcnRemove
// @Summary  Remove a ncn
// @Param    hostname  path  string  true  "hostname"
// @Tags     V2 NCNs
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /v2/ncns/{hostname} [delete]
func (u NcnController) NcnRemove(c *gin.Context) {
	c.JSON(501, "not implemented")
}
