/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022-2025 Hewlett Packard Enterprise Development LP
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
package models

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type CreateRebuildWorkflowRequest struct {
	Hosts                []string          `json:"hosts"`
	DryRun               bool              `json:"dryRun"`
	ZapOsds              bool              `json:"zapOsds,omitempty"`      // this is necessary for storage rebuilds when unable to wipe the node prior to rebuild
	WorkflowType         string            `json:"workflowType,omitempty"` // used to determine storage rebuild vs upgrade
	ImageId              string            `json:"imageId,omitempty"`
	DesiredCfsConfig     string            `json:"desiredCfsConfig,omitempty"`
	Labels               map[string]string `json:"labels,omitempty"`
	BootTimeoutInSeconds int               `json:"bootTimeoutInSeconds,omitempty"`
}

type CreateRebuildWorkflowResponse struct {
	Name       string   `json:"name"`
	TargetNcns []string `json:"targetNcns"`
}

type RebuildHooks struct {
	BeforeAll  []unstructured.Unstructured
	BeforeEach []unstructured.Unstructured
	AfterEach  []unstructured.Unstructured
	AfterAll   []unstructured.Unstructured
}
