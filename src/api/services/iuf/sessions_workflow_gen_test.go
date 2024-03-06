/*
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
package services_iuf

import (
	"encoding/json"
	"fmt"
	mocks "github.com/Cray-HPE/cray-nls/src/api/mocks/services"
	"github.com/golang/mock/gomock"
	"strconv"
	"strings"
	"testing"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	workflowmocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow/mocks"
	workflowtemplatemocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate/mocks"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fake "k8s.io/client-go/kubernetes/fake"
)

func TestWorkflowGen(t *testing.T) {
	activityName, _, iufSvc := setup(t)

	t.Run("generated workflow must have NodeSelector set to ncn-m002 when rebuilding ncn-m001", func(t *testing.T) {
		session := iuf.Session{
			Products: []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			Workflows: []iuf.SessionWorkflow{
				{
					Id:  "1",
					Url: "1",
				},
			},
			CurrentStage: "management-nodes-rollout",
			CurrentState: iuf.SessionStateInProgress,
			InputParameters: iuf.InputParameters{
				Stages: []string{"management-nodes-rollout"},
				LimitManagementNodes: []string{"ncn-m001"},
			},
			ActivityRef: activityName,
		}

		workflow, err, _ := iufSvc.workflowGen(&session)
		assert.NoError(t, err)
		assert.Equal(t, "ncn-m002", workflow.Spec.NodeSelector["kubernetes.io/hostname"])
	})

	t.Run("generated workflow must have NodeSelector set to ncn-m001 when rebuilding ncn-m002 or ncn-m003", func(t *testing.T) {
		session := iuf.Session{
			Products: []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			Workflows: []iuf.SessionWorkflow{
				{
					Id:  "1",
					Url: "1",
				},
			},
			CurrentStage: "management-nodes-rollout",
			CurrentState: iuf.SessionStateInProgress,
			InputParameters: iuf.InputParameters{
				Stages: []string{"management-nodes-rollout"},
				LimitManagementNodes: []string{"ncn-m002"},
			},
			ActivityRef: activityName,
		}

		workflow, err, _ := iufSvc.workflowGen(&session)
		assert.NoError(t, err)
		assert.Equal(t, "ncn-m001", workflow.Spec.NodeSelector["kubernetes.io/hostname"])
	})

	t.Skip("TODO generated workflow must not have NodeSelector set to ncn-m001 when NoHooks=true") /*, func(t *testing.T) {
		session := iuf.Session{
			Products: []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			Workflows: []iuf.SessionWorkflow{
				{
					Id:  "1",
					Url: "1",
				},
			},
			CurrentStage: "process-media",
			CurrentState: iuf.SessionStateInProgress,
			InputParameters: iuf.InputParameters{
				Stages: []string{"process-media", "management-m001-rollout"},
			},
			ActivityRef: activityName,
		}

		workflow, err := iufSvc.workflowGen(session)
		assert.NoError(t, err)
		assert.Equal(t, "ncn-m002", workflow.Spec.NodeSelector["kubernetes.io/hostname"])
	}*/
}

