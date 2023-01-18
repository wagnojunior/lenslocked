package models

// Session defines the session model according to the sessions SQL table
type Session struct {
	ID        int
	UserID    int
	TokenHash string
}
