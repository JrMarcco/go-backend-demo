package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"go-learning/backend-demo/util"
)

func (m *mysqlTestSuite) TestTransferTx() {
	t := m.T()

	store := NewStore(m.db)

	aid1 := m.createAccount(t, CreateAccountParams{
		AccountOwner: util.RandomString(6),
		Balance:      util.RandomInt64(100, 1000),
		Currency:     "RMB",
	})
	aid2 := m.createAccount(t, CreateAccountParams{
		AccountOwner: util.RandomString(6),
		Balance:      util.RandomInt64(100, 1000),
		Currency:     "RMB",
	})

	account1, err := m.queries.GetAccount(context.Background(), sql.NullInt64{Int64: aid1, Valid: true})
	require.NoError(t, err)
	account2, err := m.queries.GetAccount(context.Background(), sql.NullInt64{Int64: aid2, Valid: true})
	require.NoError(t, err)

	t.Log("before: ", account1.Balance, account2.Balance)

	amount := int64(10)
	n := 5

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromID: aid1,
				ToID:   aid2,
				Amount: amount,
			})

			errs <- err
			results <- res
		}()
	}

	existed := make(map[int]struct{}, n)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results

		// check transfer
		require.NotEmpty(t, res.Transfer)
		require.NotZero(t, res.Transfer.ID)
		require.Equal(t, aid1, res.Transfer.FromID)
		require.Equal(t, aid2, res.Transfer.ToID)
		require.Equal(t, amount, res.Transfer.Amount)

		// check from entry
		require.NotEmpty(t, res.FromEntry)
		require.Equal(t, aid1, res.FromEntry.AccountID)
		require.Equal(t, -amount, res.FromEntry.Amount)

		// check to entry
		require.NotEmpty(t, res.ToEntry)
		require.Equal(t, aid2, res.ToEntry.AccountID)
		require.Equal(t, amount, res.ToEntry.Amount)

		// check from account and to account
		require.NotEmpty(t, res.FromAccount)
		require.Equal(t, account1.ID, res.FromAccount.ID)
		require.NotEmpty(t, res.ToAccount)
		require.Equal(t, account2.ID, res.ToAccount.ID)

		// check the transfer amount
		diff1 := account1.Balance - res.FromAccount.Balance
		diff2 := res.ToAccount.Balance - account2.Balance
		require.True(t, diff1 == diff2 && diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = struct{}{}
	}

	updatedAccount1, err := m.queries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := m.queries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	t.Log("after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, int64(n)*amount, account1.Balance-updatedAccount1.Balance)
	require.Equal(t, int64(n)*amount, updatedAccount2.Balance-account2.Balance)
}

func (m *mysqlTestSuite) TestTransferTxDeadLock() {
	t := m.T()

	store := NewStore(m.db)

	aid1 := m.createAccount(t, CreateAccountParams{
		AccountOwner: util.RandomString(6),
		Balance:      util.RandomInt64(100, 1000),
		Currency:     "RMB",
	})
	aid2 := m.createAccount(t, CreateAccountParams{
		AccountOwner: util.RandomString(6),
		Balance:      util.RandomInt64(100, 1000),
		Currency:     "RMB",
	})

	account1, err := m.queries.GetAccount(context.Background(), sql.NullInt64{Int64: aid1, Valid: true})
	require.NoError(t, err)
	account2, err := m.queries.GetAccount(context.Background(), sql.NullInt64{Int64: aid2, Valid: true})
	require.NoError(t, err)

	t.Log("before: ", account1.Balance, account2.Balance)

	amount := int64(10)
	n := 10

	errs := make(chan error)

	for i := 0; i < n; i++ {

		fromID := aid1
		toID := aid2

		if i%2 == 0 {
			fromID = aid2
			toID = aid1
		}

		go func() {
			_, err = store.TransferTx(context.Background(), TransferTxParams{
				FromID: fromID,
				ToID:   toID,
				Amount: amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := m.queries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := m.queries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	t.Log("after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
