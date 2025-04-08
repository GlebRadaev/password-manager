// Package cmd implements the command-line interface for the password manager
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/GlebRadaev/password-manager/client/services"
)

// authService is the shared authentication service instance
var authService = services.NewAuthService()

// registerCmd handles user registration
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register new user",
	Long:  "Creates a new user account with the provided credentials",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		email, _ := cmd.Flags().GetString("email")

		_, err := authService.Register(username, password, email)
		if err != nil {
			log.Fatalf("Registration failed: %v", err)
		}
		fmt.Println("Registered user successfully")
	},
}

// loginCmd handles user authentication
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to system",
	Long:  "Authenticates user with provided credentials and stores session token",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		_, err := authService.Login(username, password)
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		}
		fmt.Println("Login successful")
	},
}

// logoutCmd terminates the current session
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from system",
	Long:  "Clears the current authentication session",
	Run: func(cmd *cobra.Command, args []string) {
		if err := authService.Logout(); err != nil {
			log.Fatalf("Logout failed: %v", err)
		}
		fmt.Println("Logged out successfully")
	},
}

// statusCmd displays current authentication status
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show auth status",
	Long:  "Displays information about current authentication state",
	Run: func(cmd *cobra.Command, args []string) {
		valid, userID, err := authService.ValidateToken()
		if err != nil {
			log.Fatalf("Status check failed: %v", err)
		}
		if valid {
			fmt.Printf("Authenticated as user ID: %s\n", userID)
		} else {
			fmt.Println("Not authenticated")
		}
	},
}

func init() {
	// Register commands with root command
	rootCmd.AddCommand(registerCmd, loginCmd, logoutCmd, statusCmd)

	// Register command flags
	registerCmd.Flags().StringP("username", "u", "", "Username for registration")
	registerCmd.Flags().StringP("password", "p", "", "Password for registration")
	registerCmd.Flags().StringP("email", "e", "", "Email for registration")
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")
	registerCmd.MarkFlagRequired("email")

	loginCmd.Flags().StringP("username", "u", "", "Username for login")
	loginCmd.Flags().StringP("password", "p", "", "Password for login")
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
}
