package token

// Storage interface of token loading and saving
type Storage interface {
	// Load loads token
	Load() (accessToken string, err error)
	// Save saves token
	Save(accessToken string) error
}
