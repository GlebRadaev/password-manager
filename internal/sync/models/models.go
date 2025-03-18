package models

import (
	"time"
)

type DataChange struct {
	ID        string
	UserID    string
	DataID    string
	Type      string
	Data      []byte
	Metadata  map[string]string
	Timestamp time.Time
}
