// Code generated by MockGen. DO NOT EDIT.
// Source: user_service.go
//
// Generated by this command:
//
//	mockgen -source=user_service.go -destination=user_service_mock_test.go -package=services
//

// Package services is a generated GoMock package.
package services

import (
	reflect "reflect"

	models "github.com/sgitwhyd/music-catalogue/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
	isgomock struct{}
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockUserRepo) Find(email, username string, id uint) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", email, username, id)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockUserRepoMockRecorder) Find(email, username, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockUserRepo)(nil).Find), email, username, id)
}

// Upsert mocks base method.
func (m *MockUserRepo) Upsert(model models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", model)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upsert indicates an expected call of Upsert.
func (mr *MockUserRepoMockRecorder) Upsert(model any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockUserRepo)(nil).Upsert), model)
}

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
	isgomock struct{}
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockUserService) Login(request models.SignInRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserServiceMockRecorder) Login(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserService)(nil).Login), request)
}

// Register mocks base method.
func (m *MockUserService) Register(request models.SignUpRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockUserServiceMockRecorder) Register(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserService)(nil).Register), request)
}
