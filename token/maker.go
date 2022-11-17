package token

import "time"

type Maker interface {
	Generate(username string, duration time.Duration) (string, error)
	Verify(token string) (*Payload, error)
}
