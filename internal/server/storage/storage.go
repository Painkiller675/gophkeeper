package storage

import (
	"context"
	"errors"
	"github.com/Painkiller675/gophkeeper/internal/server/models"
)

// Возможные ошибки при работе с хранилищем UserStorage
var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserConflict = errors.New("user conflict")
)

// UserStorage interface for user credentials saving
type UserStorage interface {
	// PutUser saves user credentials
	PutUser(ctx context.Context, user *models.User) (*models.User, error)
	// GetUser returns userID with by credentials
	GetUser(ctx context.Context, user *models.User) (*models.User, error)
}

// possible errors when working with SecretStorage storage
var (
	ErrSecretNotFound = errors.New("secret not found")
	ErrSecretConflict = errors.New("secret conflict")
)

// SecretStorage interface for users' private data saving
type SecretStorage interface {
	// GetSecret returns secret by name for the user with userID
	GetSecret(ctx context.Context, name string, userID int) (*models.Secret, error)
	// CreateSecret creates a new secret
	CreateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error)
	// UpdateSecret updates a secret
	UpdateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error)
	// DeleteSecret deletes a secret
	DeleteSecret(ctx context.Context, secret *models.Secret) error
	// ListSecrets returns all the user's secrets list by userID
	ListSecrets(ctx context.Context, userID int) ([]*models.Secret, error)
}
