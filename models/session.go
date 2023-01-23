package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/wagnojunior/lenslocked/rand"
)

const (
	// MinBytesPerToken is the minimum number of bytes per each session token
	MinBytesPerToken = 32
)

// Session defines the session model according to the `sessions` SQL table. Although this struct should map to the SQL table, some entires (such as `Token`) are not present in the DB
type Session struct {
	ID        int
	UserID    int
	Token     string // Token is only set when creating a new session
	TokenHash string
}

// SessionService defines the connection to the DB
type SessionService struct {
	DB *sql.DB
	// BytesPerToken determines how many bytes used to generate each session token. If `BytesPerToken` is not provided or is less than `MinBytesPerToken`, then `MinBytesPerToken` is used instead
	BytesPerToken int
}

// Create creates a session
func (ss *SessionService) Create(userID int) (*Session, error) {
	// Checks if the given bytes per token meets the minimum requirement
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	// Generates a session token using the specified `bytesPerToken`
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	// Hashes the session token
	tokenHash := ss.hash(token)

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
	}

	// Tries to update a user's session or creates a new session in case the user does not have one
	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		row = ss.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2)
			RETURNING id;`, session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &session, nil
}

// Delete deletes a session defined by the token stored in the cookie
func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)

	_, err := ss.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// User returns an user for a given session token
func (ss *SessionService) User(token string) (*User, error) {
	// Hashes the token string
	tokenHash := ss.hash(token)

	// Queries the DB using a given session token, gets the returning user ID, and stores it in `user`
	var user User
	row := ss.DB.QueryRow(`
		SELECT user_id
		FROM sessions
		WHERE token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	// Queries the DB using a given user ID, gets the returning user email and password hash, and stores it in `user`
	row = ss.DB.QueryRow(`
		SELECT email, password_hash
		FROM users
		WHERE id = $1;`, user.ID)
	err = row.Scan(&user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	return &user, nil
}

// hash hashes a session token
func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
