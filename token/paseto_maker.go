package token

import (
	pasetoV4 "aidanwoods.dev/go-paseto"
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"time"
)

type PasetoLocalMakerV2 struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoLocalMaker(symmetricKey string) (Maker, error) {

	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoLocalMakerV2{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

var _ Maker = (*PasetoLocalMakerV2)(nil)

func (p *PasetoLocalMakerV2) Generate(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}
	return p.paseto.Encrypt(p.symmetricKey, payload, nil)
}

func (p *PasetoLocalMakerV2) Verify(token string) (*Payload, error) {
	payload := &Payload{}

	err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if err = payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}

type PasetoLocalMakerV4 struct {
	v4SymmetricKey pasetoV4.V4SymmetricKey
	parser         pasetoV4.Parser
}

func NewPasetoLocalMarkerV4() Maker {
	return &PasetoLocalMakerV4{
		v4SymmetricKey: pasetoV4.NewV4SymmetricKey(),
		parser:         pasetoV4.NewParser(),
	}
}

var _ Maker = (*PasetoLocalMakerV4)(nil)

func (p *PasetoLocalMakerV4) Generate(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}

	token := pasetoV4.NewToken()
	token.SetString("id", payload.ID.String())
	token.SetSubject(payload.Username)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)

	return token.V4Encrypt(p.v4SymmetricKey, nil), nil
}

func (p *PasetoLocalMakerV4) Verify(token string) (*Payload, error) {
	_, err := p.parser.ParseV4Local(p.v4SymmetricKey, token, nil)
	if err != nil {
		return nil, err
	}

	return &Payload{}, nil
}

type PasetoPubMakerV4 struct {
	asymmetricSecretKey pasetoV4.V4AsymmetricSecretKey
	asymmetricPublicKey pasetoV4.V4AsymmetricPublicKey
}

func NewPasetoPubMarkerV4() Maker {
	secretKey := pasetoV4.NewV4AsymmetricSecretKey()

	return &PasetoPubMakerV4{
		asymmetricSecretKey: secretKey,
		asymmetricPublicKey: secretKey.Public(),
	}
}

var _ Maker = (*PasetoPubMakerV4)(nil)

func (p *PasetoPubMakerV4) Generate(username string, duration time.Duration) (string, error) {
	return "", nil
}

func (p *PasetoPubMakerV4) Verify(token string) (*Payload, error) {
	return nil, nil
}
