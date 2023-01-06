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
	mocks "github.com/Cray-HPE/cray-nls/src/api/mocks/services"
	"github.com/golang/mock/gomock"
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

	t.Run("generated workflow must have NodeSelector set to ncn-m001 when NoHooks=false", func(t *testing.T) {
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
				Stages: []string{"process-media", "deliver-product"},
			},
			ActivityRef: activityName,
		}

		workflow, err, _ := iufSvc.workflowGen(session)
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
	t.Run("It should get a dag task for per-product stage", func(t *testing.T) {
		session := iuf.Session{
			Products:    []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			ActivityRef: activityName,
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "product",
			Operations: []iuf.Operations{
				{Name: "this_is_an_operation_1"},
				{Name: "this_is_an_operation_2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}

		dagTasks, err := iufSvc.getDAGTasks(session, stageInfo, stages)
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
				{Name: "this_is_an_operation_1"},
				{Name: "this_is_an_operation_2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}

		dagTasks, err := iufSvc.getDAGTasks(session, stageInfo, stages)
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
				{Name: "this_is_an_operation_1"},
				{Name: "this_is_an_operation_NOT_MOCKED"},
				{Name: "this_is_an_operation_2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}

		dagTasks, err := iufSvc.getDAGTasks(session, stageInfo, stages)
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
				{Name: "this_is_an_operation_NOT_MOCKED_1"},
				{Name: "this_is_an_operation_1"},
				{Name: "this_is_an_operation_NOT_MOCKED_2"},
				{Name: "this_is_an_operation_2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
		}

		dagTasks, err := iufSvc.getDAGTasks(session, stageInfo, stages)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 2, len(dagTasks))
		assert.Equal(t, 2, len(dagTasks[0].Arguments.Parameters))
		assert.Equal(t, 0, len(dagTasks[0].Dependencies))
		assert.Equal(t, "this_is_an_operation_1", dagTasks[0].Name)
		assert.Equal(t, "this_is_an_operation_2", dagTasks[1].Name)
		assert.Equal(t, 0, len(dagTasks[0].Dependencies))
		assert.Equal(t, 1, len(dagTasks[1].Dependencies))
		assert.Equal(t, "this_is_an_operation_1", dagTasks[1].Dependencies[0])
		assert.Equal(t, dagTasks[0].Arguments.Parameters[0].Name, "auth_token")
		assert.Equal(t, dagTasks[0].Arguments.Parameters[0].Value, v1alpha1.AnyStringPtr(mockAuthToken))
	})

	t.Run("It should get DAG tasks for existing hook templates defined for product operations", func(t *testing.T) {
		session := iuf.Session{
			Products: []iuf.Product{
				{
					Name:             "cos",
					OriginalLocation: cosOriginalLocation,
					Manifest:         cosManifest,
				},
				{
					Name:             "sdu",
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
				{Name: "this_is_an_operation_1"},
				{Name: "this_is_an_operation_NOT_MOCKED"},
				{Name: "this_is_an_operation_2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
			Hooks: map[string]string{
				"master_host": "master-host-hook-script",
				"worker_host": "worker-host-hook-script",
			},
		}

		dagTasks, err := iufSvc.getDAGTasks(session, stageInfo, stages)
		assert.NoError(t, err)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 7, len(dagTasks))

		// the first task will be a hook script (see cosManifest and sduManifest constants)
		assert.Equal(t, "cos-pre-hook-pre-install-check", dagTasks[0].Name)

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

			globalParams := iufSvc.getGlobalParams(session, product, stages)
			b, _ := json.Marshal(globalParams)
			assert.Equal(t, v1alpha1.AnyStringPtr(string(b)), dagTask.Arguments.GetParameterByName("global_params").Value)

			switch dagTask.Name {
			case "cos-pre-hook-pre-install-check":
				t.Run("cos pre hook script operation exist and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "master-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 0, len(dagTask.Dependencies))
					found++
				})
			case "cos-this_is_an_operation_1":
				t.Run("cos operation 1 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this_is_an_operation_1", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.Equal(t, "cos-pre-hook-pre-install-check", dagTask.Dependencies[0])
					found++
				})
			case "cos-this_is_an_operation_2":
				t.Run("cos operation 2 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this_is_an_operation_2", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.Equal(t, "cos-this_is_an_operation_1", dagTask.Dependencies[0])
					found++
				})
			case "sdu-this_is_an_operation_1":
				t.Run("sdu operation 1 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this_is_an_operation_1", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 0, len(dagTask.Dependencies))
					found++
				})
			case "sdu-this_is_an_operation_2":
				t.Run("sdu operation 2 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this_is_an_operation_2", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.Equal(t, "sdu-this_is_an_operation_1", dagTask.Dependencies[0])
					found++
				})
			case "cos-post-hook-pre-install-check":
				t.Run("cos post hook script operation exist and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "worker-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.Equal(t, "cos-this_is_an_operation_2", dagTask.Dependencies[0])
					found++
				})
			case "sdu-post-hook-pre-install-check":
				t.Run("sdu post hook script operation exist and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "worker-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.Equal(t, "sdu-this_is_an_operation_2", dagTask.Dependencies[0])
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
					OriginalLocation: cosOriginalLocation,
					Manifest:         cosManifest,
				},
				{
					Name:             "sdu",
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
				{Name: "this_is_an_operation_NOT_MOCKED_1"},
				{Name: "this_is_an_operation_1"},
				{Name: "this_is_an_operation_NOT_MOCKED_2"},
				{Name: "this_is_an_operation_2"},
			},
		}
		stages := iuf.Stages{
			Stages: []iuf.Stage{stageInfo},
			Hooks: map[string]string{
				"master_host": "master-host-hook-script",
				"worker_host": "worker-host-hook-script",
			},
		}

		dagTasks, err := iufSvc.getDAGTasks(session, stageInfo, stages)
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
			} else if strings.Contains(dagTask.Name, "sdu") {
				product = session.Products[1]
			}

			assert.Equal(t, v1alpha1.AnyStringPtr(mockAuthToken), dagTask.Arguments.GetParameterByName("auth_token").Value)

			globalParams := iufSvc.getGlobalParams(session, product, stages)
			b, _ := json.Marshal(globalParams)
			assert.Equal(t, v1alpha1.AnyStringPtr(string(b)), dagTask.Arguments.GetParameterByName("global_params").Value)

			switch dagTask.Name {
			case "cos-pre-hook-pre-install-check":
				t.Run("cos pre hook script operation exist and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "master-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 0, len(dagTask.Dependencies))
					found++
				})
			case "sdu-pre-hook-pre-install-check":
				t.Run("sdu pre hook script operation exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "worker-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 0, len(dagTask.Dependencies))
					found++
				})
			case "this_is_an_operation_1":
				t.Run("operation 1 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this_is_an_operation_1", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 2, len(dagTask.Dependencies))
					assert.Equal(t, true, strings.Contains(dagTask.Dependencies[0], "pre-hook-pre-install-check"))
					assert.Equal(t, true, strings.Contains(dagTask.Dependencies[1], "pre-hook-pre-install-check"))
					assert.NotEqual(t, dagTask.Dependencies[0], dagTask.Dependencies[1])
					found++
				})
			case "this_is_an_operation_2":
				t.Run("operation 2 exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "this_is_an_operation_2", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.Equal(t, "this_is_an_operation_1", dagTask.Dependencies[0])
					found++
				})
			case "cos-post-hook-pre-install-check":
				t.Run("cos post hook script operation exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "worker-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.Equal(t, "this_is_an_operation_2", dagTask.Dependencies[0])
					found++
				})
			case "sdu-post-hook-pre-install-check":
				t.Run("sdu post hook script operation exists and has the right dependencies", func(t *testing.T) {
					assert.Equal(t, "worker-host-hook-script", dagTask.TemplateRef.Name)
					assert.Equal(t, "main", dagTask.TemplateRef.Template)
					assert.Equal(t, 1, len(dagTask.Dependencies))
					assert.Equal(t, "this_is_an_operation_2", dagTask.Dependencies[0])
					found++
				})
			}
		}

		assert.Equal(t, 6, found)
	})
}

func setup(t *testing.T) (string, string, iufService) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
	wfTemplateServiceClientMock := &workflowtemplatemocks.WorkflowTemplateServiceClient{}
	availableOps := []string{
		"this_is_an_operation_1", "this_is_an_operation_2",
		"extract-release-distributions",
		"loftsman-manifest-upload", "s3-upload",
		"nexus-setup", "nexus-rpm-upload",
		"nexus-setup", "nexus-rpm-upload",
		"nexus-docker-upload", "nexus-helm-upload",
		"vcs-upload", "ims-upload",
		"management-m001-rollout",
		"master-host-hook-script", "worker-host-hook-script",
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
	wt1.Name = "this_is_an_operation_1"
	wt2 := v1alpha1.WorkflowTemplate{}
	wt2.Name = "this_is_an_operation_2"
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

	mockTokenValue := "mock_token"
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
