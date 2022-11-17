package token

import (
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewPasetoLocalMaker(t *testing.T) {
	symmetricKey := util.RandomString(32)
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
			arg:     symmetricKey,
			wantErr: nil,
			wantRes: &PasetoLocalMakerV2{
				paseto:       paseto.NewV2(),
				symmetricKey: []byte(symmetricKey),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			maker, err := NewPasetoLocalMaker(tc.arg)
			if err != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}
			require.Equal(t, tc.wantRes, maker)
		})
	}
}

func TestPasetoLocalMaker_Verify(t *testing.T) {
	symmetricKey := util.RandomString(32)

	tcs := []struct {
		name     string
		username string
		duration time.Duration
		wantErr  error
		wantRes  string
	}{
		{
			name:     "Normal Case",
			username: util.RandomString(8),
			duration: time.Minute,
			wantErr:  nil,
		},
		{
			name:     "Expired Case",
			username: util.RandomString(8),
			duration: -time.Minute,
			wantErr:  ErrExpiredToken,
		},
		{
			name:     "Invalid Signature Case",
			username: util.RandomString(8),
			duration: time.Minute,
			wantErr:  ErrInvalidToken,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			maker, err := NewPasetoLocalMaker(symmetricKey)
			require.NoError(t, err)

			token, err := maker.Generate(tc.username, tc.duration)
			require.NoError(t, err)

			payload, err := maker.Verify(token)
			if err != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}
			require.Equal(t, tc.username, payload.Username)
		})
	}
}
