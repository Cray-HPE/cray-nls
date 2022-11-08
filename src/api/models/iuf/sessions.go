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
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// IufSession
type Session struct {
	InputParameters InputParameters   `json:"input_parameters"`
	CurrentState    SessionState      `json:"current_state" enums:"paused,in_progress,debug,completed"`
	CurrentStage    string            `json:"stage"`
	Workflows       []SessionWorkflow `json:"workflows"`
	Products        []Product         `json:"products" validate:"required"`
	Name            string            `json:"name"`
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

type SyncRequest struct {
	Object core_v1.ConfigMap `json:"object"`
}

type SyncResponse struct {
	ResyncAfterSeconds int `json:"resyncAfterSeconds,omitempty"`
}

type StageInputs struct {
	ProductManifest ProductManifest `json:"product_manifest"`
	InputParams     InputParams     `json:"input_params"`
	StageParams     StageParams     `json:"stage_params"`
}

type ProductManifest struct {
	Products       map[string]unstructured.Unstructured `json:"products"`
	CurrentProduct map[string]unstructured.Unstructured `json:"current_product"`
}

type InputParams struct {
	Producs        []string `json:"products"`
	MediaDir       string   `json:"media_dir"`
	SiteParameters struct {
		Global         map[string]unstructured.Unstructured `json:"global"`
		Products       map[string]unstructured.Unstructured `json:"products"`
		CurrentProduct map[string]unstructured.Unstructured `json:"current_product"`
	} `json:"site_parameters"`
	BootstrapConfigManaged struct {
		Contents string `json:"contents"`
	} `json:"bootprep_config_managed"`
	BootstrapConfigManagement struct {
		Contents string `json:"contents"`
	} `json:"bootprep_config_management"`
	LimitNodes []string `json:"limit_nodes"`
}

type StageParams map[string]struct {
	Global         map[string]unstructured.Unstructured `json:"global"`
	Products       map[string]unstructured.Unstructured `json:"products"`
	CurrentProduct map[string]unstructured.Unstructured `json:"current_product"`
}
