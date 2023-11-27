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
	"encoding/json"
	"fmt"
	services_iuf "github.com/Cray-HPE/cray-nls/src/api/services/iuf"
	"net/http"
	"time"

	"github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/gin-gonic/gin"
)

const RESYNC_TIME_IN_SECONDS = 5

// ListSessions
//
//	@Summary	List sessions of an IUF activity
//	@Param		activity_name	path	string	true	"activity name"
//	@Tags		Sessions
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]iuf.Session
//	@Failure	500	{object}	utils.ResponseError
//	@Router		/iuf/v1/activities/{activity_name}/sessions [get]
func (u IufController) ListSessions(c *gin.Context) {
	activityName := c.Param("activity_name")
	u.logger.Infof("ListSessions: received request for activity %s with params %#v", activityName, c.Request.Form)

	res, err := u.iufService.ListSessions(activityName)
	if err != nil {
		u.logger.Errorf("ListSessions: An error occurred listing sessions for activity %s: %v", activityName, err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetSession
//
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
	sessionName := c.Param("session_name")
	u.logger.Infof("GetSession: received request for session %s with params %#v", sessionName, c.Request.Form)

	res, err := u.iufService.GetSession(sessionName)
	if err != nil {
		u.logger.Errorf("GetSession: An error occurred getting session %s: %v", sessionName, err)
		errResponse := utils.ResponseError{Message: err.Error()}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (u IufController) Sync(context *gin.Context) {
	var requestBody iuf.SyncRequest
	if err := context.BindJSON(&requestBody); err != nil {
		u.logger.Errorf("Sync.1: An error occurred parsing sync request %#v: %v", context.Request.Form, err)
		context.JSON(500, utils.ResponseError{Message: err.Error()})
		return
	}
	sessionName := requestBody.Object.Name
	session, err := u.iufService.GetSession(sessionName)
	if err != nil {
		u.logger.Errorf("Sync.2: An error occurred getting session %s: %v", sessionName, err)
		context.JSON(500, utils.ResponseError{Message: err.Error()})
		return
	}

	// this is a poor-man's distributed lock. Unfortunately, in the absence of proper distributed caching, we are going
	//  to have to make do with locking using the state of the session. But note that session is stored in etcd, which
	//  is eventually consistent :\
	if session.CurrentState == iuf.SessionStateTransitioning {
		// we are already transitioning. Let's wait and try again in case the other pod/thread has died.
		response := iuf.SyncResponse{
			ResyncAfterSeconds: 60,
		}
		context.JSON(200, response)
		return
	} else if session.CurrentState == iuf.SessionStateInProgress || session.CurrentState == "" {
		u.logger.Infof("Sync.2: Setting session %s of activity %s to transitioning to prevent reentrants.", sessionName, session.ActivityRef)
		session.CurrentState = iuf.SessionStateTransitioning
		err := u.iufService.UpdateSession(session)

		// reset to anything but transitioning at the end.
		defer func() {
			u.logger.Infof("Sync.2: Trying to set session %s of activity %s out of transitioning to prevent reentrants.", sessionName, session.ActivityRef)
			session, err := u.iufService.GetSession(sessionName)
			if err != nil {
				u.logger.Errorf("Sync.defer.1: An error occurred getting session %s: %v", sessionName, err)
				return
			}

			if session.CurrentState == iuf.SessionStateTransitioning {
				u.logger.Infof("Sync.3: Session %s of activity %s is in transitioning state. Setting it to in progress to prevent reentrants.", sessionName, session.ActivityRef)
				// if no one changed the session state, then by default we assume in progress because that's what we started with
				session.CurrentState = iuf.SessionStateInProgress
				err := u.iufService.UpdateSession(session)
				if err != nil {
					u.logger.Infof("Sync.defer.3.1: Could not set session %s of activity %s to back to in_progress at the end. %#v", sessionName, session.ActivityRef, err)
				}
			} else {
				u.logger.Infof("Sync.4: Session %s of activity %s is not in transitioning state (%s). Not doing anything. Setting it to in progress to prevent reentrants.", sessionName, session.ActivityRef, session.CurrentState)
			}
		}()

		if err != nil {
			u.logger.Infof("Sync.2: Could not set session %s of activity %s to transitioning to prevent reentrants. %#v", sessionName, session.ActivityRef, err)

			// try to reset it and try again
			session.CurrentState = iuf.SessionStateInProgress
			u.iufService.UpdateSession(session)

			response := iuf.SyncResponse{
				ResyncAfterSeconds: 60,
			}
			context.JSON(200, response)
			return
		}

		// wait until the session is set to transitioning before proceeding further
		tries := 0
		for isStateTransitioning := false; isStateTransitioning == false && tries < 5; {
			time.Sleep(2 * time.Second)
			tries++
			session, err = u.iufService.GetSession(sessionName)
			isStateTransitioning = err == nil && session.CurrentState == iuf.SessionStateTransitioning
		}

		// if the session is still not set to transitioning, then give up
		if session.CurrentState != iuf.SessionStateTransitioning {
			u.logger.Infof("Sync.2: Giving up on trying to set session %s of activity %s to transitioning to prevent reentrants. %#v", sessionName, session.ActivityRef, err)
			response := iuf.SyncResponse{
				ResyncAfterSeconds: 60,
			}
			context.JSON(200, response)
			return
		}
	}

	err = u.iufService.SyncWorkflowsToSession(&session)
	if err != nil {
		u.logger.Warnf("Sync.3: State is empty, creating workflow: %s, resource version: %s, session: %s, activity: %s", session.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef)
	}

	var response iuf.SyncResponse
	switch session.CurrentState {
	case "":
		u.logger.Infof("Sync: State is empty, creating workflow: %s, resource version: %s, session: %s, activity: %s", session.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef)
		response, err, _ := u.iufService.RunNextStage(&session)
		if err != nil {
			context.JSON(500, utils.ResponseError{Message: err.Error()})
			return
		}
		context.JSON(200, response)
		return
	case iuf.SessionStateInProgress:
		activeWorkflow := u.iufService.FindLastWorkflowForCurrentStage(&session)
		if activeWorkflow == nil {
			u.restartCurrentStageFromSyncCall(context, session, requestBody, response)
			return
		}

		u.logger.Infof("Sync: Going to sync with the workflow %s for session %s in activity %s. Also, .ObjectMeta.Labels: %#v, .Labels: %#v", activeWorkflow.Name, sessionName, session.ActivityRef, activeWorkflow.ObjectMeta.Labels, activeWorkflow.Labels)

		if activeWorkflow.Status.Phase == v1alpha1.WorkflowRunning || activeWorkflow.Status.Phase == v1alpha1.WorkflowPending {
			u.logger.Infof("Sync: Workflow %s is still running for session %s in activity %s", activeWorkflow.Name, sessionName, session.ActivityRef)
			response = iuf.SyncResponse{
				ResyncAfterSeconds: RESYNC_TIME_IN_SECONDS,
			}
			context.JSON(200, response)
			return
		} else if activeWorkflow.Status.Phase == v1alpha1.WorkflowError || activeWorkflow.Status.Phase == v1alpha1.WorkflowFailed {
			u.logger.Infof("Sync: Workflow is in failed/error state. Workflow: %s, resource version: %s, session: %s, activity: %s", activeWorkflow.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef)

			// still extract the outputs from the successful steps so that if we restart we can skip over those steps.
			u.doProcessOutputs(activeWorkflow, &session, requestBody, sessionName)

			// don't do anything if session has already been aborted.
			if session.CurrentState == iuf.SessionStateAborted {
				context.JSON(200, response)
				return
			}

			var response iuf.SyncResponse

			// if this was a partial workflow, let the processing for partial workflow do the work
			if activeWorkflow.ObjectMeta.Labels[services_iuf.LABEL_PARTIAL_WORKFLOW] == "true" {
				u.logger.Infof("Sync: Stage: %s has a partial workflow that failed, moving on to the remaining products in the next workflow. Workflow failed: %s, resource version: %s, session: %s, activity: %s", session.CurrentStage, activeWorkflow.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef)
				response, err, _ = u.iufService.RunNextPartialWorkflow(&session)
				if err != nil {
					u.logger.Errorf("Sync: Unable to run the next set of products for the current stage or go to next stage. Current stage: %s, workflow: %s, resource version: %s, session: %s, activity: %s, error: %v", session.CurrentStage, activeWorkflow.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef, err)
					// note: do NOT automatically retry -- we don't know whether CurrentStage has already been updated
					//  This is the downside of using a non-transactional storage such as CRDs.
					context.JSON(500, utils.ResponseError{Message: err.Error()})
					return
				}
			} else {
				u.logger.Infof("Sync: Stage: %s a workflow that failed, and since it was not a partial workflow, setting the session state to DEBUG. Workflow failed: %s, resource version: %s, session: %s, activity: %s, .ObjectMeta.Labels: %#v, .Labels: %#v", session.CurrentStage, activeWorkflow.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef, activeWorkflow.ObjectMeta.Labels, activeWorkflow.Labels)
				session.CurrentState = iuf.SessionStateDebug
				err = u.iufService.UpdateSessionAndActivity(session, fmt.Sprintf("Failed workflow %s", activeWorkflow.Name))
				if err != nil {
					response = iuf.SyncResponse{
						ResyncAfterSeconds: RESYNC_TIME_IN_SECONDS,
					}
				} else {
					response = iuf.SyncResponse{}
				}
			}

			context.JSON(200, response)
			return
		} else if activeWorkflow.Status.Phase == v1alpha1.WorkflowSucceeded {
			u.doProcessOutputs(activeWorkflow, &session, requestBody, sessionName)

			u.logger.Infof("Sync: Stage: %s succeeded, move to the next stage. Workflow: %s, resource version: %s, session: %s, activity: %s", session.CurrentStage, activeWorkflow.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef)
			currentStage := session.CurrentStage
			var response iuf.SyncResponse

			if activeWorkflow.ObjectMeta.Labels[services_iuf.LABEL_PARTIAL_WORKFLOW] == "true" {
				u.logger.Infof("Sync: Stage: %s has a partial workflow that succeeded, moving on to the remaining products in the next workflow. Workflow completed: %s, resource version: %s, session: %s, activity: %s", session.CurrentStage, activeWorkflow.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef)
				response, err, _ = u.iufService.RunNextPartialWorkflow(&session)
				if err != nil {
					u.logger.Errorf("Sync: Unable to run the next set of products for the current stage or go to next stage. Current stage: %s, workflow: %s, resource version: %s, session: %s, activity: %s, error: %v", currentStage, activeWorkflow.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef, err)
					// note: do NOT automatically retry -- we don't know whether CurrentStage has already been updated
					//  This is the downside of using a non-transactional storage such as CRDs.
					context.JSON(500, utils.ResponseError{Message: err.Error()})
					return
				}
			} else {
				response, err, _ = u.iufService.RunNextStage(&session)
				if err != nil {
					u.logger.Errorf("Sync: Unable to go to next stage. Current stage: %s, workflow: %s, resource version: %s, session: %s, activity: %s, error: %v", currentStage, activeWorkflow.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef, err)
					// note: do NOT automatically retry -- we don't know whether CurrentStage has already been updated
					//  This is the downside of using a non-transactional storage such as CRDs.
					context.JSON(500, utils.ResponseError{Message: err.Error()})
					return
				}
			}

			context.JSON(200, response)
			return
		} else {
			context.JSON(200, response)
			return
		}
	case iuf.SessionStateTransitioning, iuf.SessionStateAborted, iuf.SessionStatePaused, iuf.SessionStateDebug, iuf.SessionStateCompleted:
		u.logger.Infof("Sync: The session %s in activity %s is in state %s and there is nothing to do", session.Name, session.ActivityRef, session.CurrentState)
		context.JSON(200, iuf.SyncResponse{})
		return
	default:
		session.CurrentState = iuf.SessionStateDebug
		err = u.iufService.UpdateSessionAndActivity(session, fmt.Sprintf("Unknown state %s", session.CurrentState))
		if err != nil {
			context.JSON(500, utils.ResponseError{Message: err.Error()})
			return
		}

		err = fmt.Errorf("sync: unknown state %s for session %s in activity %s", session.CurrentState, sessionName, session.ActivityRef)
		u.logger.Error(err)

		context.JSON(500, utils.ResponseError{Message: err.Error()})
		return
	}

	// why did we end up here? Golang really needs better static analysis.
	context.JSON(500, utils.ResponseError{Message: "Sync: Unknown code path. Shouldn't have landed here."})
	return
}

// Processes the outputs of the given workflow.
func (u IufController) doProcessOutputs(workflow *v1alpha1.Workflow, session *iuf.Session, requestBody iuf.SyncRequest, sessionName string) {
	u.logger.Infof("doProcessOutputs: About to process outputs for workflow: %s, resource version: %s, session: %s, activity: %s", workflow.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef)

	err := u.iufService.ProcessOutput(session, workflow)
	if err != nil {
		u.logger.Errorf("Sync: An error occurred processing the output for the workflow: %s, resource version: %s, session: %s, activity: %s, error: %v", workflow.Name, requestBody.Object.ObjectMeta.ResourceVersion, sessionName, session.ActivityRef, err)
		// do not return error, just continue because process output should not re-attempt stage.
	}
}

func (u IufController) restartCurrentStageFromSyncCall(context *gin.Context, session iuf.Session, requestBody iuf.SyncRequest, response iuf.SyncResponse) {
	u.logger.Infof("Sync: Restarting stage %s in session %s in activity %s", session.CurrentStage, session.Name, session.ActivityRef)

	err := u.iufService.RestartCurrentStage(&session, session.CurrentStage)
	if err != nil {
		u.logger.Errorf("Sync: Unable to restart current stage. Current stage: %s, resource version: %s, session: %s, activity: %s, error: %v", session.CurrentStage, requestBody.Object.ObjectMeta.ResourceVersion, session.Name, session.ActivityRef, err)
		// note: do NOT automatically retry -- we don't know whether CurrentStage has already been updated
		//  This is the downside of using a non-transactional storage such as CRDs.
		context.JSON(500, utils.ResponseError{Message: err.Error()})

		session.CurrentState = iuf.SessionStateDebug
		u.iufService.UpdateSessionAndActivity(session, "Unable to restart current stage")
		return
	}

	response = iuf.SyncResponse{
		ResyncAfterSeconds: RESYNC_TIME_IN_SECONDS,
	}

	context.JSON(200, response)
}

// WorkflowSync **experimental** Instead of a webhook on Session, we should have defined a webhook on Argo workflows instead
func (u IufController) WorkflowSync(context *gin.Context) {
	var requestBody iuf.WorkflowSyncRequest
	u.logger.Infof("WorkflowSync: received request with params %#v", context.Request.Form)

	if err := context.BindJSON(&requestBody); err != nil {
		u.logger.Errorf("WorkflowSync: An error occurred parsing request: %v", err)
		context.JSON(500, err.Error())
		return
	}

	bytes, _ := json.Marshal(requestBody)
	u.logger.Infof("WorkflowSync: Received the following workflow sync request: %s", string(bytes))
}
