package hasher

// Hasher main interface implemented by hash-functions
type Hasher interface {
	// Hash calculates the hash of the string
	Hash(data string) (string, error)
	// IsValid checks the hash of the string
	IsValid(data string, hash string) (bool, error)
}
