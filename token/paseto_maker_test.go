package token

import (
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPasetoLocalMaker(t *testing.T) {
	asymetricKey := util.RandomString(32)
	tcs := []struct {
		name    string
		arg     string
		wantErr error
		wantRes Maker
	}{
		{
			name:    "Short AsymetricKey Case",
			arg:     "123",
			wantErr: fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize),
			wantRes: nil,
		},
		{
			name:    "Normal Case",
			arg:     asymetricKey,
			wantErr: nil,
			wantRes: &PasetoLocalMaker{
				paseto:       paseto.NewV2(),
				asymetricKey: []byte(asymetricKey),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			maker, err := NewPasetoLocalMaker(tc.arg)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.Equal(t, tc.wantRes, maker)
		})
	}
}
