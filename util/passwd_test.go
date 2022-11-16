package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswd(t *testing.T) {

	passwd := RandomString(8)

	hashed, err := HashPasswd(passwd)
	require.NoError(t, err)

	err = CheckPasswd(passwd, hashed)
	require.NoError(t, err)

	wrongPasswd := RandomString(8)
	err = CheckPasswd(wrongPasswd, hashed)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	newHashed, err := HashPasswd(hashed)
	require.NoError(t, err)
	require.NotEqual(t, hashed, newHashed)
}
