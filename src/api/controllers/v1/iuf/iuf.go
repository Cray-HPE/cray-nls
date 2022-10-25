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
package iuf

import (
	_ "github.com/Cray-HPE/cray-nls/src/api/models/iuf/v1"
	"github.com/Cray-HPE/cray-nls/src/api/services"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/gin-gonic/gin"
)

// IufController data type
type IufController struct {
	workflowService services.WorkflowService
	iufService      services.IufService
	logger          utils.Logger
}

// NewIufController creates new Ncn controller
func NewIufController(workflowService services.WorkflowService, iufService services.IufService, logger utils.Logger) IufController {
	return IufController{
		workflowService: workflowService,
		iufService:      iufService,
		logger:          logger,
	}
}

// CreateIufActivity
// @Summary  Create an IUF activity
// @Param    activity  body  v1.IufActivity  true  "IUF activity"
// @Tags     Activities
// @Accept   json
// @Produce  json
// @Success  200  {object}  v1.IufActivity
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activity [post]
func (u IufController) CreateIufActivity(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// ListIufActivities
// @Summary  List IUF activities
// @Tags     Activities
// @Accept   json
// @Produce  json
// @Success  200  {object}  []v1.IufActivity
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities [get]
func (u IufController) ListIufActivities(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// GetIufActivity
// @Summary  Get an IUF activity
// @Param    id                path  string          true  "activity uid"
// @Tags     Activities
// @Accept   json
// @Produce  json
// @Success  200  {object}  v1.IufActivity
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{id} [get]
func (u IufController) GetIufActivity(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// PatchIufActivity
// @Summary  Patch an IUF activity
// @Param    id  path  string  true  "activity uid"
// @Param    partial_activity  body  v1.IufActivity  true  "partial IUF activity"
// @Tags     Activities
// @Accept   json
// @Produce  json
// @Success  200  {object}  v1.IufActivity
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{id} [patch]
func (u IufController) PatchIufActivity(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// ListIufSessions
// @Summary  List sessions of an IUF activity
// @Param    activity_uid  path  string         true  "activity uid"
// @Tags     Sessions
// @Accept   json
// @Produce  json
// @Success  200  {object}  []v1.IufSession
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_uid}/sessions [get]
func (u IufController) ListIufSessions(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// GetIufSession
// @Summary  Get a session of an IUF activity
// @Param    activity_uid  path  string  true  "activity uid"
// @Param    session_uid   path  string  true  "session uid"
// @Tags     Sessions
// @Accept   json
// @Produce  json
// @Success  200  {object}  v1.IufSession
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_uid}/sessions/{session_uid} [get]
func (u IufController) GetIufSession(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// CreateIufSession
// @Summary  Create a new session of an IUF activity
// @Param    activity_uid  path  string  true  "activity uid"
// @Param    session       body  v1.IufSession  true  "IUF session"
// @Tags     Sessions
// @Accept   json
// @Produce  json
// @Success  200  {object}  []v1.IufSession
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_uid}/session [post]
func (u IufController) CreateIufSession(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// StopIufSession
// @Summary  Stop a running session of an IUF activity
// @Param    activity_uid  path  string  true  "activity uid"
// @Param    session_uid   path  string  true  "session uid"
// @Tags     Sessions
// @Accept   json
// @Produce  json
// @Success  200  {object}  v1.IufSession
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_uid}/sessions/{session_uid} [delete]
func (u IufController) StopIufSession(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// ResumeIufSession
// @Summary  Resume a stopped session of an IUF activity
// @Param    activity_uid  path  string  true  "activity uid"
// @Param    session_uid   path  string  true  "session uid"
// @Tags     Sessions
// @Accept   json
// @Produce  json
// @Success  200  {object}  v1.IufSession
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_uid}/sessions/{session_uid}/resume [post]
func (u IufController) ResumeIufSession(c *gin.Context) {
	c.JSON(501, "not implemented")
}
