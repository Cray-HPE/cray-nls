package routes

import (
	"github.com/Cray-HPE/cray-nls/api/controllers"
	"github.com/Cray-HPE/cray-nls/utils"
)

// NcnRoutes struct
type NcnRoutes struct {
	logger        utils.Logger
	handler       utils.RequestHandler
	ncnController controllers.NcnController
}

// Setup Ncn routes
func (s NcnRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/api/v1")
	{
		api.POST("/ncn/:hostname/backup", s.ncnController.NcnCreateBakcup)
	}
}

// NewNcnRoutes creates new Ncn controller
func NewNcnRoutes(
	logger utils.Logger,
	handler utils.RequestHandler,
	ncnController controllers.NcnController,
) NcnRoutes {
	return NcnRoutes{
		handler:       handler,
		logger:        logger,
		ncnController: ncnController,
	}
}
