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
	"regexp"
	"testing"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	workflowmocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow/mocks"
	workflowtemplatemocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate/mocks"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fake "k8s.io/client-go/kubernetes/fake"
)

func TestCreateIufWorkflow(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTokenValue := "mock_token"
	keycloakServiceMock := mocks.NewMockKeycloakService(ctrl)
	keycloakServiceMock.EXPECT().NewKeycloakAccessToken().Return(mockTokenValue, nil).AnyTimes()

	t.Run("It can create a new iuf workflow", func(t *testing.T) {
		// setup mocks
		wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
		wfServiceClientMock.On(
			"CreateWorkflow",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.Workflow), nil)
		wfTemplateServiceClientMock := &workflowtemplatemocks.WorkflowTemplateServiceClient{}
		wfTemplateServiceClientMock.On(
			"ListWorkflowTemplates",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.WorkflowTemplateList), nil)
		fakeClient := fake.NewSimpleClientset()
		workflowSvc := iufService{
			keycloakService:        keycloakServiceMock,
			logger:                 utils.GetLogger(),
			workflowClient:         wfServiceClientMock,
			workflowTemplateClient: wfTemplateServiceClientMock,
			k8sRestClientSet:       fakeClient,
			env:                    utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "./_test_data_"},
		}
		_, err, _ := workflowSvc.CreateIufWorkflow(iuf.Session{
			CurrentStage: "process-media",
			InputParameters: iuf.InputParameters{
				Stages: []string{"process-media"},
			},
		})

		// we don't actually test the template render/upload
		// this is tested in the render package
		assert.Nil(t, err)
	})
	t.Run("It should not create a new iuf workflow with wrong stages.yaml", func(t *testing.T) {
		// setup mocks
		wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
		wfServiceClientMock.On(
			"CreateWorkflow",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.Workflow), nil)
		wfTemplateServiceClientMock := &workflowtemplatemocks.WorkflowTemplateServiceClient{}
		wfTemplateServiceClientMock.On(
			"ListWorkflowTemplates",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.WorkflowTemplateList), nil)
		fakeClient := fake.NewSimpleClientset()
		workflowSvc := iufService{
			logger:                 utils.GetLogger(),
			workflowClient:         wfServiceClientMock,
			workflowTemplateClient: wfTemplateServiceClientMock,
			keycloakService:        keycloakServiceMock,
			k8sRestClientSet:       fakeClient,
			env:                    utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "./_test_data_"},
		}
		_, err, _ := workflowSvc.CreateIufWorkflow(iuf.Session{InputParameters: iuf.InputParameters{Stages: []string{"unsupported_stage"}}})

		// we don't actually test the template render/upload
		// this is tested in the render package
		assert.NotNil(t, err)
	})
	t.Run("It should not create a new iuf workflow with MISSING stages.yaml", func(t *testing.T) {
		// setup mocks
		wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
		wfServiceClientMock.On(
			"CreateWorkflow",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.Workflow), nil)
		wfTemplateServiceClientMock := &workflowtemplatemocks.WorkflowTemplateServiceClient{}
		wfTemplateServiceClientMock.On(
			"ListWorkflowTemplates",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.WorkflowTemplateList), nil)
		fakeClient := fake.NewSimpleClientset()
		workflowSvc := iufService{
			logger:                 utils.GetLogger(),
			workflowClient:         wfServiceClientMock,
			workflowTemplateClient: wfTemplateServiceClientMock,
			keycloakService:        keycloakServiceMock,
			k8sRestClientSet:       fakeClient,
			env:                    utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "./nowhere_to_be_found"},
		}
		_, err, _ := workflowSvc.CreateIufWorkflow(iuf.Session{InputParameters: iuf.InputParameters{Stages: []string{"process-media"}}})

		// we don't actually test the template render/upload
		// this is tested in the render package
		assert.NotNil(t, err)
	})
	t.Run("It should not create a new iuf workflow with wrong stage", func(t *testing.T) {
		// setup mocks
		wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
		wfServiceClientMock.On(
			"CreateWorkflow",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.Workflow), nil)
		wfTemplateServiceClientMock := &workflowtemplatemocks.WorkflowTemplateServiceClient{}
		wfTemplateServiceClientMock.On(
			"ListWorkflowTemplates",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.WorkflowTemplateList), nil)

		workflowSvc := iufService{
			logger:                 utils.GetLogger(),
			workflowClient:         wfServiceClientMock,
			workflowTemplateClient: wfTemplateServiceClientMock,
			keycloakService:        keycloakServiceMock,
			env:                    utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "./_test_data_"},
		}
		_, err, _ := workflowSvc.CreateIufWorkflow(iuf.Session{InputParameters: iuf.InputParameters{Stages: []string{"break_it"}}})

		// we don't actually test the template render/upload
		// this is tested in the render package
		assert.NotNil(t, err)
	})
}

