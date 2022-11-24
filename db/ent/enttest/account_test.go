package enttest

import (
	"context"
	"github.com/jrmarcco/go-backend-demo/db/ent"
	"github.com/jrmarcco/go-backend-demo/db/ent/account"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/require"
)

func (e *entTestSuite) createAccount() *ent.Account {
	t := e.T()

	randAccount := ent.Account{
		AccountOwner: util.RandomString(8),
		Balance:      util.RandomInt64(100, 1000),
		Currency:     "RMB",
	}

	account, err := e.client.Account.Create().
		SetAccountOwner(randAccount.AccountOwner).
		SetBalance(randAccount.Balance).
		Save(context.Background())

	require.NoError(t, err)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.UpdatedAt)

	require.Equal(t, randAccount.AccountOwner, account.AccountOwner)
	require.Equal(t, randAccount.Balance, account.Balance)
	require.Equal(t, randAccount.Currency, account.Currency)

	return account
}

func (e *entTestSuite) TestCreateAccount() {
	_ = e.createAccount()
}

func (e *entTestSuite) TestQueryAccount() {
	t := e.T()

	account1 := e.createAccount()
	account2, err := e.client.Account.Query().
		Where(account.ID(account1.ID)).
		First(context.Background())

	require.NoError(t, err)

	require.Equal(t, account1.AccountOwner, account2.AccountOwner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
}
