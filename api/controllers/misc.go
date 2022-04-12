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
	"net/http"

	"github.com/Cray-HPE/cray-nls/api/services"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/gin-gonic/gin"
)

// use ldflags to replace this value during build:
// 		https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
const VERSION string = "development"

// MiscController data type
type MiscController struct {
	workflowService services.WorkflowService
	logger          utils.Logger
}

// NewMiscController creates new Misc controller
func NewMiscController(workflowService services.WorkflowService, logger utils.Logger) MiscController {
	return MiscController{
		workflowService: workflowService,
		logger:          logger,
	}
}

// GetVersion
// @Summary  Get version of cray-nls service
// @Tags     Misc
// @Accept   json
// @Produce  json
// @Success  200  {object}  utils.ResponseOk
// @Failure  500  {object}  utils.ResponseError
// @Router   /v1/version [get]
func (u MiscController) GetVersion(c *gin.Context) {
	c.JSON(200, utils.ResponseOk{Message: VERSION})
}

// GetReadiness
// @Summary  K8s Readiness endpoint
// @Tags     Misc
// @Accept   json
// @Produce  json
// @Success  204
// @Failure  500  {object}  utils.ResponseError
// @Router   /v1/readiness [get]
func (u MiscController) GetReadiness(c *gin.Context) {
	workflows, err := u.workflowService.GetWorkflows(c)
	if err != nil && len(workflows.Items) == 0 {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetLiveness
// @Summary  K8s Liveness endpoint
// @Tags     Misc
// @Accept   json
// @Produce  json
// @Success  204
// @Router   /v1/liveness [get]
func (u MiscController) GetLiveness(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