func TestRunNextStage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
	wfServiceClientMock.On(
		"CreateWorkflow",
		mock.Anything,
		mock.Anything,
	).Return(new(v1alpha1.Workflow), nil)
	wfTemplateServiceClientMock := &workflowtemplatemocks.WorkflowTemplateServiceClient{}
	wt1 := v1alpha1.WorkflowTemplate{}
	wt1.Name = "extract-release-distributions"
	wt2 := v1alpha1.WorkflowTemplate{}
	wt2.Name = "loftsman-manifest-upload"
	wt3 := v1alpha1.WorkflowTemplate{}
	wt3.Name = "loftsman-manifest-deploy"
	mockWorkflowTempateList := v1alpha1.WorkflowTemplateList{
		Items: v1alpha1.WorkflowTemplates{
			wt1, wt2, wt3,
		},
	}
	wfTemplateServiceClientMock.On(
		"ListWorkflowTemplates",
		mock.Anything,
		mock.Anything,
	).Return(&mockWorkflowTempateList, nil)

	mockTokenValue := "mock_token"
	keycloakServiceMock := mocks.NewMockKeycloakService(ctrl)
	keycloakServiceMock.EXPECT().NewKeycloakAccessToken().Return(mockTokenValue, nil).AnyTimes()

	fakeClient := fake.NewSimpleClientset()
	workflowSvc := iufService{
		logger:                 utils.GetLogger(),
		workflowClient:         wfServiceClientMock,
		workflowTemplateClient: wfTemplateServiceClientMock,
		k8sRestClientSet:       fakeClient,
		keycloakService:        keycloakServiceMock,
		env:                    utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "./_test_data_"},
	}
	activity, err := workflowSvc.CreateActivity(iuf.CreateActivityRequest{
		Name:          "test",
		ActivityState: iuf.ActivityStateWaitForAdmin,
	})
	if err != nil {
		t.Errorf("Unknown error occurred %v", err)
		return
	}

	type wanted struct {
		err          bool
		isCompleted  bool
		sessionState iuf.SessionState
		sessionStage string
	}
	commonProducts := []iuf.Product{
		iuf.Product{
			Name:    "cos",
			Version: "1.2.3",
		},
		iuf.Product{
			Name:    "sdu",
			Version: "2.3.4",
		},
	}
	var tests = []struct {
		name    string
		session iuf.Session
		wanted  wanted
	}{
		{
			name: "first stage",
			session: iuf.Session{
				ActivityRef:  activity.Name,
				CurrentStage: "",
				Products:     commonProducts,
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media"},
				},
			},
			wanted: wanted{
				err:          false,
				isCompleted:  false,
				sessionState: iuf.SessionStateInProgress,
				sessionStage: "process-media",
			},
		},
		{
			name: "next stage",
			session: iuf.Session{
				CurrentStage: "process-media",
				ActivityRef:  activity.Name,
				Products:     commonProducts,
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media", "deliver-product"},
				},
				Workflows: []iuf.SessionWorkflow{{Id: "asdf"}},
			},
			wanted: wanted{
				err:          false,
				isCompleted:  false,
				sessionState: iuf.SessionStateInProgress,
				sessionStage: "deliver-product",
			},
		},
		{
			name: "before last stage",
			session: iuf.Session{
				ActivityRef:  activity.Name,
				Products:     commonProducts,
				CurrentStage: "deliver-product",
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media", "deliver-product", "deploy-product"},
				},
				Workflows: []iuf.SessionWorkflow{{Id: "asdf"}},
			},
			wanted: wanted{
				err:          false,
				isCompleted:  false,
				sessionState: iuf.SessionStateInProgress,
				sessionStage: "deploy-product",
			},
		},
		{
			name: "after last stage",
			session: iuf.Session{
				ActivityRef:  activity.Name,
				Products:     commonProducts,
				CurrentStage: "deploy-product",
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media", "deliver-product", "deploy-product"},
				},
				Workflows: []iuf.SessionWorkflow{{Id: "asdf"}},
			},
			wanted: wanted{
				err:          false,
				isCompleted:  true,
				sessionState: iuf.SessionStateCompleted,
				sessionStage: "deploy-product",
			},
		},
		{
			name: "check completion where there are no stages to run",
			session: iuf.Session{
				ActivityRef:  activity.Name,
				Products:     commonProducts,
				CurrentStage: "",
				InputParameters: iuf.InputParameters{
					Stages: []string{},
				},
				Workflows: []iuf.SessionWorkflow{{Id: "asdf"}},
			},
			wanted: wanted{
				err:          false,
				isCompleted:  true,
				sessionState: iuf.SessionStateCompleted,
				sessionStage: "",
			},
		},
		{
			name: "restart if stage not found (in case input parameters were updated)",
			session: iuf.Session{
				ActivityRef:  activity.Name,
				Products:     commonProducts,
				CurrentStage: "pre-install-check",
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media", "deliver-product", "deploy-product"},
				},
				Workflows: []iuf.SessionWorkflow{{Id: "asdf"}},
			},
			wanted: wanted{
				err:          false,
				isCompleted:  false,
				sessionState: iuf.SessionStateInProgress,
				sessionStage: "process-media",
			},
		},
		{
			name: "go to next stage if there are no operations for a stage",
			session: iuf.Session{
				ActivityRef:  activity.Name,
				Products:     commonProducts,
				CurrentStage: "process-media",
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media", "pre-install-check", "deliver-product", "deploy-product"},
				},
				Workflows: []iuf.SessionWorkflow{{Id: "asdf"}},
			},
			wanted: wanted{
				err:          false,
				isCompleted:  false,
				sessionState: iuf.SessionStateInProgress,
				sessionStage: "deliver-product",
			},
		},
		{
			name: "complete session on last stages that do not have operations",
			session: iuf.Session{
				ActivityRef:  activity.Name,
				Products:     commonProducts,
				CurrentStage: "deploy-product",
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media", "deliver-product", "deploy-product", "post-install-service-check", "post-install-check"},
				},
				Workflows: []iuf.SessionWorkflow{{Id: "asdf"}},
			},
			wanted: wanted{
				err:          false,
				isCompleted:  true,
				sessionState: iuf.SessionStateCompleted,
				sessionStage: "post-install-check",
			},
		},
	}

	m1 := regexp.MustCompile(`[^a-zA-Z]`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.session.Name = utils.GenerateName(m1.ReplaceAllString(tt.name, "-"))
			_, err := workflowSvc.CreateSession(tt.session, tt.session.Name, activity)
			if err != nil {
				t.Errorf("got unexpted error while creating session %v", err)
				return
			}

			_, err, completed := workflowSvc.RunNextStage(&tt.session)
			if (err != nil) != tt.wanted.err {
				t.Errorf("got %v, wantErr %v", err, tt.wanted.err)
				return
			}
			assert.Equal(t, tt.wanted.isCompleted, completed)
			assert.Equal(t, tt.wanted.sessionState, tt.session.CurrentState)
			assert.Equal(t, tt.wanted.sessionStage, tt.session.CurrentStage)
		})
	}
}

