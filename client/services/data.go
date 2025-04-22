// Package services provides business logic and API integration for data operations
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/GlebRadaev/password-manager/client/models"
	"github.com/GlebRadaev/password-manager/client/storage"
)

const (
	dataSyncPath = "/v1/data/sync"
)

// DataService handles data operations with local storage and remote server
type DataService struct {
	storage StorageInterface
	baseURL string
	client  HTTPClientInterface
}

// NewDataService creates a new DataService instance with default configuration
func NewDataService() *DataService {
	return &DataService{
		storage: storage.NewLocalStorage(),
		baseURL: "http://localhost:8079",
		client: &http.Client{
			Timeout: timeOut,
		},
	}
}

// Add saves a new data entry to local storage
// Returns error if local save fails
func (s *DataService) Add(entry *models.DataEntry) error {
	if err := s.storage.Add(entry); err != nil {
		return fmt.Errorf("failed to save locally: %w", err)
	}
	return nil
}

// List retrieves all data entries from local storage
// Returns slice of entries or error if retrieval fails
func (s *DataService) List() ([]*models.DataEntry, error) {
	entries, err := s.storage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get local data: %w", err)
	}
	return entries, nil
}

// Get retrieves a single data entry by ID from local storage
// Returns entry or error if not found
func (s *DataService) Get(id string) (*models.DataEntry, error) {
	entry, err := s.storage.Get(id)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: entry not found", err)
		}
		return nil, fmt.Errorf("failed to get entry %s: %w", id, err)
	}
	return entry, nil
}

// Delete removes a data entry by ID from local storage
// Returns error if deletion fails
func (s *DataService) Delete(id string) error {
	if err := s.storage.Delete(id); err != nil {
		return fmt.Errorf("failed to delete entry %s: %w", id, err)
	}
	return nil
}

// SyncWithServer synchronizes local data with remote server
// Returns error if sync operation fails
func (s *DataService) SyncWithServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	entries, err := s.storage.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get entries for sync: %w", err)
	}

	url := s.baseURL + dataSyncPath
	reqBody := map[string]interface{}{
		"entries": entries,
	}

	resp, err := s.doAuthenticatedRequest(ctx, "POST", url, reqBody)
	if err != nil {
		return fmt.Errorf("sync request failed: %w", err)
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to decode sync response: %w", err)
	}

	if !result.Success {
		return errors.New("sync failed: server returned unsuccessful status")
	}

	return nil
}

// doAuthenticatedRequest performs authenticated HTTP requests with context
// Handles token management, request building and response parsing
// Returns response body or error if request fails
func (s *DataService) doAuthenticatedRequest(ctx context.Context, method, url string, body interface{}) ([]byte, error) {
	token, err := s.storage.GetAuthToken()
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Message != "" {
			return nil, fmt.Errorf("server error: %s (status %d)", errResp.Message, resp.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return respBody, nil
}
