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

// NcnCreateBakcup perform backup action on a NCN
// @Summary      Create a NCN backup
// @Description  Create a NCN backup before rebuild
// @Tags         NCN
// @Param		 hostname path string true "Hostname"
// @Accept       json
// @Produce      json
// @Header       200  {string}  Token  "qwerty"
// @Failure      400  {object}  utils.ResponseError
// @Failure      404  {object}  utils.ResponseError
// @Failure      500  {object}  utils.ResponseError
// @Router       /ncn/{hostname}/bakcup [post]
func (u NcnController) NcnCreateBakcup(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnCreateBakcup perform backup action on a NCN
// @Summary      Create a NCN backup
// @Description  Create a NCN backup before rebuild
// @Tags         NCN
// @Param		 hostname path string true "Hostname"
// @Accept       json
// @Produce      json
// @Header       200  {string}  Token  "qwerty"
// @Failure      400  {object}  utils.ResponseError
// @Failure      404  {object}  utils.ResponseError
// @Failure      500  {object}  utils.ResponseError
// @Router       /ncn/{hostname}/bakcup [get]
func (u NcnController) NcnGetBakcup(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnMoveFirstMaster 	move first master to a master ncn
// @Summary      		Move first master to a master ncn
// @Description  		Move first master to a master ncn
// @Param		 		hostname body string true "Hostname of target first master"
// @Tags         		NCN
// @Accept       		json
// @Produce      		json
// @Header      		 200  {string}  Token  "qwerty"
// @Failure      		400  {object}  utils.ResponseError
// @Failure      		404  {object}  utils.ResponseError
// @Failure      		500  {object}  utils.ResponseError
// @Router       		/ncn/first-master [put]
func (u NcnController) NcnMoveFirstMaster(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnGetFirstMaster 	get hostname of first master
// @Summary      		Get hostname of first master
// @Description  		Get hostname of first master
// @Tags         		NCN
// @Accept       		json
// @Produce      		json
// @Header      		 200  {string}  Token  "qwerty"
// @Failure      		400  {object}  utils.ResponseError
// @Failure      		404  {object}  utils.ResponseError
// @Failure      		500  {object}  utils.ResponseError
// @Router       		/ncn/first-master [get]
func (u NcnController) NcnGetFirstMaster(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}
