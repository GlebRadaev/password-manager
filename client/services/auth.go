package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/GlebRadaev/password-manager/client/models"
)

const (
	tokenFileMode = 0600 // -rw-------
	tokenFileName = ".pm_token"
)

// AuthService provides authentication operations for the password manager client.
// It handles user registration, login, logout, and token validation.
type AuthService struct {
	baseURL   string
	client    HTTPClientInterface
	tokenPath func() (string, error)
}

// NewAuthService creates a new AuthService instance with default base URL.
func NewAuthService() *AuthService {
	return &AuthService{
		baseURL: "http://localhost:8079",
		client:  &http.Client{},
		tokenPath: func() (string, error) {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get home directory: %w", err)
			}
			return filepath.Join(home, tokenFileName), nil
		},
	}
}

// Register creates a new user account with the provided credentials.
// Returns AuthResponse containing access token on success.
func (s *AuthService) Register(username, password, email string) (*models.RegisterResponse, error) {
	url := fmt.Sprintf("%s/v1/auth/register", s.baseURL)

	reqBody := map[string]string{
		"username": username,
		"password": password,
		"email":    email,
	}

	resp, err := s.doRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	var result models.RegisterResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Login authenticates a user and stores the access token locally.
// Returns AuthResponse containing access token on success.
func (s *AuthService) Login(username, password string) (*models.AuthResponse, error) {
	url := fmt.Sprintf("%s/v1/auth/login", s.baseURL)

	reqBody := map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := s.doRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	var result models.AuthResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.AccessToken == "" {
		return nil, fmt.Errorf("server returned empty token")
	}

	if err := s.saveToken(result.AccessToken); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	return &result, nil
}

// Logout removes the locally stored access token.
func (s *AuthService) Logout() error {
	return s.clearToken()
}

// ValidateToken checks if the stored token is valid.
// Returns validation status and associated user ID if valid.
func (s *AuthService) ValidateToken() (bool, string, error) {
	token, err := s.loadToken()
	if err != nil {
		if os.IsNotExist(err) {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to load token: %w", err)
	}
	if token == "" {
		return false, "", nil
	}

	url := fmt.Sprintf("%s/v1/auth/validate-token", s.baseURL)
	reqBody := map[string]string{"token": token}

	resp, err := s.doRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return false, "", fmt.Errorf("validation request failed: %w", err)
	}

	var result struct {
		Valid  bool   `json:"valid"`
		UserID string `json:"user_id"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return false, "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Valid, result.UserID, nil
}

// saveToken stores the access token in the user's home directory.
func (s *AuthService) saveToken(token string) error {
	path, err := s.tokenPath()
	if err != nil {
		return fmt.Errorf("failed to get token path: %w", err)
	}
	return os.WriteFile(path, []byte(token), tokenFileMode)
}

// loadToken retrieves the stored access token from the user's home directory.
func (s *AuthService) loadToken() (string, error) {
	path, err := s.tokenPath()
	if err != nil {
		return "", fmt.Errorf("failed to get token path: %w", err)
	}

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// clearToken removes the stored access token file.
func (s *AuthService) clearToken() error {
	path, err := s.tokenPath()
	if err != nil {
		return fmt.Errorf("failed to get token path: %w", err)
	}
	return os.Remove(path)
}

// doRequest performs an HTTP request with JSON body and handles the response.
func (s *AuthService) doRequest(method, url string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

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

	return respBody, nil
}
