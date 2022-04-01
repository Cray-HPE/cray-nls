//
//  MIT License
//
//  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
//
//  Permission is hereby granted, free of charge, to any person obtaining a
//  copy of this software and associated documentation files (the "Software"),
//  to deal in the Software without restriction, including without limitation
//  the rights to use, copy, modify, merge, publish, distribute, sublicense,
//  and/or sell copies of the Software, and to permit persons to whom the
//  Software is furnished to do so, subject to the following conditions:
//
//  The above copyright notice and this permission notice shall be included
//  in all copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
//  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
//  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
//  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
//  OTHER DEALINGS IN THE SOFTWARE.
//
package routes

import (
	"github.com/Cray-HPE/cray-nls/api/controllers"
	"github.com/Cray-HPE/cray-nls/utils"
)

// NcnRoutes struct
type NcnRoutes struct {
	logger         utils.Logger
	handler        utils.RequestHandler
	ncnsController controllers.NcnController
}

// Setup Ncn routes
func (s NcnRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/apis/nls/v1")
	{
		api.POST("/ncns/:hostname/reboot", s.ncnsController.NcnCreateRebootWorkflow)
		api.POST("/ncns/:hostname/rebuild", s.ncnsController.NcnCreateRebuildWorkflow)

	}
}

// NewNcnRoutes creates new Ncn controller
func NewNcnRoutes(
	logger utils.Logger,
	handler utils.RequestHandler,
	ncnsController controllers.NcnController,
) NcnRoutes {
	return NcnRoutes{
		handler:        handler,
		logger:         logger,
		ncnsController: ncnsController,
	}
}
