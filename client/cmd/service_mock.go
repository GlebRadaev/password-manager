// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=service_mock.go -package=cmd
//

// Package cmd is a generated GoMock package.
package cmd

import (
	reflect "reflect"

	models "github.com/GlebRadaev/password-manager/client/models"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthServiceInterface is a mock of AuthServiceInterface interface.
type MockAuthServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceInterfaceMockRecorder
	isgomock struct{}
}

// MockAuthServiceInterfaceMockRecorder is the mock recorder for MockAuthServiceInterface.
type MockAuthServiceInterfaceMockRecorder struct {
	mock *MockAuthServiceInterface
}

// NewMockAuthServiceInterface creates a new mock instance.
func NewMockAuthServiceInterface(ctrl *gomock.Controller) *MockAuthServiceInterface {
	mock := &MockAuthServiceInterface{ctrl: ctrl}
	mock.recorder = &MockAuthServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServiceInterface) EXPECT() *MockAuthServiceInterfaceMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuthServiceInterface) Login(username, password string) (*models.AuthResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", username, password)
	ret0, _ := ret[0].(*models.AuthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthServiceInterfaceMockRecorder) Login(username, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthServiceInterface)(nil).Login), username, password)
}

// Logout mocks base method.
func (m *MockAuthServiceInterface) Logout() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout")
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockAuthServiceInterfaceMockRecorder) Logout() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthServiceInterface)(nil).Logout))
}

// Register mocks base method.
func (m *MockAuthServiceInterface) Register(username, password, email string) (*models.RegisterResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", username, password, email)
	ret0, _ := ret[0].(*models.RegisterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockAuthServiceInterfaceMockRecorder) Register(username, password, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthServiceInterface)(nil).Register), username, password, email)
}

// ValidateToken mocks base method.
func (m *MockAuthServiceInterface) ValidateToken() (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockAuthServiceInterfaceMockRecorder) ValidateToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockAuthServiceInterface)(nil).ValidateToken))
}

// MockDataServiceInterface is a mock of DataServiceInterface interface.
type MockDataServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDataServiceInterfaceMockRecorder
	isgomock struct{}
}

// MockDataServiceInterfaceMockRecorder is the mock recorder for MockDataServiceInterface.
type MockDataServiceInterfaceMockRecorder struct {
	mock *MockDataServiceInterface
}

// NewMockDataServiceInterface creates a new mock instance.
func NewMockDataServiceInterface(ctrl *gomock.Controller) *MockDataServiceInterface {
	mock := &MockDataServiceInterface{ctrl: ctrl}
	mock.recorder = &MockDataServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataServiceInterface) EXPECT() *MockDataServiceInterfaceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockDataServiceInterface) Add(entry *models.DataEntry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", entry)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockDataServiceInterfaceMockRecorder) Add(entry any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockDataServiceInterface)(nil).Add), entry)
}

// Delete mocks base method.
func (m *MockDataServiceInterface) Delete(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDataServiceInterfaceMockRecorder) Delete(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDataServiceInterface)(nil).Delete), id)
}

// Get mocks base method.
func (m *MockDataServiceInterface) Get(id string) (*models.DataEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*models.DataEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDataServiceInterfaceMockRecorder) Get(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDataServiceInterface)(nil).Get), id)
}

// List mocks base method.
func (m *MockDataServiceInterface) List() ([]*models.DataEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]*models.DataEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockDataServiceInterfaceMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDataServiceInterface)(nil).List))
}

// SyncWithServer mocks base method.
func (m *MockDataServiceInterface) SyncWithServer() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncWithServer")
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncWithServer indicates an expected call of SyncWithServer.
func (mr *MockDataServiceInterfaceMockRecorder) SyncWithServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncWithServer", reflect.TypeOf((*MockDataServiceInterface)(nil).SyncWithServer))
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
