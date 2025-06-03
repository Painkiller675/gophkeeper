package token

import (
	"errors"
)

// MinKeySize minimum key size
const MinKeySize = 16

var (
	// ErrInvalidKeySize insufficient length of the key
	ErrInvalidKeySize = errors.New("invalid key size")
	// ErrInvalidToken invalid token
	ErrInvalidToken = errors.New("token is invalid")
	// ErrExpiredToken token is expired
	ErrExpiredToken = errors.New("token has expired")
)

// Manager interface for token generating and token checking for authentication
type Manager interface {
	// Create creates token for given userID
	Create(userID int) (token string, err error)
	// Validate checks if token is valid
	Validate(accessToken string) (payload *Payload, err error)
}
