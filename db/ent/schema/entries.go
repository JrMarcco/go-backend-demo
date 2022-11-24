package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// Entries holds the schema definition for the Entries entity.
type Entries struct {
	ent.Schema
}

// Fields of the Entries.
func (Entries) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("account_id").Positive(),
		field.Int64("amount").Positive(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Entries.
func (Entries) Edges() []ent.Edge {
	return nil
}
