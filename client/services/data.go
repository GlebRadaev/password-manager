// Package services provides business logic and API integration for data operations
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/GlebRadaev/password-manager/client/models"
	"github.com/GlebRadaev/password-manager/client/storage"
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
		client:  &http.Client{},
	}
}

// Add saves a new data entry to local storage
// Returns error if local save fails
func (s *DataService) Add(entry *models.DataEntry) error {
	if err := s.storage.Add(entry); err != nil {
		return fmt.Errorf("failed to save locally: %v", err)
	}
	return nil
}

// List retrieves all data entries from local storage
// Returns slice of entries or error if retrieval fails
func (s *DataService) List() ([]*models.DataEntry, error) {
	entries, err := s.storage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get local data: %v", err)
	}
	return entries, nil
}

// Get retrieves a single data entry by ID from local storage
// Returns entry or error if not found
func (s *DataService) Get(id string) (*models.DataEntry, error) {
	return s.storage.Get(id)
}

// Delete removes a data entry by ID from local storage
// Returns error if deletion fails
func (s *DataService) Delete(id string) error {
	return s.storage.Delete(id)
}

// SyncWithServer synchronizes local data with remote server
// Returns error if sync operation fails
func (s *DataService) SyncWithServer() error {
	entries, err := s.storage.GetAll()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/v1/data/sync", s.baseURL)
	reqBody := map[string]interface{}{
		"entries": entries,
	}

	resp, err := s.doAuthenticatedRequest("POST", url, reqBody)
	if err != nil {
		return err
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !result.Success {
		return fmt.Errorf("sync failed")
	}

	return nil
}

// doAuthenticatedRequest performs authenticated HTTP requests
// Handles token management, request building and response parsing
// Returns response body or error if request fails
func (s *DataService) doAuthenticatedRequest(method, url string, body interface{}) ([]byte, error) {
	token, err := s.storage.GetAuthToken()
	if err != nil {
		return nil, fmt.Errorf("authentication required: %v", err)
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(respBody, &errResp)
		if errResp.Message != "" {
			return nil, fmt.Errorf("server error: %s", errResp.Message)
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return respBody, nil
}
