package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// User defines the user model according to the `users` SQL table
type User struct {
	ID           int
	Email        string
	PasswordHash string
}

// UserService defines the connection to the users DB
type UserService struct {
	DB *sql.DB
}

// Create creates a new user
func (us *UserService) Create(email, password string) (*User, error) {
	// Makes sure all emails are lower case
	email = strings.ToLower(email)

	// Generates a []byte hashed password and converts it to string
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHash := string(hashedBytes)

	// Creates a new user object with the information that is being written to the DB
	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	// Creates a new user using the `email` and the `passwordHash`
	row := us.DB.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2) RETURNING id`,
		email, passwordHash)

	// Scans for the id that was returned and saves if into `user.ID`
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil

}

// Authenticate authenticates a user who is signin in to the server. The parameters `email` and `password` are received from the GUI. The user ID and password hash are queried from the DB and `bcrypt` is used to to match the entered password hash to the DB-stored value.
func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)

	user := User{
		Email: email,
	}

	row := us.DB.QueryRow(`
		SELECT id, password_hash
		FROM users
		WHERE email = $1`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	return &user, nil

}
