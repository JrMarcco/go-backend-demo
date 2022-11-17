package token

import (
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"time"
)

type PasetoLocalMakerV2 struct {
	paseto       *paseto.V2
	asymetricKey []byte
}

func NewPasetoLocalMaker(asymetricKey string) (Maker, error) {

	if len(asymetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoLocalMakerV2{
		paseto:       paseto.NewV2(),
		asymetricKey: []byte(asymetricKey),
	}, nil
}

var _ Maker = (*PasetoLocalMakerV2)(nil)

func (p *PasetoLocalMakerV2) Generate(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}
	return p.paseto.Encrypt(p.asymetricKey, payload, nil)
}

func (p *PasetoLocalMakerV2) Verify(token string) (*Payload, error) {
	payload := &Payload{}

	err := p.paseto.Decrypt(token, p.asymetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if err = payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}

// ssh-keygen -t rsa -b 2048 -m PEM -f paseto.key
// openssl rsa -in paseto.key -pubout -outform PEM -out paseto.key.pub
