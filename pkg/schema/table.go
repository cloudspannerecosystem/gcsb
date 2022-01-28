// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schema

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/spanner"
	"github.com/cloudspannerecosystem/gcsb/pkg/schema/information"
)

type (
	Table interface {
		SetName(string)
		Name() string
		SetType(string)
		Type() string
		HasParent() bool
		SetParentName(string)
		ParentName() string
		SetParent(Table)
		Parent() Table
		HasChild() bool
		SetChildName(string)
		ChildName() string
		SetChild(Table)
		Child() Table

		SetSpanenrState(string)
		SpannerState() string

		AddColumn(Column)
		AddIndex(Index)
		Columns() Columns
		ColumnNames() []string

		PrimaryKeys() Columns
		PrimaryKeyNames() []string
		PointInsertStatement() (string, error)
		PointReadStatement(...string) (string, error)
		TableSample(float64) (string, error)

		IsView() bool
		// IsInterleaved will return true if the table has a parent or child
		IsInterleaved() bool
		// IsApex will return true if the table is the apex of an interleaved relationship
		IsApex() bool
		// IsBottom will return true if table is the bottom of an interleaved relationship
		IsBottom() bool
		// GetApex will return the top level parent or nil if it does not exist
		GetApex() Table
		GetAllRelationNames() []string
	}

	table struct {
		n            string
		t            string
		p            string // parent name
		parent       Table
		c            string // child name
		child        Table
		spannerState string
		columns      Columns
		indexes      Indexes
	}
)

func NewTable() Table {
	return &table{
		columns: NewColumns(),
		indexes: NewIndexes(),
	}
}

