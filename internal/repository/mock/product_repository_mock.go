// Code generated by MockGen. DO NOT EDIT.
// Source: ../../domain/repository/product_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	models "github.com/smthjapanese/avito_pvz/internal/domain/models"
)

// MockProductRepository is a mock of ProductRepository interface.
type MockProductRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProductRepositoryMockRecorder
}

// MockProductRepositoryMockRecorder is the mock recorder for MockProductRepository.
type MockProductRepositoryMockRecorder struct {
	mock *MockProductRepository
}

// NewMockProductRepository creates a new mock instance.
func NewMockProductRepository(ctrl *gomock.Controller) *MockProductRepository {
	mock := &MockProductRepository{ctrl: ctrl}
	mock.recorder = &MockProductRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductRepository) EXPECT() *MockProductRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProductRepository) Create(ctx context.Context, product *models.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, product)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockProductRepositoryMockRecorder) Create(ctx, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProductRepository)(nil).Create), ctx, product)
}

// Delete mocks base method.
func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockProductRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProductRepository)(nil).Delete), ctx, id)
}

// GetByID mocks base method.
func (m *MockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockProductRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockProductRepository)(nil).GetByID), ctx, id)
}

// GetLastByReceptionID mocks base method.
func (m *MockProductRepository) GetLastByReceptionID(ctx context.Context, receptionID uuid.UUID) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastByReceptionID", ctx, receptionID)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastByReceptionID indicates an expected call of GetLastByReceptionID.
func (mr *MockProductRepositoryMockRecorder) GetLastByReceptionID(ctx, receptionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastByReceptionID", reflect.TypeOf((*MockProductRepository)(nil).GetLastByReceptionID), ctx, receptionID)
}

// ListByReceptionID mocks base method.
func (m *MockProductRepository) ListByReceptionID(ctx context.Context, receptionID uuid.UUID) ([]*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByReceptionID", ctx, receptionID)
	ret0, _ := ret[0].([]*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByReceptionID indicates an expected call of ListByReceptionID.
func (mr *MockProductRepositoryMockRecorder) ListByReceptionID(ctx, receptionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByReceptionID", reflect.TypeOf((*MockProductRepository)(nil).ListByReceptionID), ctx, receptionID)
}
