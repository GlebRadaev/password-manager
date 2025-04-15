// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -destination=service_mock.go -source=service.go -package=services
//

// Package services is a generated GoMock package.
package services

import (
	http "net/http"
	reflect "reflect"

	models "github.com/GlebRadaev/password-manager/client/models"
	gomock "go.uber.org/mock/gomock"
)

// MockStorageInterface is a mock of StorageInterface interface.
type MockStorageInterface struct {
	ctrl     *gomock.Controller
	recorder *MockStorageInterfaceMockRecorder
	isgomock struct{}
}

// MockStorageInterfaceMockRecorder is the mock recorder for MockStorageInterface.
type MockStorageInterfaceMockRecorder struct {
	mock *MockStorageInterface
}

// NewMockStorageInterface creates a new mock instance.
func NewMockStorageInterface(ctrl *gomock.Controller) *MockStorageInterface {
	mock := &MockStorageInterface{ctrl: ctrl}
	mock.recorder = &MockStorageInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageInterface) EXPECT() *MockStorageInterfaceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockStorageInterface) Add(entry *models.DataEntry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", entry)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockStorageInterfaceMockRecorder) Add(entry any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockStorageInterface)(nil).Add), entry)
}

// ClearPendingSync mocks base method.
func (m *MockStorageInterface) ClearPendingSync() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearPendingSync")
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearPendingSync indicates an expected call of ClearPendingSync.
func (mr *MockStorageInterfaceMockRecorder) ClearPendingSync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearPendingSync", reflect.TypeOf((*MockStorageInterface)(nil).ClearPendingSync))
}

// Delete mocks base method.
func (m *MockStorageInterface) Delete(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStorageInterfaceMockRecorder) Delete(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorageInterface)(nil).Delete), id)
}

// Get mocks base method.
func (m *MockStorageInterface) Get(id string) (*models.DataEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*models.DataEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStorageInterfaceMockRecorder) Get(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStorageInterface)(nil).Get), id)
}

// GetAll mocks base method.
func (m *MockStorageInterface) GetAll() ([]*models.DataEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*models.DataEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockStorageInterfaceMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockStorageInterface)(nil).GetAll))
}

// GetAuthToken mocks base method.
func (m *MockStorageInterface) GetAuthToken() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthToken")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuthToken indicates an expected call of GetAuthToken.
func (mr *MockStorageInterfaceMockRecorder) GetAuthToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthToken", reflect.TypeOf((*MockStorageInterface)(nil).GetAuthToken))
}

// GetPendingSyncEntries mocks base method.
func (m *MockStorageInterface) GetPendingSyncEntries() ([]*models.DataEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingSyncEntries")
	ret0, _ := ret[0].([]*models.DataEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPendingSyncEntries indicates an expected call of GetPendingSyncEntries.
func (mr *MockStorageInterfaceMockRecorder) GetPendingSyncEntries() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingSyncEntries", reflect.TypeOf((*MockStorageInterface)(nil).GetPendingSyncEntries))
}

// UpdateSyncStatus mocks base method.
func (m *MockStorageInterface) UpdateSyncStatus(entries []*models.DataEntry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSyncStatus", entries)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSyncStatus indicates an expected call of UpdateSyncStatus.
func (mr *MockStorageInterfaceMockRecorder) UpdateSyncStatus(entries any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSyncStatus", reflect.TypeOf((*MockStorageInterface)(nil).UpdateSyncStatus), entries)
}

// MockHTTPClientInterface is a mock of HTTPClientInterface interface.
type MockHTTPClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPClientInterfaceMockRecorder
	isgomock struct{}
}

// MockHTTPClientInterfaceMockRecorder is the mock recorder for MockHTTPClientInterface.
type MockHTTPClientInterfaceMockRecorder struct {
	mock *MockHTTPClientInterface
}

// NewMockHTTPClientInterface creates a new mock instance.
func NewMockHTTPClientInterface(ctrl *gomock.Controller) *MockHTTPClientInterface {
	mock := &MockHTTPClientInterface{ctrl: ctrl}
	mock.recorder = &MockHTTPClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHTTPClientInterface) EXPECT() *MockHTTPClientInterfaceMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockHTTPClientInterface) Do(req *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", req)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockHTTPClientInterfaceMockRecorder) Do(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockHTTPClientInterface)(nil).Do), req)
}

// MockSyncServiceInterface is a mock of SyncServiceInterface interface.
type MockSyncServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockSyncServiceInterfaceMockRecorder
	isgomock struct{}
}

// MockSyncServiceInterfaceMockRecorder is the mock recorder for MockSyncServiceInterface.
type MockSyncServiceInterfaceMockRecorder struct {
	mock *MockSyncServiceInterface
}

// NewMockSyncServiceInterface creates a new mock instance.
func NewMockSyncServiceInterface(ctrl *gomock.Controller) *MockSyncServiceInterface {
	mock := &MockSyncServiceInterface{ctrl: ctrl}
	mock.recorder = &MockSyncServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSyncServiceInterface) EXPECT() *MockSyncServiceInterfaceMockRecorder {
	return m.recorder
}

// Resolve mocks base method.
func (m *MockSyncServiceInterface) Resolve(conflictID, strategy string) (*models.ResolutionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resolve", conflictID, strategy)
	ret0, _ := ret[0].(*models.ResolutionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Resolve indicates an expected call of Resolve.
func (mr *MockSyncServiceInterfaceMockRecorder) Resolve(conflictID, strategy any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resolve", reflect.TypeOf((*MockSyncServiceInterface)(nil).Resolve), conflictID, strategy)
}

// Sync mocks base method.
func (m *MockSyncServiceInterface) Sync() (*models.SyncResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync")
	ret0, _ := ret[0].(*models.SyncResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sync indicates an expected call of Sync.
func (mr *MockSyncServiceInterfaceMockRecorder) Sync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockSyncServiceInterface)(nil).Sync))
}
