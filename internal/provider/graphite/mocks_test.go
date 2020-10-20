// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package graphite_test is a generated GoMock package.
package graphite_test

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockClient) Do(arg0 *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockClientMockRecorder) Do(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockClient)(nil).Do), arg0)
}

// MockProgressListener is a mock of ProgressListener interface
type MockProgressListener struct {
	ctrl     *gomock.Controller
	recorder *MockProgressListenerMockRecorder
}

// MockProgressListenerMockRecorder is the mock recorder for MockProgressListener
type MockProgressListenerMockRecorder struct {
	mock *MockProgressListener
}

// NewMockProgressListener creates a new mock instance
func NewMockProgressListener(ctrl *gomock.Controller) *MockProgressListener {
	mock := &MockProgressListener{ctrl: ctrl}
	mock.recorder = &MockProgressListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProgressListener) EXPECT() *MockProgressListenerMockRecorder {
	return m.recorder
}

// OnStepDone mocks base method
func (m *MockProgressListener) OnStepDone() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnStepDone")
}

// OnStepDone indicates an expected call of OnStepDone
func (mr *MockProgressListenerMockRecorder) OnStepDone() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnStepDone", reflect.TypeOf((*MockProgressListener)(nil).OnStepDone))
}

// OnStepQueued mocks base method
func (m *MockProgressListener) OnStepQueued() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnStepQueued")
}

// OnStepQueued indicates an expected call of OnStepQueued
func (mr *MockProgressListenerMockRecorder) OnStepQueued() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnStepQueued", reflect.TypeOf((*MockProgressListener)(nil).OnStepQueued))
}
