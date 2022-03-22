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

// K8sController data type
type K8sController struct {
	service services.K8sService
	logger  utils.Logger
}

// NewK8sController creates new K8s controller
func NewK8sController(K8sService services.K8sService, logger utils.Logger) K8sController {
	return K8sController{
		service: K8sService,
		logger:  logger,
	}
}

// K8sPreRebuild 	pre rebuild action
// @Summary               Kubernetes node pre rebuild action
// @description.markdown  k8s-pre-rebuild
// @Param                 hostname  path  string  true  "Hostname"
// @Tags                  Kubernetes
// @Accept                json
// @Produce               json
// @Failure               400  {object}  utils.ResponseError
// @Failure               401  {object}  utils.ResponseError
// @Failure               403  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /kubernetes/{hostname}/pre-rebuild [post]
// @Security              OAuth2Application[admin]
func (u K8sController) K8sPreRebuild(c *gin.Context) {
	c.JSON(200, gin.H{"data": "K8s updated"})
}

// K8sDrain 			drain a k8s node
// @Summary               Drain a Kubernetes node
// @description.markdown  k8s-drain
// @Tags                  Kubernetes
// @Accept                json
// @Produce               json
// @Param                 hostname  path      string  true  "Hostname"
// @Failure               400       {object}  utils.ResponseError
// @Failure               404       {object}  utils.ResponseError
// @Failure               500       {object}  utils.ResponseError
// @Router                /kubernetes/{hostname}/drain [post]
// @Security              OAuth2Application[admin]
func (u K8sController) K8sDrain(c *gin.Context) {
	c.JSON(200, gin.H{"data": "K8s updated"})
}

// K8sPostRebuild 		Post rebuild action
// @Summary               Kubernetes node post rebuild action
// @description.markdown  k8s-post-rebuild
// @Tags                  Kubernetes
// @Accept                json
// @Produce               json
// @Param                 hostname  path      string  true  "Hostname"
// @Failure               400       {object}  utils.ResponseError
// @Failure               404       {object}  utils.ResponseError
// @Failure               500       {object}  utils.ResponseError
// @Router                /kubernetes/{hostname}/post-rebuild [post]
// @Security              OAuth2Application[admin]
func (u K8sController) K8sPostRebuild(c *gin.Context) {
	c.JSON(200, gin.H{"data": "K8s updated"})
}
