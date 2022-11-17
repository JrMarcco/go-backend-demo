package db

import (
	"context"
	"github.com/stretchr/testify/assert"
)

func (m *mysqlTestSuite) TestTransferTx() {
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
		assert.NoError(t, err)

		res := <-results

		// check transfer
		assert.NotEmpty(t, res.Transfer)
		assert.NotZero(t, res.Transfer.ID)
		assert.Equal(t, account1.ID.Int64, res.Transfer.FromID)
		assert.Equal(t, account2.ID.Int64, res.Transfer.ToID)
		assert.Equal(t, amount, res.Transfer.Amount)

		// check from entry
		assert.NotEmpty(t, res.FromEntry)
		assert.Equal(t, account1.ID.Int64, res.FromEntry.AccountID)
		assert.Equal(t, -amount, res.FromEntry.Amount)

		// check to entry
		assert.NotEmpty(t, res.ToEntry)
		assert.Equal(t, account2.ID.Int64, res.ToEntry.AccountID)
		assert.Equal(t, amount, res.ToEntry.Amount)

		// check from account and to account
		assert.NotEmpty(t, res.FromAccount)
		assert.Equal(t, account1.ID, res.FromAccount.ID)
		assert.NotEmpty(t, res.ToAccount)
		assert.Equal(t, account2.ID, res.ToAccount.ID)

		// check the transfer amount
		diff1 := account1.Balance - res.FromAccount.Balance
		diff2 := res.ToAccount.Balance - account2.Balance
		assert.True(t, diff1 == diff2 && diff1%amount == 0)

		k := int(diff1 / amount)
		assert.True(t, k >= 1 && k <= n)
		assert.NotContains(t, existed, k)
		existed[k] = struct{}{}
	}

	updatedAccount1, err := m.queries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	updatedAccount2, err := m.queries.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)

	t.Log("after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	assert.Equal(t, int64(n)*amount, account1.Balance-updatedAccount1.Balance)
	assert.Equal(t, int64(n)*amount, updatedAccount2.Balance-account2.Balance)
}

func (m *mysqlTestSuite) TestTransferTxDeadLock() {
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
		assert.NoError(t, err)
	}

	updatedAccount1, err := m.queries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	updatedAccount2, err := m.queries.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)

	t.Log("after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	assert.Equal(t, account1.Balance, updatedAccount1.Balance)
	assert.Equal(t, account2.Balance, updatedAccount2.Balance)
}
