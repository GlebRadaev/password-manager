// Package storage provides local file-based storage implementation for password manager data.
// It handles secure storage of data entries and synchronization status.
package storage

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/GlebRadaev/password-manager/client/models"
)

// LocalStorage implements thread-safe local file storage for data entries
type LocalStorage struct {
	path         string     // Base directory path for data storage
	syncFilePath string     // Path to sync status tracking file
	mu           sync.Mutex // Mutex for concurrent access protection
}

// NewLocalStorage creates new LocalStorage instance with default paths (/tmp/.pm_data)
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		path:         "/tmp/.pm_data",
		syncFilePath: "/tmp/.pm_data/.sync_status",
	}
}

// Add saves a new data entry to local storage as JSON file
// Returns error if file operations fail
func (s *LocalStorage) Add(entry *models.DataEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fullPath := filepath.Join(s.path, entry.ID)
	log.Printf("Saving to: %s", fullPath)

	if err := os.MkdirAll(s.path, 0700); err != nil {
		return err
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(s.path, entry.ID), data, 0600)
}

// Get retrieves single data entry by ID
// Returns entry or ErrNotExist if not found
func (s *LocalStorage) Get(id string) (*models.DataEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("Looking for file: %s", filepath.Join(s.path, id))

	files, err := os.ReadDir(s.path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), id) {
			data, err := os.ReadFile(filepath.Join(s.path, file.Name()))
			if err != nil {
				return nil, err
			}
			var entry models.DataEntry
			if err := json.Unmarshal(data, &entry); err != nil {
				return nil, err
			}
			return &entry, nil
		}
	}
	return nil, os.ErrNotExist
}

// GetAll returns all stored data entries
// Returns empty slice if storage directory doesn't exist
func (s *LocalStorage) GetAll() ([]*models.DataEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	files, err := os.ReadDir(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.DataEntry{}, nil
		}
		return nil, err
	}

	var entries []*models.DataEntry
	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(s.path, file.Name()))
		if err != nil {
			continue
		}

		var entry models.DataEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		entries = append(entries, &entry)
	}

	return entries, nil
}

// Delete removes data entry by ID
// Returns error if file doesn't exist or can't be removed
func (s *LocalStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return os.Remove(filepath.Join(s.path, id))
}

// GetPendingSyncEntries returns entries that need synchronization
// Compares local update times with last sync status
func (s *LocalStorage) GetPendingSyncEntries() ([]*models.DataEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	allEntries, err := s.getAllEntries()
	if err != nil {
		return nil, err
	}

	syncStatus, err := s.loadSyncStatus()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var pendingEntries []*models.DataEntry
	for _, entry := range allEntries {
		lastSynced, exists := syncStatus[entry.ID]
		if !exists || entry.UpdatedAt > lastSynced {
			pendingEntries = append(pendingEntries, entry)
		}
	}

	return pendingEntries, nil
}

// ClearPendingSync resets synchronization status tracking
func (s *LocalStorage) ClearPendingSync() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return os.Remove(s.syncFilePath)
}

// UpdateSyncStatus updates last sync timestamps for given entries
func (s *LocalStorage) UpdateSyncStatus(entries []*models.DataEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	syncStatus, err := s.loadSyncStatus()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	for _, entry := range entries {
		syncStatus[entry.ID] = entry.UpdatedAt
	}

	return s.saveSyncStatus(syncStatus)
}

// loadSyncStatus reads sync status from tracking file
// Returns empty map if file doesn't exist
func (s *LocalStorage) loadSyncStatus() (map[string]int64, error) {
	data, err := os.ReadFile(s.syncFilePath)
	if err != nil {
		return make(map[string]int64), err
	}

	var status map[string]int64
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, err
	}
	return status, nil
}

// saveSyncStatus writes sync status to tracking file
func (s *LocalStorage) saveSyncStatus(status map[string]int64) error {
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(s.path, 0700); err != nil {
		return err
	}

	return os.WriteFile(s.syncFilePath, data, 0600)
}

// getAllEntries retrieves all entries from storage (internal helper)
func (s *LocalStorage) getAllEntries() ([]*models.DataEntry, error) {
	files, err := os.ReadDir(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.DataEntry{}, nil
		}
		return nil, err
	}

	var entries []*models.DataEntry
	for _, file := range files {
		if file.Name()[0] == '.' {
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.path, file.Name()))
		if err != nil {
			continue
		}

		var entry models.DataEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		entries = append(entries, &entry)
	}

	return entries, nil
}

// GetAuthToken retrieves stored authentication token from ~/.pm_token
// Returns empty string if token file doesn't exist
func (s *LocalStorage) GetAuthToken() (string, error) {
	home, _ := os.UserHomeDir()
	tokenPath := filepath.Join(home, ".pm_token")

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
