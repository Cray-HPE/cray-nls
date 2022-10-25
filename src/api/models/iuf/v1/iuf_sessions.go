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
// +groupName=iuf.hpe.com
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IufSyncRequest struct {
	Parent IufSession `json:"parent"`
}

type IufSyncResponse struct {
	Status             IufSessionStatus `json:"status,omitempty"`
	ResyncAfterSeconds int              `json:"resyncAfterSeconds,omitempty"`
}

// IufSession
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=iufsessions,scope=Namespaced
// +kubebuilder:storageversion
type IufSession struct {
	metav1.TypeMeta   `json:",inline" swaggerignore:"true"`
	metav1.ObjectMeta `json:"metadata" swaggerignore:"true"`
	Spec              IufSessionSpec   `json:"spec"`
	Status            IufSessionStatus `json:"status,omitempty" swaggerignore:"true"`
}

type IufSessionCurrentState struct {
	Type    IufSessionStageState `json:"type" validate:"required"`
	Comment string               `json:"comment"  validate:"optional"`
} // @name IufSession.CurrentState

// An IUF session represents the intent of an Admin to initiate an install-upgrade workflow. It contains both input data, as well as any intermediary data that is needed to generate the final Argo workflow.
type IufSessionSpec struct {
	SharedInput `json:",inline"`
	// The stages that need to be executed.
	// This is either explicitly specified by the Admin, or it is computed from the workflow type.
	// An Stage is a group of Operations. Stages represent the overall workflow at a high-level, and executing a stage means executing a bunch of Operations in a predefined manner.  An Admin can specify the stages that must be executed for an install-upgrade workflow. And Product Developers can extend each stage with custom hook scripts that they would like to run before and after the stage's execution.  The high-level stages allow their configuration would revealing too many details to the consumers of IUF.
	// if not specified, we apply all stages
	Stages []string `json:"stages"`

	Products []IufProduct `json:"products" validate:"required"`
} // @name IufSession.Spec

// +kubebuilder:validation:Enum=not_started;in_progress;error;complete
type IufSessionStageState string

// Node types
const (
	IufSessionStageNotStarted IufSessionStageState = "not_started"
	IufSessionStageInProgress IufSessionStageState = "in_progres"
	IufSessionStageError      IufSessionStageState = "error"
	IufSessionStageComplete   IufSessionStageState = "complete"
)

type IufSessionStage struct {
	Name          string               `json:"products" validate:"required"`
	State         IufSessionStageState `json:"state" validate:"required"`
	WorkflowId    string               `json:"workflou_id" validate:"required"`
	WorkflowOuput map[string]string    `json:"workflou_output" validate:"optional"`
} // @name IufSession.Stage

type IufSessionStatus struct {
	CurrentState IufSessionCurrentState `json:"current_state"`
	Stages       []IufSessionStage      `json:"products" validate:"optional"`
}
