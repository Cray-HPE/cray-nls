//
//  MIT License
//
//  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
//
//  Permission is hereby granted, free of charge, to any person obtaining a
//  copy of this software and associated documentation files (the "Software"),
//  to deal in the Software without restriction, including without limitation
//  the rights to use, copy, modify, merge, publish, distribute, sublicense,
//  and/or sell copies of the Software, and to permit persons to whom the
//  Software is furnished to do so, subject to the following conditions:
//
//  The above copyright notice and this permission notice shall be included
//  in all copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
//  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
//  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
//  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
//  OTHER DEALINGS IN THE SOFTWARE.
//
package argo_templates

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestRenderWorkerRebuildTemplate(t *testing.T) {
	t.Run("It should render a workflow template for a worker node", func(t *testing.T) {
		targetNcns := []string{"ncn-w006", "ncn-w005"}
		a, _ := GetWorkerRebuildWorkflow(targetNcns)
		t.Errorf("%s", string(a))
		assert.Equal(t, true, strings.Contains(string(a), "targetNcn"))
	})
	t.Run("It should fail when host is not a worker node", func(t *testing.T) {
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
		}
		for _, tt := range tests {
			t.Run(tt.hostnames[0], func(t *testing.T) {
				_, err := GetWorkerRebuildWorkflow(tt.hostnames)
				if (err != nil) != tt.wantErr {
					t.Errorf("got %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		}

	})
	t.Run("It should select nodes that is not being rebuilt", func(t *testing.T) {
		targetNcn := "ncn-w99999"
		workerRebuildWorkflow, _ := GetWorkerRebuildWorkflow([]string{targetNcn})
		workerRebuildWorkflowJson, _ := yaml.YAMLToJSON(workerRebuildWorkflow)
		var myWorkflow v1alpha1.Workflow
		json.Unmarshal(workerRebuildWorkflowJson, &myWorkflow)
		assert.Equal(t, targetNcn, myWorkflow.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values[0])
	})
}

func TestRenderStorageRebuildTemplate(t *testing.T) {
	t.Skip("skip: future feature")
	t.Run("It should render a workflow template for a storage node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should fail when host is not a storage node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

func TestRenderMasterRebuildTemplate(t *testing.T) {
	t.Skip("skip: future feature")
	t.Run("It should render a workflow template for a master node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})

	t.Run("It should fail when host is not a master node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

func TestRenderRebuildTemplate(t *testing.T) {
	t.Skip("skip: future feature")
	t.Run("It should render a workflow template for a node", func(t *testing.T) {
		// loop test 3 types: master/worker/storage
		assert.Fail(t, "NOT IMPLEMENTED")
	})

	t.Run("It should fail when parameters are invalid", func(t *testing.T) {
		// loop test: hostname, xname, image version
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}
