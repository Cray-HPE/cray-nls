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

type IufActivitiesSyncRequest struct {
	Parent IufActivity `json:"parent"`
}

type IufActivitiesSyncResponse struct {
	Status             IufActivitiesStatus `json:"status,omitempty"`
	ResyncAfterSeconds int                 `json:"resyncAfterSeconds,omitempty"`
}

// IufActivity
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=iufactivities,scope=Namespaced
// +kubebuilder:storageversion
type IufActivity struct {
	metav1.TypeMeta   `json:",inline" swaggerignore:"true"`
	metav1.ObjectMeta `json:"metadata" swaggerignore:"true"`
	Spec              IufActivitiesSpec   `json:"spec"`
	Status            IufActivitiesStatus `json:"status,omitempty"`
} // @name Activity

type IufActivityCurrentState struct {
	Type    string `json:"type" validate:"required"`
	Comment string `json:"comment"  validate:"optional"`
} // @name Activity.CurrentState

// An IUF session represents the intent of an Admin to initiate an install-upgrade workflow. It contains both input data, as well as any intermediary data that is needed to generate the final Argo workflow.
type IufActivitiesSpec struct {
	SharedInput    `json:",inline"`
	IsBlocked      bool   `json:"is_blocked"`
	IsCompleted    bool   `json:"is_completed"`
	CurrentComment string `json:"current_comment"`
} // @name Activity.Spec

type IufActivitiesStatus struct {
	SharedInput `json:",inline"`
	Sessions    []IufSession `json:"sessions" validate:"optional"`
	Products    []IufProduct `json:"products" validate:"optional"`
} // @name Activity.Status
