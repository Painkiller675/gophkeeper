package models

// User includes users' credentials
type User struct {
	ID           int
	Email        string
	PasswordHash string
}
