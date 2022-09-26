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
	misc_controllers "github.com/Cray-HPE/cray-nls/src/api/controllers/v1/misc"
	"github.com/Cray-HPE/cray-nls/src/utils"
)

// MiscRoutes struct
type MiscRoutes struct {
	logger         utils.Logger
	handler        utils.RequestHandler
	miscController misc_controllers.MiscController
}

// Setup Misc routes
func (s MiscRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/apis/nls/v1")
	{
		api.GET("/liveness", s.miscController.GetLiveness)
		api.GET("/readiness", s.miscController.GetReadiness)
		api.GET("/version", s.miscController.GetVersion)

	}
}

// NewMiscRoutes creates new Misc controller
func NewMiscRoutes(
	logger utils.Logger,
	handler utils.RequestHandler,
	miscController misc_controllers.MiscController,
) MiscRoutes {
	return MiscRoutes{
		handler:        handler,
		logger:         logger,
		miscController: miscController,
	}
}
