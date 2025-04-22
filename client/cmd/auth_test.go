package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/GlebRadaev/password-manager/client/models"
)

func TestRegisterCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := NewMockAuthServiceInterface(ctrl)
	mockAuthService.EXPECT().
		Register("testuser", "testpass", "test@example.com").
		Return(&models.RegisterResponse{}, nil)

	originalAuthService := authService
	authService = mockAuthService
	defer func() { authService = originalAuthService }()

	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "Register new user",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			email, _ := cmd.Flags().GetString("email")

			_, err := authService.Register(username, password, email)
			if err != nil {
				cmd.PrintErrln("Registration failed:", err)
				return
			}
			cmd.Println("Registered user successfully")
		},
	}
	registerCmd.Flags().StringP("username", "u", "", "Username")
	registerCmd.Flags().StringP("password", "p", "", "Password")
	registerCmd.Flags().StringP("email", "e", "", "Email")

	buf := new(bytes.Buffer)
	registerCmd.SetOut(buf)
	registerCmd.SetErr(buf)

	registerCmd.SetArgs([]string{"--username", "testuser", "--password", "testpass", "--email", "test@example.com"})
	err := registerCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Registered user successfully")
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

	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to system",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")

			_, err := authService.Login(username, password)
			if err != nil {
				cmd.PrintErrln("Login failed:", err)
				return
			}
			cmd.Println("Login successful")
		},
	}
	loginCmd.Flags().StringP("username", "u", "", "Username")
	loginCmd.Flags().StringP("password", "p", "", "Password")

	buf := new(bytes.Buffer)
	loginCmd.SetOut(buf)
	loginCmd.SetErr(buf)

	loginCmd.SetArgs([]string{"--username", "testuser", "--password", "testpass"})
	err := loginCmd.Execute()

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

	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from system",
		Run: func(cmd *cobra.Command, args []string) {
			if err := authService.Logout(); err != nil {
				cmd.PrintErrln("Logout failed:", err)
				return
			}
			cmd.Println("Logged out successfully")
		},
	}

	buf := new(bytes.Buffer)
	logoutCmd.SetOut(buf)
	logoutCmd.SetErr(buf)

	logoutCmd.SetArgs([]string{})
	err := logoutCmd.Execute()

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

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show auth status",
		Run: func(cmd *cobra.Command, args []string) {
			valid, userID, err := authService.ValidateToken()
			if err != nil {
				cmd.PrintErrln("Status check failed:", err)
				return
			}
			if valid {
				cmd.Printf("Authenticated as user ID: %s\n", userID)
			} else {
				cmd.Println("Not authenticated")
			}
		},
	}

	buf := new(bytes.Buffer)
	statusCmd.SetOut(buf)
	statusCmd.SetErr(buf)

	statusCmd.SetArgs([]string{})
	err := statusCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Authenticated as user ID: user123")
}
