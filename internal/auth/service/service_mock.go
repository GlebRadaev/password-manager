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
	time "time"

	models "github.com/GlebRadaev/password-manager/internal/auth/models"
	pg "github.com/GlebRadaev/password-manager/internal/common/pg"
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

// CheckExists mocks base method.
func (m *MockRepo) CheckExists(ctx context.Context, username, email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckExists", ctx, username, email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckExists indicates an expected call of CheckExists.
func (mr *MockRepoMockRecorder) CheckExists(ctx, username, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckExists", reflect.TypeOf((*MockRepo)(nil).CheckExists), ctx, username, email)
}

// CreateOTP mocks base method.
func (m *MockRepo) CreateOTP(ctx context.Context, otp models.OTP) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOTP", ctx, otp)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOTP indicates an expected call of CreateOTP.
func (mr *MockRepoMockRecorder) CreateOTP(ctx, otp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOTP", reflect.TypeOf((*MockRepo)(nil).CreateOTP), ctx, otp)
}

// CreateUser mocks base method.
func (m *MockRepo) CreateUser(ctx context.Context, user models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepoMockRecorder) CreateUser(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepo)(nil).CreateUser), ctx, user)
}

// GetOTP mocks base method.
func (m *MockRepo) GetOTP(ctx context.Context, userID, otpCode string) (models.OTP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOTP", ctx, userID, otpCode)
	ret0, _ := ret[0].(models.OTP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOTP indicates an expected call of GetOTP.
func (mr *MockRepoMockRecorder) GetOTP(ctx, userID, otpCode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOTP", reflect.TypeOf((*MockRepo)(nil).GetOTP), ctx, userID, otpCode)
}

// GetUserByUsername mocks base method.
func (m *MockRepo) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", ctx, username)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockRepoMockRecorder) GetUserByUsername(ctx, username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockRepo)(nil).GetUserByUsername), ctx, username)
}

// MockTxManager is a mock of TxManager interface.
type MockTxManager struct {
	ctrl     *gomock.Controller
	recorder *MockTxManagerMockRecorder
	isgomock struct{}
}

// MockTxManagerMockRecorder is the mock recorder for MockTxManager.
type MockTxManagerMockRecorder struct {
	mock *MockTxManager
}

// NewMockTxManager creates a new mock instance.
func NewMockTxManager(ctrl *gomock.Controller) *MockTxManager {
	mock := &MockTxManager{ctrl: ctrl}
	mock.recorder = &MockTxManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxManager) EXPECT() *MockTxManagerMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockTxManager) Begin(ctx context.Context, fn pg.TransactionalFn) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Begin indicates an expected call of Begin.
func (mr *MockTxManagerMockRecorder) Begin(ctx, fn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockTxManager)(nil).Begin), ctx, fn)
}

// MockAuthManager is a mock of AuthManager interface.
type MockAuthManager struct {
	ctrl     *gomock.Controller
	recorder *MockAuthManagerMockRecorder
	isgomock struct{}
}

// MockAuthManagerMockRecorder is the mock recorder for MockAuthManager.
type MockAuthManagerMockRecorder struct {
	mock *MockAuthManager
}

// NewMockAuthManager creates a new mock instance.
func NewMockAuthManager(ctrl *gomock.Controller) *MockAuthManager {
	mock := &MockAuthManager{ctrl: ctrl}
	mock.recorder = &MockAuthManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthManager) EXPECT() *MockAuthManagerMockRecorder {
	return m.recorder
}

// Compare mocks base method.
func (m *MockAuthManager) Compare(hashedPassword, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compare", hashedPassword, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Compare indicates an expected call of Compare.
func (mr *MockAuthManagerMockRecorder) Compare(hashedPassword, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compare", reflect.TypeOf((*MockAuthManager)(nil).Compare), hashedPassword, password)
}

// GenerateOTP mocks base method.
func (m *MockAuthManager) GenerateOTP() (string, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateOTP")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateOTP indicates an expected call of GenerateOTP.
func (mr *MockAuthManagerMockRecorder) GenerateOTP() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateOTP", reflect.TypeOf((*MockAuthManager)(nil).GenerateOTP))
}

// GenerateToken mocks base method.
func (m *MockAuthManager) GenerateToken(userID string) (string, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockAuthManagerMockRecorder) GenerateToken(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockAuthManager)(nil).GenerateToken), userID)
}

// Hash mocks base method.
func (m *MockAuthManager) Hash(password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hash", password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hash indicates an expected call of Hash.
func (mr *MockAuthManagerMockRecorder) Hash(password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hash", reflect.TypeOf((*MockAuthManager)(nil).Hash), password)
}

// ValidateOTP mocks base method.
func (m *MockAuthManager) ValidateOTP(storedOTPCode string, storedExpiresAt time.Time, providedOTPCode string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateOTP", storedOTPCode, storedExpiresAt, providedOTPCode)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateOTP indicates an expected call of ValidateOTP.
func (mr *MockAuthManagerMockRecorder) ValidateOTP(storedOTPCode, storedExpiresAt, providedOTPCode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateOTP", reflect.TypeOf((*MockAuthManager)(nil).ValidateOTP), storedOTPCode, storedExpiresAt, providedOTPCode)
}

// ValidateToken mocks base method.
func (m *MockAuthManager) ValidateToken(token string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockAuthManagerMockRecorder) ValidateToken(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockAuthManager)(nil).ValidateToken), token)
}
