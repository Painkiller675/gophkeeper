package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Payload includes payload token data
type Payload struct {
	jwt.StandardClaims

	UserID int `json:"user_id"`
}

// NewPayload create new payload with given userID and token expiration date
func NewPayload(userID int, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	payload := &Payload{
		StandardClaims: jwt.StandardClaims{
			Id:        tokenID.String(),
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(duration).Unix(),
		},
		UserID: userID,
	}
	return payload, nil
}

// Valid checks the token
func (p *Payload) Valid() error {
	if time.Now().Unix() > p.ExpiresAt {
		return ErrExpiredToken
	}
	return nil
}
