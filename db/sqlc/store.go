package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) withTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	if err = fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v\n", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromID int64 `json:"fromID"`
	ToID   int64 `json:"toID"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"fromAccount"`
	ToAccount   Account  `json:"toAccount"`
	FromEntry   Entry    `json:"fromEntry"`
	ToEntry     Entry    `json:"toEntry"`
}

// TransferTx do transfer in transaction.
// Pay attention to that don't call the Store's query func,
// it's out of transaction.
func (s *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var txRes TransferTxResult
	err := s.withTx(ctx, func(queries *Queries) error {
		var err error

		// create transfer record
		txRes.Transfer, err = doWitTransfer(ctx, CreateTransferParams{
			FromID: args.FromID,
			ToID:   args.ToID,
			Amount: args.Amount,
		}, queries)
		if err != nil {
			return err
		}

		// create from entry
		txRes.FromEntry, err = doWithEntry(ctx, CreateEntryParams{
			AccountID: args.FromID,
			Amount:    -args.Amount,
		}, queries)
		if err != nil {
			return err
		}

		// create to entry
		txRes.ToEntry, err = doWithEntry(ctx, CreateEntryParams{
			AccountID: args.ToID,
			Amount:    args.Amount,
		}, queries)
		if err != nil {
			return err
		}

		// ensure that fromID lt toID
		if args.FromID > args.ToID {
			args.FromID, args.ToID = args.ToID, args.FromID
			args.Amount = -args.Amount
		}

		// update from account balance
		txRes.FromAccount, err = doWithAccount(ctx, args.FromID, -args.Amount, queries)
		if err != nil {
			return err
		}

		// update to account balance
		txRes.ToAccount, err = doWithAccount(ctx, args.ToID, args.Amount, queries)
		if err != nil {
			return err
		}

		return err
	})

	return txRes, err
}

func doWitTransfer(ctx context.Context, args CreateTransferParams, q *Queries) (Transfer, error) {
	res, err := q.CreateTransfer(ctx, args)
	if err != nil {
		return Transfer{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Transfer{}, err
	}

	return Transfer{
		ID: sql.NullInt64{
			Int64: id,
			Valid: true,
		},
		FromID: args.FromID,
		ToID:   args.ToID,
		Amount: args.Amount,
	}, nil
}

func doWithEntry(ctx context.Context, args CreateEntryParams, q *Queries) (Entry, error) {
	res, err := q.CreateEntry(ctx, CreateEntryParams{
		AccountID: args.AccountID,
		Amount:    args.Amount,
	})
	if err != nil {
		return Entry{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Entry{}, err
	}

	return Entry{
		ID: sql.NullInt64{
			Int64: id,
			Valid: true,
		},
		AccountID: args.AccountID,
		Amount:    args.Amount,
	}, nil
}

func doWithAccount(ctx context.Context, aid int64, amount int64, q *Queries) (Account, error) {
	// update account balance and return after update
	sqlId := sql.NullInt64{Int64: aid, Valid: true}
	err := q.AddBalance(ctx, AddBalanceParams{
		ID:     sqlId,
		Amount: amount,
	})
	if err != nil {
		return Account{}, err
	}
	return q.GetAccount(ctx, sqlId)
}
