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

type (
	// ColumnOptions is a collection of ColumnOption
	ColumnOptions []*ColumnOption

	// ColumnOption is a row from information_schema.column_options (see: https://cloud.google.com/spanner/docs/information-schema#information_schemacolumn_options)
	ColumnOption struct {
		// The name of the catalog. Always an empty string.
		TableCatalog string `spanner:"TABLE_CATALOG"`
		// The name of the schema. The name is empty for the default schema and non-empty for other schemas (for example, the INFORMATION_SCHEMA itself). This column is never null.
		TableSchema string `spanner:"TABLE_SCHEMA"`
		// The name of the table.
		TableName string `spanner:"TABLE_NAME"`
		// The name of the column.
		ColumnName string `spanner:"COLUMN_NAME"`
		// A SQL identifier that uniquely identifies the option. This identifier is the key of the OPTIONS clause in DDL.
		OptionName string `spanner:"OPTION_NAME"`
		// A data type name that is the type of this option value.
		OptionType string `spanner:"OPTION_TYPE"`
		// A SQL literal describing the value of this option. The value of this column must be parsable as part of a query. The expression resulting from parsing the value must be castable to OPTION_TYPE. This column is never null.
		OptionValue string `spanner:"OPTION_VALUE"`
	}
)
