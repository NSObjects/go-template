// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_ucase is a generated GoMock package.
package mock_ucase

import (
	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo"
	reflect "reflect"
)

// MockUserUsecase is a mock of UserUsecase interface
type MockUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseMockRecorder
}

// MockUserUsecaseMockRecorder is the mock recorder for MockUserUsecase
type MockUserUsecaseMockRecorder struct {
	mock *MockUserUsecase
}

// NewMockUserUsecase creates a new mock instance
func NewMockUserUsecase(ctrl *gomock.Controller) *MockUserUsecase {
	mock := &MockUserUsecase{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserUsecase) EXPECT() *MockUserUsecaseMockRecorder {
	return m.recorder
}

// GetUser mocks base method
func (m *MockUserUsecase) GetUser(c echo.Context) error {
	ret := m.ctrl.Call(m, "GetUser", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetUser indicates an expected call of GetUser
func (mr *MockUserUsecaseMockRecorder) GetUser(c interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserUsecase)(nil).GetUser), c)
}

// GetUserDetail mocks base method
func (m *MockUserUsecase) GetUserDetail(c echo.Context) error {
	ret := m.ctrl.Call(m, "GetUserDetail", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetUserDetail indicates an expected call of GetUserDetail
func (mr *MockUserUsecaseMockRecorder) GetUserDetail(c interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDetail", reflect.TypeOf((*MockUserUsecase)(nil).GetUserDetail), c)
}

// CreateUser mocks base method
func (m *MockUserUsecase) CreateUser(c echo.Context) error {
	ret := m.ctrl.Call(m, "CreateUser", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockUserUsecaseMockRecorder) CreateUser(c interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserUsecase)(nil).CreateUser), c)
}

// DeleteUser mocks base method
func (m *MockUserUsecase) DeleteUser(c echo.Context) error {
	ret := m.ctrl.Call(m, "DeleteUser", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockUserUsecaseMockRecorder) DeleteUser(c interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserUsecase)(nil).DeleteUser), c)
}

// UpdateUser mocks base method
func (m *MockUserUsecase) UpdateUser(c echo.Context) error {
	ret := m.ctrl.Call(m, "UpdateUser", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser
func (mr *MockUserUsecaseMockRecorder) UpdateUser(c interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserUsecase)(nil).UpdateUser), c)
}