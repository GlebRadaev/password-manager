package models

import (
	"time"
)

type Metadata struct {
	Key   string
	Value string
}

type DataType int

const (
	LoginPassword DataType = iota
	Text
	Binary
	Card
)

type Operation int

const (
	Add Operation = iota
	Update
	Delete
)

type ResolutionStrategy int

const (
	UseClientVersion ResolutionStrategy = iota
	UseServerVersion
	MergeVersions
)

type Conflict struct {
	ID         string
	UserID     string
	DataID     string
	ClientData []byte
	ServerData []byte
	Resolved   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
type ClientData struct {
	DataID    string
	Type      DataType
	Data      []byte
	UpdatedAt time.Time
	Metadata  []Metadata
	Operation Operation
}
