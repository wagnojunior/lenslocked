package models

import "errors"

var (
	// SESSION
	ErrCreateSession = errors.New("models: could not create a session")

	// USER
	ErrEmailTaken  = errors.New("models: email address is already in use")
	ErrInvalidUser = errors.New("models: failed to retrieve user from the database")
	ErrInvalidPW   = errors.New("models: failed to match the password with the stored password-hash")

	// GALLERY
	ErrInvalidGallery = errors.New("models: failed to retrieve gallery from the databse")

	// IMAGE
	ErrImageNotFound = errors.New("models: failed to query for image")
)
