package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("token si invalid")
	ErrExpiredToken = errors.New("token is expired")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
