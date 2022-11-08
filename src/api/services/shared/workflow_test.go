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
package services_shared

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	argo_templates "github.com/Cray-HPE/cray-nls/src/api/argo-templates"
	models_nls "github.com/Cray-HPE/cray-nls/src/api/models/nls"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	workflowmocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow/mocks"
	wftemplatemocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate/mocks"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

func TestInitializeWorkflowTemplate(t *testing.T) {
	// setup mocks
	wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
	wftServiceSclientMock := &wftemplatemocks.WorkflowTemplateServiceClient{}
	wftServiceSclientMock.On(
		"ListWorkflowTemplates",
		mock.Anything,
		mock.Anything,
	).Return(new(v1alpha1.WorkflowTemplateList), nil)
	wftServiceSclientMock.On(
		"CreateWorkflowTemplate",
		mock.Anything,
		mock.Anything,
	).Return(nil, nil)

	workflowSvc := workflowService{
		logger:                utils.GetLogger(),
		ctx:                   context.Background(),
		workflowCient:         wfServiceClientMock,
		workflowTemplateCient: wftServiceSclientMock,
		env:                   utils.Env{},
	}
	t.Run("It should initialize workflow template", func(t *testing.T) {
		workflowTemplates, _ := argo_templates.GetWorkflowTemplate()
		for _, workflowTemplate := range workflowTemplates {
			err := workflowSvc.InitializeWorkflowTemplate(workflowTemplate)
			assert.Nil(t, err)
		}
		wftServiceSclientMock.AssertExpectations(t)
	})
}

func TestCreateRebuildWorkflow(t *testing.T) {

	t.Run("It can create a new workflow", func(t *testing.T) {
		// setup mocks
		wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
		wftServiceSclientMock := &wftemplatemocks.WorkflowTemplateServiceClient{}
		wfServiceClientMock.On(
			"ListWorkflows",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.WorkflowList), nil)
		wfServiceClientMock.On(
			"CreateWorkflow",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.Workflow), nil)

		workflowSvc := workflowService{
			logger:                utils.GetLogger(),
			ctx:                   context.Background(),
			workflowCient:         wfServiceClientMock,
			workflowTemplateCient: wftServiceSclientMock,
			env:                   utils.Env{WorkerRebuildWorkflowFiles: "../argo-templates/_test_data_"},
		}
		req := models_nls.CreateRebuildWorkflowRequest{
			Hosts: []string{"ncn-w001"},
		}
		_, err := workflowSvc.CreateRebuildWorkflow(req)

		// we don't actually test the template render/upload
		// this is tested in the render package
		assert.Contains(t, err.Error(), "template: pattern matches no files: `**/*.yaml`")
	})
	t.Run("It should NOT create a new workflow when there is a running one of same type", func(t *testing.T) {
		// setup mocks
		wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
		wftServiceSclientMock := &wftemplatemocks.WorkflowTemplateServiceClient{}
		wfServiceClientMock.On(
			"ListWorkflows",
			mock.Anything,
			mock.Anything,
		).Return(&v1alpha1.WorkflowList{Items: make(v1alpha1.Workflows, 2)}, nil)

		workflowSvc := workflowService{
			logger:                utils.GetLogger(),
			ctx:                   context.Background(),
			workflowCient:         wfServiceClientMock,
			workflowTemplateCient: wftServiceSclientMock,
			env:                   utils.Env{},
		}
		req := models_nls.CreateRebuildWorkflowRequest{
			Hosts: []string{"ncn-w001"},
		}
		_, err := workflowSvc.CreateRebuildWorkflow(req)

		// we don't actually test the template render/upload
		// this is tested in the render package
		assert.Contains(t, err.Error(), "another ncn rebuild workflow (type: worker) is running/failed")
		wfServiceClientMock.AssertExpectations(t)
	})
	t.Run("It should NOT create a new workflow when there is a failed/error one of same type", func(t *testing.T) {
		t.Skip("skip: TODO")
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should NOT create a new workflow when request has mixed type", func(t *testing.T) {
		workflowSvc := workflowService{
			logger:                utils.GetLogger(),
			ctx:                   context.Background(),
			workflowCient:         nil,
			workflowTemplateCient: nil,
			env:                   utils.Env{},
		}
		req := models_nls.CreateRebuildWorkflowRequest{
			Hosts: []string{"ncn-w001", "ncn-s001"},
		}
		_, err := workflowSvc.CreateRebuildWorkflow(req)
		assert.Contains(t, err.Error(), "hostnames cannot contain both worker and storage nodes. Only one node type is supported at a time")
	})
	t.Run("It should NOT create a new workflow when request has wrong hostname", func(t *testing.T) {
		workflowSvc := workflowService{
			logger:                utils.GetLogger(),
			ctx:                   context.Background(),
			workflowCient:         nil,
			workflowTemplateCient: nil,
			env:                   utils.Env{},
		}
		req := models_nls.CreateRebuildWorkflowRequest{
			Hosts: []string{"ncn-ws1", "ncn-s001"},
		}
		_, err := workflowSvc.CreateRebuildWorkflow(req)
		assert.Contains(t, err.Error(), "invalid worker or storage node hostname")
	})
}

func TestGetWorkflows(t *testing.T) {
	// setup mocks
	wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
	wftServiceSclientMock := &wftemplatemocks.WorkflowTemplateServiceClient{}
	wfServiceClientMock.On(
		"ListWorkflows",
		mock.Anything,
		mock.Anything,
	).Return(nil, nil)

	workflowSvc := workflowService{
		logger:                utils.GetLogger(),
		ctx:                   context.Background(),
		workflowCient:         wfServiceClientMock,
		workflowTemplateCient: wftServiceSclientMock,
		env:                   utils.Env{},
	}
	t.Run("It should get workflows", func(t *testing.T) {
		response := httptest.NewRecorder()
		context, _ := gin.CreateTestContext(response)

		context.Request, _ = http.NewRequest("GET", "/", nil)
		_, err := workflowSvc.GetWorkflows(context)
		assert.Nil(t, err)
		wfServiceClientMock.AssertExpectations(t)
	})

}

func TestGetWorkflowByName(t *testing.T) {
	// setup mocks
	wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
	wftServiceSclientMock := &wftemplatemocks.WorkflowTemplateServiceClient{}
	wfServiceClientMock.On(
		"GetWorkflow",
		mock.Anything,
		mock.Anything,
	).Return(nil, nil)

	workflowSvc := workflowService{
		logger:                utils.GetLogger(),
		ctx:                   context.Background(),
		workflowCient:         wfServiceClientMock,
		workflowTemplateCient: wftServiceSclientMock,
		env:                   utils.Env{},
	}
	t.Run("It should get workflow", func(t *testing.T) {
		response := httptest.NewRecorder()
		context, _ := gin.CreateTestContext(response)

		context.Request, _ = http.NewRequest("GET", "/", nil)
		_, err := workflowSvc.GetWorkflowByName("test", context)
		assert.Nil(t, err)
		wfServiceClientMock.AssertExpectations(t)
	})

}
