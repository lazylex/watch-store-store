// Code generated by MockGen. DO NOT EDIT.
// Source: http.go

// Package mock_http is a generated GoMock package.
package mock_http

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetricsInterface is a mock of MetricsInterface interface.
type MockMetricsInterface struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsInterfaceMockRecorder
}

// MockMetricsInterfaceMockRecorder is the mock recorder for MockMetricsInterface.
type MockMetricsInterfaceMockRecorder struct {
	mock *MockMetricsInterface
}

// NewMockMetricsInterface creates a new mock instance.
func NewMockMetricsInterface(ctrl *gomock.Controller) *MockMetricsInterface {
	mock := &MockMetricsInterface{ctrl: ctrl}
	mock.recorder = &MockMetricsInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsInterface) EXPECT() *MockMetricsInterfaceMockRecorder {
	return m.recorder
}

// RequestsDurationObserve mocks base method.
func (m *MockMetricsInterface) RequestsDurationObserve(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RequestsDurationObserve", arg0)
}

// RequestsDurationObserve indicates an expected call of RequestsDurationObserve.
func (mr *MockMetricsInterfaceMockRecorder) RequestsDurationObserve(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestsDurationObserve", reflect.TypeOf((*MockMetricsInterface)(nil).RequestsDurationObserve), arg0)
}

// RequestsTotalInc mocks base method.
func (m *MockMetricsInterface) RequestsTotalInc(arg0 map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RequestsTotalInc", arg0)
}

// RequestsTotalInc indicates an expected call of RequestsTotalInc.
func (mr *MockMetricsInterfaceMockRecorder) RequestsTotalInc(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestsTotalInc", reflect.TypeOf((*MockMetricsInterface)(nil).RequestsTotalInc), arg0)
}
