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
	"fmt"
	"net/http"

	_ "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/gin-gonic/gin"
)

// ListSessions
// @Summary  List sessions of an IUF activity
// @Param    activity_name  path  string  true  "activity name"
// @Tags     Sessions
// @Accept   json
// @Produce  json
// @Success  200  {object}  []iuf.Session
// @Failure  500  {object}  utils.ResponseError
// @Router   /iuf/v1/activities/{activity_name}/sessions [get]
func (u IufController) ListSessions(c *gin.Context) {
	res, err := u.iufService.ListSessions(c.Param("activity_name"))
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(http.StatusInternalServerError, errResponse)
	}
	c.JSON(http.StatusOK, res)
}

// GetSession
// @Summary  Get a session of an IUF activity
// @Param    activity_name  path  string  true  "activity name"
// @Param    session_name   path  string  true  "session name"
// @Tags     Sessions
// @Accept   json
// @Produce  json
// @Success  200  {object}  iuf.Session
// @Failure  500  {object}  utils.ResponseError
// @Router   /iuf/v1/activities/{activity_name}/sessions/{session_name} [get]
func (u IufController) GetSession(c *gin.Context) {
	res, err := u.iufService.GetSession(c.Param("activity_name"), c.Param("session_name"))
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(http.StatusInternalServerError, errResponse)
	}
	c.JSON(http.StatusOK, res)
}
