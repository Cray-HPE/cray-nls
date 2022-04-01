// Code generated by MockGen. DO NOT EDIT.
// Source: workflow.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	v1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockWorkflowService is a mock of WorkflowService interface.
type MockWorkflowService struct {
	ctrl     *gomock.Controller
	recorder *MockWorkflowServiceMockRecorder
}

// MockWorkflowServiceMockRecorder is the mock recorder for MockWorkflowService.
type MockWorkflowServiceMockRecorder struct {
	mock *MockWorkflowService
}

// NewMockWorkflowService creates a new mock instance.
func NewMockWorkflowService(ctrl *gomock.Controller) *MockWorkflowService {
	mock := &MockWorkflowService{ctrl: ctrl}
	mock.recorder = &MockWorkflowServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorkflowService) EXPECT() *MockWorkflowServiceMockRecorder {
	return m.recorder
}

// CreateWorkflow mocks base method.
func (m *MockWorkflowService) CreateWorkflow(hostname string) (*v1alpha1.Workflow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWorkflow", hostname)
	ret0, _ := ret[0].(*v1alpha1.Workflow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWorkflow indicates an expected call of CreateWorkflow.
func (mr *MockWorkflowServiceMockRecorder) CreateWorkflow(hostname interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWorkflow", reflect.TypeOf((*MockWorkflowService)(nil).CreateWorkflow), hostname)
}

// GetWorkflows mocks base method.
func (m *MockWorkflowService) GetWorkflows(ctx *gin.Context) (*v1alpha1.WorkflowList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkflows", ctx)
	ret0, _ := ret[0].(*v1alpha1.WorkflowList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkflows indicates an expected call of GetWorkflows.
func (mr *MockWorkflowServiceMockRecorder) GetWorkflows(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkflows", reflect.TypeOf((*MockWorkflowService)(nil).GetWorkflows), ctx)
}

// InitializeWorkflowTemplate mocks base method.
func (m *MockWorkflowService) InitializeWorkflowTemplate(template []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitializeWorkflowTemplate", template)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitializeWorkflowTemplate indicates an expected call of InitializeWorkflowTemplate.
func (mr *MockWorkflowServiceMockRecorder) InitializeWorkflowTemplate(template interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitializeWorkflowTemplate", reflect.TypeOf((*MockWorkflowService)(nil).InitializeWorkflowTemplate), template)
}
