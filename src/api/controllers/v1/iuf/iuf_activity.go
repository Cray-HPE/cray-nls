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

	v1 "github.com/Cray-HPE/cray-nls/src/api/models/iuf/v1"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

// IufActivitySync
func (u IufController) IufActivitySync(c *gin.Context) {
	var requestBody v1.IufActivitiesSyncRequest
	var response v1.IufActivitiesSyncResponse
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(400, errResponse)
		return
	}
	response.Status = requestBody.Parent.Status
	// assuming we will always resync until it reaches end state
	response.ResyncAfterSeconds = 10

	if requestBody.Parent.Spec.IsCompleted {
		// activity is marked as done, no-op
		response.ResyncAfterSeconds = 0
	} else {
		// fetch all sessions belong to current acitivy
		var err error
		response.Status.Sessions, err = u.iufService.GetSessionsByActivityName(requestBody.Parent.Metadata.Name)
		if err != nil {
			u.logger.Error(err)
			errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
			c.JSON(400, errResponse)
			return
		}

		if len(response.Status.Sessions) > 0 {
			lastSession := response.Status.Sessions[len(response.Status.Sessions)-1]
			if lastSession.Status.CurrentState.Type == v1.IufSessionStageInProgress {
				u.logger.Warnf("Session: %s is in progress, sync again", lastSession.Metadata.Name)
				c.JSON(200, response)
				return
			}
		}

		if !cmp.Equal(requestBody.Parent.Spec.SharedInput, response.Status.SharedInput) {
			u.logger.Info("input changed, reprocessing artifacts")
			//TODO: ^
			response.Status.SharedInput = requestBody.Parent.Spec.SharedInput
		}
		u.logger.Warn("Sync again")
		c.JSON(200, response)
	}

}
