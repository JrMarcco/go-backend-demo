package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func (m *mysqlTestSuite) createUser(t *testing.T) User {

	hashedPasswd, err := util.HashPasswd(util.RandomString(8))
	require.NoError(t, err)

	createUserArgs := CreateUserParams{
		Username:     util.RandomString(6),
		Email:        fmt.Sprintf("%s@email.com", util.RandomString(6)),
		HashedPasswd: hashedPasswd,
	}
	res, err := m.queries.CreateUser(context.Background(), createUserArgs)

	require.NoError(t, err)
	id, _ := res.LastInsertId()
	require.NotZero(t, id)

	user, err := m.queries.GetUser(context.Background(), sql.NullInt64{
		Int64: id,
		Valid: true,
	})

	require.NoError(t, err)
	require.Equal(t, createUserArgs.Username, user.Username)
	require.Equal(t, createUserArgs.Email, user.Email)
	require.Equal(t, createUserArgs.HashedPasswd, user.HashedPasswd)

	return user
}

func (m *mysqlTestSuite) TestCreateUser() {
	_ = m.createUser(m.T())
}
