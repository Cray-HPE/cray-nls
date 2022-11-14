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
	"testing"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	workflowmocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow/mocks"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes"
	fake "k8s.io/client-go/rest/fake"
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
	workflowSvc := iufService{
		logger:        utils.GetLogger(),
		workflowCient: wfServiceClientMock,
		env:           utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "/_test_data_"},
	}
	t.Run("It should get a dag task for per-product stage", func(t *testing.T) {
		session := iuf.Session{
			Products: []iuf.Product{{Name: "product_A"}, {Name: "product_B"}},
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

}

func TestRunNextStage(t *testing.T) {
	wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
	wfServiceClientMock.On(
		"CreateWorkflow",
		mock.Anything,
		mock.Anything,
	).Return(new(v1alpha1.Workflow), nil)
	fakeClient := fake.RESTClient{}
	workflowSvc := iufService{
		logger:           utils.GetLogger(),
		workflowCient:    wfServiceClientMock,
		k8sRestClientSet: kubernetes.New(&fakeClient),
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
			_, err := workflowSvc.RunNextStage(&tt.session, "test")
			if (err != nil) != tt.wanted.err {
				t.Errorf("got %v, wantErr %v", err, tt.wanted.err)
				return
			}
			assert.Equal(t, tt.wanted.sessionState, tt.session.CurrentState)
			assert.Equal(t, tt.wanted.sessionStage, tt.session.CurrentStage)
		})
	}
}
