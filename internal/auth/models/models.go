package models

import (
	"time"
)

type User struct {
	ID           string
	Username     string
	PasswordHash string
	Email        string
	CreatedAt    time.Time
}

type OTP struct {
	UserID    string
	OTPCode   string
	ExpiresAt time.Time
}
