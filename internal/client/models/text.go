package models

import "fmt"

var _ Secret = (*Text)(nil)

// Text any text data
type Text struct {
	Data string
}

// Type returns saved data type
func (t Text) Type() SecretType {
	return secretTypeText
}

// String show private data
func (t Text) String() string {
	return fmt.Sprintf("TextData: %s", t.Data)
}
