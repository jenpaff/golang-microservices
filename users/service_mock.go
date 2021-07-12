// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jenpaff/golang-microservices/users (interfaces: Service)

// Package users is a generated GoMock package.
package users

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	common "github.com/jenpaff/golang-microservices/common"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockService) CreateUser(arg0 context.Context, arg1, arg2, arg3 string) (*common.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*common.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockServiceMockRecorder) CreateUser(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockService)(nil).CreateUser), arg0, arg1, arg2, arg3)
}

// CreateUserWithNewFeature mocks base method.
func (m *MockService) CreateUserWithNewFeature(arg0 context.Context, arg1, arg2, arg3 string) (*common.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserWithNewFeature", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*common.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserWithNewFeature indicates an expected call of CreateUserWithNewFeature.
func (mr *MockServiceMockRecorder) CreateUserWithNewFeature(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserWithNewFeature", reflect.TypeOf((*MockService)(nil).CreateUserWithNewFeature), arg0, arg1, arg2, arg3)
}

// GetUser mocks base method.
func (m *MockService) GetUser(arg0 context.Context, arg1 string) (*common.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(*common.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockServiceMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockService)(nil).GetUser), arg0, arg1)
}
