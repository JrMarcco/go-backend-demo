package db

import (
	"context"
	"database/sql"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func (m *mysqlTestSuite) createAccount(t *testing.T) Account {
	user := m.createUser(t)

	createAccountArgs := CreateAccountParams{
		AccountOwner: user.Username,
		Balance:      util.RandomInt64(100, 10000),
		Currency:     "RMB",
	}

	res, err := m.queries.CreateAccount(context.Background(), createAccountArgs)

	assert.NoError(t, err)
	accountID, _ := res.LastInsertId()

	account, err := m.queries.GetAccount(context.Background(), sql.NullInt64{Int64: accountID, Valid: true})

	assert.NoError(t, err)
	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	return account
}

func (m *mysqlTestSuite) TestCreateAccount() {
	_ = m.createAccount(m.T())
}

func (m *mysqlTestSuite) TestGetAccount() {
	t := m.T()

	account1 := m.createAccount(t)
	account2, err := m.queries.GetAccount(context.Background(), account1.ID)

	assert.NoError(t, err)
	assert.Equal(t, account1.AccountOwner, account2.AccountOwner)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
}

func (m *mysqlTestSuite) TestDeleteAccount() {
	t := m.T()

	account := m.createAccount(t)

	err := m.queries.DeleteAccount(context.Background(), sql.NullInt64{
		Int64: account.ID.Int64,
		Valid: true,
	})

	assert.NoError(t, err)

	deletedAccount, err := m.queries.GetAccount(context.Background(), account.ID)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, deletedAccount)
}