func LoadTable(ctx context.Context, client *spanner.Client, s Schema, t string) error {
	iter := client.Single().Query(ctx, information.GetTableQuery(t))
	defer iter.Stop()
	err := iter.Do(func(row *spanner.Row) error {
		var ti information.Table
		if err := row.ToStruct(&ti); err != nil {
			return err
		}

		tp := NewTableFromSchema(ti)
		s.SetTable(tp)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func LoadTables(ctx context.Context, client *spanner.Client, s Schema) error {
	iter := client.Single().Query(ctx, information.ListTablesQuery())
	defer iter.Stop()

	err := iter.Do(func(row *spanner.Row) error {
		var ti information.Table
		if err := row.ToStruct(&ti); err != nil {
			return err
		}

		tp := NewTableFromSchema(ti)
		s.AddTable(tp)

		return nil
	})

	if err != nil {
		return fmt.Errorf("iterating tables: %s", err.Error())
	}

	return nil
}

func NewTableFromSchema(x information.Table) Table {
	t := NewTable()

	if x.TableName != nil {
		t.SetName(*x.TableName)
	}

	if x.TableType != nil {
		t.SetType(*x.TableType)
	}

	if x.SpannerState != nil {
		t.SetSpanenrState(*x.SpannerState)
	}

	if x.ParentTableName != nil {
		t.SetParentName(*x.ParentTableName)
	}

	return t
}

func (t *table) SetName(x string) {
	t.n = x
}

func (t *table) Name() string {
	return t.n
}

func (t *table) SetType(x string) {
	t.t = x
}

func (t *table) Type() string {
	return t.t
}

func (t *table) SetParentName(x string) {
	t.p = x
}

func (t *table) ParentName() string {
	return t.p
}

func (t *table) HasParent() bool {
	return t.p != ""
}

func (t *table) SetParent(x Table) {
	t.parent = x
}

func (t *table) Parent() Table {
	return t.parent
}

func (t *table) HasChild() bool {
	return t.c != ""
}

func (t *table) SetChildName(x string) {
	t.c = x
}

func (t *table) ChildName() string {
	return t.c
}

func (t *table) SetChild(x Table) {
	t.child = x
}

func (t *table) Child() Table {
	return t.child
}

func (t *table) SetSpanenrState(x string) {
	t.spannerState = x
}

func (t *table) SpannerState() string {
	return t.spannerState
}

func (t *table) AddColumn(x Column) {
	t.columns.AddColumn(x)
}

func (t *table) AddIndex(x Index) {
	t.indexes.AddIndex(x)
}

func (t *table) Columns() Columns {
	return t.columns
}

func (t *table) PointInsertStatement() (string, error) {
	var b strings.Builder

	cols := t.columns.ColumnNames()
	if len(cols) <= 0 {
		return "", errors.New("no columns associated with table")
	}

	fmt.Fprintf(&b, "INSERT INTO %s(%s) VALUES(", t.Name(), strings.Join(cols, ", "))

	fmt.Fprintf(&b, "@%s", cols[0])
	if len(cols) > 1 {
		for _, s := range cols[1:] {
			fmt.Fprintf(&b, ", @%s", s)
		}
	}
	b.WriteString(")")

	return b.String(), nil
}

func (t *table) PointReadStatement(predicates ...string) (string, error) {
	// TODO: Return error if predicate is not a valid column
	if len(predicates) <= 0 {
		return "", errors.New("can not generate point read without predicates")
	}

	var b strings.Builder

	cols := t.columns.ColumnNames()
	if len(cols) <= 0 {
		return "", errors.New("no columns associated with table")
	}

	fmt.Fprintf(&b, "SELECT %s FROM %s WHERE ", strings.Join(cols, ", "), t.Name())
	fmt.Fprintf(&b, "%s = @%s", predicates[0], predicates[0])
	if len(cols) > 1 {
		for _, s := range predicates[1:] {
			fmt.Fprintf(&b, " AND %s = @%s", s, s)
		}
	}

	return b.String(), nil
}

func (t *table) TableSample(x float64) (string, error) {
	pkeys := t.PrimaryKeyNames()

	if len(pkeys) <= 0 {
		return "", errors.New("no primary keys associated with table")
	}

	var b strings.Builder

	fmt.Fprintf(&b, "SELECT %s FROM %s TABLESAMPLE BERNOULLI (%f PERCENT)", strings.Join(pkeys, ", "), t.Name(), x)

	return b.String(), nil
}

func (t *table) PrimaryKeys() Columns {
	return t.columns.PrimaryKeys()
}

func (t *table) PrimaryKeyNames() []string {
	cols := t.columns.PrimaryKeys()
	ret := make([]string, 0, cols.Len())
	for cols.HasNext() {
		col := cols.GetNext()
		ret = append(ret, col.Name())
	}

	return ret
}

func (t *table) ColumnNames() []string {
	ret := make([]string, 0, t.columns.Len())
	for t.columns.HasNext() {
		c := t.columns.GetNext()
		ret = append(ret, c.Name())
	}

	t.columns.ResetIterator()

	return ret
}

func (t *table) IsView() bool {
	return t.t == "VIEW"
}

// IsInterleaved will return true if the table has a parent or child
func (t *table) IsInterleaved() bool {
	return t.HasChild() || t.HasParent()
}

// IsApex will return true if the table is the apex of an interleaved relationship
func (t *table) IsApex() bool {
	return t.Parent() == nil
}

// IsBottom will return true if the table is the bottom of an interleaved relationship (has no children)
func (t *table) IsBottom() bool {
	return t.Child() == nil
}

// GetApex returns
func (t *table) GetApex() Table {
	// If this table is the apex return it
	if t.IsInterleaved() && t.IsApex() {
		return t
	}

	p := t.Parent()
	if p == nil {
		return nil
	}

	last := p
	for ok := true; ok; ok = (p != nil) {
		last = p
		p = p.Parent()
	}

	return last
}

func (t *table) GetAllRelationNames() []string {
	apex := t.GetApex()
	ret := []string{apex.Name()}

	child := apex.Child()
	for ok := true; ok; ok = (child != nil) {
		ret = append(ret, child.Name())
		child = child.Child()
	}

	return ret
}
