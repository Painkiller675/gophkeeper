package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/Painkiller675/gophkeeper/pkg/hasher"
)

var _ hasher.Hasher = (*Hasher)(nil)

// MinKeySize minimum key size
const MinKeySize = 16

var (
	// ErrInvalidKeySize error which occurs when the Hasher with a short key is created
	ErrInvalidKeySize = errors.New("invalid key size")
)

// Hasher hmac implementation of hasher.Hasher interface
type Hasher struct {
	key []byte
}

// New returns a new hmac Hasher
func New(key string) (*Hasher, error) {
	if len(key) < MinKeySize {
		return nil, ErrInvalidKeySize
	}
	return &Hasher{
		key: []byte(key),
	}, nil
}

// Hash calculates the hash of the string
func (h *Hasher) Hash(data string) (string, error) {
	mac := hmac.New(sha256.New, h.key)
	if _, err := mac.Write([]byte(data)); err != nil {
		return "", err
	}
	sum := mac.Sum(nil)
	return hex.EncodeToString(sum), nil
}

// IsValid check the hash of given string
func (h *Hasher) IsValid(data, hash string) (bool, error) {
	mac := hmac.New(sha256.New, h.key)
	if _, err := mac.Write([]byte(data)); err != nil {
		return false, err
	}
	expected, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	actual := mac.Sum(nil)
	return hmac.Equal(expected, actual), nil
}
