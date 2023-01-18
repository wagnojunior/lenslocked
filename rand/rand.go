package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Bytes returns a random byte slice of size `n` and an error using crypto/rand package
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("Bytes: %w", err)
	}

	return b, nil
}

// String returns returns a random string of size `n` and an error using the crypto/rand package
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("String: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

const SessionTokenBytes = 32

// SessionToken returns a random string of size 32 and and error using the crypto/rand package
func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}
