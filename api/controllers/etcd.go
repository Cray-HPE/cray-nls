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
// @Failure               500  {object}  utils.ResponseError
// @Router                /etcd/{hostname}/prepare [put]
func (u EtcdController) EtcdPrepare(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Etcd updated"})
}
