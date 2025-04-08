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

// SyncService handles synchronization between local storage and remote server.
// It manages data sync operations and conflict resolution.
type SyncService struct {
	baseURL string
	storage *storage.LocalStorage
}

// NewSyncService creates a new SyncService instance with default configuration.
func NewSyncService() *SyncService {
	return &SyncService{
		baseURL: "http://localhost:8079",
		storage: storage.NewLocalStorage(),
	}
}

// Sync performs synchronization of pending entries with the remote server.
// Returns SyncResponse containing sync results and any conflicts.
func (s *SyncService) Sync() (*models.SyncResponse, error) {
	token, err := storage.GetAuthToken()
	if err != nil {
		return nil, fmt.Errorf("authentication required: %v", err)
	}

	userID, err := s.validateTokenAndGetUserID(token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %v", err)
	}

	entries, err := s.storage.GetPendingSyncEntries()
	if err != nil {
		return nil, fmt.Errorf("failed to get pending sync entries: %v", err)
	}

	var clientData []*models.ClientData
	for _, entry := range entries {
		clientData = append(clientData, &models.ClientData{
			ID:        entry.ID,
			Type:      entry.Type.String(),
			Data:      entry.Data,
			UpdatedAt: entry.UpdatedAt,
		})
	}

	url := fmt.Sprintf("%s/v1/sync/data", s.baseURL)
	reqBody := map[string]interface{}{
		"user_id":     userID,
		"client_data": clientData,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
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

	var result models.SyncResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result.Conflicts) == 0 {
		if err := s.storage.UpdateSyncStatus(entries); err != nil {
			return nil, fmt.Errorf("failed to update sync status: %v", err)
		}
	}

	return &result, nil
}

// Resolve handles conflict resolution using specified strategy.
// Returns ResolutionResponse with resolution results.
func (s *SyncService) Resolve(conflictID, strategy string) (*models.ResolutionResponse, error) {
	token, err := storage.GetAuthToken()
	if err != nil {
		return nil, fmt.Errorf("authentication required: %v", err)
	}

	url := fmt.Sprintf("%s/v1/sync/resolve", s.baseURL)
	reqBody := map[string]string{
		"conflict_id": conflictID,
		"strategy":    strategy,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
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

	var result models.ResolutionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}

// validateTokenAndGetUserID checks token validity and retrieves associated user ID.
// Returns user ID if token is valid, error otherwise.
func (s *SyncService) validateTokenAndGetUserID(token string) (string, error) {
	url := fmt.Sprintf("%s/v1/auth/validate-token", s.baseURL)
	reqBody := map[string]string{"token": token}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to encode request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("validation request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("invalid token status: %d, response: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Valid  bool   `json:"valid"`
		UserID string `json:"UserID"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if !result.Valid {
		return "", fmt.Errorf("invalid token")
	}

	if result.UserID == "" {
		return "", fmt.Errorf("server returned empty user_id")
	}

	return result.UserID, nil
}
