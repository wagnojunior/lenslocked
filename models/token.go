package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/wagnojunior/lenslocked/rand"
)

// New returns a new token, token hash, and an error given the number of bytes
func New(bytes int) (token, tokenHash string, err error) {
	// Checks if thegiven bytes per token meets the minimum requirement
	if bytes < MinBytesPerToken {
		bytes = MinBytesPerToken
	}

	// Generates a token using the specified number of bytes
	token, err = rand.String(bytes)
	if err != nil {
		return "", "", fmt.Errorf("create: %w", err)
	}

	// Hashes the token
	tokenHash = Hash(token)

	return token, tokenHash, nil

}

// hash hashes a session token
func Hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
