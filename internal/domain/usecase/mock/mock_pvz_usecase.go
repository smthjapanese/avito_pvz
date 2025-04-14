// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/usecase/pvz_usecase.go
//
// Generated by this command:
//
//	mockgen -source=internal/domain/usecase/pvz_usecase.go -destination=internal/domain/usecase/mock/mock_pvz_usecase.go -package=mock_usecase
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"
	time "time"

	uuid "github.com/google/uuid"
	models "github.com/smthjapanese/avito_pvz/internal/domain/models"
	usecase "github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	gomock "go.uber.org/mock/gomock"
)

// MockPVZUseCase is a mock of PVZUseCase interface.
type MockPVZUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockPVZUseCaseMockRecorder
	isgomock struct{}
}

// MockPVZUseCaseMockRecorder is the mock recorder for MockPVZUseCase.
type MockPVZUseCaseMockRecorder struct {
	mock *MockPVZUseCase
}

// NewMockPVZUseCase creates a new mock instance.
func NewMockPVZUseCase(ctrl *gomock.Controller) *MockPVZUseCase {
	mock := &MockPVZUseCase{ctrl: ctrl}
	mock.recorder = &MockPVZUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPVZUseCase) EXPECT() *MockPVZUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPVZUseCase) Create(ctx context.Context, city models.City) (*models.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, city)
	ret0, _ := ret[0].(*models.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPVZUseCaseMockRecorder) Create(ctx, city any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPVZUseCase)(nil).Create), ctx, city)
}

// GetAll mocks base method.
func (m *MockPVZUseCase) GetAll(ctx context.Context) ([]*models.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*models.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPVZUseCaseMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPVZUseCase)(nil).GetAll), ctx)
}

// GetByID mocks base method.
func (m *MockPVZUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockPVZUseCaseMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockPVZUseCase)(nil).GetByID), ctx, id)
}

// List mocks base method.
func (m *MockPVZUseCase) List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*usecase.PVZWithReceptions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, startDate, endDate, page, limit)
	ret0, _ := ret[0].([]*usecase.PVZWithReceptions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockPVZUseCaseMockRecorder) List(ctx, startDate, endDate, page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPVZUseCase)(nil).List), ctx, startDate, endDate, page, limit)
}
