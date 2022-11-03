/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package routes

import (
	"github.com/Cray-HPE/cray-nls/src/api/controllers/v1/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
)

// IufRoutes struct
type IufRoutes struct {
	logger        utils.Logger
	handler       utils.RequestHandler
	iufController iuf.IufController
}

// Setup Iuf routes
func (s IufRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/apis/iuf/v1")
	{
		// activities CRUD
		api.POST("/activities", s.iufController.CreateActivity)
		api.GET("/activities", s.iufController.ListActivities)
		api.GET("/activities/:activity_name", s.iufController.GetActivity)
		api.PATCH("/activities/:activity_name", s.iufController.PatchActivity)
		// history CRUD
		api.GET("/activities/:activity_name/history", s.iufController.ListHistory)
		api.POST("/activities/:activity_name/history/run", s.iufController.HistoryRunAction)
		// session CRUD
		api.GET("/activities/:activity_name/sessions", s.iufController.ListSessions)
		api.GET("/activities/:activity_name/sessions/:session_name", s.iufController.GetSession)
	}
}

// NewIufRoutes creates new Iuf controller
func NewIufRoutes(
	logger utils.Logger,
	handler utils.RequestHandler,
	iufController iuf.IufController,
) IufRoutes {
	return IufRoutes{
		handler:       handler,
		logger:        logger,
		iufController: iufController,
	}
}
