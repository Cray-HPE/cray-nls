package routes

import (
	"github.com/Cray-HPE/cray-nls/api/controllers"
	"github.com/Cray-HPE/cray-nls/utils"
)

// UserRoutes struct
type UserRoutes struct {
	logger         utils.Logger
	handler        utils.RequestHandler
	userController controllers.UserController
}

// Setup user routes
func (s UserRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/api/v1")
	{
		api.POST("/user/:id", s.userController.UpdateUser)
	}
}

// NewUserRoutes creates new user controller
func NewUserRoutes(
	logger utils.Logger,
	handler utils.RequestHandler,
	userController controllers.UserController,
) UserRoutes {
	return UserRoutes{
		handler:        handler,
		logger:         logger,
		userController: userController,
	}
}