func TestGetDagTasks(t *testing.T) {
	activityName, mockAuthToken, iufSvc := setup(t)
	globalParamsPerProduct := map[string]string{
		"product":   "product",
		"cos-1-2-3": "cos-1-2-3",
		"sdu-3-4-5": "sdu-3-4-5",
	}

	t.Run("It should get a dag task for per-product stage", func(t *testing.T) {
		session := iuf.Session{
			Products:    []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			ActivityRef: activityName,
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "product",
			Operations: []iuf.Operations{
				{Name: "this-is-an-operation-1"},
				{Name: "this-is-an-operation-2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}
		workflow := v1alpha1.Workflow{}

		dagTasks, _, err := iufSvc.getDAGTasks(&session, stageInfo, stages, globalParamsPerProduct, "global_params", "auth_token", &workflow)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 4, len(dagTasks))
		assert.Equal(t, 2, len(dagTasks[0].Arguments.Parameters))
		assert.Equal(t, dagTasks[0].Arguments.Parameters[0].Name, "auth_token")
		assert.Equal(t, dagTasks[0].Arguments.Parameters[0].Value, v1alpha1.AnyStringPtr(mockAuthToken))
	})
	t.Run("It should get a dag task for global stage", func(t *testing.T) {
		session := iuf.Session{
			Products:    []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			ActivityRef: activityName,
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "global",
			Operations: []iuf.Operations{
				{Name: "this-is-an-operation-1"},
				{Name: "this-is-an-operation-2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}
		workflow := v1alpha1.Workflow{}

		dagTasks, _, err := iufSvc.getDAGTasks(&session, stageInfo, stages, globalParamsPerProduct, "global_params", "auth_token", &workflow)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 2, len(dagTasks))
	})
	t.Run("It should not get a dag task for per-product operations not defined in Argo", func(t *testing.T) {
		session := iuf.Session{
			Products:    []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			ActivityRef: activityName,
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "product",
			Operations: []iuf.Operations{
				{Name: "this-is-an-operation-1"},
				{Name: "this-is-an-operation-NOT_MOCKED"},
				{Name: "this-is-an-operation-2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}
		workflow := v1alpha1.Workflow{}

		dagTasks, _, err := iufSvc.getDAGTasks(&session, stageInfo, stages, globalParamsPerProduct, "global_params", "auth_token", &workflow)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 4, len(dagTasks))
		assert.Equal(t, dagTasks[0].Arguments.Parameters[0].Name, "auth_token")
		assert.Equal(t, dagTasks[0].Arguments.Parameters[0].Value, v1alpha1.AnyStringPtr(mockAuthToken))
	})
	t.Run("It should not get a dag task for global operations not defined in Argo", func(t *testing.T) {
		session := iuf.Session{
			Products:    []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			ActivityRef: activityName,
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "global",
			Operations: []iuf.Operations{
				{Name: "this-is-an-operation-NOT-MOCKED-1"},
				{Name: "this-is-an-operation-1"},
				{Name: "this-is-an-operation-NOT-MOCKED-2"},
				{Name: "this-is-an-operation-2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}
		workflow := v1alpha1.Workflow{}

		dagTasks, _, err := iufSvc.getDAGTasks(&session, stageInfo, stages, globalParamsPerProduct, "global_params", "auth_token", &workflow)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 2, len(dagTasks))
		assert.Equal(t, 2, len(dagTasks[0].Arguments.Parameters))
		assert.Equal(t, 0, len(dagTasks[0].Dependencies))
		assert.Equal(t, "this-is-an-operation-1", dagTasks[0].Name)
		assert.Equal(t, "this-is-an-operation-2", dagTasks[1].Name)
		assert.Equal(t, 0, len(dagTasks[0].Dependencies))
		assert.Equal(t, 1, len(dagTasks[1].Dependencies))
		assert.Equal(t, "this-is-an-operation-1", dagTasks[1].Dependencies[0])
		assert.Equal(t, dagTasks[0].Arguments.Parameters[0].Name, "auth_token")
		assert.Equal(t, dagTasks[0].Arguments.Parameters[0].Value, v1alpha1.AnyStringPtr(mockAuthToken))
	})

	t.Run("It should get DAG tasks for existing hook templates defined for product operations", func(t *testing.T) {
		session := iuf.Session{
			Products: []iuf.Product{
				{
					Name:             "cos",
					Version:          "1.2.3",
					OriginalLocation: cosOriginalLocation,
					Manifest:         cosManifest,
				},
				{
					Name:             "sdu",
					Version:          "3.4.5",
					OriginalLocation: sduOriginalLocation,
					Manifest:         sduManifest,
				},
			},
			ActivityRef: activityName,
		}
		stageInfo := iuf.Stage{
			Name: "pre-install-check",
			Type: "product",
			Operations: []iuf.Operations{
				{Name: "this-is-an-operation-1"},
				{Name: "this-is-an-operation-NOT_MOCKED"},
				{Name: "this-is-an-operation-2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
			Hooks: map[string]string{
				"master_host": "master-host-hook-script",
				"worker_host": "worker-host-hook-script",
			},
		}
		workflow := v1alpha1.Workflow{}

		dagTasks, _, err := iufSvc.getDAGTasks(&session, stageInfo, stages, globalParamsPerProduct, "global_params", "auth_token", &workflow)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 7, len(dagTasks))

		// the first task will be a hook script (see cosManifest and sduManifest constants)
		assert.True(t, strings.HasPrefix(dagTasks[0].Name, "cos-1-2-3-pre-hook-pre-install-check"))

		var found int
		for _, dagTask := range dagTasks {
			// first determine what product this task is for
			var product iuf.Product
			if strings.Contains(dagTask.Name, "cos") {
				product = session.Products[0]
			} else {
				product = session.Products[1]
			}

			assert.Equal(t, v1alpha1.AnyStringPtr(mockAuthToken), dagTask.Arguments.GetParameterByName("auth_token").Value)

			productKey := iufSvc.getProductVersionKey(product)
			assert.Equal(t, v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.%s}}", productKey)), dagTask.Arguments.GetParameterByName("global_params").Value)

			if strings.HasPrefix(dagTask.Name, "cos-1-2-3-pre-hook-pre-install-check") {
				t.Run("cos pre hook script operation exist and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "master-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 0, len(dagTask.Dependencies))
					found++
				})
			} else if strings.HasPrefix(dagTask.Name, "cos-1-2-3-this-is-an-operation-1") {
				t.Run("cos operation 1 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this-is-an-operation-1", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.True(t, strings.HasPrefix(dagTask.Dependencies[0], "cos-1-2-3-pre-hook-pre-install-check"))
					found++
				})
			} else if strings.HasPrefix(dagTask.Name, "cos-1-2-3-this-is-an-operation-2") {
				t.Run("cos operation 2 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this-is-an-operation-2", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.True(t, strings.HasPrefix(dagTask.Dependencies[0], "cos-1-2-3-this-is-an-operation-1"))
					found++
				})
			} else if strings.HasPrefix(dagTask.Name, "sdu-3-4-5-this-is-an-operation-1") {
				t.Run("sdu operation 1 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this-is-an-operation-1", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 0, len(dagTask.Dependencies))
					found++
				})
			} else if strings.HasPrefix(dagTask.Name, "sdu-3-4-5-this-is-an-operation-2") {
				t.Run("sdu operation 2 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this-is-an-operation-2", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.True(t, strings.HasPrefix(dagTask.Dependencies[0], "sdu-3-4-5-this-is-an-operation-1"))
					found++
				})
			} else if strings.HasPrefix(dagTask.Name, "cos-1-2-3-post-hook-pre-install-check") {
				t.Run("cos post hook script operation exist and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "worker-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.True(t, strings.HasPrefix(dagTask.Dependencies[0], "cos-1-2-3-this-is-an-operation-2"))
					found++
				})
			} else if strings.HasPrefix(dagTask.Name, "sdu-3-4-5-post-hook-pre-install-check") {
				t.Run("sdu post hook script operation exist and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "worker-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.True(t, strings.HasPrefix(dagTask.Dependencies[0], "sdu-3-4-5-this-is-an-operation-2"))
					found++
				})
			}
		}

		assert.Equal(t, 7, found)
	})

	t.Run("It should get DAG tasks for existing hook templates defined for global operations", func(t *testing.T) {
		session := iuf.Session{
			Products: []iuf.Product{
				{
					Name:             "cos",
					Version:          "1.2.3",
					OriginalLocation: cosOriginalLocation,
					Manifest:         cosManifest,
				},
				{
					Name:             "sdu",
					Version:          "3.4.5",
					OriginalLocation: sduOriginalLocation,
					Manifest:         sduManifest_alt,
				},
			},
			ActivityRef: activityName,
		}
		stageInfo := iuf.Stage{
			Name: "pre-install-check",
			Type: "global",
			Operations: []iuf.Operations{
				{Name: "this-is-an-operation-NOT-MOCKED-1"},
				{Name: "this-is-an-operation-1"},
				{Name: "this-is-an-operation-NOT-MOCKED-2"},
				{Name: "this-is-an-operation-2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
			Hooks: map[string]string{
				"master_host": "master-host-hook-script",
				"worker_host": "worker-host-hook-script",
			},
		}
		workflow := v1alpha1.Workflow{}

		dagTasks, _, err := iufSvc.getDAGTasks(&session, stageInfo, stages, globalParamsPerProduct, "global_params", "auth_token", &workflow)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 6, len(dagTasks))

		// the first two tasks will be hook scripts (see cosManifest and sduManifest_alt constants)
		assert.Equal(t, true, strings.Contains(dagTasks[0].Name, "pre-hook-pre-install-check"))
		assert.Equal(t, true, strings.Contains(dagTasks[1].Name, "pre-hook-pre-install-check"))

		var found int
		for _, dagTask := range dagTasks {
			// first determine what product this task is for
			var product iuf.Product
			if strings.Contains(dagTask.Name, "cos") {
				product = session.Products[0]
				productKey := iufSvc.getProductVersionKey(product)
				assert.Equal(t, v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.%s}}", productKey)), dagTask.Arguments.GetParameterByName("global_params").Value)
			} else if strings.Contains(dagTask.Name, "sdu") {
				product = session.Products[1]
				productKey := iufSvc.getProductVersionKey(product)
				assert.Equal(t, v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.%s}}", productKey)), dagTask.Arguments.GetParameterByName("global_params").Value)
			} else {
				assert.Equal(t, v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.global_params}}")), dagTask.Arguments.GetParameterByName("global_params").Value)
			}

			assert.Equal(t, v1alpha1.AnyStringPtr(mockAuthToken), dagTask.Arguments.GetParameterByName("auth_token").Value)

			if strings.HasPrefix(dagTask.Name, "cos-1-2-3-pre-hook-pre-install-check") {
				t.Run("cos pre hook script operation exist and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "master-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 0, len(dagTask.Dependencies))
					found++
				})
			} else if strings.HasPrefix(dagTask.Name, "sdu-3-4-5-pre-hook-pre-install-check") {
				t.Run("sdu pre hook script operation exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "worker-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 0, len(dagTask.Dependencies))
					found++
				})
			} else if strings.HasPrefix(dagTask.Name, "this-is-an-operation-1") {
				t.Run("operation 1 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this-is-an-operation-1", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 2, len(dagTask.Dependencies))
					assert.Equal(t, true, strings.Contains(dagTask.Dependencies[0], "pre-hook-pre-install-check"))
					assert.Equal(t, true, strings.Contains(dagTask.Dependencies[1], "pre-hook-pre-install-check"))
					assert.NotEqual(t, dagTask.Dependencies[0], dagTask.Dependencies[1])
					found++
				})
			} else if strings.HasPrefix(dagTask.Name, "this-is-an-operation-2") {
				t.Run("operation 2 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this-is-an-operation-2", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.Equal(t, "this-is-an-operation-1", dagTask.Dependencies[0])
					found++
				})
			} else if strings.HasPrefix(dagTask.Name, "cos-1-2-3-post-hook-pre-install-check") {
				t.Run("cos post hook script operation exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "worker-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.Equal(t, "this-is-an-operation-2", dagTask.Dependencies[0])
					found++
				})
			} else if strings.HasPrefix(dagTask.Name, "sdu-3-4-5-post-hook-pre-install-check") {
				t.Run("sdu post hook script operation exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "worker-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.Equal(t, "this-is-an-operation-2", dagTask.Dependencies[0])
					found++
				})
			}
		}

		assert.Equal(t, 6, found)
	})
	t.Run("It should get correct template for Management-nodes-rollout if worker hostname is provided", func(t *testing.T) {
		session := iuf.Session{
			Products:        []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			ActivityRef:     activityName,
			InputParameters: iuf.InputParameters{LimitManagementNodes: []string{"ncn-w002"}},
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "global",
			Operations: []iuf.Operations{
				{Name: "this-is-an-operation-1"},
				{Name: "management-nodes-rollout"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}
		workflow := v1alpha1.Workflow{}

		dagTasks, _, err := iufSvc.getDAGTasks(&session, stageInfo, stages, globalParamsPerProduct, "global_params", "auth_token", &workflow)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 2, len(dagTasks))
		assert.Equal(t, "this-is-an-operation-1", dagTasks[0].TemplateRef.Name)
		assert.Equal(t, "management-worker-nodes-rollout", dagTasks[1].TemplateRef.Name)
	})
	t.Run("It should get correct template for Management-nodes-rollout if Storage HSM role_subrole is provided", func(t *testing.T) {
		session := iuf.Session{
			Products:        []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			ActivityRef:     activityName,
			InputParameters: iuf.InputParameters{LimitManagementNodes: []string{"Management_Storage"}},
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "global",
			Operations: []iuf.Operations{
				{Name: "this-is-an-operation-1"},
				{Name: "management-nodes-rollout"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}
		workflow := v1alpha1.Workflow{}

		dagTasks, _, err := iufSvc.getDAGTasks(&session, stageInfo, stages, globalParamsPerProduct, "global_params", "auth_token", &workflow)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 2, len(dagTasks))
		assert.Equal(t, "management-storage-nodes-rollout", dagTasks[1].TemplateRef.Name)
	})
	t.Run("It should get correct template for Management-nodes-rollout if worker hostname is provided", func(t *testing.T) {
		session := iuf.Session{
			Products:        []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			ActivityRef:     activityName,
			InputParameters: iuf.InputParameters{LimitManagementNodes: []string{"ncn-w002"}},
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "global",
			Operations: []iuf.Operations{
				{Name: "this-is-an-operation-1"},
				{Name: "management-nodes-rollout"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}
		workflow := v1alpha1.Workflow{}

		dagTasks, _, err := iufSvc.getDAGTasks(&session, stageInfo, stages, globalParamsPerProduct, "global_params", "auth_token", &workflow)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 2, len(dagTasks))
		assert.Equal(t, "this-is-an-operation-1", dagTasks[0].TemplateRef.Name)
		assert.Equal(t, "management-worker-nodes-rollout", dagTasks[1].TemplateRef.Name)
	})
	t.Run("It should not get an error if --limit-management-rollout has an invalid input. Should still return echoTemplate.", func(t *testing.T) {
		session := iuf.Session{
			Products:        []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			ActivityRef:     activityName,
			InputParameters: iuf.InputParameters{LimitManagementNodes: []string{"ncn-bad"}},
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "global",
			Operations: []iuf.Operations{
				{Name: "this-is-an-operation-1"},
				{Name: "management-nodes-rollout"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}
		workflow := v1alpha1.Workflow{}

		dagTasks, _, err := iufSvc.getDAGTasks(&session, stageInfo, stages, globalParamsPerProduct, "global_params", "auth_token", &workflow)
		assert.NotEmpty(t, dagTasks)
		assert.NoError(t, err)
	})
	t.Run("It should skip operations which do not have the appropriate required attributes in the manifest.", func(t *testing.T) {
		session := iuf.Session{
			Products: []iuf.Product{
				{
					Name: "product_A",
					Manifest: `
---
iuf_version: ^0.5.0
name: product_A

content:
 op1: hello
 op2:
  - op2_A
  - op2_B
`,
				},
				{
					Name: "product_B",
					Manifest: `
---
iuf_version: ^0.5.0
name: product_B

content:
 op1: hello
`,
				},
				{
					Name: "product_C",
					Manifest: `
---
iuf_version: ^0.5.0
name: product_C

content:
 op2:
  - op2_A
  - op2_B
`,
				},
			},
			ActivityRef: activityName,
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "product",
			Operations: []iuf.Operations{
				{
					Name:                       "this-is-an-operation-1",
					RequiredManifestAttributes: []string{"content.op1"},
				},
				{
					Name:                       "this-is-an-operation-2",
					RequiredManifestAttributes: []string{"content.op2"},
				},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}
		workflow := v1alpha1.Workflow{}

		dagTasks, _, err := iufSvc.getDAGTasks(&session, stageInfo, stages, globalParamsPerProduct, "global_params", "auth_token",&workflow)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 6, len(dagTasks))
		assert.True(t, strings.Contains(dagTasks[0].Name, "this-is-an-operation-1"))
		assert.True(t, strings.Contains(dagTasks[0].Name, "product-A"))
		assert.True(t, dagTasks[0].TemplateRef.Name == "this-is-an-operation-1")

		assert.True(t, strings.Contains(dagTasks[1].Name, "this-is-an-operation-2"))
		assert.True(t, strings.Contains(dagTasks[1].Name, "product-A"))
		assert.True(t, dagTasks[1].TemplateRef.Name == "this-is-an-operation-2")

		assert.True(t, strings.Contains(dagTasks[2].Name, "this-is-an-operation-1"))
		assert.True(t, strings.Contains(dagTasks[2].Name, "product-B"))
		assert.True(t, dagTasks[2].TemplateRef.Name == "this-is-an-operation-1")

		assert.True(t, strings.Contains(dagTasks[3].Name, "this-is-an-operation-2"))
		assert.True(t, strings.Contains(dagTasks[3].Name, "product-B"))
		assert.True(t, dagTasks[3].TemplateRef.Name == "echo-template")

		assert.True(t, strings.Contains(dagTasks[4].Name, "this-is-an-operation-1"))
		assert.True(t, strings.Contains(dagTasks[4].Name, "product-C"))
		assert.True(t, dagTasks[4].TemplateRef.Name == "echo-template")

		assert.True(t, strings.Contains(dagTasks[5].Name, "this-is-an-operation-2"))
		assert.True(t, strings.Contains(dagTasks[5].Name, "product-C"))
		assert.True(t, dagTasks[5].TemplateRef.Name == "this-is-an-operation-2")
	})

	t.Run("It should split up large workflows into smaller workflows -- two large workflows", func(t *testing.T) {
		var products []iuf.Product
		for i := 0; i < 30; i++ {
			products = append(products, iuf.Product{Name: "product_" + strconv.Itoa(i)})
		}

		session := iuf.Session{
			Products:     products,
			ActivityRef:  activityName,
			CurrentStage: "deliver-product",
			InputParameters: iuf.InputParameters{
				Force:  true,
				Stages: []string{"deliver-product"},
			},
		}

		numOperations := 2
		stagesMetadata, err := iufSvc.GetStages()
		for _, stage := range stagesMetadata.Stages {
			if stage.Name == "deliver-product" {
				numOperations = len(stage.Operations)
				break
			}
		}

		// this is a predetermined number from running this test. Remember that the workflow is split as per the size of the JSON serialized form of the workflow.
		expectedProductsToProcessInFirstWorkflow := 15

		workflowRes, err, _ := iufSvc.workflowGen(&session)
		assert.NoError(t, err, "Should not have had an error when generating first workflow")
		assert.NotEmpty(t, workflowRes.Spec.Templates[0].DAG.Tasks, "Tasks should not be empty for first workflow")
		assert.Equal(t, workflowRes.Labels[LABEL_PARTIAL_WORKFLOW], "true", "partial workflow expected")
		assert.Equal(t, expectedProductsToProcessInFirstWorkflow, len(session.ProcessedProductsByStage["deliver-product"]), "Unexpected number of processed products for first workflow")
		assert.Equal(t, expectedProductsToProcessInFirstWorkflow*numOperations, len(workflowRes.Spec.Templates[0].DAG.Tasks), "Unexpected number of total operations for first workflow") // (number of products to process) * (2 ops per product as per above)

		// now let's try to get the next set of tasks
		workflowRes, err, _ = iufSvc.workflowGen(&session)
		assert.NoError(t, err, "Should not have had an error when generating second workflow")
		assert.NotEmpty(t, workflowRes.Spec.Templates[0].DAG.Tasks, "Tasks should not be empty for second workflow")
		assert.Equal(t, len(products), len(session.ProcessedProductsByStage["deliver-product"]), "Unexpected number of processed products for second workflow")
		assert.Equal(t, (len(products)-expectedProductsToProcessInFirstWorkflow)*numOperations, len(workflowRes.Spec.Templates[0].DAG.Tasks), "Unexpected number of total operations for second workflow")
	})

	t.Run("It should split up large workflows into smaller workflows -- one large workflow and one small", func(t *testing.T) {
		var products []iuf.Product
		for i := 0; i < 20; i++ {
			products = append(products, iuf.Product{Name: "product_" + strconv.Itoa(i)})
		}

		session := iuf.Session{
			Products:     products,
			ActivityRef:  activityName,
			CurrentStage: "deliver-product",
			InputParameters: iuf.InputParameters{
				Force:  true,
				Stages: []string{"deliver-product"},
			},
		}

		numOperations := 2
		stagesMetadata, err := iufSvc.GetStages()
		for _, stage := range stagesMetadata.Stages {
			if stage.Name == "deliver-product" {
				numOperations = len(stage.Operations)
				break
			}
		}

		// this is a predetermined number from running this test. Remember that the workflow is split as per the size of the JSON serialized form of the workflow.
		expectedProductsToProcessInFirstWorkflow := 15

		workflowRes, err, _ := iufSvc.workflowGen(&session)
		assert.NoError(t, err, "Should not have had an error when generating first workflow")
		assert.NotEmpty(t, workflowRes.Spec.Templates[0].DAG.Tasks, "Tasks should not be empty for first workflow")
		assert.Equal(t, workflowRes.Labels[LABEL_PARTIAL_WORKFLOW], "true", "partial workflow expected")
		assert.Equal(t, expectedProductsToProcessInFirstWorkflow, len(session.ProcessedProductsByStage["deliver-product"]), "Unexpected number of processed products for first workflow")
		assert.Equal(t, expectedProductsToProcessInFirstWorkflow*numOperations, len(workflowRes.Spec.Templates[0].DAG.Tasks), "Unexpected number of total operations for first workflow") // (number of products to process) * (2 ops per product as per above)

		// now let's try to get the next set of tasks
		workflowRes, err, _ = iufSvc.workflowGen(&session)
		assert.NoError(t, err, "Should not have had an error when generating second workflow")
		assert.NotEmpty(t, workflowRes.Spec.Templates[0].DAG.Tasks, "Tasks should not be empty for second workflow")
		assert.Equal(t, len(products), len(session.ProcessedProductsByStage["deliver-product"]), "Unexpected number of processed products for second workflow")
		assert.Equal(t, (len(products)-expectedProductsToProcessInFirstWorkflow)*numOperations, len(workflowRes.Spec.Templates[0].DAG.Tasks), "Unexpected number of total operations for second workflow")
	})

	t.Run("It should not split up large workflows into smaller workflows when products are containable in a single workflow", func(t *testing.T) {
		var products []iuf.Product
		for i := 0; i < 15; i++ {
			products = append(products, iuf.Product{Name: "product_" + strconv.Itoa(i)})
		}

		session := iuf.Session{
			Products:     products,
			ActivityRef:  activityName,
			CurrentStage: "deliver-product",
			InputParameters: iuf.InputParameters{
				Force:  true,
				Stages: []string{"deliver-product"},
			},
		}

		numOperations := 2
		stagesMetadata, err := iufSvc.GetStages()
		for _, stage := range stagesMetadata.Stages {
			if stage.Name == "deliver-product" {
				numOperations = len(stage.Operations)
				break
			}
		}

		workflowRes, err, _ := iufSvc.workflowGen(&session)
		assert.NoError(t, err, "Should not have had an error when generating first workflow")
		assert.NotEmpty(t, workflowRes.Spec.Templates[0].DAG.Tasks, "Tasks should not be empty for first workflow")
		assert.Empty(t, workflowRes.Labels[LABEL_PARTIAL_WORKFLOW], "partial workflow unexpected")
		assert.Equal(t, len(products)*numOperations, len(workflowRes.Spec.Templates[0].DAG.Tasks), "Unexpected number of total operations for first workflow") // (number of products to process) * (2 ops per product as per above)
	})

	t.Run("It should not create on exit tasks when there are no products with onExit hooks defined", func(t *testing.T) {
		session := iuf.Session{
			Products: []iuf.Product{
				{
					Name:             "cos",
					Version:          "2.5.1",
					OriginalLocation: cosOriginalLocation,
					Manifest:         cosManifest,
				},
			},
			CurrentStage: "deliver-product",
			CurrentState: iuf.SessionStateInProgress,
			InputParameters: iuf.InputParameters{
				Stages: []string{"deliver-product"},
			},
			ActivityRef: activityName,
		}

		workflow, err, _ := iufSvc.workflowGen(&session)
		assert.NoError(t, err)

		assert.Empty(t, workflow.Spec.OnExit)
	})

	t.Run("It should create on exit tasks for products with onExit hooks defined", func(t *testing.T) {
		session := iuf.Session{
			Products: []iuf.Product{
				{
					Name:             "csm",
					Version:          "1.6.0",
					OriginalLocation: csmOriginalLocation,
					Manifest:         csmManifest,
				},
			},
			CurrentStage: "deliver-product",
			CurrentState: iuf.SessionStateInProgress,
			InputParameters: iuf.InputParameters{
				Stages: []string{"deliver-product"},
			},
			ActivityRef: activityName,
		}

		workflow, err, _ := iufSvc.workflowGen(&session)
		assert.NoError(t, err)

		assert.Equal(t, "onExitHandlers", workflow.Spec.OnExit)
		assert.Equal(t, 2, len(workflow.Spec.Templates))
		assert.Equal(t, "onExitHandlers", workflow.Spec.Templates[1].Name)
		assert.Equal(t, 1, len(workflow.Spec.Templates[1].Steps))
		assert.Equal(t, 1, len(workflow.Spec.Templates[1].Steps[0].Steps))
		assert.Equal(t, 3, len(workflow.Spec.Templates[1].Steps[0].Steps[0].Arguments.Parameters))
		assert.Equal(t, "script_path", workflow.Spec.Templates[1].Steps[0].Steps[0].Arguments.Parameters[2].Name)
		assert.Equal(t, v1alpha1.AnyStringPtr("/etc/cray/upgrade/csm/test-activity/csm-160/on_exit/upgrade_k8s.sh"), workflow.Spec.Templates[1].Steps[0].Steps[0].Arguments.Parameters[2].Value)
	})

	t.Run("It should create on exit tasks for products with onExit hooks defined only for the last partial workflow and no other partial workflow", func(t *testing.T) {
		var products []iuf.Product

		products = append(products, iuf.Product{
			Name:             "csm",
			Version:          "1.6.0",
			OriginalLocation: csmOriginalLocation,
			Manifest:         csmManifest,
		})

		for i := 0; i < 19; i++ {
			products = append(products, iuf.Product{Name: "product_" + strconv.Itoa(i)})
		}

		session := iuf.Session{
			Products:     products,
			ActivityRef:  activityName,
			CurrentStage: "deliver-product",
			InputParameters: iuf.InputParameters{
				Force:  true,
				Stages: []string{"deliver-product"},
			},
		}

		// this is a predetermined number from running this test. Remember that the workflow is split as per the size of the JSON serialized form of the workflow.
		expectedProductsToProcessInFirstWorkflow := 15

		workflowRes, err, _ := iufSvc.workflowGen(&session)
		assert.NoError(t, err)
		assert.Equal(t, len(session.ProcessedProductsByStage["deliver-product"]), expectedProductsToProcessInFirstWorkflow)
		assert.Empty(t, workflowRes.Spec.OnExit)

		// now let's try to get the next set of tasks
		workflowRes, err, _ = iufSvc.workflowGen(&session)
		assert.NoError(t, err)
		assert.Equal(t, len(session.ProcessedProductsByStage["deliver-product"]), len(products))

		assert.Equal(t, "onExitHandlers", workflowRes.Spec.OnExit)
		assert.Equal(t, 2, len(workflowRes.Spec.Templates))
		assert.Equal(t, "onExitHandlers", workflowRes.Spec.Templates[1].Name)
		assert.Equal(t, 1, len(workflowRes.Spec.Templates[1].Steps))
		assert.Equal(t, 1, len(workflowRes.Spec.Templates[1].Steps[0].Steps))
		assert.Equal(t, 3, len(workflowRes.Spec.Templates[1].Steps[0].Steps[0].Arguments.Parameters))
		assert.Equal(t, "script_path", workflowRes.Spec.Templates[1].Steps[0].Steps[0].Arguments.Parameters[2].Name)
		assert.Equal(t, v1alpha1.AnyStringPtr("/etc/cray/upgrade/csm/test-activity/csm-160/on_exit/upgrade_k8s.sh"), workflowRes.Spec.Templates[1].Steps[0].Steps[0].Arguments.Parameters[2].Value)
	})
}

func setup(t *testing.T) (string, string, iufService) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
	wfServiceClientMock.On(
		"GetWorkflow",
		mock.Anything,
		mock.Anything,
	).Return(new(v1alpha1.Workflow), nil)
	wfServiceClientMock.On(
		"ListWorkflows",
		mock.Anything,
		mock.Anything,
	).Return(new(v1alpha1.WorkflowList), nil)
	wfTemplateServiceClientMock := &workflowtemplatemocks.WorkflowTemplateServiceClient{}
	availableOps := []string{
		"this-is-an-operation-1", "this-is-an-operation-2",
		"extract-release-distributions",
		"loftsman-manifest-upload", "s3-upload",
		"nexus-setup", "nexus-rpm-upload",
		"nexus-setup", "nexus-rpm-upload",
		"nexus-docker-upload", "nexus-helm-upload",
		"vcs-upload", "ims-upload",
		"management-m001-rollout",
		"master-host-hook-script", "worker-host-hook-script",
		"management-worker-nodes-rollout",
		"management-storage-nodes-rollout",
		"management-two-master-nodes-rollout",
	}
	var availableTemplates v1alpha1.WorkflowTemplates

	for _, op := range availableOps {
		template := v1alpha1.WorkflowTemplate{}
		template.Name = op
		availableTemplates = append(availableTemplates, template)
	}
	templateList := v1alpha1.WorkflowTemplateList{
		Items: availableTemplates,
	}
	wt1 := v1alpha1.WorkflowTemplate{}
	wt1.Name = "this-is-an-operation-1"
	wt2 := v1alpha1.WorkflowTemplate{}
	wt2.Name = "this-is-an-operation-2"
	wtMasterHook := v1alpha1.WorkflowTemplate{}
	wtMasterHook.Name = "master-host-hook-script"
	wt2WorkerHook := v1alpha1.WorkflowTemplate{}
	wt2WorkerHook.Name = "worker-host-hook-script"
	wfTemplateServiceClientMock.On(
		"ListWorkflowTemplates",
		mock.Anything,
		mock.Anything,
	).Return(&templateList, nil)
	name := uuid.NewString()
	activity := iuf.Activity{
		Name: name,
	}
	reqBytes, _ := json.Marshal(activity)
	configmap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: DEFAULT_NAMESPACE,
			Labels: map[string]string{
				"type": LABEL_ACTIVITY,
			},
		},
		Data: map[string]string{LABEL_ACTIVITY: string(reqBytes)},
	}
	fakeClient := fake.NewSimpleClientset(&configmap)

	mockTokenValue := "{{workflow.parameters.auth_token}}"
	keycloakServiceMock := mocks.NewMockKeycloakService(ctrl)
	keycloakServiceMock.EXPECT().NewKeycloakAccessToken().Return(mockTokenValue, nil).AnyTimes()

	iufSvc := iufService{
		logger:                 utils.GetLogger(),
		workflowClient:         wfServiceClientMock,
		workflowTemplateClient: wfTemplateServiceClientMock,
		k8sRestClientSet:       fakeClient,
		keycloakService:        keycloakServiceMock,
		env:                    utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "./_test_data_"},
	}
	return name, mockTokenValue, iufSvc
}
