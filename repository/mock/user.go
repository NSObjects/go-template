// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	gomock "github.com/golang/mock/gomock"
	models "go-template/models"
	reflect "reflect"
)

// MockUserRepository is a mock of UserRepository interface
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// GetUserById mocks base method
func (m *MockUserRepository) GetUserById(id int64) (models.User, error) {
	ret := m.ctrl.Call(m, "GetUserById", id)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById
func (mr *MockUserRepositoryMockRecorder) GetUserById(id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockUserRepository)(nil).GetUserById), id)
}

// FindUser mocks base method
func (m *MockUserRepository) FindUser(param models.UserParam) ([]models.User, int64, error) {
	ret := m.ctrl.Call(m, "FindUser", param)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindUser indicates an expected call of FindUser
func (mr *MockUserRepositoryMockRecorder) FindUser(param interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*MockUserRepository)(nil).FindUser), param)
}

// DeleteUserById mocks base method
func (m *MockUserRepository) DeleteUserById(id int64) error {
	ret := m.ctrl.Call(m, "DeleteUserById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserById indicates an expected call of DeleteUserById
func (mr *MockUserRepositoryMockRecorder) DeleteUserById(id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserById", reflect.TypeOf((*MockUserRepository)(nil).DeleteUserById), id)
}

// UpdateUser mocks base method
func (m *MockUserRepository) UpdateUser(param models.UserParam, id int64) error {
	ret := m.ctrl.Call(m, "UpdateUser", param, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser
func (mr *MockUserRepositoryMockRecorder) UpdateUser(param, id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepository)(nil).UpdateUser), param, id)
}

// CreateUser mocks base method
func (m *MockUserRepository) CreateUser(param models.UserParam) (int64, error) {
	ret := m.ctrl.Call(m, "CreateUser", param)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockUserRepositoryMockRecorder) CreateUser(param interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), param)
}
