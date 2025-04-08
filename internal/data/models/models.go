// Package models defines core domain types for data management
package models

import (
	"time"
)

// DataEntry represents a single data record stored in the system
type DataEntry struct {
	ID        string     // Unique identifier for the data entry
	UserID    string     // Owner user ID
	Type      DataType   // Type of data stored
	Data      []byte     // Actual data content
	Metadata  []Metadata // Additional metadata key-value pairs
	CreatedAt time.Time  // When the entry was created
	UpdatedAt time.Time  // Last modification time
}

// Metadata represents additional key-value attributes for data entries
type Metadata struct {
	Key   string // Metadata attribute name
	Value string // Metadata attribute value
}

// DataType categorizes different kinds of stored data
type DataType int

const (
	// LoginPassword credentials (username/password)
	LoginPassword DataType = iota
	// Text plain data
	Text
	// Binary data (files, images)
	Binary
	// Card information
	Card
)

// Operation defines possible actions on data entries
type Operation int

const (
	// Add new entry
	Add Operation = iota
	// Update existing entry
	Update
	// Delete entry
	Delete
)

// ResolutionStrategy defines approaches for resolving data conflicts
type ResolutionStrategy int

const (
	// UseClientVersion prefer client's version
	UseClientVersion ResolutionStrategy = iota
	// UseServerVersion prefer server's version
	UseServerVersion
	// MergeVersions both versions
	MergeVersions
)
