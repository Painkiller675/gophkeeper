package models

import (
	"fmt"
)

var _ Secret = (*Card)(nil)

// Card bank cards data
type Card struct {
	Number       string
	ExpiryDate   string
	SecurityCode string
	Holder       string
}

// Type returns saved type info
func (c Card) Type() SecretType {
	return secretTypeCard
}

// String shows private info
func (c Card) String() string {
	return fmt.Sprintf("Number: %s, ExpiryDate: %s, SecurityCode: %s, Holder: %s",
		c.Number, c.ExpiryDate, c.SecurityCode, c.Holder)
}
