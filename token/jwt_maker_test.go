package token

import (
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewJwtMaker(t *testing.T) {
	secretKey := util.RandomString(32)
	tcs := []struct {
		name    string
		arg     string
		wantErr error
		wantRes Maker
	}{
		{
			name:    "Short SecretKey Case",
			arg:     "util",
			wantErr: ErrInvalidKeySize,
			wantRes: nil,
		},
		{
			name:    "Normal Case",
			arg:     secretKey,
			wantErr: nil,
			wantRes: &JwtMaker{
				secretKey: secretKey,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			maker, err := NewJwtMaker(tc.arg)
			if err != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.Equal(t, tc.wantRes, maker)
		})
	}
}

func TestJwtMaker_Verify(t *testing.T) {

	secretKey := util.RandomString(32)

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
			maker, err := NewJwtMaker(secretKey)
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
