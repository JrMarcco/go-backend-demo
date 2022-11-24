// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/jrmarcco/go-backend-demo/db/ent/entries"
	"github.com/jrmarcco/go-backend-demo/db/ent/predicate"
)

// EntriesQuery is the builder for querying Entries entities.
type EntriesQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.Entries
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the EntriesQuery builder.
func (eq *EntriesQuery) Where(ps ...predicate.Entries) *EntriesQuery {
	eq.predicates = append(eq.predicates, ps...)
	return eq
}

// Limit adds a limit step to the query.
func (eq *EntriesQuery) Limit(limit int) *EntriesQuery {
	eq.limit = &limit
	return eq
}

// Offset adds an offset step to the query.
func (eq *EntriesQuery) Offset(offset int) *EntriesQuery {
	eq.offset = &offset
	return eq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (eq *EntriesQuery) Unique(unique bool) *EntriesQuery {
	eq.unique = &unique
	return eq
}

// Order adds an order step to the query.
func (eq *EntriesQuery) Order(o ...OrderFunc) *EntriesQuery {
	eq.order = append(eq.order, o...)
	return eq
}

// First returns the first Entries entity from the query.
// Returns a *NotFoundError when no Entries was found.
func (eq *EntriesQuery) First(ctx context.Context) (*Entries, error) {
	nodes, err := eq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{entries.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (eq *EntriesQuery) FirstX(ctx context.Context) *Entries {
	node, err := eq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Entries ID from the query.
// Returns a *NotFoundError when no Entries ID was found.
func (eq *EntriesQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = eq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{entries.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (eq *EntriesQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := eq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Entries entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Entries entity is found.
// Returns a *NotFoundError when no Entries entities are found.
func (eq *EntriesQuery) Only(ctx context.Context) (*Entries, error) {
	nodes, err := eq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{entries.Label}
	default:
		return nil, &NotSingularError{entries.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (eq *EntriesQuery) OnlyX(ctx context.Context) *Entries {
	node, err := eq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Entries ID in the query.
// Returns a *NotSingularError when more than one Entries ID is found.
// Returns a *NotFoundError when no entities are found.
func (eq *EntriesQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = eq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{entries.Label}
	default:
		err = &NotSingularError{entries.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (eq *EntriesQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := eq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of EntriesSlice.
func (eq *EntriesQuery) All(ctx context.Context) ([]*Entries, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return eq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (eq *EntriesQuery) AllX(ctx context.Context) []*Entries {
	nodes, err := eq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Entries IDs.
func (eq *EntriesQuery) IDs(ctx context.Context) ([]uint64, error) {
	var ids []uint64
	if err := eq.Select(entries.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (eq *EntriesQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := eq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (eq *EntriesQuery) Count(ctx context.Context) (int, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return eq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (eq *EntriesQuery) CountX(ctx context.Context) int {
	count, err := eq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (eq *EntriesQuery) Exist(ctx context.Context) (bool, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return eq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (eq *EntriesQuery) ExistX(ctx context.Context) bool {
	exist, err := eq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the EntriesQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (eq *EntriesQuery) Clone() *EntriesQuery {
	if eq == nil {
		return nil
	}
	return &EntriesQuery{
		config:     eq.config,
		limit:      eq.limit,
		offset:     eq.offset,
		order:      append([]OrderFunc{}, eq.order...),
		predicates: append([]predicate.Entries{}, eq.predicates...),
		// clone intermediate query.
		sql:    eq.sql.Clone(),
		path:   eq.path,
		unique: eq.unique,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		AccountID uint64 `json:"account_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Entries.Query().
//		GroupBy(entries.FieldAccountID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (eq *EntriesQuery) GroupBy(field string, fields ...string) *EntriesGroupBy {
	grbuild := &EntriesGroupBy{config: eq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return eq.sqlQuery(ctx), nil
	}
	grbuild.label = entries.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		AccountID uint64 `json:"account_id,omitempty"`
//	}
//
//	client.Entries.Query().
//		Select(entries.FieldAccountID).
//		Scan(ctx, &v)
func (eq *EntriesQuery) Select(fields ...string) *EntriesSelect {
	eq.fields = append(eq.fields, fields...)
	selbuild := &EntriesSelect{EntriesQuery: eq}
	selbuild.label = entries.Label
	selbuild.flds, selbuild.scan = &eq.fields, selbuild.Scan
	return selbuild
}

// Aggregate returns a EntriesSelect configured with the given aggregations.
func (eq *EntriesQuery) Aggregate(fns ...AggregateFunc) *EntriesSelect {
	return eq.Select().Aggregate(fns...)
}

func (eq *EntriesQuery) prepareQuery(ctx context.Context) error {
	for _, f := range eq.fields {
		if !entries.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if eq.path != nil {
		prev, err := eq.path(ctx)
		if err != nil {
			return err
		}
		eq.sql = prev
	}
	return nil
}

func (eq *EntriesQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Entries, error) {
	var (
		nodes = []*Entries{}
		_spec = eq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Entries).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Entries{config: eq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, eq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (eq *EntriesQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := eq.querySpec()
	_spec.Node.Columns = eq.fields
	if len(eq.fields) > 0 {
		_spec.Unique = eq.unique != nil && *eq.unique
	}
	return sqlgraph.CountNodes(ctx, eq.driver, _spec)
}

func (eq *EntriesQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := eq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (eq *EntriesQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   entries.Table,
			Columns: entries.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: entries.FieldID,
			},
		},
		From:   eq.sql,
		Unique: true,
	}
	if unique := eq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := eq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, entries.FieldID)
		for i := range fields {
			if fields[i] != entries.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := eq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := eq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := eq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := eq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (eq *EntriesQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(eq.driver.Dialect())
	t1 := builder.Table(entries.Table)
	columns := eq.fields
	if len(columns) == 0 {
		columns = entries.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if eq.sql != nil {
		selector = eq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if eq.unique != nil && *eq.unique {
		selector.Distinct()
	}
	for _, p := range eq.predicates {
		p(selector)
	}
	for _, p := range eq.order {
		p(selector)
	}
	if offset := eq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := eq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// EntriesGroupBy is the group-by builder for Entries entities.
type EntriesGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (egb *EntriesGroupBy) Aggregate(fns ...AggregateFunc) *EntriesGroupBy {
	egb.fns = append(egb.fns, fns...)
	return egb
}

// Scan applies the group-by query and scans the result into the given value.
func (egb *EntriesGroupBy) Scan(ctx context.Context, v any) error {
	query, err := egb.path(ctx)
	if err != nil {
		return err
	}
	egb.sql = query
	return egb.sqlScan(ctx, v)
}

func (egb *EntriesGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range egb.fields {
		if !entries.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := egb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := egb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (egb *EntriesGroupBy) sqlQuery() *sql.Selector {
	selector := egb.sql.Select()
	aggregation := make([]string, 0, len(egb.fns))
	for _, fn := range egb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(egb.fields)+len(egb.fns))
		for _, f := range egb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(egb.fields...)...)
}

// EntriesSelect is the builder for selecting fields of Entries entities.
type EntriesSelect struct {
	*EntriesQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (es *EntriesSelect) Aggregate(fns ...AggregateFunc) *EntriesSelect {
	es.fns = append(es.fns, fns...)
	return es
}

// Scan applies the selector query and scans the result into the given value.
func (es *EntriesSelect) Scan(ctx context.Context, v any) error {
	if err := es.prepareQuery(ctx); err != nil {
		return err
	}
	es.sql = es.EntriesQuery.sqlQuery(ctx)
	return es.sqlScan(ctx, v)
}

func (es *EntriesSelect) sqlScan(ctx context.Context, v any) error {
	aggregation := make([]string, 0, len(es.fns))
	for _, fn := range es.fns {
		aggregation = append(aggregation, fn(es.sql))
	}
	switch n := len(*es.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		es.sql.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		es.sql.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := es.sql.Query()
	if err := es.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
