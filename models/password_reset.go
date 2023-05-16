package models

import (
	"database/sql"
	"fmt"
	"strings"
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
	// Normalize the email
	email = strings.ToLower(email)

	var userID int
	row := prs.DB.QueryRow(`
		SELECT id
		FROM users
		WHERE email = $1`, email)
	err := row.Scan(&userID)
	if err != nil {
		// TODO: consider returning a specific error when the user does not
		// exists
		return nil, fmt.Errorf("create: %w", err)
	}

	// Gets a token and token hash
	token, tokenHash, err := New(prs.BytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	// Handles the case where a duration is not provided
	duration := prs.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(duration),
	}

	// Creates a new session with the given value, or updates an existing
	// session (ON CONFLICT clause)
	row = prs.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &pwReset, nil
}

// Consume takes an existing password reset token and uses it
func (prs *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := Hash(token)
	var user User
	var pwReset PasswordReset

	// Checks if the provided token has a corresponsing hash stored in the DB
	row := prs.DB.QueryRow(`
		SELECT password_resets.id, password_resets.expires_at,
			   users.id, users.email, users.password_hash
		FROM password_resets
		JOIN users ON users.id = password_resets.user_id
		WHERE password_resets.token_hash = $1`, tokenHash)
	err := row.Scan(
		&pwReset.ID, &pwReset.ExpiresAt,
		&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	// Checks if the token has expired
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}

	// Consumes the token
	err = prs.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	return &user, nil
}

func (prs *PasswordResetService) delete(id int) error {
	_, err := prs.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
