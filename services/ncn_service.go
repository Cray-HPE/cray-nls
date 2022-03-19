package services

import (
	"github.com/Cray-HPE/cray-nls/utils"
	"gorm.io/gorm"
)

// NcnService service layer
type NcnService struct {
	logger utils.Logger
}

// NewNcnService creates a new Ncnservice
func NewNcnService(logger utils.Logger) NcnService {
	return NcnService{
		logger: logger,
	}
}

// WithTrx delegates transaction to repository database
func (s NcnService) WithTrx(trxHandle *gorm.DB) NcnService {
	return s
}
