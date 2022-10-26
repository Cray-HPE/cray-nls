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
	Name             string                 `json:"name"`                                  // Name of activity
	InputParameters  InputParameters        `json:"input_parameters" validate:"required"`  // Input parameters by admin
	OperationOutputs map[string]interface{} `json:"operation_outputs" validate:"required"` // Operation outputs from argo
	Products         []Product              `json:"products" validate:"required"`          // List of products included in an activity
	ActivityStates   []ActivityState        `json:"activity_states" validate:"required"`   // History of states
} // @name Activity

type ActivityState struct {
	State       string `json:"state" validate:"required"`
	SessionName string `json:"session_name" validate:"required"`
	StartTime   string `json:"start_time" validate:"required"`
	Comment     string `json:"comment" validate:"optional"`
} // @name Activity.State

type InputParameters struct {
	MediaDir                 string   `json:"media_dir"`                  // Location of media
	SiteParameters           string   `json:"site_parameters"`            // The inline contents of the site_parameters.yaml file.
	LimitNodes               []string `json:"limit_nodes"`                // Each item is the xname of a node
	BootprepConfigManaged    []string `json:"bootprep_config_managed"`    // Each item is a path of the bootprep files
	BootprepConfigManagement []string `json:"bootprep_config_management"` // Each item is a path of the bootprep files
	Stages                   []string `json:"stages"`                     // Execution of the specified stages
	Force                    bool     `json:"force"`                      // Force re-execution of stage operations
} // @name Activity.InputParameters

type CreateActivityRequest struct {
	Name            string          `json:"name" validate:"required"`             // Name of activity
	InputParameters InputParameters `json:"input_parameters" validate:"required"` // Input parameters by admin
} // @name Activity.CreateActivityRequest

type PatchActivityRequest struct {
	InputParameters InputParameters `json:"input_parameters" validate:"required"`
} // @name Activity.PatchActivityRequest
