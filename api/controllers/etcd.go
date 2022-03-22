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

// EtcdController data type
type EtcdController struct {
	service services.EtcdService
	logger  utils.Logger
}

// NewEtcdController creates new Etcd controller
func NewEtcdController(EtcdService services.EtcdService, logger utils.Logger) EtcdController {
	return EtcdController{
		service: EtcdService,
		logger:  logger,
	}
}

// EtcdPrepare 	prepare baremetal etcd for a master node to rejoin
// @Summary               Prepare baremetal etcd for a master node to rejoin
// @description.markdown  etcd-prepare
// @Param                 hostname  path  string  true  "Hostname of target ncn"
// @Tags                  Etcd
// @Accept                json
// @Produce               json
// @Success               200  {string}  string  "ok"
// @Failure               400  {object}  utils.ResponseError
// @Failure               401  {object}  utils.ResponseError
// @Failure               403  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /etcd/{hostname}/prepare [put]
// @Security              OAuth2Application[admin]
func (u EtcdController) EtcdPrepare(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Etcd updated"})
}
