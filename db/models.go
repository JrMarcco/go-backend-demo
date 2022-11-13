// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"database/sql"
	"time"
)

// account
type Account struct {
	// id
	ID           sql.NullInt64 `json:"id"`
	AccountOwner string        `json:"accountOwner"`
	Balance      int64         `json:"balance"`
	Currency     string        `json:"currency"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

// entries
type Entry struct {
	// id
	ID        sql.NullInt64 `json:"id"`
	AccountID int64         `json:"accountID"`
	Amount    int64         `json:"amount"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

// transfer
type Transfer struct {
	// id
	ID        sql.NullInt64 `json:"id"`
	FromID    int64         `json:"fromID"`
	ToID      int64         `json:"toID"`
	Amount    int64         `json:"amount"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}
