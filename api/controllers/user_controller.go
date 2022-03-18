package controllers

import (
	"github.com/Cray-HPE/cray-nls/services"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/gin-gonic/gin"
)

// UserController data type
type UserController struct {
	service services.UserService
	logger  utils.Logger
}

// NewUserController creates new user controller
func NewUserController(userService services.UserService, logger utils.Logger) UserController {
	return UserController{
		service: userService,
		logger:  logger,
	}
}

// UpdateUser updates user
func (u UserController) UpdateUser(c *gin.Context) {
	c.JSON(200, gin.H{"data": "user updated"})
}
