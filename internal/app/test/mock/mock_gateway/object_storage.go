// Code generated by MockGen. DO NOT EDIT.
// Source: object_storage.go

// Package mock_gateway is a generated GoMock package.
package mock_gateway

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockObjectStorage is a mock of ObjectStorage interface
type MockObjectStorage struct {
	ctrl     *gomock.Controller
	recorder *MockObjectStorageMockRecorder
}

// MockObjectStorageMockRecorder is the mock recorder for MockObjectStorage
type MockObjectStorageMockRecorder struct {
	mock *MockObjectStorage
}

// NewMockObjectStorage creates a new mock instance
func NewMockObjectStorage(ctrl *gomock.Controller) *MockObjectStorage {
	mock := &MockObjectStorage{ctrl: ctrl}
	mock.recorder = &MockObjectStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockObjectStorage) EXPECT() *MockObjectStorageMockRecorder {
	return m.recorder
}

// Save mocks base method
func (m *MockObjectStorage) Save(ctx context.Context, object []byte, key, bucketName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, object, key, bucketName)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockObjectStorageMockRecorder) Save(ctx, object, key, bucketName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockObjectStorage)(nil).Save), ctx, object, key, bucketName)
}

// Find mocks base method
func (m *MockObjectStorage) Find(ctx context.Context, key, bucketName string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, key, bucketName)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockObjectStorageMockRecorder) Find(ctx, key, bucketName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockObjectStorage)(nil).Find), ctx, key, bucketName)
}
