// Code generated by MockGen. DO NOT EDIT.
// Source: source.go

// Package pgsql is a generated GoMock package.
package pgsql

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSource is a mock of Source interface.
type MockSource struct {
	ctrl     *gomock.Controller
	recorder *MockSourceMockRecorder
}

// MockSourceMockRecorder is the mock recorder for MockSource.
type MockSourceMockRecorder struct {
	mock *MockSource
}

// NewMockSource creates a new mock instance.
func NewMockSource(ctrl *gomock.Controller) *MockSource {
	mock := &MockSource{ctrl: ctrl}
	mock.recorder = &MockSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSource) EXPECT() *MockSourceMockRecorder {
	return m.recorder
}

// AddAnalytics mocks base method.
func (m *MockSource) AddAnalytics(ctx context.Context, data *AnalyticsData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAnalytics", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAnalytics indicates an expected call of AddAnalytics.
func (mr *MockSourceMockRecorder) AddAnalytics(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAnalytics", reflect.TypeOf((*MockSource)(nil).AddAnalytics), ctx, data)
}