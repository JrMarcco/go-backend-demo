package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func (m *mysqlTestSuite) createUser(t *testing.T, user CreateUserParams) int64 {

	res, err := m.queries.CreateUser(context.Background(), user)

	require.NoError(t, err)
	id, err := res.LastInsertId()

	require.NoError(t, err)
	require.NotZero(t, id)

	return id
}

func (m *mysqlTestSuite) TestCreateUser() {
	t := m.T()

	args := CreateUserParams{
		Username:     util.RandomString(6),
		Email:        fmt.Sprintf("%s@email.com", util.RandomString(6)),
		HashedPasswd: "secret",
	}

	id := m.createUser(t, args)

	user, err := m.queries.GetUser(context.Background(), sql.NullInt64{
		Int64: id,
		Valid: true,
	})

	require.NoError(t, err)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.HashedPasswd, user.HashedPasswd)
}
