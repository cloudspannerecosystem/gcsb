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
	"fmt"

	"github.com/cloudspannerecosystem/gcsb/pkg/config"
)

type (
	Schema interface {
		SetTable(Table) // TODO: Remove this
		Table() Table
		AddTable(Table)
		Tables() Tables
		GetTable(string) Table

		// Traverse should be called after a schema load operation to detect interleaved relationships
		Traverse() error
	}

	schema struct {
		table  Table // TODO: Remove this
		tables Tables
	}
)

func NewSchema() Schema {
	return &schema{
		tables: NewTables(),
	}
}

func LoadSchema(ctx context.Context, cfg *config.Config) (Schema, error) {
	client, err := cfg.Client(ctx)
	if err != nil {
		return nil, err
	}

	s := NewSchema()

	// Load Tables
	err = LoadTables(ctx, client, s)
	if err != nil {
		return nil, err
	}

	// Load Columns
	ts := s.Tables()
	for ts.HasNext() {
		t := ts.GetNext()

		// Load Columns
		err := LoadColumns(ctx, client, t)
		if err != nil {
			return nil, fmt.Errorf("loading columns for table '%s': %s", t.Name(), err.Error())
		}

		// Load Indexes
		// Load Indexes
		err = LoadIndexes(ctx, client, t)
		if err != nil {
			return nil, fmt.Errorf("loading indexes for table '%s': %s", t.Name(), err.Error())
		}
	}

	// reset iterator
	ts.ResetIterator()

	// Traverse the schema to setup parent/child relationships
	s.Traverse()

	return s, nil
}

// TODO: Get rid of this in milestone 2 (multi-table)
func LoadSingleTableSchema(ctx context.Context, cfg *config.Config, t string) (Schema, error) {
	client, err := cfg.Client(ctx)
	if err != nil {
		return nil, err
	}

	s := NewSchema()

	// Load Table
	err = LoadTable(ctx, client, s, t)
	if err != nil {
		return nil, err
	}

	tab := s.Table()

	// Load Columns
	err = LoadColumns(ctx, client, tab)
	if err != nil {
		return nil, fmt.Errorf("loading columns for table '%s': %s", tab.Name(), err.Error())
	}

	// Load Indexes
	err = LoadIndexes(ctx, client, tab)
	if err != nil {
		return nil, fmt.Errorf("loading indexes for table '%s': %s", tab.Name(), err.Error())
	}

	return s, nil
}

func (s *schema) SetTable(x Table) {
	s.table = x
}

func (s *schema) Table() Table {
	return s.table
}

func (s *schema) AddTable(x Table) {
	s.tables.AddTable(x)
}

func (s *schema) Tables() Tables {
	return s.tables
}

func (s *schema) GetTable(x string) Table {
	return s.tables.GetTable(x)
}

// Traverse should be called after a schema load operation to detect interleaved relationships
func (s *schema) Traverse() error {
	return s.tables.Traverse()
}
