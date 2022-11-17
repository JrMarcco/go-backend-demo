package token

import (
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/require"
	"testing"
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

func TestJwtMaker_Generate(t *testing.T) {

}

func TestJwtMaker_Verify(t *testing.T) {
}
