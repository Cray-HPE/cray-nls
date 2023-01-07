// Code generated by MockGen. DO NOT EDIT.
// Source: src/api/services/iuf/iuf.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	v1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
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

// ConfigMapDataToSession mocks base method.
func (m *MockIufService) ConfigMapDataToSession(data string) (iuf.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfigMapDataToSession", data)
	ret0, _ := ret[0].(iuf.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConfigMapDataToSession indicates an expected call of ConfigMapDataToSession.
func (mr *MockIufServiceMockRecorder) ConfigMapDataToSession(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigMapDataToSession", reflect.TypeOf((*MockIufService)(nil).ConfigMapDataToSession), data)
}

// CreateActivity mocks base method.
func (m *MockIufService) CreateActivity(req iuf.CreateActivityRequest) (iuf.Activity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateActivity", req)
	ret0, _ := ret[0].(iuf.Activity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateActivity indicates an expected call of CreateActivity.
func (mr *MockIufServiceMockRecorder) CreateActivity(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateActivity", reflect.TypeOf((*MockIufService)(nil).CreateActivity), req)
}

// CreateIufWorkflow mocks base method.
func (m *MockIufService) CreateIufWorkflow(req iuf.Session) (*v1alpha1.Workflow, error, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIufWorkflow", req)
	ret0, _ := ret[0].(*v1alpha1.Workflow)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// CreateIufWorkflow indicates an expected call of CreateIufWorkflow.
func (mr *MockIufServiceMockRecorder) CreateIufWorkflow(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIufWorkflow", reflect.TypeOf((*MockIufService)(nil).CreateIufWorkflow), req)
}

// GetActivity mocks base method.
func (m *MockIufService) GetActivity(name string) (iuf.Activity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActivity", name)
	ret0, _ := ret[0].(iuf.Activity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActivity indicates an expected call of GetActivity.
func (mr *MockIufServiceMockRecorder) GetActivity(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActivity", reflect.TypeOf((*MockIufService)(nil).GetActivity), name)
}

// GetActivityHistory mocks base method.
func (m *MockIufService) GetActivityHistory(activityName string, startTime int32) (iuf.History, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActivityHistory", activityName, startTime)
	ret0, _ := ret[0].(iuf.History)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActivityHistory indicates an expected call of GetActivityHistory.
func (mr *MockIufServiceMockRecorder) GetActivityHistory(activityName, startTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActivityHistory", reflect.TypeOf((*MockIufService)(nil).GetActivityHistory), activityName, startTime)
}

// GetSession mocks base method.
func (m *MockIufService) GetSession(sessionName string) (iuf.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", sessionName)
	ret0, _ := ret[0].(iuf.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockIufServiceMockRecorder) GetSession(sessionName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockIufService)(nil).GetSession), sessionName)
}

// GetStages mocks base method.
func (m *MockIufService) GetStages() (iuf.Stages, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStages")
	ret0, _ := ret[0].(iuf.Stages)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStages indicates an expected call of GetStages.
func (mr *MockIufServiceMockRecorder) GetStages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStages", reflect.TypeOf((*MockIufService)(nil).GetStages))
}

// HistoryRunAction mocks base method.
func (m *MockIufService) HistoryRunAction(activityName string, req iuf.HistoryRunActionRequest) (iuf.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HistoryRunAction", activityName, req)
	ret0, _ := ret[0].(iuf.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HistoryRunAction indicates an expected call of HistoryRunAction.
func (mr *MockIufServiceMockRecorder) HistoryRunAction(activityName, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HistoryRunAction", reflect.TypeOf((*MockIufService)(nil).HistoryRunAction), activityName, req)
}

// ListActivities mocks base method.
func (m *MockIufService) ListActivities() ([]iuf.Activity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListActivities")
	ret0, _ := ret[0].([]iuf.Activity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListActivities indicates an expected call of ListActivities.
func (mr *MockIufServiceMockRecorder) ListActivities() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListActivities", reflect.TypeOf((*MockIufService)(nil).ListActivities))
}

// ListActivityHistory mocks base method.
func (m *MockIufService) ListActivityHistory(activityName string) ([]iuf.History, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListActivityHistory", activityName)
	ret0, _ := ret[0].([]iuf.History)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListActivityHistory indicates an expected call of ListActivityHistory.
func (mr *MockIufServiceMockRecorder) ListActivityHistory(activityName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListActivityHistory", reflect.TypeOf((*MockIufService)(nil).ListActivityHistory), activityName)
}

// ListSessions mocks base method.
func (m *MockIufService) ListSessions(activityName string) ([]iuf.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSessions", activityName)
	ret0, _ := ret[0].([]iuf.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSessions indicates an expected call of ListSessions.
func (mr *MockIufServiceMockRecorder) ListSessions(activityName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSessions", reflect.TypeOf((*MockIufService)(nil).ListSessions), activityName)
}

// PatchActivity mocks base method.
func (m *MockIufService) PatchActivity(activity iuf.Activity, req iuf.PatchActivityRequest) (iuf.Activity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchActivity", activity, req)
	ret0, _ := ret[0].(iuf.Activity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchActivity indicates an expected call of PatchActivity.
func (mr *MockIufServiceMockRecorder) PatchActivity(activity, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchActivity", reflect.TypeOf((*MockIufService)(nil).PatchActivity), activity, req)
}

// ProcessOutput mocks base method.
func (m *MockIufService) ProcessOutput(session *iuf.Session, workflow *v1alpha1.Workflow) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessOutput", session, workflow)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessOutput indicates an expected call of ProcessOutput.
func (mr *MockIufServiceMockRecorder) ProcessOutput(session, workflow interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessOutput", reflect.TypeOf((*MockIufService)(nil).ProcessOutput), session, workflow)
}

// ReplaceHistoryComment mocks base method.
func (m *MockIufService) ReplaceHistoryComment(activityName string, startTime int32, req iuf.ReplaceHistoryCommentRequest) (iuf.History, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceHistoryComment", activityName, startTime, req)
	ret0, _ := ret[0].(iuf.History)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplaceHistoryComment indicates an expected call of ReplaceHistoryComment.
func (mr *MockIufServiceMockRecorder) ReplaceHistoryComment(activityName, startTime, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceHistoryComment", reflect.TypeOf((*MockIufService)(nil).ReplaceHistoryComment), activityName, startTime, req)
}

// RunNextStage mocks base method.
func (m *MockIufService) RunNextStage(session *iuf.Session) (iuf.SyncResponse, error, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunNextStage", session)
	ret0, _ := ret[0].(iuf.SyncResponse)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// RunNextStage indicates an expected call of RunNextStage.
func (mr *MockIufServiceMockRecorder) RunNextStage(session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunNextStage", reflect.TypeOf((*MockIufService)(nil).RunNextStage), session)
}

// UpdateActivityStateFromSessionState mocks base method.
func (m *MockIufService) UpdateActivityStateFromSessionState(session iuf.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateActivityStateFromSessionState", session)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateActivityStateFromSessionState indicates an expected call of UpdateActivityStateFromSessionState.
func (mr *MockIufServiceMockRecorder) UpdateActivityStateFromSessionState(session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateActivityStateFromSessionState", reflect.TypeOf((*MockIufService)(nil).UpdateActivityStateFromSessionState), session)
}

// UpdateSession mocks base method.
func (m *MockIufService) UpdateSession(session iuf.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSession", session)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSession indicates an expected call of UpdateSession.
func (mr *MockIufServiceMockRecorder) UpdateSession(session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSession", reflect.TypeOf((*MockIufService)(nil).UpdateSession), session)
}

// UpdateSessionAndActivity mocks base method.
func (m *MockIufService) UpdateSessionAndActivity(session iuf.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSessionAndActivity", session)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSessionAndActivity indicates an expected call of UpdateSessionAndActivity.
func (mr *MockIufServiceMockRecorder) UpdateSessionAndActivity(session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSessionAndActivity", reflect.TypeOf((*MockIufService)(nil).UpdateSessionAndActivity), session)
}
