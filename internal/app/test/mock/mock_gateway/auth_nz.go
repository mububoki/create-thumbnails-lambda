// Code generated by MockGen. DO NOT EDIT.
// Source: auth_nz.go

// Package mock_gateway is a generated GoMock package.
package mock_gateway

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAuthNZ is a mock of AuthNZ interface
type MockAuthNZ struct {
	ctrl     *gomock.Controller
	recorder *MockAuthNZMockRecorder
}

// MockAuthNZMockRecorder is the mock recorder for MockAuthNZ
type MockAuthNZMockRecorder struct {
	mock *MockAuthNZ
}

// NewMockAuthNZ creates a new mock instance
func NewMockAuthNZ(ctrl *gomock.Controller) *MockAuthNZ {
	mock := &MockAuthNZ{ctrl: ctrl}
	mock.recorder = &MockAuthNZMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthNZ) EXPECT() *MockAuthNZMockRecorder {
	return m.recorder
}

// CreateRole mocks base method
func (m *MockAuthNZ) CreateRole(ctx context.Context, roleName, serviceName string, actions []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRole", ctx, roleName, serviceName, actions)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRole indicates an expected call of CreateRole
func (mr *MockAuthNZMockRecorder) CreateRole(ctx, roleName, serviceName, actions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRole", reflect.TypeOf((*MockAuthNZ)(nil).CreateRole), ctx, roleName, serviceName, actions)
}
