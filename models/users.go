package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// User defines the user model
type User struct {
	ID           int
	Email        string
	PasswordHash string
}

// UserService defines the connection to the DB
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

	// Creates a new user object with the information that is being writted to the DB
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
