package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// /////////////////////////////////////////////////////////////////////////////
// CONSTANTS
// /////////////////////////////////////////////////////////////////////////////
const (
	// MinBytesPerToken is the minimum number of bytes per each session token
	MinBytesPerToken = 32
)

// Session defines the session model according to the `sessions` SQL table.
// Although this struct should map to the SQL table, some entires (such as
// `Token`) are not present in the DB
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

// type TokenManager struct {
// 	sessionService *SessionService
// }

// func (tm TokenManager) New() (token, tokenHash string, err error) {
// 	// Checks if the given bytes per token meets the minimum requirement
// 	bytesPerToken := tm.sessionService.BytesPerToken
// 	if bytesPerToken < MinBytesPerToken {
// 		bytesPerToken = MinBytesPerToken
// 	}

// 	// Generates a session token using the specified `bytesPerToken`
// 	token, err = rand.String(bytesPerToken)
// 	if err != nil {
// 		return "", "", fmt.Errorf("create: %w", err)
// 	}

// 	// Hashes the session token
// 	tokenHash = tm.sessionService.hash(token)

// 	return token, tokenHash, nil
// }

// Create creates a session
func (ss *SessionService) Create(userID int) (*Session, error) {
	// Gets a token and token hash
	token, tokenHash, err := New(ss.BytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
	}

	// Creates a new session with the given value or updates an existing session (ON CONFLICT clause)
	row := ss.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2
		RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		// Check if the error is of type `*pgconn.PgError`
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			// Checks if the error is of type `UndefinedTable`
			if pgError.Code == pgerrcode.UndefinedTable {
				return nil, ErrCreateSession
			}
		}
		return nil, fmt.Errorf("create: %w", err)
	}

	return &session, nil
}

// Delete deletes a session defined by the token stored in the cookie
func (ss *SessionService) Delete(token string) error {
	tokenHash := Hash(token)

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
	tokenHash := Hash(token)

	// Queries the DB for the user that corresponds to a token hash
	var user User
	row := ss.DB.QueryRow(`
		SELECT users.id, users.email, users.password_hash
		FROM users
		JOIN sessions ON users.id = sessions.user_id
		WHERE token_hash = $1`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	return &user, nil
}

// // hash hashes a session token
// func (ss *SessionService) hash(token string) string {
// 	tokenHash := sha256.Sum256([]byte(token))
// 	return base64.URLEncoding.EncodeToString(tokenHash[:])
// }
