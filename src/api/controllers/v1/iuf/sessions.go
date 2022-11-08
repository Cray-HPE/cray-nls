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
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
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
		return
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
	res, _, err := u.iufService.GetSession(c.Param("session_name"))
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: fmt.Sprint(err)}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (u IufController) Sync(c *gin.Context) {
	var requestBody iuf.SyncRequest
	// var response v1.SyncResponse
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		c.JSON(500, fmt.Sprint(err))
		return
	}
	session, activityRef, err := u.iufService.GetSession(requestBody.Object.Name)
	if err != nil {
		u.logger.Error(err)
		c.JSON(500, fmt.Sprint(err))
		return
	}
	var response iuf.SyncResponse
	switch session.CurrentState {
	case "":
		u.logger.Infof("State is empty, creating workflow: %s, resoure version: %s", session.Name, requestBody.Object.ObjectMeta.ResourceVersion)
		// get list of stages
		stages := session.InputParameters.Stages
		workflow, err := u.iufService.CreateIufWorkflow(session)
		if err != nil {
			u.logger.Error(err)
			c.JSON(500, fmt.Sprint(err))
			return
		}
		u.logger.Infof("workflow: %s has been created", workflow.Name)

		session.Workflows = append(session.Workflows, iuf.SessionWorkflow{Id: workflow.Name})
		session.CurrentStage = stages[0]
		session.CurrentState = iuf.SessionStateInProgress
		u.logger.Infof("Update activity state, session state: %s", session.CurrentState)
		err = u.iufService.UpdateActivityStateFromSessionState(session, activityRef)
		if err != nil {
			u.logger.Error(err)
			c.JSON(500, fmt.Sprint(err))
			return
		}
		u.logger.Infof("Update session: %v", session)
		err = u.iufService.UpdateSession(session, activityRef)
		if err != nil {
			u.logger.Error(err)
			c.JSON(500, fmt.Sprint(err))
			return
		}
		response = iuf.SyncResponse{
			ResyncAfterSeconds: 5,
		}
		c.JSON(200, response)
		return
	case iuf.SessionStateInProgress:
		activeWorkflowInfo := session.Workflows[len(session.Workflows)-1]
		activeWorkflow, _ := u.workflowService.GetWorkflowByName(activeWorkflowInfo.Id, c)
		if activeWorkflow.Status.Phase == v1alpha1.WorkflowRunning {
			u.logger.Infof("Workflow is still running: %s", activeWorkflowInfo.Id)
			response = iuf.SyncResponse{
				ResyncAfterSeconds: 5,
			}
			c.JSON(200, response)
			return
		}
		if activeWorkflow.Status.Phase == v1alpha1.WorkflowError || activeWorkflow.Status.Phase == v1alpha1.WorkflowFailed {
			u.logger.Infof("Workflow is in failed/error state: %s,resource version: %s", activeWorkflowInfo.Id, requestBody.Object.ObjectMeta.ResourceVersion)
			session.CurrentState = iuf.SessionStateDebug
			u.iufService.UpdateActivityStateFromSessionState(session, activityRef)
			u.iufService.UpdateSession(session, activityRef)
			response = iuf.SyncResponse{}
			c.JSON(200, response)
			return
		}
		// todo: move to next stage or complete
	case iuf.SessionStatePaused, iuf.SessionStateDebug, iuf.SessionStateCompleted:
		u.logger.Infof("session state: %s", session.CurrentState)
		response = iuf.SyncResponse{}
		c.JSON(200, response)
		return
	default:
		err := fmt.Errorf("unknow state: session.CurrentState")
		u.logger.Error(err)
		c.JSON(500, utils.ResponseError{Message: fmt.Sprint(err)})
		return
	}
}
