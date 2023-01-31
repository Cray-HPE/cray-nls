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
	"net/http"

	"github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/gin-gonic/gin"
)

// CreateActivity
//	@Summary	Create an IUF activity
//	@Param		activity	body	iuf.CreateActivityRequest	true	"IUF activity"
//	@Tags		Activities
//	@Accept		json
//	@Produce	json
//	@Success	201	{object}	iuf.Activity
//	@Failure	400	{object}	utils.ResponseError
//	@Failure	500	{object}	utils.ResponseError
//	@Router		/iuf/v1/activities [post]
func (u IufController) CreateActivity(c *gin.Context) {
	u.logger.Infof("CreateActivity: received request with params %#v", c.Request.Form)
	var requestBody iuf.CreateActivityRequest
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Errorf("CreateActivity: An error occurred while parsing parameters for an activity named %s: %v", requestBody.Name, err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	res, err := u.iufService.CreateActivity(requestBody)
	if err != nil {
		u.logger.Errorf("CreateActivity: An error occurred while creating an activity named %s: %v", requestBody.Name, err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// ListActivities
//	@Summary	List IUF activities
//	@Tags		Activities
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]iuf.Activity
//	@Failure	500	{object}	utils.ResponseError
//	@Router		/iuf/v1/activities [get]
func (u IufController) ListActivities(c *gin.Context) {
	u.logger.Infof("ListActivities: received request with params %#v", c.Request.Form)
	res, err := u.iufService.ListActivities()
	if err != nil {
		u.logger.Errorf("ListActivities: An error occurred while listing activities: %v", err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetActivity
//	@Summary	Get an IUF activity
//	@Param		activity_name	path	string	true	"activity name"
//	@Tags		Activities
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	iuf.Activity
//	@Failure	404	{object}	utils.ResponseError
//	@Failure	500	{object}	utils.ResponseError
//	@Router		/iuf/v1/activities/{activity_name} [get]
func (u IufController) GetActivity(c *gin.Context) {
	activityName := c.Param("activity_name")
	u.logger.Infof("GetActivity: received request for activity %s with params %#v", activityName, c.Request.Form)
	res, err := u.iufService.GetActivity(activityName)
	if err != nil {
		u.logger.Errorf("GetActivity: An error occurred while fetching activity %s: %v", activityName, err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusNotFound, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}

// PatchActivity
//	@Summary	Patches an existing IUF activity
//	@Param		activity_name	path	string	true	"activity name"
//	@Tags		Activities
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	iuf.Activity
//	@Failure	400	{object}	utils.ResponseError
//	@Failure	404	{object}	utils.ResponseError
//	@Failure	500	{object}	utils.ResponseError
//	@Router		/iuf/v1/activities/{activity_name} [patch]
func (u IufController) PatchActivity(c *gin.Context) {
	var requestBody iuf.PatchActivityRequest
	name := c.Param("activity_name")
	u.logger.Infof("PatchActivity: received request for activity %s with params %#v", name, c.Request.Form)
	activity, err := u.iufService.GetActivity(name)
	if err != nil {
		u.logger.Errorf("PatchActivity: An error occurred while fetching activity %s: %v", name, err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusNotFound, errResponse)
		return
	}

	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Errorf("PatchActivity: An error occurred parsing request parameters for activity %s: %v", name, err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	res, err := u.iufService.PatchActivity(activity, requestBody)
	if err != nil {
		u.logger.Errorf("PatchActivity: An error occurred patching activity %s: %v", name, err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}
