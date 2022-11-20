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

	token, err := buildPasetoToken(username, duration)
	if err != nil {
		return "", err
	}

	return token.V4Encrypt(p.v4SymmetricKey, nil), nil
}

func (p *PasetoLocalMakerV4) Verify(token string) (*Payload, error) {
	pasetoToken, err := p.parser.ParseV4Local(p.v4SymmetricKey, token, nil)
	if err != nil {
		return nil, err
	}
	return payload2PasetoToken(pasetoToken)
}

type PasetoPubMakerV4 struct {
	asymmetricSecretKey pasetoV4.V4AsymmetricSecretKey
	asymmetricPublicKey pasetoV4.V4AsymmetricPublicKey
	parser              pasetoV4.Parser
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
	token, err := buildPasetoToken(username, duration)
	if err != nil {
		return "", err
	}

	return token.V4Sign(p.asymmetricSecretKey, nil), nil
}

func (p *PasetoPubMakerV4) Verify(token string) (*Payload, error) {
	pasetoToken, err := p.parser.ParseV4Public(p.asymmetricPublicKey, token, nil)
	if err != nil {
		return nil, err
	}
	return payload2PasetoToken(pasetoToken)
}

func buildPasetoToken(username string, duration time.Duration) (pasetoV4.Token, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return pasetoV4.Token{}, err
	}

	token := pasetoV4.NewToken()
	token.SetString("id", payload.ID)
	token.SetSubject(payload.Username)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)

	return token, nil
}

func payload2PasetoToken(token *pasetoV4.Token) (*Payload, error) {
	id, err := token.GetString("id")
	if err != nil {
		return nil, err
	}
	username, err := token.GetSubject()
	if err != nil {
		return nil, err
	}
	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		return nil, err
	}
	expireAt, err := token.GetExpiration()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expireAt,
	}
	if err = payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}
