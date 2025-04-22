// Package services provides synchronization services between client and server.
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/GlebRadaev/password-manager/client/models"
	"github.com/GlebRadaev/password-manager/client/storage"
)

const timeOut = 30 * time.Second

// SyncService handles synchronization between local storage and remote server.
// It manages data sync operations and conflict resolution.
type SyncService struct {
	baseURL string
	storage StorageInterface
	client  HTTPClientInterface
}

// NewSyncService creates a new SyncService instance with default configuration.
func NewSyncService() *SyncService {
	return &SyncService{
		baseURL: "http://localhost:8079",
		storage: storage.NewLocalStorage(),
		client:  &http.Client{Timeout: timeOut},
	}
}

// Sync performs synchronization of pending entries with the remote server.
// Returns SyncResponse containing sync results and any conflicts.
func (s *SyncService) Sync() (*models.SyncResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	token, err := s.storage.GetAuthToken()
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	userID, err := s.validateTokenAndGetUserID(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	entries, err := s.storage.GetPendingSyncEntries()
	if err != nil {
		return nil, fmt.Errorf("failed to get pending sync entries: %w", err)
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

	url := s.baseURL + "/v1/sync/data"
	reqBody := map[string]interface{}{
		"user_id":     userID,
		"client_data": clientData,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
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
		_ = json.Unmarshal(respBody, &errResp)
		if errResp.Message != "" {
			return nil, fmt.Errorf("server error: %s", errResp.Message)
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result models.SyncResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Conflicts) == 0 {
		if err := s.storage.UpdateSyncStatus(entries); err != nil {
			return nil, fmt.Errorf("failed to update sync status: %w", err)
		}
	}

	return &result, nil
}

// Resolve handles conflict resolution using specified strategy.
// Returns ResolutionResponse with resolution results.
func (s *SyncService) Resolve(conflictID, strategy string) (*models.ResolutionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	token, err := s.storage.GetAuthToken()
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	url := s.baseURL + "/v1/sync/resolve"
	reqBody := map[string]string{
		"conflict_id": conflictID,
		"strategy":    strategy,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
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
		_ = json.Unmarshal(respBody, &errResp)
		if errResp.Message != "" {
			return nil, fmt.Errorf("server error: %s", errResp.Message)
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result models.ResolutionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// validateTokenAndGetUserID checks token validity and retrieves associated user ID.
// Returns user ID if token is valid, error otherwise.
func (s *SyncService) validateTokenAndGetUserID(ctx context.Context, token string) (string, error) {
	url := s.baseURL + "/v1/auth/validate-token"
	reqBody := map[string]string{"token": token}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("validation request failed: %w", err)
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
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Valid {
		return "", fmt.Errorf("invalid token")
	}

	if result.UserID == "" {
		return "", fmt.Errorf("server returned empty user_id")
	}

	return result.UserID, nil
}
