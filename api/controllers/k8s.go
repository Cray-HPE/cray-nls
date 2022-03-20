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

// K8sMoveFirstMaster 	move first master to a master k8s
// @Summary               Move first master to a master k8s
// @description.markdown  k8s-move-first-master
// @Param                 hostname  body  string  true  "Hostname of target first master"
// @Tags                  Kubernetes
// @Accept                json
// @Produce               json
// @Header                200  {string}  Token  "qwerty"
// @Failure               400  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /kubernetes/first-master [put]
func (u K8sController) K8sMoveFirstMaster(c *gin.Context) {
	c.JSON(200, gin.H{"data": "K8s updated"})
}

// K8sDrain 			drain a k8s node
// @Summary               Drain a Kubernetes node
// @description.markdown  k8s-drain
// @Tags                  Kubernetes
// @Accept                json
// @Produce               json
// @Header                200       {string}  Token   "qwerty"
// @Param                 hostname  path      string  true  "Hostname"
// @Failure               400       {object}  utils.ResponseError
// @Failure               404       {object}  utils.ResponseError
// @Failure               500       {object}  utils.ResponseError
// @Router                /kubernetes/{hostname}/drain [post]
func (u K8sController) K8sDrain(c *gin.Context) {
	c.JSON(200, gin.H{"data": "K8s updated"})
}

// K8sPostRebuild 		Post rebuild action
// @Summary               Kubernetes node post rebuild action
// @description.markdown  k8s-post-rebuild
// @Tags                  Kubernetes
// @Accept                json
// @Produce               json
// @Header                200       {string}  Token   "qwerty"
// @Param                 hostname  path      string  true  "Hostname"
// @Failure               400       {object}  utils.ResponseError
// @Failure               404       {object}  utils.ResponseError
// @Failure               500       {object}  utils.ResponseError
// @Router                /kubernetes/{hostname}/post-rebuild [post]
func (u K8sController) K8sPostRebuild(c *gin.Context) {
	c.JSON(200, gin.H{"data": "K8s updated"})
}
