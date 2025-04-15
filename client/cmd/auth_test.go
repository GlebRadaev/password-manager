package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/GlebRadaev/password-manager/client/models"
)

func TestRegisterCmd_Success(t *testing.T) {
	originalAuthService := authService
	defer func() { authService = originalAuthService }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := NewMockAuthServiceInterface(ctrl)
	mockAuthService.EXPECT().
		Register("testuser", "testpass", "test@example.com").
		Return(&models.RegisterResponse{}, nil)

	authService = mockAuthService

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	cmd := &cobra.Command{Use: "pm"}
	cmd.AddCommand(registerCmd)
	cmd.SetArgs([]string{"register", "--username", "testuser", "--password", "testpass", "--email", "test@example.com"})
	err := cmd.Execute()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.NoError(t, err)
	assert.Contains(t, output, "Registered user successfully")
}

func TestLoginCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := NewMockAuthServiceInterface(ctrl)
	mockAuthService.EXPECT().
		Login("testuser", "testpass").
		Return(&models.AuthResponse{}, nil)

	originalAuthService := authService
	authService = mockAuthService
	defer func() { authService = originalAuthService }()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "pm"}
	cmd.AddCommand(loginCmd)
	cmd.SetArgs([]string{"login", "--username", "testuser", "--password", "testpass"})
	err := cmd.Execute()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Login successful")
}

func TestLogoutCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := NewMockAuthServiceInterface(ctrl)
	mockAuthService.EXPECT().
		Logout().
		Return(nil)

	originalAuthService := authService
	authService = mockAuthService
	defer func() { authService = originalAuthService }()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "pm"}
	cmd.AddCommand(logoutCmd)
	cmd.SetArgs([]string{"logout"})
	err := cmd.Execute()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Logged out successfully")
}

func TestStatusCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := NewMockAuthServiceInterface(ctrl)
	mockAuthService.EXPECT().
		ValidateToken().
		Return(true, "user123", nil)

	originalAuthService := authService
	authService = mockAuthService
	defer func() { authService = originalAuthService }()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "pm"}
	cmd.AddCommand(statusCmd)
	cmd.SetArgs([]string{"status"})
	err := cmd.Execute()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Authenticated as user ID: user123")
}
