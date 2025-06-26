package models

import "fmt"

var _ Secret = (*Credentials)(nil)

// Credentials login & password
type Credentials struct {
	Login    string
	Password string
}

// Type returns saved info type
func (c Credentials) Type() SecretType {
	return secretTypeCredentials
}

// String shows private info
func (c Credentials) String() string {
	return fmt.Sprintf("Login: %s, Password: %s", c.Login, c.Password)
}
