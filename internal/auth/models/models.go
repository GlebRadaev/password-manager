// Package models defines core domain entities for the authentication system
package models

import "time"

// User represents an application user account
type User struct {
	ID           string    // Unique user identifier
	Username     string    // Login username
	PasswordHash string    // Hashed password
	Email        string    // User email address
	CreatedAt    time.Time // Account creation timestamp
}

// OTP represents a one-time password for multi-factor authentication
type OTP struct {
	UserID    string    // Associated user ID
	OTPCode   string    // One-time password value
	ExpiresAt time.Time // Expiration timestamp
	DeviceID  string    // Device identifier
}

// Session represents an active user authentication session
type Session struct {
	SessionID  string    // Unique session identifier
	UserID     string    // Associated user ID
	DeviceInfo string    // Device/browser information
	CreatedAt  time.Time // Session creation timestamp
	ExpiresAt  time.Time // Session expiration timestamp
}
