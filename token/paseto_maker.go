package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoLocalMaker struct {
	paseto       *paseto.V2
	asymetricKey []byte
}

func NewPasetoLocalMaker(asymetricKey string) (Maker, error) {

	if len(asymetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoLocalMaker{
		paseto:       paseto.NewV2(),
		asymetricKey: []byte(asymetricKey),
	}, nil
}

var _ Maker = (*PasetoLocalMaker)(nil)

func (p *PasetoLocalMaker) Generate(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}
	return p.paseto.Encrypt(p.asymetricKey, payload, nil)
}

func (p *PasetoLocalMaker) Verify(token string) (*Payload, error) {
	return nil, nil
}

type PasetoPublicMaker struct {
	paseto *paseto.V2
}
