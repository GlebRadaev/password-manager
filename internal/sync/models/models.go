// Package models defines domain models and types for the application
package models

import (
	"time"
)

// Metadata represents key-value pairs for additional data attributes
type Metadata struct {
	Key   string // Metadata key
	Value string // Metadata value
}

// DataType enumerates different types of stored data
type DataType int

const (
	// LoginPassword represents login/password credentials
	LoginPassword DataType = iota
	// Text represents plain text data
	Text
	// Binary represents binary data
	Binary
	// Card represents payment card information
	Card
)

// Operation defines possible data operations
type Operation int

const (
	// Add represents an add operation
	Add Operation = iota
	// Update represents an update operation
	Update
	// Delete represents a delete operation
	Delete
)

// ResolutionStrategy defines conflict resolution approaches
type ResolutionStrategy int

const (
	// UseClientVersion indicates using client's version in conflict resolution
	UseClientVersion ResolutionStrategy = iota
	// UseServerVersion indicates using server's version in conflict resolution
	UseServerVersion
	// MergeVersions indicates merging both versions in conflict resolution
	MergeVersions
)

// Conflict represents a data synchronization conflict
type Conflict struct {
	ID         string    // Unique conflict ID
	UserID     string    // Associated user ID
	DataID     string    // Related data entry ID
	ClientData []byte    // Client's version of data
	ServerData []byte    // Server's version of data
	Resolved   bool      // Resolution status
	CreatedAt  time.Time // Creation timestamp
	UpdatedAt  time.Time // Last update timestamp
}

// ClientData represents data sent from client to server
type ClientData struct {
	DataID    string     // Unique data identifier
	Type      DataType   // Data type
	Data      []byte     // Actual data content
	UpdatedAt time.Time  // Last update timestamp
	Metadata  []Metadata // Additional metadata
	Operation Operation  // Operation type (add/update/delete)
}
