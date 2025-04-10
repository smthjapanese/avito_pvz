// Code generated by MockGen. DO NOT EDIT.
// Source: ../../domain/repository/reception_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	models "github.com/smthjapanese/avito_pvz/internal/domain/models"
)

// MockReceptionRepository is a mock of ReceptionRepository interface.
type MockReceptionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockReceptionRepositoryMockRecorder
}

// MockReceptionRepositoryMockRecorder is the mock recorder for MockReceptionRepository.
type MockReceptionRepositoryMockRecorder struct {
	mock *MockReceptionRepository
}

// NewMockReceptionRepository creates a new mock instance.
func NewMockReceptionRepository(ctrl *gomock.Controller) *MockReceptionRepository {
	mock := &MockReceptionRepository{ctrl: ctrl}
	mock.recorder = &MockReceptionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReceptionRepository) EXPECT() *MockReceptionRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockReceptionRepository) Create(ctx context.Context, reception *models.Reception) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, reception)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockReceptionRepositoryMockRecorder) Create(ctx, reception interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockReceptionRepository)(nil).Create), ctx, reception)
}

// GetByID mocks base method.
func (m *MockReceptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockReceptionRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockReceptionRepository)(nil).GetByID), ctx, id)
}

// GetLastByPVZID mocks base method.
func (m *MockReceptionRepository) GetLastByPVZID(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastByPVZID", ctx, pvzID)
	ret0, _ := ret[0].(*models.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastByPVZID indicates an expected call of GetLastByPVZID.
func (mr *MockReceptionRepositoryMockRecorder) GetLastByPVZID(ctx, pvzID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastByPVZID", reflect.TypeOf((*MockReceptionRepository)(nil).GetLastByPVZID), ctx, pvzID)
}

// GetLastOpenByPVZID mocks base method.
func (m *MockReceptionRepository) GetLastOpenByPVZID(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastOpenByPVZID", ctx, pvzID)
	ret0, _ := ret[0].(*models.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastOpenByPVZID indicates an expected call of GetLastOpenByPVZID.
func (mr *MockReceptionRepositoryMockRecorder) GetLastOpenByPVZID(ctx, pvzID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastOpenByPVZID", reflect.TypeOf((*MockReceptionRepository)(nil).GetLastOpenByPVZID), ctx, pvzID)
}

// ListByPVZID mocks base method.
func (m *MockReceptionRepository) ListByPVZID(ctx context.Context, pvzID uuid.UUID) ([]*models.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByPVZID", ctx, pvzID)
	ret0, _ := ret[0].([]*models.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByPVZID indicates an expected call of ListByPVZID.
func (mr *MockReceptionRepositoryMockRecorder) ListByPVZID(ctx, pvzID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByPVZID", reflect.TypeOf((*MockReceptionRepository)(nil).ListByPVZID), ctx, pvzID)
}

// Update mocks base method.
func (m *MockReceptionRepository) Update(ctx context.Context, reception *models.Reception) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, reception)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockReceptionRepositoryMockRecorder) Update(ctx, reception interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockReceptionRepository)(nil).Update), ctx, reception)
}
