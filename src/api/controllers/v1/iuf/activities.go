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

	"github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/gin-gonic/gin"
)

// CreateActivity
// @Summary  Create an IUF activity
// @Param    activity  body  iuf.CreateActivityRequest  true  "IUF activity"
// @Tags     Activities
// @Accept   json
// @Produce  json
// @Success  201  {object}  iuf.Activity
// @Failure  400  {object}  utils.ResponseError
// @Failure  500  {object}  utils.ResponseError
// @Router   /iuf/v1/activities [post]
func (u IufController) CreateActivity(c *gin.Context) {
	var requestBody iuf.CreateActivityRequest
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	res, err := u.iufService.CreateActivity(requestBody)
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// ListActivities
// @Summary  List IUF activities
// @Tags     Activities
// @Accept   json
// @Produce  json
// @Success  200  {object}  []iuf.Activity
// @Failure  500  {object}  utils.ResponseError
// @Router   /iuf/v1/activities [get]
func (u IufController) ListActivities(c *gin.Context) {
	res, err := u.iufService.ListActivities()
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetActivity
// @Summary  Get an IUF activity
// @Param    activity_name  path  string  true  "activity name"
// @Tags     Activities
// @Accept   json
// @Produce  json
// @Success  200  {object}  iuf.Activity
// @Failure  500  {object}  utils.ResponseError
// @Router   /iuf/v1/activities/{activity_name} [get]
func (u IufController) GetActivity(c *gin.Context) {
	res, err := u.iufService.GetActivity(c.Param("activity_name"))
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}
