package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

// PasswordReset defines the password reset model according to the
// `password_reset` SQL table. Although this struct should map to the SQL
// table, some entires (such as `Token`) are not present in the DB
type PasswordReset struct {
	ID        int
	UserID    int
	Token     string // Token is only set when creating a new session
	TokenHash string
	ExpiresAt time.Time
}

// PasswordResetService defines the connection to the DB
type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken determines how many bytes used to generate each password
	// reset token. If `BytesPerToken` is not provided or is less than
	// `MinBytesPerToken`, then `MinBytesPerToken` is used instead
	BytesPerToken int
	// Duration is the amount of time during which a PasswordReset is valid.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

// Create creates a new `PasswordReset` type
func (prs *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Create")
}

// Consume takes an existing password reset token and uses it
func (psr *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Consume")
}
