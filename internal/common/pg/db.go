// Package pg provides PostgreSQL database abstractions
package pg

// db wraps a Database implementation
type db struct {
	Database
}

// New creates a new Database wrapper
func New(dbase Database) Database {
	return &db{
		Database: dbase,
	}
}
