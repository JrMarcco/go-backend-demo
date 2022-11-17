package token

import "time"

type PasetoMaker struct {
}

var _ Maker = (*PasetoMaker)(nil)

func (p *PasetoMaker) Generate(username string, duration time.Duration) (string, error) {
	return "", nil
}

func (p *PasetoMaker) Verify(token string) (*Payload, error) {
	return nil, nil
}
