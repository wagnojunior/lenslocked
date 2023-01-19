package models

import "database/sql"

// Session defines the session model according to the `sessions` SQL table. Although this struct should map to the SQL table, some entires (such as `Token`) are not present in the DB
type Session struct {
	ID        int
	UserID    int
	Token     string // Token is only set when creating a new session
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

// Create creates a session
func (ss *SessionService) Create(userID int) (*Session, error) {
	// TODO: create the session token
	// TODO: implement SessionService.Create
	return nil, nil
}

// User returns an user for a given session token
func (ss *SessionService) User(token string) (*User, error) {
	// TODO: implement
	return nil, nil
}
