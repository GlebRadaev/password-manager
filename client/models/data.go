// Package models defines core domain types and data structures
// for the password manager application.
package models

import "strings"

// DataType represents different categories of stored data
type DataType int

const (
	// Login represents authentication credentials (username/password)
	Login DataType = iota

	// Note represents text notes and secure messages
	Note

	// Card represents payment card information
	Card

	// Binary represents binary data (files, images)
	Binary
)

// String returns the string representation of DataType
func (t DataType) String() string {
	return [...]string{"login", "note", "card", "binary"}[t]
}

// DataTypeFromString converts string to DataType
// Returns -1 for unknown types
func DataTypeFromString(s string) DataType {
	switch strings.ToLower(s) {
	case "login":
		return Login
	case "note":
		return Note
	case "card":
		return Card
	case "binary":
		return Binary
	default:
		return -1
	}
}

// DataEntry represents a single encrypted data record
type DataEntry struct {
	ID        string   // Unique identifier
	Type      DataType // Data category
	Data      []byte   // Encrypted content
	CreatedAt int64    // Creation timestamp (unix)
	UpdatedAt int64    // Last update timestamp (unix)
}

// AuthResponse contains authentication tokens from the server
type AuthResponse struct {
	AccessToken  string // Bearer token for API access
	RefreshToken string // Token for obtaining new access tokens
	ExpiresIn    string // Token expiration duration
}

// ClientData represents data structure for client-server communication
type ClientData struct {
	ID        string `json:"id"`         // Entry ID
	Type      string `json:"type"`       // Data type as string
	Data      []byte `json:"data"`       // Encrypted data
	UpdatedAt int64  `json:"updated_at"` // Last update timestamp
}
