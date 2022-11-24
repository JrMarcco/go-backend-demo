// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/jrmarcco/go-backend-demo/db/ent/entries"
)

// Entries is the model entity for the Entries schema.
type Entries struct {
	config `json:"-"`
	// ID of the ent.
	ID uint64 `json:"id,omitempty"`
	// AccountID holds the value of the "account_id" field.
	AccountID uint64 `json:"account_id,omitempty"`
	// Amount holds the value of the "amount" field.
	Amount int64 `json:"amount,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Entries) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case entries.FieldID, entries.FieldAccountID, entries.FieldAmount:
			values[i] = new(sql.NullInt64)
		case entries.FieldCreatedAt, entries.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Entries", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Entries fields.
func (e *Entries) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case entries.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			e.ID = uint64(value.Int64)
		case entries.FieldAccountID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field account_id", values[i])
			} else if value.Valid {
				e.AccountID = uint64(value.Int64)
			}
		case entries.FieldAmount:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field amount", values[i])
			} else if value.Valid {
				e.Amount = value.Int64
			}
		case entries.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				e.CreatedAt = value.Time
			}
		case entries.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				e.UpdatedAt = value.Time
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Entries.
// Note that you need to call Entries.Unwrap() before calling this method if this Entries
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Entries) Update() *EntriesUpdateOne {
	return (&EntriesClient{config: e.config}).UpdateOne(e)
}

// Unwrap unwraps the Entries entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Entries) Unwrap() *Entries {
	_tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("ent: Entries is not a transactional entity")
	}
	e.config.driver = _tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Entries) String() string {
	var builder strings.Builder
	builder.WriteString("Entries(")
	builder.WriteString(fmt.Sprintf("id=%v, ", e.ID))
	builder.WriteString("account_id=")
	builder.WriteString(fmt.Sprintf("%v", e.AccountID))
	builder.WriteString(", ")
	builder.WriteString("amount=")
	builder.WriteString(fmt.Sprintf("%v", e.Amount))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(e.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(e.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// EntriesSlice is a parsable slice of Entries.
type EntriesSlice []*Entries

func (e EntriesSlice) config(cfg config) {
	for _i := range e {
		e[_i].config = cfg
	}
}
