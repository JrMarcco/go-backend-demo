package models

import (
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func (m *modelTestSuite) createAccount() *Account {
	t := m.T()

	addAccountParams := AddAccountParam{
		AccountOwner: util.RandomString(8),
		Balance:      util.RandomInt64(100, 10000),
		Currency:     "RMB",
	}

	account, err := AddAccount(addAccountParams)

	require.NoError(t, err)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.UpdatedAt)

	require.Equal(t, addAccountParams.AccountOwner, account.AccountOwner)
	require.Equal(t, addAccountParams.Balance, account.Balance)
	require.Equal(t, addAccountParams.Currency, account.Currency)

	return account
}

func (m *modelTestSuite) TestAddAccount() {
	_ = m.createAccount()
}

func (m *modelTestSuite) TestGetAccount() {
	t := m.T()
	account1 := m.createAccount()

	account2, err := GetAccount(account1.ID)
	require.NoError(t, err)

	require.Equal(t, account1.AccountOwner, account2.AccountOwner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
}

func (m *modelTestSuite) TestDelAccount() {
	t := m.T()

	account1 := m.createAccount()

	err := DelAccount(account1.ID)
	require.NoError(t, err)

	account2, err := GetAccount(account1.ID)
	require.EqualError(t, gorm.ErrRecordNotFound, err.Error())
	require.Empty(t, account2)
}
