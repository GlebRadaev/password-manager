// Package state provides thread-safe global state management for authentication
package state

import (
	"sync"
)

// authState holds the authentication state with thread-safe access
var (
	authState struct {
		accessToken string       // Current access token
		mu          sync.RWMutex // Mutex for concurrent access protection
	}
)

// SetAccessToken updates the current access token
// Thread-safe operation using mutex lock
func SetAccessToken(token string) {
	authState.mu.Lock()
	defer authState.mu.Unlock()
	authState.accessToken = token
}

// GetAccessToken retrieves the current access token
// Thread-safe operation using read lock
func GetAccessToken() string {
	authState.mu.RLock()
	defer authState.mu.RUnlock()
	return authState.accessToken
}

// IsAuthenticated checks if there's an active access token
// Returns true if access token is not empty
func IsAuthenticated() bool {
	return GetAccessToken() != ""
}
