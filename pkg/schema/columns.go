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

type (
	ColumnIterator interface {
		ResetIterator()
		HasNext() bool
		GetNext() Column
	}

	Columns interface {
		ColumnIterator
		Columns() []Column
		AddColumn(Column)
		ColumnNames() []string
		PrimaryKeys() Columns
		Len() int
	}

	columns struct {
		iteratorIndex int
		columns       []Column
	}
)

func NewColumns() Columns {
	return &columns{
		columns: make([]Column, 0),
	}
}

func (c *columns) Len() int {
	return len(c.columns)
}

func (c *columns) ResetIterator() {
	c.iteratorIndex = 0
}

func (c *columns) HasNext() bool {
	return c.iteratorIndex < len(c.columns)
}

func (c *columns) GetNext() Column {
	if c.HasNext() {
		to := c.columns[c.iteratorIndex]
		c.iteratorIndex++
		return to
	}

	return nil
}

func (c *columns) Columns() []Column {
	return c.columns
}

func (c *columns) AddColumn(x Column) {
	c.columns = append(c.columns, x)
}

func (c *columns) ColumnNames() []string {
	ret := make([]string, 0)
	for _, col := range c.columns {
		ret = append(ret, col.Name())
	}

	return ret
}

func (c *columns) PrimaryKeys() Columns {
	ret := NewColumns()

	for _, col := range c.columns {
		if col.PrimaryKey() {
			ret.AddColumn(col)
		}
	}

	return ret
}