func TestProcessOutput(t *testing.T) {
	wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
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
	mySvc := iufService{
		logger:           utils.GetLogger(),
		workflowClient:   wfServiceClientMock,
		k8sRestClientSet: fakeClient,
	}

	var tests = []struct {
		name     string
		session  iuf.Session
		workflow *v1alpha1.Workflow
		wantErr  bool
	}{
		{
			name: "invalid stage type",
			session: iuf.Session{
				ActivityRef: name,
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media"},
				},
			},
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"stage_type": "INVALID_STAGE_TYPE",
					},
				},
				Spec: v1alpha1.WorkflowSpec{
					Templates: []v1alpha1.Template{
						{
							DAG: &v1alpha1.DAGTemplate{
								Tasks: []v1alpha1.DAGTask{},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "product stage type: no outputs",
			session: iuf.Session{
				ActivityRef: name,
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media"},
				},
			},
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"stage_type": "product",
					},
				},
				Spec: v1alpha1.WorkflowSpec{
					Templates: []v1alpha1.Template{
						{
							DAG: &v1alpha1.DAGTemplate{
								Tasks: []v1alpha1.DAGTask{
									{
										Name: "product-name-dash-operation-name",
										TemplateRef: &v1alpha1.TemplateRef{
											Name: "this-is-a-name-of-templateRef",
										},
									},
								},
							},
						},
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Nodes: v1alpha1.Nodes{
						"this-is-a-name-of-templateRef": v1alpha1.NodeStatus{
							DisplayName: "product-name-dash-operation-name",
							Inputs: &v1alpha1.Inputs{
								Parameters: []v1alpha1.Parameter{
									{
										Name: "global_params",
										Value: v1alpha1.AnyStringPtr(`
										{
											"product_manifest": {
											  "products": {
												"cos": {
												  "manifest": {
													"iuf_version": "^0.1.0",
													"name": "cos",
													"description": "The Cray Operating System (COS).\n",
													"version": "2.5.34-20221012230953",
													"content": {
													  "docker": [
														{
														  "path": "docker/cray"
														}
													  ],
													  "rest of the file snipped in this example": []
													}
												  }
												},
												"sdu": {
												  "manifest": {
													"iuf_version": "^0.1.0",
													"name": "sdu",
													"rest of the file snipped in this example": {}
												  }
												}
											  },
											  "current_product": {
												"manifest": {
												  "iuf_version": "^0.1.0",
												  "name": "cos",
												  "description": "The Cray Operating System (COS).\n",
												  "version": "2.5.34-20221012230953",
												  "content": {
													"docker": [
													  {
														"path": "docker/cray"
													  }
													],
													"rest of the file snipped in this example": []
												  }
												}
											  }
											},
											"input_params": {
											  "products": ["cos", "sdu"],
											  "media_dir": "/iuf/alice_230116",
											  "bootprep_config_managed": [
												{
												  "contents": "boot prep file contents as a string"
												}
											  ],
											  "bootprep_config_management": [
												{
												  "contents": "boot prep file contents as a string"
												}
											  ],
											  "limit_nodes": ["x12413515", "x15464574"]
											},
											"site_params": {
											  "global": {
												"some_global_site_parameter": "lorem ipsum"
											  },
											  "products": {
												"cos": {
												  "vcs_branch": "integration-2.5.34",
												  "some_cos_site_parameter": "lorem ipsum"
												},
												"sdu": {
												  "vcs_branch": "integration-1.2.3",
												  "some_sdu_site_parameter": "lorem ipsum"
												}
											  },
											  "current_product": {
												"vcs_branch": "integration-2.5.34",
												"some_cos_site_parameter": "lorem ipsum"
											  }
											},
											"stage_params": {
											  "process-media": {
												"global": {},
												"products": {
												  "cos": {
													"parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
													"output_of_cos": "whatever"
												  },
												  "sdu": {
													"parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
													"output_of_sdu": "whatever"
												  }
												},
												"current_product": {
												  "parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
												  "output_of_cos": "whatever"
												}
											  },
											  "pre-install-check": {
												"global": {}
											  }
											}
										  }
										  `),
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "product stage type: with outputs",
			session: iuf.Session{
				ActivityRef: name,
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media"},
				},
			},
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"stage_type": "product",
					},
				},
				Spec: v1alpha1.WorkflowSpec{
					Templates: []v1alpha1.Template{
						{
							DAG: &v1alpha1.DAGTemplate{
								Tasks: []v1alpha1.DAGTask{
									{
										Name: "product-name-dash-operation-name",
										TemplateRef: &v1alpha1.TemplateRef{
											Name: "this-is-a-name-of-templateRef",
										},
									},
								},
							},
						},
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Nodes: v1alpha1.Nodes{
						"this-is-a-name-of-templateRef": v1alpha1.NodeStatus{
							DisplayName: "product-name-dash-operation-name",
							Inputs: &v1alpha1.Inputs{
								Parameters: []v1alpha1.Parameter{
									{
										Name: "global_params",
										Value: v1alpha1.AnyStringPtr(`
										{
											"product_manifest": {
											  "products": {
												"cos": {
												  "manifest": {
													"iuf_version": "^0.1.0",
													"name": "cos",
													"description": "The Cray Operating System (COS).\n",
													"version": "2.5.34-20221012230953",
													"content": {
													  "docker": [
														{
														  "path": "docker/cray"
														}
													  ],
													  "rest of the file snipped in this example": []
													}
												  }
												},
												"sdu": {
												  "manifest": {
													"iuf_version": "^0.1.0",
													"name": "sdu",
													"rest of the file snipped in this example": {}
												  }
												}
											  },
											  "current_product": {
												"manifest": {
												  "iuf_version": "^0.1.0",
												  "name": "cos",
												  "description": "The Cray Operating System (COS).\n",
												  "version": "2.5.34-20221012230953",
												  "content": {
													"docker": [
													  {
														"path": "docker/cray"
													  }
													],
													"rest of the file snipped in this example": []
												  }
												}
											  }
											},
											"input_params": {
											  "products": ["cos", "sdu"],
											  "media_dir": "/iuf/alice_230116",
											  "bootprep_config_managed": [
												{
												  "contents": "boot prep file contents as a string"
												}
											  ],
											  "bootprep_config_management": [
												{
												  "contents": "boot prep file contents as a string"
												}
											  ],
											  "limit_nodes": ["x12413515", "x15464574"]
											},
											"site_params": {
											  "global": {
												"some_global_site_parameter": "lorem ipsum"
											  },
											  "products": {
												"cos": {
												  "vcs_branch": "integration-2.5.34",
												  "some_cos_site_parameter": "lorem ipsum"
												},
												"sdu": {
												  "vcs_branch": "integration-1.2.3",
												  "some_sdu_site_parameter": "lorem ipsum"
												}
											  },
											  "current_product": {
												"vcs_branch": "integration-2.5.34",
												"some_cos_site_parameter": "lorem ipsum"
											  }
											},
											"stage_params": {
											  "process-media": {
												"global": {},
												"products": {
												  "cos": {
													"parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
													"output_of_cos": "whatever"
												  },
												  "sdu": {
													"parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
													"output_of_sdu": "whatever"
												  }
												},
												"current_product": {
												  "parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
												  "output_of_cos": "whatever"
												}
											  },
											  "pre-install-check": {
												"global": {}
											  }
											}
										  }
										  `),
									},
								},
							},
							Outputs: &v1alpha1.Outputs{
								Parameters: []v1alpha1.Parameter{
									{
										Name:  "this_is_the_name_of_an_output",
										Value: v1alpha1.AnyStringPtr("this_is_the_value_of_an_output"),
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "global stage type: with outputs",
			session: iuf.Session{
				ActivityRef: name,
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media"},
				},
			},
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"stage_type": "global",
					},
				},
				Spec: v1alpha1.WorkflowSpec{
					Templates: []v1alpha1.Template{
						{
							DAG: &v1alpha1.DAGTemplate{
								Tasks: []v1alpha1.DAGTask{
									{
										Name: "product-name-dash-operation-name",
										TemplateRef: &v1alpha1.TemplateRef{
											Name: "this-is-a-name-of-templateRef",
										},
									},
								},
							},
						},
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Nodes: v1alpha1.Nodes{
						"this-is-a-name-of-templateRef": v1alpha1.NodeStatus{
							DisplayName: "product-name-dash-operation-name",
							Inputs: &v1alpha1.Inputs{
								Parameters: []v1alpha1.Parameter{
									{
										Name: "global_params",
										Value: v1alpha1.AnyStringPtr(`
										{
											"product_manifest": {
											  "products": {
												"cos": {
												  "manifest": {
													"iuf_version": "^0.1.0",
													"name": "cos",
													"description": "The Cray Operating System (COS).\n",
													"version": "2.5.34-20221012230953",
													"content": {
													  "docker": [
														{
														  "path": "docker/cray"
														}
													  ],
													  "rest of the file snipped in this example": []
													}
												  }
												},
												"sdu": {
												  "manifest": {
													"iuf_version": "^0.1.0",
													"name": "sdu",
													"rest of the file snipped in this example": {}
												  }
												}
											  },
											  "current_product": {
												"manifest": {
												  "iuf_version": "^0.1.0",
												  "name": "cos",
												  "description": "The Cray Operating System (COS).\n",
												  "version": "2.5.34-20221012230953",
												  "content": {
													"docker": [
													  {
														"path": "docker/cray"
													  }
													],
													"rest of the file snipped in this example": []
												  }
												}
											  }
											},
											"input_params": {
											  "products": ["cos", "sdu"],
											  "media_dir": "/iuf/alice_230116",
											  "bootprep_config_managed": [
												{
												  "contents": "boot prep file contents as a string"
												}
											  ],
											  "bootprep_config_management": [
												{
												  "contents": "boot prep file contents as a string"
												}
											  ],
											  "limit_nodes": ["x12413515", "x15464574"]
											},
											"site_params": {
											  "global": {
												"some_global_site_parameter": "lorem ipsum"
											  },
											  "products": {
												"cos": {
												  "vcs_branch": "integration-2.5.34",
												  "some_cos_site_parameter": "lorem ipsum"
												},
												"sdu": {
												  "vcs_branch": "integration-1.2.3",
												  "some_sdu_site_parameter": "lorem ipsum"
												}
											  },
											  "current_product": {
												"vcs_branch": "integration-2.5.34",
												"some_cos_site_parameter": "lorem ipsum"
											  }
											},
											"stage_params": {
											  "process-media": {
												"global": {},
												"products": {
												  "cos": {
													"parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
													"output_of_cos": "whatever"
												  },
												  "sdu": {
													"parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
													"output_of_sdu": "whatever"
												  }
												},
												"current_product": {
												  "parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
												  "output_of_cos": "whatever"
												}
											  },
											  "pre-install-check": {
												"global": {}
											  }
											}
										  }
										  `),
									},
								},
							},
							Outputs: &v1alpha1.Outputs{
								Parameters: []v1alpha1.Parameter{
									{
										Name:  "this_is_the_name_of_an_output",
										Value: v1alpha1.AnyStringPtr("this_is_the_value_of_an_output"),
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "global stage - process-media",
			session: iuf.Session{
				ActivityRef: name,
				InputParameters: iuf.InputParameters{
					Stages: []string{"process-media"},
				},
			},
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"stage_type": "global",
						"stage":      "process-media",
					},
				},
				Spec: v1alpha1.WorkflowSpec{
					Templates: []v1alpha1.Template{
						{
							DAG: &v1alpha1.DAGTemplate{
								Tasks: []v1alpha1.DAGTask{
									{
										Name: "product-name-dash-operation-name",
										TemplateRef: &v1alpha1.TemplateRef{
											Name: "this-is-a-name-of-templateRef",
										},
									},
								},
							},
						},
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Nodes: v1alpha1.Nodes{
						"this-is-a-name-of-templateRef": v1alpha1.NodeStatus{
							DisplayName: "product-name-dash-operation-name",
							Inputs: &v1alpha1.Inputs{
								Parameters: []v1alpha1.Parameter{
									{
										Name: "global_params",
										Value: v1alpha1.AnyStringPtr(`
										{
											"product_manifest": {
											  "products": {
												"cos": {
												  "manifest": {
													"iuf_version": "^0.1.0",
													"name": "cos",
													"description": "The Cray Operating System (COS).\n",
													"version": "2.5.34-20221012230953",
													"content": {
													  "docker": [
														{
														  "path": "docker/cray"
														}
													  ],
													  "rest of the file snipped in this example": []
													}
												  }
												},
												"sdu": {
												  "manifest": {
													"iuf_version": "^0.1.0",
													"name": "sdu",
													"rest of the file snipped in this example": {}
												  }
												}
											  },
											  "current_product": {
												"manifest": {
												  "iuf_version": "^0.1.0",
												  "name": "cos",
												  "description": "The Cray Operating System (COS).\n",
												  "version": "2.5.34-20221012230953",
												  "content": {
													"docker": [
													  {
														"path": "docker/cray"
													  }
													],
													"rest of the file snipped in this example": []
												  }
												}
											  }
											},
											"input_params": {
											  "products": ["cos", "sdu"],
											  "media_dir": "/iuf/alice_230116",
											  "bootprep_config_managed": [
												{
												  "contents": "boot prep file contents as a string"
												}
											  ],
											  "bootprep_config_management": [
												{
												  "contents": "boot prep file contents as a string"
												}
											  ],
											  "limit_nodes": ["x12413515", "x15464574"]
											},
											"site_params": {
											  "global": {
												"some_global_site_parameter": "lorem ipsum"
											  },
											  "products": {
												"cos": {
												  "vcs_branch": "integration-2.5.34",
												  "some_cos_site_parameter": "lorem ipsum"
												},
												"sdu": {
												  "vcs_branch": "integration-1.2.3",
												  "some_sdu_site_parameter": "lorem ipsum"
												}
											  },
											  "current_product": {
												"vcs_branch": "integration-2.5.34",
												"some_cos_site_parameter": "lorem ipsum"
											  }
											},
											"stage_params": {
											  "process-media": {
												"global": {},
												"products": {
												  "cos": {
													"parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
													"output_of_cos": "whatever"
												  },
												  "sdu": {
													"parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
													"output_of_sdu": "whatever"
												  }
												},
												"current_product": {
												  "parent_directory": "path to parent directory containing COS's iuf-product-manifest.yaml",
												  "output_of_cos": "whatever"
												}
											  },
											  "pre-install-check": {
												"global": {}
											  }
											}
										  }
										  `),
									},
								},
							},
							Outputs: &v1alpha1.Outputs{
								Parameters: []v1alpha1.Parameter{
									{
										Name:  "this_is_the_name_of_an_output",
										Value: v1alpha1.AnyStringPtr("iuf: test-me"),
									},
									{
										Name:  "this_is_the_name_of_an_output2",
										Value: v1alpha1.AnyStringPtr("this_is_the_value_of_an_output"),
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mySvc.ProcessOutput(&tt.session, tt.workflow)
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestProcessOutputOfProcessMedia(t *testing.T) {
	wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
	name := uuid.NewString()
	activity := iuf.Activity{
		Name:     name,
		Products: []iuf.Product{},
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
	mySvc := iufService{
		logger:           utils.GetLogger(),
		workflowClient:   wfServiceClientMock,
		k8sRestClientSet: fakeClient,
	}

	var tests = []struct {
		name         string
		workflow     *v1alpha1.Workflow
		wantActivity iuf.Activity
		wantErr      bool
	}{
		{
			name: "no outputs",
			workflow: &v1alpha1.Workflow{
				Status: v1alpha1.WorkflowStatus{
					Nodes: v1alpha1.Nodes{
						"this-is-a-name-of-templateRef": v1alpha1.NodeStatus{
							DisplayName: "product-name-dash-operation-name",
							Outputs:     &v1alpha1.Outputs{},
						},
					},
				},
			},
			wantActivity: activity,
			wantErr:      false,
		},
		{
			name: "only one output",
			workflow: &v1alpha1.Workflow{
				Status: v1alpha1.WorkflowStatus{
					Nodes: v1alpha1.Nodes{
						"this-is-a-name-of-templateRef": v1alpha1.NodeStatus{
							DisplayName: "product-name-dash-operation-name",
							Outputs: &v1alpha1.Outputs{
								Parameters: []v1alpha1.Parameter{
									{
										Name:  "this_is_the_name_of_an_output",
										Value: v1alpha1.AnyStringPtr("this_is_the_value_of_an_output"),
									},
								},
							},
						},
					},
				},
			},
			wantActivity: activity,
			wantErr:      false,
		},
		{
			name: "two outputs - invalid yaml",
			workflow: &v1alpha1.Workflow{
				Status: v1alpha1.WorkflowStatus{
					Nodes: v1alpha1.Nodes{
						"this-is-a-name-of-templateRef": v1alpha1.NodeStatus{
							DisplayName: "product-name-dash-operation-name",
							Outputs: &v1alpha1.Outputs{
								Parameters: []v1alpha1.Parameter{
									{
										Name:  "this_is_the_name_of_an_output",
										Value: v1alpha1.AnyStringPtr("this_is_the_value_of_an_output"),
									},
									{
										Name:  "this_is_the_name_of_another_output",
										Value: v1alpha1.AnyStringPtr("this_is_the_value_of_another_output"),
									},
								},
							},
						},
					},
				},
			},
			wantActivity: iuf.Activity{
				Name: activity.Name,
				OperationOutputs: map[string]interface{}{
					"stage_params": map[string]interface{}{
						"process-media": map[string]interface{}{
							"products": map[string]interface{}{},
						},
					},
				},
				Products: activity.Products,
			},
			wantErr: true,
		},
		{
			name: "two outputs - valid yaml but invalid iuf-manifest",
			workflow: &v1alpha1.Workflow{
				Status: v1alpha1.WorkflowStatus{
					Nodes: v1alpha1.Nodes{
						"this-is-a-name-of-templateRef": v1alpha1.NodeStatus{
							DisplayName: "product-name-dash-operation-name",
							Outputs: &v1alpha1.Outputs{
								Parameters: []v1alpha1.Parameter{
									{
										Name:  "this_is_the_name_of_an_output",
										Value: v1alpha1.AnyStringPtr("valid: yaml"),
									},
									{
										Name:  "this_is_the_name_of_another_output",
										Value: v1alpha1.AnyStringPtr("this_is_the_value_of_another_output"),
									},
								},
							},
						},
					},
				},
			},
			wantActivity: iuf.Activity{
				Name: activity.Name,
				OperationOutputs: map[string]interface{}{
					"stage_params": map[string]interface{}{
						"process-media": map[string]interface{}{
							"products": map[string]interface{}{},
						},
					},
				},
				Products: activity.Products,
			},
			wantErr: false,
		},
		{
			name: "two outputs - valid yaml, invalid manifest with only product name and version",
			workflow: &v1alpha1.Workflow{
				Status: v1alpha1.WorkflowStatus{
					Nodes: v1alpha1.Nodes{
						"this-is-a-name-of-templateRef": v1alpha1.NodeStatus{
							DisplayName: "product-name-dash-operation-name",
							Outputs: &v1alpha1.Outputs{
								Parameters: []v1alpha1.Parameter{
									{
										Name:  "this_is_the_name_of_product",
										Value: v1alpha1.AnyStringPtr("name: this-is-a-name\nversion: this-is-a-version"),
									},
									{
										Name:  "this_is_the_parent-directory",
										Value: v1alpha1.AnyStringPtr("this_is_the_value_of_parent-directory"),
									},
								},
							},
						},
					},
				},
			},
			wantActivity: iuf.Activity{
				Name: activity.Name,
				OperationOutputs: map[string]interface{}{
					"stage_params": map[string]interface{}{
						"process-media": map[string]interface{}{
							"products": map[string]interface{}{
								"this-is-a-name": map[string]interface{}{
									"parent_directory": "this_is_the_value_of_parent-directory",
								},
							},
						},
					},
				},
				Products: []iuf.Product{
					{
						Name:             "this-is-a-name",
						Version:          "this-is-a-version",
						OriginalLocation: "this_is_the_value_of_parent-directory",
						Validated:        false,
						Manifest:         `{"name":"this-is-a-name","version":"this-is-a-version"}`,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mySvc.processOutputOfProcessMedia(&activity, tt.workflow)
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(activity, tt.wantActivity) {
				t.Errorf("Wrong object received, got=%s", cmp.Diff(tt.wantActivity, activity))
			}
		})
	}
}
