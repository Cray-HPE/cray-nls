package services

import (
	"github.com/Cray-HPE/cray-nls/utils"
	"gorm.io/gorm"
)

// K8sService service layer
type K8sService struct {
	logger utils.Logger
}

// NewK8sService creates a new K8sservice
func NewK8sService(logger utils.Logger) K8sService {
	return K8sService{
		logger: logger,
	}
}

// WithTrx delegates transaction to repository database
func (s K8sService) WithTrx(trxHandle *gorm.DB) K8sService {
	return s
}
