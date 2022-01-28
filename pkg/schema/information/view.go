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
	// Views is a collection of View
	Views []*View

	// View is a row in information_schema.views (see: https://cloud.google.com/spanner/docs/information-schema#information_schemaviews)
	View struct {
		// The name of the catalog. Always an empty string.
		TableCatalog string `spanner:"TABLE_CATALOG"`
		// The name of the schema. An empty string if unnamed.
		TableSchema string `spanner:"TABLE_SCHEMA"`
		// The name of the view.
		TableName string `spanner:"TABLE_NAME"`
		// The SQL text of the query that defines the view.
		ViewDefinition string `spanner:"VIEW_DEFINITION"`
	}
)
