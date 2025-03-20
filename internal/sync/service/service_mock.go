// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -destination=service_mock.go -source=service.go -package=service
//

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	models "github.com/GlebRadaev/password-manager/internal/sync/models"
	data "github.com/GlebRadaev/password-manager/pkg/data"
	gomock "go.uber.org/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
	isgomock struct{}
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// AddConflicts mocks base method.
func (m *MockRepo) AddConflicts(ctx context.Context, conflicts []models.Conflict) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddConflicts", ctx, conflicts)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddConflicts indicates an expected call of AddConflicts.
func (mr *MockRepoMockRecorder) AddConflicts(ctx, conflicts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddConflicts", reflect.TypeOf((*MockRepo)(nil).AddConflicts), ctx, conflicts)
}

// DeleteConflicts mocks base method.
func (m *MockRepo) DeleteConflicts(ctx context.Context, conflictIDs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteConflicts", ctx, conflictIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteConflicts indicates an expected call of DeleteConflicts.
func (mr *MockRepoMockRecorder) DeleteConflicts(ctx, conflictIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConflicts", reflect.TypeOf((*MockRepo)(nil).DeleteConflicts), ctx, conflictIDs)
}

// GetConflict mocks base method.
func (m *MockRepo) GetConflict(ctx context.Context, conflictID string) (models.Conflict, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConflict", ctx, conflictID)
	ret0, _ := ret[0].(models.Conflict)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConflict indicates an expected call of GetConflict.
func (mr *MockRepoMockRecorder) GetConflict(ctx, conflictID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConflict", reflect.TypeOf((*MockRepo)(nil).GetConflict), ctx, conflictID)
}

// GetUnresolvedConflicts mocks base method.
func (m *MockRepo) GetUnresolvedConflicts(ctx context.Context, userID string) ([]models.Conflict, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnresolvedConflicts", ctx, userID)
	ret0, _ := ret[0].([]models.Conflict)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnresolvedConflicts indicates an expected call of GetUnresolvedConflicts.
func (mr *MockRepoMockRecorder) GetUnresolvedConflicts(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnresolvedConflicts", reflect.TypeOf((*MockRepo)(nil).GetUnresolvedConflicts), ctx, userID)
}

// ResolveConflict mocks base method.
func (m *MockRepo) ResolveConflict(ctx context.Context, conflictID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveConflict", ctx, conflictID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResolveConflict indicates an expected call of ResolveConflict.
func (mr *MockRepoMockRecorder) ResolveConflict(ctx, conflictID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveConflict", reflect.TypeOf((*MockRepo)(nil).ResolveConflict), ctx, conflictID)
}

// MockDataClient is a mock of DataClient interface.
type MockDataClient struct {
	ctrl     *gomock.Controller
	recorder *MockDataClientMockRecorder
	isgomock struct{}
}

// MockDataClientMockRecorder is the mock recorder for MockDataClient.
type MockDataClientMockRecorder struct {
	mock *MockDataClient
}

// NewMockDataClient creates a new mock instance.
func NewMockDataClient(ctrl *gomock.Controller) *MockDataClient {
	mock := &MockDataClient{ctrl: ctrl}
	mock.recorder = &MockDataClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataClient) EXPECT() *MockDataClientMockRecorder {
	return m.recorder
}

// BatchProcess mocks base method.
func (m *MockDataClient) BatchProcess(ctx context.Context, userID string, operations []*data.DataOperation) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchProcess", ctx, userID, operations)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchProcess indicates an expected call of BatchProcess.
func (mr *MockDataClientMockRecorder) BatchProcess(ctx, userID, operations any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchProcess", reflect.TypeOf((*MockDataClient)(nil).BatchProcess), ctx, userID, operations)
}

// CreateAddOperation mocks base method.
func (m *MockDataClient) CreateAddOperation(userID string, entry models.ClientData) *data.DataOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAddOperation", userID, entry)
	ret0, _ := ret[0].(*data.DataOperation)
	return ret0
}

// CreateAddOperation indicates an expected call of CreateAddOperation.
func (mr *MockDataClientMockRecorder) CreateAddOperation(userID, entry any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAddOperation", reflect.TypeOf((*MockDataClient)(nil).CreateAddOperation), userID, entry)
}

// CreateDeleteOperation mocks base method.
func (m *MockDataClient) CreateDeleteOperation(userID, dataID string) *data.DataOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDeleteOperation", userID, dataID)
	ret0, _ := ret[0].(*data.DataOperation)
	return ret0
}

// CreateDeleteOperation indicates an expected call of CreateDeleteOperation.
func (mr *MockDataClientMockRecorder) CreateDeleteOperation(userID, dataID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDeleteOperation", reflect.TypeOf((*MockDataClient)(nil).CreateDeleteOperation), userID, dataID)
}

// CreateUpdateOperation mocks base method.
func (m *MockDataClient) CreateUpdateOperation(userID string, entry models.ClientData) *data.DataOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUpdateOperation", userID, entry)
	ret0, _ := ret[0].(*data.DataOperation)
	return ret0
}

// CreateUpdateOperation indicates an expected call of CreateUpdateOperation.
func (mr *MockDataClientMockRecorder) CreateUpdateOperation(userID, entry any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUpdateOperation", reflect.TypeOf((*MockDataClient)(nil).CreateUpdateOperation), userID, entry)
}

// ListData mocks base method.
func (m *MockDataClient) ListData(ctx context.Context, userID string) ([]models.DataEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListData", ctx, userID)
	ret0, _ := ret[0].([]models.DataEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListData indicates an expected call of ListData.
func (mr *MockDataClientMockRecorder) ListData(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListData", reflect.TypeOf((*MockDataClient)(nil).ListData), ctx, userID)
}

// UpdateData mocks base method.
func (m *MockDataClient) UpdateData(ctx context.Context, userID string, entry models.ClientData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateData", ctx, userID, entry)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateData indicates an expected call of UpdateData.
func (mr *MockDataClientMockRecorder) UpdateData(ctx, userID, entry any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateData", reflect.TypeOf((*MockDataClient)(nil).UpdateData), ctx, userID, entry)
}
