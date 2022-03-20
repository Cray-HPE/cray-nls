package services

import (
	"github.com/Cray-HPE/cray-nls/utils"
	"gorm.io/gorm"
)

// EtcdService service layer
type EtcdService struct {
	logger utils.Logger
}

// NewEtcdService creates a new Etcdservice
func NewEtcdService(logger utils.Logger) EtcdService {
	return EtcdService{
		logger: logger,
	}
}

// WithTrx delegates transaction to repository database
func (s EtcdService) WithTrx(trxHandle *gorm.DB) EtcdService {
	return s
}
