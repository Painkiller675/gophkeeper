// gcm is used for encrypting and decrypting data
package gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"

	"github.com/rs/zerolog/log"

	blockCipher "github.com/Painkiller675/gophkeeper/pkg/cipher"
)

const (
	minPasswordSize = 32
	keySize         = 32
	nonceSize       = 12
)

var (
	// ErrInvalidPasswordSize the error of an unacceptable length
	ErrInvalidPasswordSize = errors.New("invalid password size")
)

var _ blockCipher.BlockCipher = (*Cipher)(nil)

// Cipher block cipher AES in the GCM mode
type Cipher struct {
	key   []byte
	nonce []byte
}

// New creates a Cipher instance
func New(password string) (*Cipher, error) {
	if len(password) < minPasswordSize {
		return nil, ErrInvalidPasswordSize
	}
	key := sha256.Sum256([]byte(password))
	return &Cipher{
		key:   key[:],
		nonce: key[len(key)-nonceSize:],
	}, nil
}

// Encrypt encrypt byte sequence
func (c Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	// create aes block cipher
	aesCipher, err := aes.NewCipher(c.key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create aes block cipher")
		return nil, err
	}
	// to cipher a message of any length we use the GCM (Galois/Counter Mode) algorithm
	aesGCM, err := cipher.NewGCM(aesCipher)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create aes block cipher in gcm mode")
		return nil, err
	}
	// seal - to encrypt
	return aesGCM.Seal(nil, c.nonce, plaintext, nil), nil
}

// Decrypt decrypt byte sequence
func (c Cipher) Decrypt(ciphertext []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(c.key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create aes block cipher")
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(aesCipher)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create aes block cipher in gcm mode")
		return nil, err
	}

	return aesGCM.Open(nil, c.nonce, ciphertext, nil)
}
