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

package information

import (
	"cloud.google.com/go/spanner"
	"cloud.google.com/go/spanner/spansql"
)

type (
	// Tables is a collection of tables
	Tables []*Table

	// Table is a row in information_schema.tables (see: https://cloud.google.com/spanner/docs/information-schema#information_schematables)
	Table struct {
		// The name of the catalog. Always an empty string.
		TableCatalog *string `spanner:"TABLE_CATALOG"`
		// The name of the schema. An empty string if unnamed.
		TableSchema *string `spanner:"TABLE_SCHEMA"`
		// The name of the table or view.
		TableName *string `spanner:"TABLE_NAME"`
		// The type of the table. For tables it has the value BASE TABLE; for views it has the value VIEW.
		TableType *string `spanner:"TABLE_TYPE"`
		// The name of the parent table if this table is interleaved, or NULL.
		ParentTableName *string `spanner:"PARENT_TABLE_NAME"`
		// This is set to CASCADE or NO ACTION for interleaved tables, and NULL otherwise. See TABLE statements for more information.
		OnDeleteAction *string `spanner:"ON_DELETE_ACTION"`
		// A table can go through multiple states during creation, if bulk operations are involved. For example, when the table is created with a foreign key that requires backfilling of its indexes. Possible states are:
		//   ADDING_FOREIGN_KEY: Adding the table's foreign keys.
		//   WAITING_FOR_COMMIT: Finalizing the schema change.
		//   COMMITTED: The schema change to create the table has been committed. You cannot write to the table until the change is committed.
		SpannerState *string `spanner:"SPANNER_STATE"`
	}
)

var (
	// listTablesQuery renders a query for fetching all tables from information_schema.tables
	listTablesQuery = spansql.Query{
		Select: spansql.Select{
			List: []spansql.Expr{
				spansql.ID("TABLE_CATALOG"),
				spansql.ID("TABLE_SCHEMA"),
				spansql.ID("TABLE_NAME"),
				spansql.ID("TABLE_TYPE"),
				spansql.ID("PARENT_TABLE_NAME"),
				spansql.ID("ON_DELETE_ACTION"),
				spansql.ID("SPANNER_STATE"),
			},
			From: []spansql.SelectFrom{
				spansql.SelectFromTable{
					Table: "information_schema.tables",
				},
			},
			Where: spansql.LogicalOp{
				Op: spansql.And,
				LHS: spansql.ComparisonOp{
					Op:  spansql.Eq,
					LHS: spansql.ID("table_catalog"),
					RHS: spansql.StringLiteral(""),
				},
				RHS: spansql.ComparisonOp{
					Op:  spansql.Eq,
					LHS: spansql.ID("table_schema"),
					RHS: spansql.StringLiteral(""),
				},
			},
		},
		Order: []spansql.Order{
			{Expr: spansql.ID("table_catalog")},
			{Expr: spansql.ID("table_schema")},
			{Expr: spansql.ID("table_name")},
		},
	}

	// getTableQuery fetches a table from information_schema.tables
	getTableQuery = spansql.Query{
		Select: spansql.Select{
			List: []spansql.Expr{
				spansql.ID("TABLE_CATALOG"),
				spansql.ID("TABLE_SCHEMA"),
				spansql.ID("TABLE_NAME"),
				spansql.ID("TABLE_TYPE"),
				spansql.ID("PARENT_TABLE_NAME"),
				spansql.ID("ON_DELETE_ACTION"),
				spansql.ID("SPANNER_STATE"),
			},
			From: []spansql.SelectFrom{
				spansql.SelectFromTable{
					Table: "information_schema.tables",
				},
			},
			Where: spansql.LogicalOp{
				Op: spansql.And,
				LHS: spansql.ComparisonOp{
					Op:  spansql.Eq,
					LHS: spansql.ID("table_schema"),
					RHS: spansql.StringLiteral(""),
				},
				RHS: spansql.ComparisonOp{
					Op:  spansql.Eq,
					LHS: spansql.ID("table_name"),
					RHS: spansql.Param("table_name"),
				},
			},
		},
	}
)

// ListTablesQuery will construct a query for fetching all tables from information_schema
// TODO: Support schema?
func ListTablesQuery() spanner.Statement {
	return spanner.NewStatement(listTablesQuery.SQL())
}

func GetTableQuery(table string) spanner.Statement {
	st := spanner.NewStatement(getTableQuery.SQL())
	st.Params["table_name"] = table

	return st
}
