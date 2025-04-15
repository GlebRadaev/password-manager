package models

import (
	"testing"
)

func TestDataType(t *testing.T) {
	tests := []struct {
		name     string
		input    DataType
		expected string
	}{
		{"Login", Login, "login"},
		{"Note", Note, "note"},
		{"Card", Card, "card"},
		{"Binary", Binary, "binary"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.expected && tt.input != -1 {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDataTypeFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected DataType
	}{
		{"Login lowercase", "login", Login},
		{"Login uppercase", "LOGIN", Login},
		{"Note mixed case", "NoTe", Note},
		{"Card", "card", Card},
		{"Binary", "binary", Binary},
		{"Unknown type", "unknown", -1},
		{"Empty string", "", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DataTypeFromString(tt.input); got != tt.expected {
				t.Errorf("DataTypeFromString() = %v, want %v", got, tt.expected)
			}
		})
	}
}
