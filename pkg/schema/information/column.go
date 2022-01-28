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
	// Columns is a collection of Collumns
	Columns []*Column

	// Column is a row from information_schema.columns (see: https://cloud.google.com/spanner/docs/information-schema#information_schemacolumns)
	Column struct {
		// The name of the catalog. Always an empty string.
		TableCatalog *string `spanner:"TABLE_CATALOG"`
		// The name of the schema. An empty string if unnamed.
		TableSchema *string `spanner:"TABLE_SCHEMA"`
		// The name of the table.
		TableName *string `spanner:"TABLE_NAME"`
		// The name of the column.
		ColumnName *string `spanner:"COLUMN_NAME"`
		// The ordinal position of the column in the table, starting with a value of 1.
		OrdinalPosition int64 `spanner:"ORDINAL_POSITION"`
		// Included to satisfy the SQL standard. Always NULL.
		ColumnDefault []byte `spanner:"COLUMN_DEFAULT"`
		// Included to satisfy the SQL standard. Always NULL.
		DataType *string `spanner:"DATA_TYPE"`
		// A string that indicates whether the column is nullable. In accordance with the SQL standard, the string is either YES or NO, rather than a Boolean value.
		IsNullable *string `spanner:"IS_NULLABLE"`
		// The data type of the column.
		SpannerType *string `spanner:"SPANNER_TYPE"`
		// A string that indicates whether the column is generated. The string is either ALWAYS for a generated column or NEVER for a non-generated column.
		IsGenerated bool `spanner:"IS_GENERATED"`
		// A string representing the SQL expression of a generated column. NULL if the column is not a generated column.
		GenerationExpression *string `spanner:"GENERATION_EXPRESSION"`
		// A string that indicates whether the generated column is stored. The string is always YES for generated columns, and NULL for non-generated columns.
		IsStored *string `spanner:"IS_STORED"`
		// The current state of the column. A new stored generated column added to an existing table may go through multiple user-observable states before it is fully usable. Possible values are:
		//   WRITE_ONLY: The column is being backfilled. No read is allowed.
		// 	 COMMITTED: The column is fully usable.
		SpannerState         *string `spanner:"SPANNER_STATE"`
		IsPrimaryKey         bool    `spanner:"IS_PRIMARY_KEY"`
		AllowCommitTimestamp bool    `spanner:"ALLOW_COMMIT_TIMESTAMP"`
	}
)

const getColSqlStr = `
SELECT
  c.COLUMN_NAME,
  c.ORDINAL_POSITION,
  c.IS_NULLABLE,
  c.SPANNER_TYPE,
  c.SPANNER_STATE,
  EXISTS (
  SELECT
    1
  FROM
    INFORMATION_SCHEMA.INDEX_COLUMNS ic
  WHERE
    ic.TABLE_SCHEMA = ""
    AND ic.TABLE_NAME = c.TABLE_NAME
    AND ic.COLUMN_NAME = c.COLUMN_NAME
    AND ic.INDEX_NAME = "PRIMARY_KEY" ) IS_PRIMARY_KEY,
  IS_GENERATED = "ALWAYS" AS IS_GENERATED,
  EXISTS (
    SELECT
        1
    FROM
         INFORMATION_SCHEMA.COLUMN_OPTIONS co
    WHERE
        co.TABLE_SCHEMA = ""
        AND co.TABLE_NAME = c.TABLE_NAME
        AND co.COLUMN_NAME = c.COLUMN_NAME
        AND co.OPTION_NAME = "allow_commit_timestamp"
        AND co.OPTION_VALUE = "TRUE"
  ) ALLOW_COMMIT_TIMESTAMP
FROM
  INFORMATION_SCHEMA.COLUMNS c
WHERE
  c.TABLE_SCHEMA = ""
  AND c.TABLE_NAME = @table_name
ORDER BY
  c.ORDINAL_POSITION
`

var (
	// getColumnsForTableQuery renders a query for fetching all columns for a table from information_schema.columns
	getColumnsForTableQuery = spansql.Query{
		Select: spansql.Select{
			List: []spansql.Expr{
				spansql.ID("TABLE_CATALOG"),
				spansql.ID("TABLE_SCHEMA"),
				spansql.ID("TABLE_NAME"),
				spansql.ID("COLUMN_NAME"),
				spansql.ID("ORDINAL_POSITION"),
				spansql.ID("COLUMN_DEFAULT"),
				spansql.ID("DATA_TYPE"),
				spansql.ID("IS_NULLABLE"),
				spansql.ID("SPANNER_TYPE"),
				spansql.ID("IS_GENERATED"),
				spansql.ID("GENERATION_EXPRESSION"),
				spansql.ID("IS_STORED"),
				spansql.ID("SPANNER_STATE"),
			},
			From: []spansql.SelectFrom{
				spansql.SelectFromTable{
					Table: "information_schema.columns",
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
			{Expr: spansql.ID("ORDINAL_POSITION")},
			{Expr: spansql.ID("COLUMN_NAME")},
		},
	}
)

// GetColumnsQuery returns a spanner statement for fetching column information for a table
func GetColumnsQuery(table string) spanner.Statement {
	// st := spanner.NewStatement(getColumnsForTableQuery.SQL())
	st := spanner.NewStatement(getColSqlStr)
	st.Params["table_name"] = table

	return st
}
