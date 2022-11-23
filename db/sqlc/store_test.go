package db

import (
	"context"
	"github.com/stretchr/testify/require"
)

func (m *sqlcTestSuite) TestTransferTx() {
	t := m.T()

	store := NewStore(m.conn)

	account1 := m.createAccount(t)
	account2 := m.createAccount(t)

	t.Log("before: ", account1.Balance, account2.Balance)

	amount := int64(10)
	n := 5

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromID: account1.ID.Int64,
				ToID:   account2.ID.Int64,
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
		require.Equal(t, account1.ID.Int64, res.Transfer.FromID)
		require.Equal(t, account2.ID.Int64, res.Transfer.ToID)
		require.Equal(t, amount, res.Transfer.Amount)

		// check from entry
		require.NotEmpty(t, res.FromEntry)
		require.Equal(t, account1.ID.Int64, res.FromEntry.AccountID)
		require.Equal(t, -amount, res.FromEntry.Amount)

		// check to entry
		require.NotEmpty(t, res.ToEntry)
		require.Equal(t, account2.ID.Int64, res.ToEntry.AccountID)
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

func (m *sqlcTestSuite) TestTransferTxDeadLock() {
	t := m.T()

	store := NewStore(m.conn)

	account1 := m.createAccount(t)
	account2 := m.createAccount(t)

	t.Log("before: ", account1.Balance, account2.Balance)

	amount := int64(10)
	n := 10

	errs := make(chan error)

	for i := 0; i < n; i++ {

		fromID := account1.ID.Int64
		toID := account2.ID.Int64

		if i%2 == 0 {
			fromID = account2.ID.Int64
			toID = account1.ID.Int64
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
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
