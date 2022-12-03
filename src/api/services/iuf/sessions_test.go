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
package services_iuf

import (
	"encoding/json"
	"testing"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	workflowmocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow/mocks"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fake "k8s.io/client-go/kubernetes/fake"
)

func TestCreateIufWorkflow(t *testing.T) {

	t.Run("It can create a new iuf workflow", func(t *testing.T) {
		// setup mocks
		wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
		wfServiceClientMock.On(
			"CreateWorkflow",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.Workflow), nil)

		workflowSvc := iufService{
			logger:        utils.GetLogger(),
			workflowCient: wfServiceClientMock,
			env:           utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "./_test_data_"},
		}
		_, err := workflowSvc.CreateIufWorkflow(iuf.Session{InputParameters: iuf.InputParameters{Stages: []string{"process_media"}}})

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

		workflowSvc := iufService{
			logger:        utils.GetLogger(),
			workflowCient: wfServiceClientMock,
			env:           utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "/_test_data_"},
		}
		_, err := workflowSvc.CreateIufWorkflow(iuf.Session{InputParameters: iuf.InputParameters{Stages: []string{"process_media"}}})

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

		workflowSvc := iufService{
			logger:        utils.GetLogger(),
			workflowCient: wfServiceClientMock,
			env:           utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "/_test_data_"},
		}
		_, err := workflowSvc.CreateIufWorkflow(iuf.Session{InputParameters: iuf.InputParameters{Stages: []string{"break_it"}}})

		// we don't actually test the template render/upload
		// this is tested in the render package
		assert.NotNil(t, err)
	})
}

func TestGetDagTasks(t *testing.T) {
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
	workflowSvc := iufService{
		logger:           utils.GetLogger(),
		workflowCient:    wfServiceClientMock,
		k8sRestClientSet: fakeClient,
		env:              utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "/_test_data_"},
	}
	t.Run("It should get a dag task for per-product stage", func(t *testing.T) {
		session := iuf.Session{
			Products:    []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			ActivityRef: name,
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "product",
			Operations: []iuf.Operations{
				{Name: "this_is_an_operationr_1"},
				{Name: "this_is_an_operationr_2"},
			},
		}

		dagTasks := workflowSvc.getDagTasks(session, stageInfo)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 4, len(dagTasks))
	})
	t.Run("It should get a dag task for global stage", func(t *testing.T) {
		session := iuf.Session{
			Products:    []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
			ActivityRef: name,
		}
		stageInfo := iuf.Stage{
			Name: "this_is_a_stage_name",
			Type: "global",
			Operations: []iuf.Operations{
				{Name: "this_is_an_operationr_1"},
				{Name: "this_is_an_operationr_2"},
			},
		}

		dagTasks := workflowSvc.getDagTasks(session, stageInfo)
		assert.NotEmpty(t, dagTasks)
		assert.Equal(t, 2, len(dagTasks))
	})

}

func TestRunNextStage(t *testing.T) {
	wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
	wfServiceClientMock.On(
		"CreateWorkflow",
		mock.Anything,
		mock.Anything,
	).Return(new(v1alpha1.Workflow), nil)
	fakeClient := fake.NewSimpleClientset()
	workflowSvc := iufService{
		logger:           utils.GetLogger(),
		workflowCient:    wfServiceClientMock,
		k8sRestClientSet: fakeClient,
		env:              utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "./_test_data_"},
	}
	type wanted struct {
		err          bool
		sessionState iuf.SessionState
		sessionStage string
	}
	var tests = []struct {
		name    string
		session iuf.Session
		wanted  wanted
	}{
		{
			name: "first stage",
			session: iuf.Session{
				ActivityRef: "test",
				InputParameters: iuf.InputParameters{
					Stages: []string{"process_media"},
				},
			},
			wanted: wanted{
				err:          true,
				sessionState: iuf.SessionStateInProgress,
				sessionStage: "process_media",
			},
		},
		{
			name: "next stage",
			session: iuf.Session{
				ActivityRef: "test",
				InputParameters: iuf.InputParameters{
					Stages: []string{"process_media", "deliver_product"},
				},
				Workflows: []iuf.SessionWorkflow{{Id: "asdf"}},
			},
			wanted: wanted{
				err:          true,
				sessionState: iuf.SessionStateInProgress,
				sessionStage: "deliver_product",
			},
		},
		{
			name: "last stage",
			session: iuf.Session{
				ActivityRef: "test",
				InputParameters: iuf.InputParameters{
					Stages: []string{"process_media", "deliver_product"},
				},
				Workflows: []iuf.SessionWorkflow{{Id: "asdf"}},
			},
			wanted: wanted{
				err:          true,
				sessionState: iuf.SessionStateInProgress,
				sessionStage: "deliver_product",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := workflowSvc.RunNextStage(&tt.session)
			if (err != nil) != tt.wanted.err {
				t.Errorf("got %v, wantErr %v", err, tt.wanted.err)
				return
			}
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
		workflowCient:    wfServiceClientMock,
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
					Stages: []string{"process_media"},
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
					Stages: []string{"process_media"},
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
											  "process_media": {
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
											  "pre_install_check": {
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
					Stages: []string{"process_media"},
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
											  "process_media": {
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
											  "pre_install_check": {
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
					Stages: []string{"process_media"},
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
											  "process_media": {
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
											  "pre_install_check": {
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
					Stages: []string{"process_media"},
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
											  "process_media": {
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
											  "pre_install_check": {
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
		workflowCient:    wfServiceClientMock,
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
