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
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	core_v1 "k8s.io/api/core/v1"
)

// IufSession
type Session struct {
	InputParameters InputParameters   `json:"input_parameters"`
	SiteParameters  SiteParameters    `json:"site_parameters"`
	CurrentState    SessionState      `json:"current_state" enums:"paused,in_progress,debug,completed,aborted"`
	CurrentStage    string            `json:"stage"`
	Workflows       []SessionWorkflow `json:"workflows"`
	Products        []Product         `json:"products" validate:"required"`
	Name            string            `json:"name"`
	ActivityRef     string            `json:"activityRef" swaggerignore:"true"`
} //	@name	Session

type SessionState string

const (
	SessionStateInProgress    SessionState = "in_progress"
	SessionStateTransitioning SessionState = "transitioning"
	SessionStatePaused        SessionState = "paused"
	SessionStateDebug         SessionState = "debug"
	SessionStateCompleted     SessionState = "completed"
	SessionStateAborted       SessionState = "aborted"
)

type SessionWorkflow struct {
	Id  string `json:"id"`  // id of argo workflow
	Url string `json:"url"` // url to the argo workflow
} //	@name	Session.Workflow

type SyncRequest struct {
	Object core_v1.ConfigMap `json:"object"`
}
type WorkflowSyncRequest struct {
	Object v1alpha1.Workflow `json:"object"`
}

type SyncResponse struct {
	ResyncAfterSeconds int `json:"resyncAfterSeconds,omitempty"`
}

type ManifestHookScript struct {
	ScriptPath       string `json:"script_path"`
	ExecutionContext string `json:"execution_context"`
}

type ManifestStageHooks struct {
	PreHook  ManifestHookScript `json:"pre"`
	PostHook ManifestHookScript `json:"post"`
}
