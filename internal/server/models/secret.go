package models

import (
	"github.com/google/uuid"
)

// Secret includes users' secret data
type Secret struct {
	Name    string
	Content []byte
	Version uuid.UUID
	OwnerID int
}
