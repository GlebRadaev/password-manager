package models

import (
	"time"
)

type DataType int

const (
	LoginPassword DataType = iota
	Text
	Binary
	BankCard
)

type Data struct {
	ID        string
	UserID    string
	Type      DataType
	Data      []byte
	Metadata  map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
}
