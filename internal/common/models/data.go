package models

import (
	"time"
)

type DataType int

const (
	LoginPassword DataType = iota
	Text
	Binary
	Card
)

type Metadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DataEntry struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Type      DataType   `json:"type"`
	Data      []byte     `json:"data"`
	Metadata  []Metadata `json:"metadata"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
