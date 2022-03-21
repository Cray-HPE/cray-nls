package controllers

import (
	"github.com/Cray-HPE/cray-nls/services"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/gin-gonic/gin"
)

// NcnController data type
type NcnController struct {
	service services.NcnService
	logger  utils.Logger
}

// NewNcnController creates new Ncn controller
func NewNcnController(NcnService services.NcnService, logger utils.Logger) NcnController {
	return NcnController{
		service: NcnService,
		logger:  logger,
	}
}

// NcnCreateBakcup 	perform backup action on a NCN
// @Summary      Create a NCN backup
// @Description.markdown

// @Tags     NCN
// @Param    hostname  path  string  true  "Hostname"
// @Accept   json
// @Produce  json
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /ncn/{hostname}/backup [post]
func (u NcnController) NcnCreateBakcup(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnRestoreBakcup 	restore a NCN backup
// @Summary      Restore a NCN backup
// @Description.markdown

// @Tags     NCN
// @Param    hostname  path  string  true  "Hostname"
// @Accept   json
// @Produce  json
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /ncn/{hostname}/restore [post]
func (u NcnController) NcnRestoreBakcup(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnWipe 		perform disk wipe on a NCN
// @Summary      Perform disk wipe on a NCN
// @Description.markdown

// @Tags     NCN
// @Param    hostname  path  string  true  "Hostname"
// @Accept   json
// @Produce  json
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /ncn/{hostname}/wipe [post]
func (u NcnController) NcnWipe(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnReboot 		perform reboot on a NCN
// @Summary      Perform reboot on a NCN
// @Description.markdown

// @Tags     NCN
// @Param    hostname  path  string  true  "Hostname"
// @Accept   json
// @Produce  json
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /ncn/{hostname}/reboot [post]
func (u NcnController) NcnReboot(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnPostRebuild 		perform post rebuild action on a NCN
// @Summary      Perform post rebuild action on a NCN
// @Description.markdown

// @Tags     NCN
// @Param    hostname  path  string  true  "Hostname"
// @Accept   json
// @Produce  json
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /ncn/{hostname}/post-rebuild [post]
func (u NcnController) NcnPostRebuild(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnValidate 		perform validation on a NCN
// @Summary      Perform validation on a NCN
// @Description.markdown

// @Tags     NCN
// @Param    hostname  path  string  true  "Hostname"
// @Accept   json
// @Produce  json
// @Failure  400  {object}  utils.ResponseError
// @Failure  404  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /ncn/{hostname}/validate [post]
func (u NcnController) NcnValidate(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}
