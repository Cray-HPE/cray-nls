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
	_ "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/gin-gonic/gin"
)

// ListSessions
// @Summary  List sessions of an IUF activity
// @Param    activity_uid  path  string  true  "activity uid"
// @Tags     Sessions
// @Accept   json
// @Produce  json
// @Success  200  {object}  []iuf.Session
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_uid}/sessions [get]
func (u IufController) ListSessions(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// GetSession
// @Summary  Get a session of an IUF activity
// @Param    activity_uid  path  string  true  "activity uid"
// @Param    session_uid   path  string  true  "session uid"
// @Tags     Sessions
// @Accept   json
// @Produce  json
// @Success  200  {object}  iuf.Session
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_uid}/sessions/{session_uid} [get]
func (u IufController) GetSession(c *gin.Context) {
	c.JSON(501, "not implemented")
}
