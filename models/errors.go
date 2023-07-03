package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

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

// /////////////////////////////////////////////////////////////////////////////
// FILEERROR
// /////////////////////////////////////////////////////////////////////////////

// Server-side validation to ensure that the user does not update files that are
// not an image file. The validation has two aspects: 1. check whether the
// extension is supported and 2. check whether the content type matches that of
// an image
type FileError struct {
	Issue string
}

// Error prints the string format of the file error
func (fe FileError) Error() string {
	return fmt.Sprintf("invalid file: %v", fe.Issue)
}

// ReadSeeke reades a file up to a certain point and restarts reading it from
// the beginning. This is important to ensure that the bytes already read are
// not left behing.
func checkContentType(r io.ReadSeeker, allowedTypes []string) error {
	// The algorithm that checks the content type only needs 512 bytes to assert
	// the type
	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes)
	if err != nil {
		return fmt.Errorf("checking content type: %w", err)
	}

	// Returns to the beginning of the file
	_, err = r.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("checking content type: %w", err)
	}

	// Checks whether the content type matches any of the allowed types
	contentType := http.DetectContentType(testBytes)
	var isAllowed bool = false
	for _, t := range allowedTypes {
		if contentType == t {
			isAllowed = true
		}
	}

	if !isAllowed {
		return FileError{
			Issue: fmt.Sprintf("invalid content type: %s", contentType),
		}
	}

	return nil

}

// checkExtension checks whether the given filename has one of the allowed
// extensions. It returns a FileError if it doesn't, and nil if it does
func checkExtension(filename string, allowedExtensions []string) error {
	if !hasExtension(filename, allowedExtensions) {
		return FileError{
			Issue: fmt.Sprintf("invalid extension: %s", filepath.Ext(filename)),
		}
	}

	return nil
}
