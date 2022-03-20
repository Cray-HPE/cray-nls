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

// EtcdPrepare 	prepare etcd for a master node
// @Summary               Prepare etcd on a master node
// @description.markdown  etcd-prepare
// @Param                 hostname  body  string  true  "Hostname of target first master"
// @Tags                  Etcd
// @Accept                json
// @Produce               json
// @Header                200  {string}  Token  "qwerty"
// @Failure               400  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /etcd/{hostname}/prepare [put]
func (u EtcdController) EtcdPrepare(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Etcd updated"})
}
