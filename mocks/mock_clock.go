// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ansd/driving-time/clock (interfaces: Nower)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockNower is a mock of Nower interface
type MockNower struct {
	ctrl     *gomock.Controller
	recorder *MockNowerMockRecorder
}

// MockNowerMockRecorder is the mock recorder for MockNower
type MockNowerMockRecorder struct {
	mock *MockNower
}

// NewMockNower creates a new mock instance
func NewMockNower(ctrl *gomock.Controller) *MockNower {
	mock := &MockNower{ctrl: ctrl}
	mock.recorder = &MockNowerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNower) EXPECT() *MockNowerMockRecorder {
	return m.recorder
}

// Now mocks base method
func (m *MockNower) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now
func (mr *MockNowerMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockNower)(nil).Now))
}
