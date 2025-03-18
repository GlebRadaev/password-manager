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
	DeviceID  string
}

type Session struct {
	SessionID  string
	UserID     string
	DeviceInfo string
	CreatedAt  time.Time
	ExpiresAt  time.Time
}
