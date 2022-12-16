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

// Activity
type Activity struct {
	Name             string                 `json:"name"`                                                                                      // Name of activity
	InputParameters  InputParameters        `json:"input_parameters" binding:"required"`                                                       // Input parameters by admin
	OperationOutputs map[string]interface{} `json:"operation_outputs" binding:"required"`                                                      // Operation outputs from argo
	Products         []Product              `json:"products" binding:"required"`                                                               // List of products included in an activity
	ActivityState    ActivityState          `json:"activity_state" binding:"required" enums:"paused,in_progress,debug,blocked,wait_for_admin"` // State of activity
} // @name Activity

type CreateActivityRequest struct {
	Name          string        `json:"name" binding:"required"` // Name of activity
	ActivityState ActivityState `json:"activity_state" swaggerignore:"true"`
} // @name Activity.CreateActivityRequest

type PatchActivityRequest struct {
	InputParameters InputParameters `json:"input_parameters" binding:"required"`
} // @name Activity.PatchActivityRequest

type ActivityState string

const (
	ActivityStateInProgress   ActivityState = "in_progress"
	ActivityStatePaused       ActivityState = "paused"
	ActivityStateDebug        ActivityState = "debug"
	ActivityStateBlocked      ActivityState = "blocked"
	ActivityStateWaitForAdmin ActivityState = "wait_for_admin"
)
