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

// NcnController data type
type NcnController struct {
	workflowService services.WorkflowService
	logger          utils.Logger
}

// NewNcnController creates new Ncn controller
func NewNcnController(workflowService services.WorkflowService, logger utils.Logger) NcnController {
	return NcnController{
		workflowService: workflowService,
		logger:          logger,
	}
}

// NcnCreateRebuildWorkflow
// @Summary   End to end rebuild of a single ncn
// @Param     hostname  path  string  true  "hostname"
// @Tags      NCN
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v1/ncns/{hostname}/rebuild [post]
// @Security  OAuth2Application[admin]
func (u NcnController) NcnCreateRebuildWorkflow(c *gin.Context) {
	hostname := c.Param("hostname")
	u.logger.Infof("Hostname: %s", hostname)
	u.workflowService.CreateWorkflow(hostname)

	c.JSON(200, gin.H{"data": "work flow created"})
}

// NcnCreateRebootWorkflow
// @Summary   End to end reboot of a single ncn
// @Param     hostname  path  string  true  "hostname"
// @Tags      NCN
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v1/ncns/{hostname}/reboot [post]
// @Security  OAuth2Application[admin]
func (u NcnController) NcnCreateRebootWorkflow(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnsCreateRebootsWorkflow
// @Summary   End to end rolling reboot request
// @Tags      NCN v2
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v2/ncns/reboot [post]
// @Security  OAuth2Application[admin]
func (u NcnController) NcnsCreateRebootWorkflow(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnsCreateRebuildWorkflow
// @Summary   End to end rolling rebuild request
// @Tags      NCN v2
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v2/ncns/rebuild [post]
// @Security  OAuth2Application[admin]
func (u NcnController) NcnsCreateRebuildWorkflow(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnAdd
// @Summary   Add a ncn
// @Tags      NCN v2
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v2/ncn [post]
// @Security  OAuth2Application[admin]
func (u NcnController) NcnAdd(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnRemove
// @Summary   Remove a ncn
// @Param     hostname  path  string  true  "hostname"
// @Tags      NCN v2
// @Accept    json
// @Produce   json
// @Failure   501  "Not Implemented"
// @Router    /v2/ncns/{hostname} [delete]
// @Security  OAuth2Application[admin]
func (u NcnController) NcnRemove(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}
