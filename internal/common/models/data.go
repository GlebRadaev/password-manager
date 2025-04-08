// Package models defines core data structures for the application
package models

import (
	"time"
)

// DataType represents different categories of stored data
type DataType int

const (
	// LoginPassword credentials (username/password)
	LoginPassword DataType = iota
	// Text plain content
	Text
	// Binary data (files, images)
	Binary
	// Card information
	Card
)

// Metadata stores additional key-value attributes for data entries
// Uses JSON tags for API serialization
type Metadata struct {
	Key   string `json:"key"`   // Attribute name
	Value string `json:"value"` // Attribute value
}

// DataEntry represents a single encrypted data record
// Uses JSON tags for API serialization
type DataEntry struct {
	ID        string     `json:"id"`         // Unique identifier (UUID)
	UserID    string     `json:"user_id"`    // Owner reference
	Type      DataType   `json:"type"`       // Data category
	Data      []byte     `json:"data"`       // Encrypted content
	Metadata  []Metadata `json:"metadata"`   // Additional attributes
	CreatedAt time.Time  `json:"created_at"` // Creation timestamp
	UpdatedAt time.Time  `json:"updated_at"` // Last modification time
}
