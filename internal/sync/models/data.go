// Package models defines the core data structures for the application
package models

import "time"

// DataEntry represents a single data record in the system
type DataEntry struct {
	DataID    string     // Unique identifier for the data entry
	Type      DataType   // Type of data (login, card, text, etc.)
	Data      []byte     // The encrypted data content
	Metadata  []Metadata // Additional key-value metadata
	CreatedAt time.Time  // When the entry was created
	UpdatedAt time.Time  // When the entry was last updated
}
