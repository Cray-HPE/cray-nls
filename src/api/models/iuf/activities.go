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
	// Name of activity
	Name string `json:"name" validate:"required"`
	// location of media
	MediaDir string `json:"media_dir" validate:"required"`
	// The inline contents of the site_parameters.yaml file.
	SiteParameters string `json:"site_parameters" validate:"required"`
	// Each item is the xname of a node
	LimitNodes []string `json:"limit_nodes" validate:"optional"`
	// Each item is a path of the bootprep files
	BootprepConfigManaged []string `json:"bootprep_config_managed" validate:"required"`
	// Each item is a path of the bootprep files
	BootprepConfigManagement []string `json:"bootprep_config_management" validate:"required"`
	// Operation outputs from argo
	OperationOutputs map[string]interface{} `json:"operation_outputs" validate:"required"`
	// Comment provided by admin
	CurrentComment string `json:"current_comment" validate:"optional"`
	// List of products included in an activity
	Products []Product `json:"products" validate:"required"`
	// History of states
	ActivityStates []ActivityState `json:"activity_states" validate:"required"`
} // @name Activity

type ActivityState struct {
	State       string `json:"state" validate:"required"`
	SessionName string `json:"session_name" validate:"required"`
	StartTime   string `json:"start_time" validate:"required"`
	Comment     string `json:"comment" validate:"optional"`
} // @name Activity.State

type CreateActivityRequest struct {
	// TODO
} // @name Activity.CreateActivityRequest

type PatchActivityRequest struct {
	// TODO
} // @name Activity.PatchActivityRequest
