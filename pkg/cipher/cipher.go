package cipher

// BlockCipher interface with the methods of encrypting & decrypting of any binary data
type BlockCipher interface {
	// Encrypt encrypting function
	Encrypt(plaintext []byte) (ciphertext []byte, err error)
	// Decrypt decrypting function
	Decrypt(ciphertext []byte) (plaintext []byte, err error)
}
