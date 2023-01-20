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
	// TODO: store the session in DB
	return &session, nil
}

// User returns an user for a given session token
func (ss *SessionService) User(token string) (*User, error) {
	// TODO: implement
	return nil, nil
}

// hash hashes a session token
func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
