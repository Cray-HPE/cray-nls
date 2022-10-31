// Code generated by MockGen. DO NOT EDIT.
// Source: src/api//services/iuf.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	gomock "github.com/golang/mock/gomock"
)

// MockIufService is a mock of IufService interface.
type MockIufService struct {
	ctrl     *gomock.Controller
	recorder *MockIufServiceMockRecorder
}

// MockIufServiceMockRecorder is the mock recorder for MockIufService.
type MockIufServiceMockRecorder struct {
	mock *MockIufService
}

// NewMockIufService creates a new mock instance.
func NewMockIufService(ctrl *gomock.Controller) *MockIufService {
	mock := &MockIufService{ctrl: ctrl}
	mock.recorder = &MockIufServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIufService) EXPECT() *MockIufServiceMockRecorder {
	return m.recorder
}

// CreateActivity mocks base method.
func (m *MockIufService) CreateActivity(req iuf.CreateActivityRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateActivity", req)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateActivity indicates an expected call of CreateActivity.
func (mr *MockIufServiceMockRecorder) CreateActivity(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateActivity", reflect.TypeOf((*MockIufService)(nil).CreateActivity), req)
}

// GetSessionsByActivityName mocks base method.
func (m *MockIufService) GetSessionsByActivityName(activityName string) ([]iuf.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionsByActivityName", activityName)
	ret0, _ := ret[0].([]iuf.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionsByActivityName indicates an expected call of GetSessionsByActivityName.
func (mr *MockIufServiceMockRecorder) GetSessionsByActivityName(activityName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionsByActivityName", reflect.TypeOf((*MockIufService)(nil).GetSessionsByActivityName), activityName)
}
