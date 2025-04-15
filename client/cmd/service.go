// Package cmd implements the command-line interface for the password manager.
package cmd

//go:generate mockgen -source=service.go -destination=service_mock.go -package=cmd
import (
	"github.com/GlebRadaev/password-manager/client/models"
)

// AuthServiceInterface defines the contract for authentication-related operations.
type AuthServiceInterface interface {
	// Register creates a new user account with the provided credentials.
	Register(username, password, email string) (*models.RegisterResponse, error)

	// Login authenticates a user with the provided credentials.
	Login(username, password string) (*models.AuthResponse, error)

	// Logout terminates the current user session.
	Logout() error

	// ValidateToken checks if the current session token is valid.
	ValidateToken() (bool, string, error)
}

// DataServiceInterface defines operations for managing password entries.
type DataServiceInterface interface {
	// Add creates a new data entry in the password manager.
	Add(entry *models.DataEntry) error

	// List retrieves all stored data entries.
	List() ([]*models.DataEntry, error)

	// Get retrieves a specific entry by its ID.
	Get(id string) (*models.DataEntry, error)

	// Delete removes an entry by its ID.
	Delete(id string) error

	// SyncWithServer synchronizes local data with the remote server.
	SyncWithServer() error
}

// SyncServiceInterface defines operations for data synchronization and conflict resolution.
type SyncServiceInterface interface {
	// Sync performs a two-way synchronization between client and server.
	Sync() (*models.SyncResponse, error)

	// Resolve handles conflict resolution using the specified strategy.
	Resolve(conflictID, strategy string) (*models.ResolutionResponse, error)
}
