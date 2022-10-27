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

// IufSession
type Session struct {
	InputParameters InputParameters   `json:"input_parameters"`
	CurrentState    SessionState      `json:"current_state" enums:"paused,in_progress,debug,completed"`
	CurrentStage    string            `json:"stage"`
	Workflows       []SessionWorkflow `json:"workflows"`
	Products        []Product         `json:"products" validate:"required"`
} // @name Session

type SessionState string

const (
	SessionStateInProgress SessionState = "in_progress"
	SessionStatePaused     SessionState = "paused"
	SessionStateDebug      SessionState = "debug"
	SessionStateCompleted  SessionState = "completed"
)

type SessionWorkflow struct {
	Id  string `json:"id"`  // id of argo workflow
	Url string `json:"url"` // url to the argo workflow
} // @name Session.Workflow
