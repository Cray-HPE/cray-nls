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
//	@Summary	List sessions of an IUF activity
//	@Param		activity_name	path	string	true	"activity name"
//	@Tags		Sessions
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]iuf.Session
//	@Failure	500	{object}	utils.ResponseError
//	@Router		/iuf/v1/activities/{activity_name}/sessions [get]
func (u IufController) ListSessions(c *gin.Context) {
	res, err := u.iufService.ListSessions(c.Param("activity_name"))
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetSession
//	@Summary	Get a session of an IUF activity
//	@Param		activity_name	path	string	true	"activity name"
//	@Param		session_name	path	string	true	"session name"
//	@Tags		Sessions
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	iuf.Session
//	@Failure	500	{object}	utils.ResponseError
//	@Router		/iuf/v1/activities/{activity_name}/sessions/{session_name} [get]
func (u IufController) GetSession(c *gin.Context) {
	res, err := u.iufService.GetSession(c.Param("session_name"))
	if err != nil {
		u.logger.Error(err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (u IufController) Sync(c *gin.Context) {
	var requestBody iuf.SyncRequest
	if err := c.BindJSON(&requestBody); err != nil {
		u.logger.Error(err)
		c.JSON(500, err.Error())
		return
	}
	session, err := u.iufService.GetSession(requestBody.Object.Name)
	if err != nil {
		u.logger.Error(err)
		c.JSON(500, err.Error())
		return
	}
	var response iuf.SyncResponse
	switch session.CurrentState {
	case "":
		u.logger.Infof("State is empty, creating workflow: %s, resource version: %s", session.Name, requestBody.Object.ObjectMeta.ResourceVersion)
		response, err, _ := u.iufService.RunNextStage(&session)
		if err != nil {
			c.JSON(500, err.Error())
		}
		c.JSON(200, response)
		return
	case iuf.SessionStateInProgress:
		if len(session.Workflows) == 0 {
			break
		}

		activeWorkflowInfo := session.Workflows[len(session.Workflows)-1]
		activeWorkflow, _ := u.workflowService.GetWorkflowByName(activeWorkflowInfo.Id, c)
		if activeWorkflow.Status.Phase == v1alpha1.WorkflowRunning {
			// set the session back to in progress if the workflow is running.
			if session.CurrentState != iuf.SessionStateInProgress {
				session.CurrentState = iuf.SessionStateInProgress
				u.iufService.UpdateSessionAndActivity(session)

				// note: if there was an error in UpdateSession above, then we would resync anyway below after x seconds
			}

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
			err = u.iufService.UpdateSessionAndActivity(session)
			var response iuf.SyncResponse
			if err != nil {
				response = iuf.SyncResponse{
					ResyncAfterSeconds: 5,
				}
			} else {
				response = iuf.SyncResponse{}
			}
			c.JSON(200, response)
			return
		}

		if activeWorkflow.Status.Phase == v1alpha1.WorkflowSucceeded {
			err := u.iufService.ProcessOutput(&session, activeWorkflow)
			if err != nil {
				u.logger.Error(err)
				c.JSON(500, err.Error())
			}

			u.logger.Infof("Stage: %s is Succeeded, move to next stage", session.CurrentStage)
			response, err, _ := u.iufService.RunNextStage(&session)
			if err != nil {
				u.logger.Errorf("Unable to go to next stage: %v", err)
				// note: do NOT automatically retry -- we don't know whether CurrentStage has already been updated
				//  This is the downside of using a non-transactional storage such as CRDs.
				c.JSON(500, iuf.SyncResponse{})
				return
			}

			c.JSON(200, response)
			return
		}
	case iuf.SessionStatePaused, iuf.SessionStateDebug, iuf.SessionStateCompleted:
		u.logger.Infof("The session %s is in state: %s", session.Name, session.CurrentState)
		response = iuf.SyncResponse{}
		c.JSON(200, response)
		return
	default:
		err := fmt.Errorf("unknow state: session.CurrentState")
		u.logger.Error(err)
		c.JSON(500, utils.ResponseError{Message: err.Error()})
		return
	}
}
