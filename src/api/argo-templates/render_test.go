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
package argo_templates

import (
	"embed"
	"encoding/json"
	"testing"

	models_nls "github.com/Cray-HPE/cray-nls/src/api/models/nls"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

const doDryRun bool = true

//go:embed _test_data_/*
var rebuildWorkflowFS embed.FS

func TestRenderWorkerRebuildTemplate(t *testing.T) {
	t.Run("It should render a workflow template for a group of worker nodes", func(t *testing.T) {
		req := models_nls.CreateRebuildWorkflowRequest{
			Hosts:  []string{"ncn-w006", "ncn-w005"},
			DryRun: doDryRun,
		}
		_, err := GetWorkerRebuildWorkflow(rebuildWorkflowFS, req, models_nls.RebuildHooks{})
		assert.Equal(t, true, err == nil)
	})
	t.Run("Render with valid/invalid hostnames", func(t *testing.T) {
		var tests = []struct {
			hostnames []string
			wantErr   bool
		}{
			{[]string{"ncn-m001"}, true},
			{[]string{"ncn-w001"}, false},
			{[]string{"ncn-s001"}, true},
			{[]string{"ncn-m011"}, true},
			{[]string{"ncn-x001"}, true},
			{[]string{"sccn-m001"}, true},
			{[]string{"ncn-x001"}, true},
			{[]string{"ncn-m001asdf"}, true},
			{[]string{"ncn-w001", "ncn-m001asdf"}, true},
		}
		for _, tt := range tests {
			t.Run(tt.hostnames[0], func(t *testing.T) {
				req := models_nls.CreateRebuildWorkflowRequest{
					Hosts:  tt.hostnames,
					DryRun: doDryRun,
				}
				_, err := GetWorkerRebuildWorkflow(rebuildWorkflowFS, req, models_nls.RebuildHooks{})
				if (err != nil) != tt.wantErr {
					t.Errorf("got %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		}

	})
	t.Run("It should select nodes that is not being rebuilt", func(t *testing.T) {
		req := models_nls.CreateRebuildWorkflowRequest{
			Hosts:  []string{"ncn-w99999"},
			DryRun: doDryRun,
		}
		workerRebuildWorkflow, _ := GetWorkerRebuildWorkflow(rebuildWorkflowFS, req, models_nls.RebuildHooks{})
		workerRebuildWorkflowJson, _ := yaml.YAMLToJSON(workerRebuildWorkflow)
		var myWorkflow v1alpha1.Workflow
		json.Unmarshal(workerRebuildWorkflowJson, &myWorkflow)
		assert.Equal(t, "ncn-w99999", myWorkflow.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values[0])
	})
}

// ---- Storage Testing ----//
func TestRenderStorageRebuildTemplate(t *testing.T) {
	t.Run("It should render a workflow template for a group of storage nodes", func(t *testing.T) {
		req := models_nls.CreateRebuildWorkflowRequest{
			Hosts:            []string{"ncn-s006", "ncn-s005"},
			DryRun:           doDryRun,
			ZapOsds:          false,
			WorkflowType:     "rebuild",
			ImageId:          "",
			DesiredCfsConfig: "",
		}
		_, err := GetStorageRebuildWorkflow(rebuildWorkflowFS, req)
		assert.Equal(t, true, err == nil)
	})
	t.Run("Render with valid/invalid hostnames", func(t *testing.T) {
		var tests = []struct {
			hostnames []string
			wantErr   bool
		}{
			{[]string{"ncn-m001"}, true},
			{[]string{"ncn-w001"}, true},
			{[]string{"ncn-s001"}, false},
			{[]string{"ncn-m011"}, true},
			{[]string{"ncn-x001"}, true},
			{[]string{"sccn-m001"}, true},
			{[]string{"ncn-x001"}, true},
			{[]string{"ncn-m001asdf"}, true},
			{[]string{"ncn-w001", "ncn-m001asdf"}, true},
		}
		for _, tt := range tests {
			t.Run(tt.hostnames[0], func(t *testing.T) {
				req := models_nls.CreateRebuildWorkflowRequest{
					Hosts:  tt.hostnames,
					DryRun: doDryRun,
				}
				_, err := GetStorageRebuildWorkflow(rebuildWorkflowFS, req)
				if (err != nil) != tt.wantErr {
					t.Errorf("got %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		}

	})
}
