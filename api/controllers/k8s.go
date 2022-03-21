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

// K8sMoveFirstMaster 	move first master from a master k8s node
// @Summary               Move first master from a master k8s node
// @description.markdown  k8s-move-first-master
// @Param                 hostname  path  string  true  "Hostname"
// @Tags                  Kubernetes
// @Accept                json
// @Produce               json
// @Failure               400  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /kubernetes/{hostname}/move-first-master [post]
func (u K8sController) K8sMoveFirstMaster(c *gin.Context) {
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
func (u K8sController) K8sPostRebuild(c *gin.Context) {
	c.JSON(200, gin.H{"data": "K8s updated"})
}
