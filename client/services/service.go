// Package services provides the core business logic implementations for
package services

//go:generate mockgen -destination=service_mock.go -source=service.go -package=services
import (
	"net/http"

	"github.com/GlebRadaev/password-manager/client/models"
)

// StorageInterface defines the contract for persistent data storage operations.
type StorageInterface interface {
	// Add stores a new data entry in the local storage.
	Add(entry *models.DataEntry) error

	// Get retrieves a specific entry by its unique identifier.
	Get(id string) (*models.DataEntry, error)

	// GetAll retrieves all stored data entries from local storage.
	GetAll() ([]*models.DataEntry, error)

	// Delete removes an entry from local storage by its ID.
	Delete(id string) error

	// GetAuthToken retrieves the current authentication token.
	GetAuthToken() (string, error)

	// GetPendingSyncEntries retrieves all entries marked for synchronization.
	GetPendingSyncEntries() ([]*models.DataEntry, error)

	// UpdateSyncStatus updates the synchronization status of multiple entries.
	UpdateSyncStatus(entries []*models.DataEntry) error

	// ClearPendingSync resets the synchronization status of all entries.
	ClearPendingSync() error
}

// HTTPClientInterface abstracts HTTP client operations.
type HTTPClientInterface interface {
	// Do executes an HTTP request and returns the response.
	Do(req *http.Request) (*http.Response, error)
}

// SyncServiceInterface defines operations for client-server synchronization.
type SyncServiceInterface interface {
	// Sync performs bidirectional synchronization between client and server.
	Sync() (*models.SyncResponse, error)

	// Resolve handles a synchronization conflict using the specified strategy.
	Resolve(conflictID, strategy string) (*models.ResolutionResponse, error)
}
