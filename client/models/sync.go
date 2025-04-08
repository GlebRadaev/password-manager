// Package models defines core domain types and data structures for synchronization
package models

// SyncResponse represents the response from a synchronization operation
type SyncResponse struct {
	Success   bool       // Indicates if sync was successful
	Conflicts []Conflict // List of conflicts found during sync
}

// Conflict represents a data synchronization conflict between client and server
type Conflict struct {
	ConflictID string `json:"conflict_id"` // Unique identifier for the conflict
	DataID     string `json:"data_id"`     // ID of the conflicting data entry
}

// ResolutionResponse represents the result of a conflict resolution attempt
type ResolutionResponse struct {
	Success bool   `json:"success"` // Whether resolution was successful
	Message string `json:"message"` // Additional status message
}

// ResolutionStrategy defines approaches for resolving data conflicts
type ResolutionStrategy string

const (
	// UseClient indicates using client's version in conflict resolution
	UseClient ResolutionStrategy = "client"
	// UseServer indicates using server's version in conflict resolution
	UseServer ResolutionStrategy = "server"
	// Merge indicates merging both versions in conflict resolution
	Merge ResolutionStrategy = "merge"
)
