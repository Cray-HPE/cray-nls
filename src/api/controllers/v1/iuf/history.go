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

// ListHistory
// @Summary  List history of an iuf activity
// @Tags     History
// @Accept   json
// @Produce  json
// @Success  200  {object}  []iuf.History
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_id}/history [get]
func (u IufController) ListHistory(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// GetHistory
// @Summary  Get a history item of an iuf activity
// @Param    start_time  path  string                            true  "start time of a history item"
// @Tags     History
// @Accept   json
// @Produce  json
// @Success  200  {object}  iuf.History
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_id}/history/{start_time} [get]
func (u IufController) GetHistory(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// ReplaceHistoryComment
// @Summary  Get a history item of an iuf activity
// @Param    start_time  path  string  true  "start time of a history item"
// @Param    activity    body  iuf.ReplaceHistoryCommentRequest  true  "Modify comment of a history"
// @Tags     History
// @Accept   json
// @Produce  json
// @Success  200  {object}  iuf.History
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_id}/history/{start_time} [patch]
func (u IufController) ReplaceHistoryComment(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// HistoryRunAction
// @Summary  Run a session
// @Param    activity_name   path  string                    true  "activity name"
// @Param    action_request  body  iuf.HistoryActionRequest  true  "Action Request"
// @Tags     History
// @Accept   json
// @Produce  json
// @Success  201  "Created"
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_name}/history/run [post]
func (u IufController) HistoryRunAction(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// HistoryBlockedAction
// @Summary  Mark a session blocked
// @Param    activity_name   path  string                    true  "activity name"
// @Param    action_request  body  iuf.HistoryActionRequest  true  "Action Request"
// @Tags     History
// @Accept   json
// @Produce  json
// @Success  201  "Created"
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_name}/history/blocked [post]
func (u IufController) HistoryBlockedAction(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// HistoryResumeAction
// @Summary  Resume an activity
// @Param    activity_name   path  string                    true  "activity name"
// @Param    action_request  body  iuf.HistoryActionRequest  true  "Action Request"
// @Tags     History
// @Accept   json
// @Produce  json
// @Success  201  "Created"
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_name}/history/resume [post]
func (u IufController) HistoryResumeAction(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// HistoryPausedAction
// @Summary  Pause a session
// @Param    activity_name   path  string                    true  "activity name"
// @Param    action_request  body  iuf.HistoryActionRequest  true  "Action Request"
// @Tags     History
// @Accept   json
// @Produce  json
// @Success  201  "Created"
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_name}/history/paused [post]
func (u IufController) HistoryPausedAction(c *gin.Context) {
	c.JSON(501, "not implemented")
}

// HistoryAbortAction
// @Summary  Abort a session
// @Param    activity_name   path  string                    true  "activity name"
// @Param    action_request  body  iuf.HistoryActionRequest  true  "Action Request"
// @Tags     History
// @Accept   json
// @Produce  json
// @Success  201  "Created"
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/activities/{activity_name}/history/abort [post]
func (u IufController) HistoryAbortAction(c *gin.Context) {
	c.JSON(501, "not implemented")
}