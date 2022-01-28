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
	// Indexes is a collection of Index
	Indexes []*Index

	// Index is a row in information_schema.indexes (see: https://cloud.google.com/spanner/docs/information-schema#indexes)
	Index struct {
		// The name of the catalog. Always an empty string.
		TableCatalog string `spanner:"TABLE_CATALOG"`
		// The name of the schema. An empty string if unnamed.
		TableSchema string `spanner:"TABLE_SCHEMA"`
		// The name of the table.
		TableName string `spanner:"TABLE_NAME"`
		// The name of the index. Tables with a PRIMARY KEY specification have a pseudo-index entry generated with the name PRIMARY_KEY, which allows the fields of the primary key to be determined.
		IndexName string `spanner:"INDEX_NAME"`
		// The type of the index. The type is INDEX or PRIMARY_KEY.
		IndexType string `spanner:"INDEX_TYPE"`
		// Secondary indexes can be interleaved in a parent table, as discussed in Creating a secondary index. This column holds the name of that parent table, or NULL if the index is not interleaved.
		ParentTableName string `spanner:"PARENT_TABLE_NAME"`
		// Whether the index keys must be unique.
		IsUnique bool `spanner:"IS_UNIQUE"`
		// Whether the index includes entries with NULL values.
		IsNullFiltered bool `spanner:"IS_NULL_FILTERED"`
		// The current state of the index. Possible values and the states they represent are:
		//   PREPARE: creating empty tables for a new index.
		//   WRITE_ONLY: backfilling data for a new index.
		//   WRITE_ONLY_CLEANUP: cleaning up a new index.
		//   WRITE_ONLY_VALIDATE_UNIQUE: checking uniqueness of data in a new index.
		//   READ_WRITE: normal index operation.
		IndexState string `spanner:"INDEX_STATE"`
		// True if the index is managed by Cloud Spanner; Otherwise, False. Secondary backing indexes for foreign keys are managed by Cloud Spanner.
		SpannerIsManaged bool `spanner:"SPANNER_IS_MANAGED"`
	}
)

// sql query
const getIndexSqlstr = `SELECT ` +
	`INDEX_NAME, IS_UNIQUE, IS_NULL_FILTERED, INDEX_STATE ` +
	`FROM INFORMATION_SCHEMA.INDEXES ` +
	`WHERE TABLE_SCHEMA = "" ` +
	`AND INDEX_NAME != "PRIMARY_KEY" ` +
	`AND TABLE_NAME = @table_name ` +
	`AND SPANNER_IS_MANAGED = FALSE `

var (
	// getIndexesForTableQuery renders a query for fetching all indexes for a table from information_schema.indexes
	getIndexesForTableQuery = spansql.Query{
		Select: spansql.Select{
			List: []spansql.Expr{
				spansql.ID("TABLE_CATALOG"),
				spansql.ID("TABLE_SCHEMA"),
				spansql.ID("TABLE_NAME"),
				spansql.ID("INDEX_NAME"),
				spansql.ID("INDEX_TYPE"),
				spansql.ID("PARENT_TABLE_NAME"),
				spansql.ID("IS_UNIQUE"),
				spansql.ID("IS_NULL_FILTERED"),
				spansql.ID("INDEX_STATE"),
				spansql.ID("SPANNER_IS_MANAGED"),
			},
			From: []spansql.SelectFrom{
				spansql.SelectFromTable{
					Table: "information_schema.indexes",
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
		Order: []spansql.Order{
			{Expr: spansql.ID("INDEX_NAME")},
		},
	}
)

// GetIndexesQuery returns a spanner statement for fetching index information for a table
func GetIndexesQuery(table string) spanner.Statement {
	// st := spanner.NewStatement(getIndexesForTableQuery.SQL())
	st := spanner.NewStatement(getIndexSqlstr)
	st.Params["table_name"] = table

	return st
}
