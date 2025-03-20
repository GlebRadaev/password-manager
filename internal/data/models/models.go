package models

import (
	"time"
)

type DataEntry struct {
	ID        string
	UserID    string
	Type      DataType
	Data      []byte
	Metadata  []Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
