// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=../../services/spotify/service_mock_test.go -package=spotify
//

// Package spotify is a generated GoMock package.
package spotify

import (
	context "context"
	reflect "reflect"

	spotify "github.com/sgitwhyd/music-catalogue/internal/models/spotify"
	spotify0 "github.com/sgitwhyd/music-catalogue/internal/repositorys/spotify"
	gomock "go.uber.org/mock/gomock"
)

// MockSpotifyOutbond is a mock of SpotifyOutbond interface.
type MockSpotifyOutbond struct {
	ctrl     *gomock.Controller
	recorder *MockSpotifyOutbondMockRecorder
	isgomock struct{}
}

// MockSpotifyOutbondMockRecorder is the mock recorder for MockSpotifyOutbond.
type MockSpotifyOutbondMockRecorder struct {
	mock *MockSpotifyOutbond
}

// NewMockSpotifyOutbond creates a new mock instance.
func NewMockSpotifyOutbond(ctrl *gomock.Controller) *MockSpotifyOutbond {
	mock := &MockSpotifyOutbond{ctrl: ctrl}
	mock.recorder = &MockSpotifyOutbondMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpotifyOutbond) EXPECT() *MockSpotifyOutbondMockRecorder {
	return m.recorder
}

// Search mocks base method.
func (m *MockSpotifyOutbond) Search(ctx context.Context, query string, limit, offset int) (*spotify0.SpotifySearchResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, query, limit, offset)
	ret0, _ := ret[0].(*spotify0.SpotifySearchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockSpotifyOutbondMockRecorder) Search(ctx, query, limit, offset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockSpotifyOutbond)(nil).Search), ctx, query, limit, offset)
}

// MockSpotifyRepository is a mock of SpotifyRepository interface.
type MockSpotifyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSpotifyRepositoryMockRecorder
	isgomock struct{}
}

// MockSpotifyRepositoryMockRecorder is the mock recorder for MockSpotifyRepository.
type MockSpotifyRepositoryMockRecorder struct {
	mock *MockSpotifyRepository
}

// NewMockSpotifyRepository creates a new mock instance.
func NewMockSpotifyRepository(ctrl *gomock.Controller) *MockSpotifyRepository {
	mock := &MockSpotifyRepository{ctrl: ctrl}
	mock.recorder = &MockSpotifyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpotifyRepository) EXPECT() *MockSpotifyRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSpotifyRepository) Create(ctx context.Context, model spotify.TrackActivity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, model)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockSpotifyRepositoryMockRecorder) Create(ctx, model any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSpotifyRepository)(nil).Create), ctx, model)
}

// Get mocks base method.
func (m *MockSpotifyRepository) Get(ctx context.Context, UserID uint, spotifyID string) (*spotify.TrackActivity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, UserID, spotifyID)
	ret0, _ := ret[0].(*spotify.TrackActivity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSpotifyRepositoryMockRecorder) Get(ctx, UserID, spotifyID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSpotifyRepository)(nil).Get), ctx, UserID, spotifyID)
}

// GetBulkSpotifyIDs mocks base method.
func (m *MockSpotifyRepository) GetBulkSpotifyIDs(ctx context.Context, UserID uint, spotifyIDs []string) (map[string]spotify.TrackActivity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBulkSpotifyIDs", ctx, UserID, spotifyIDs)
	ret0, _ := ret[0].(map[string]spotify.TrackActivity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBulkSpotifyIDs indicates an expected call of GetBulkSpotifyIDs.
func (mr *MockSpotifyRepositoryMockRecorder) GetBulkSpotifyIDs(ctx, UserID, spotifyIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBulkSpotifyIDs", reflect.TypeOf((*MockSpotifyRepository)(nil).GetBulkSpotifyIDs), ctx, UserID, spotifyIDs)
}

// Update mocks base method.
func (m *MockSpotifyRepository) Update(ctx context.Context, model spotify.TrackActivity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, model)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockSpotifyRepositoryMockRecorder) Update(ctx, model any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSpotifyRepository)(nil).Update), ctx, model)
}
