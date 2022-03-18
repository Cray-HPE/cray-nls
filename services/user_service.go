package services

import (
	"github.com/Cray-HPE/cray-nls/utils"
	"gorm.io/gorm"
)

// UserService service layer
type UserService struct {
	logger utils.Logger
}

// NewUserService creates a new userservice
func NewUserService(logger utils.Logger) UserService {
	return UserService{
		logger: logger,
	}
}

// WithTrx delegates transaction to repository database
func (s UserService) WithTrx(trxHandle *gorm.DB) UserService {
	return s
}
