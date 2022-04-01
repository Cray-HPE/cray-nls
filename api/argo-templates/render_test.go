package argo_templates

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestRenderRebuildTemplate(t *testing.T) {
	t.Run("It should render a workflow template for a node", func(t *testing.T) {
		// loop test 3 types: master/worker/storage
		assert.Fail(t, "NOT IMPLEMENTED")
	})

	t.Run("It should fail when parameters are invalid", func(t *testing.T) {
		// loop test: hostname, xname, image version
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

func TestRenderMasterRebuildTemplate(t *testing.T) {
	t.Run("It should render a workflow template for a master node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})

	t.Run("It should fail when host is not a master node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

func TestRenderWorkerRebuildTemplate(t *testing.T) {
	t.Run("It should render a workflow template for a worker node", func(t *testing.T) {
		targetNcn := "ncn-w99999"
		a, _ := GetWrokerRebuildWorkflow(targetNcn, "")
		assert.Equal(t, true, strings.Contains(string(a), targetNcn))
	})
	t.Run("It should fail when host is not a worker node", func(t *testing.T) {
		var tests = []struct {
			hostname string
			wantErr  bool
		}{
			{"ncn-m001", true},
			{"ncn-w001", false},
			{"ncn-s001", true},
			{"ncn-m011", true},
			{"ncn-x001", true},
			{"sccn-m001", true},
			{"ncn-x001", true},
			{"ncn-m001asdf", true},
		}
		for _, tt := range tests {
			t.Run(tt.hostname, func(t *testing.T) {
				_, err := GetWrokerRebuildWorkflow(tt.hostname, "")
				if (err != nil) != tt.wantErr {
					t.Errorf("got %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		}

	})
	t.Run("It should select nodes that is not being rebuilt", func(t *testing.T) {
		targetNcn := "ncn-w99999"
		workerRebuildWorkflow, _ := GetWrokerRebuildWorkflow(targetNcn, "")
		workerRebuildWorkflowJson, _ := yaml.YAMLToJSON(workerRebuildWorkflow)
		var myWorkflow v1alpha1.Workflow
		json.Unmarshal(workerRebuildWorkflowJson, &myWorkflow)
		assert.Equal(t, targetNcn, myWorkflow.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values[0])
	})
}

func TestRenderStorageRebuildTemplate(t *testing.T) {
	t.Run("It should render a workflow template for a storage node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should fail when host is not a storage node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}
