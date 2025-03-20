package models

import "time"

type DataEntry struct {
	DataID    string
	Type      DataType
	Data      []byte
	Metadata  []Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}
