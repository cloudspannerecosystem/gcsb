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

import "fmt"

type (
	TableIterator interface {
		ResetIterator()
		HasNext() bool
		GetNext() Table
	}

	Tables interface {
		TableIterator
		Len() int
		Tables() []Table
		AddTable(Table)
		GetTable(string) Table
		Traverse() error
	}

	tables struct {
		iteratorIndex int
		tables        []Table
	}
)

func NewTables() Tables {
	return &tables{
		tables: make([]Table, 0),
	}
}

func (t *tables) Len() int {
	return len(t.tables)
}

func (t *tables) ResetIterator() {
	t.iteratorIndex = 0
}

func (t *tables) HasNext() bool {
	return t.iteratorIndex < len(t.tables)
}

func (t *tables) GetNext() Table {
	if t.HasNext() {
		to := t.tables[t.iteratorIndex]
		t.iteratorIndex++
		return to
	}

	return nil
}

func (t *tables) Tables() []Table {
	return t.tables
}

func (t *tables) AddTable(x Table) {
	t.tables = append(t.tables, x)
}

func (t *tables) GetTable(x string) Table {
	for _, tab := range t.tables {
		if tab.Name() == x {
			return tab
		}
	}

	return nil
}

func (t *tables) Traverse() error {
	// Iterate over tables setting parental relationships
	for _, child := range t.tables {
		if child.ParentName() != "" {
			// fetch the parent table
			parent := t.GetTable(child.ParentName())
			if parent == nil {
				return fmt.Errorf("table '%s' references a parent table '%s' that is not in information schema", child.Name(), child.ParentName())
			}

			// Set parent as this tables parent
			child.SetParent(parent)

			// Set parents child
			parent.SetChildName(child.Name())
			parent.SetChild(child)
		}
	}

	return nil
}
