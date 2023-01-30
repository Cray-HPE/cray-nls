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
	"strconv"

	"github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/gin-gonic/gin"
)

// ListHistory
//	@Summary	List history of an iuf activity
//	@Param		activity_name	path	string	true	"activity name"
//	@Tags		History
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]iuf.History
//	@Failure	500	{object}	utils.ResponseError
//	@Router		/iuf/v1/activities/{activity_name}/history [get]
func (u IufController) ListHistory(c *gin.Context) {
	res, err := u.iufService.ListActivityHistory(c.Param("activity_name"))
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetHistory
//	@Summary	Get a history item of an iuf activity
//	@Param		start_time	path	string	true	"start time of a history item"
//	@Tags		History
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	iuf.History
//	@Failure	400	{object}	utils.ResponseError
//	@Failure	404	{object}	utils.ResponseError
//	@Failure	500	{object}	utils.ResponseError
//	@Router		/iuf/v1/activities/{activity_name}/history/{start_time} [get]
func (u IufController) GetHistory(c *gin.Context) {
	startTime, err := strconv.Atoi(c.Param("start_time"))
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	res, err := u.iufService.GetActivityHistory(c.Param("activity_name"), int32(startTime))
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	if res.StartTime == 0 {
		err := fmt.Errorf("history with start_time:%d not found", startTime)
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusNotFound, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}

// ReplaceHistoryComment
//	@Summary	replace comment of a history item of an iuf activity
//	@Param		activity_name	path	string								true	"activity name"
//	@Param		start_time		path	string								true	"start time of a history item"
//	@Param		activity		body	iuf.ReplaceHistoryCommentRequest	true	"Modify comment of a history"
//	@Tags		History
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	iuf.History
//	@Failure	400	{object}	utils.ResponseError
//	@Failure	404	{object}	utils.ResponseError
//	@Failure	500	{object}	utils.ResponseError
//	@Router		/iuf/v1/activities/{activity_name}/history/{start_time} [patch]
func (u IufController) ReplaceHistoryComment(c *gin.Context) {
	startTime, err := strconv.Atoi(c.Param("start_time"))
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	var requestBody iuf.ReplaceHistoryCommentRequest
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	res, err := u.iufService.ReplaceHistoryComment(c.Param("activity_name"), int32(startTime), requestBody)
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}

// HistoryRunAction
//	@Summary	Run a session
//	@Param		activity_name	path	string						true	"activity name"
//	@Param		action_request	body	iuf.HistoryRunActionRequest	true	"Action Request"
//	@Tags		History
//	@Accept		json
//	@Produce	json
//	@Success	201	{object}	iuf.Session
//	@Failure	500	{object}	utils.ResponseError
//	@Router		/iuf/v1/activities/{activity_name}/history/run [post]
func (u IufController) HistoryRunAction(c *gin.Context) {
	var requestBody iuf.HistoryRunActionRequest
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	res, err := u.iufService.HistoryRunAction(c.Param("activity_name"), requestBody)
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// HistoryBlockedAction
//	@Summary	Mark a session blocked
//	@Param		activity_name	path	string						true	"activity name"
//	@Param		action_request	body	iuf.HistoryActionRequest	true	"Action Request"
//	@Tags		History
//	@Accept		json
//	@Produce	json
//	@Success	201	"Created"
//	@Router		/iuf/v1/activities/{activity_name}/history/blocked [post]
func (u IufController) HistoryBlockedAction(c *gin.Context) {
	var requestBody iuf.HistoryActionRequest
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	res, err := u.iufService.HistoryBlockedAction(c.Param("activity_name"), requestBody)
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// HistoryResumeAction
//	@Summary	Resume an activity
//	@Param		activity_name	path	string						true	"activity name"
//	@Param		action_request	body	iuf.HistoryActionRequest	true	"Action Request"
//	@Tags		History
//	@Accept		json
//	@Produce	json
//	@Success	201	"Created"
//	@Router		/iuf/v1/activities/{activity_name}/history/resume [post]
func (u IufController) HistoryResumeAction(c *gin.Context) {
	var requestBody iuf.HistoryActionRequest
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	res, err := u.iufService.HistoryResumeAction(c.Param("activity_name"), requestBody)
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// HistoryPausedAction
//	@Summary	Pause a session
//	@Param		activity_name	path	string						true	"activity name"
//	@Param		action_request	body	iuf.HistoryActionRequest	true	"Action Request"
//	@Tags		History
//	@Accept		json
//	@Produce	json
//	@Success	201	"Created"
//	@Router		/iuf/v1/activities/{activity_name}/history/paused [post]
func (u IufController) HistoryPausedAction(c *gin.Context) {
	var requestBody iuf.HistoryActionRequest
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	res, err := u.iufService.HistoryPausedAction(c.Param("activity_name"), requestBody)
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// HistoryAbortAction
//	@Summary	Abort a session
//	@Param		activity_name	path	string						true	"activity name"
//	@Param		action_request	body	iuf.HistoryActionRequest	true	"Action Request"
//	@Tags		History
//	@Accept		json
//	@Produce	json
//	@Success	201	"Created"
//	@Router		/iuf/v1/activities/{activity_name}/history/abort [post]
func (u IufController) HistoryAbortAction(c *gin.Context) {
	var requestBody iuf.HistoryAbortRequest
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	res, err := u.iufService.HistoryAbortAction(c.Param("activity_name"), requestBody)
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusCreated, res)
}
