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
	// SpannerStatistics is a collection of SpannerStatistic
	SpannerStatistics []*SpannerStatistic

	// SpannerStatistic is a row in information_schema.spanner_statistics (see https://cloud.google.com/spanner/docs/information-schema#information_schemaspanner_statistics)
	SpannerStatistic struct {
		// The name of the catalog. Always an empty string.
		CatalogName string `spanner:"CATALOG_NAME"`
		// The name of the schema. The name is empty for the default schema and non-empty for other schemas (for example, the INFORMATION_SCHEMA itself). This column is never null.
		SchemaName string `spanner:"SCHEMA_NAME"`
		// The name of the statistics package.
		PackageName string `spanner:"PACKAGE_NAME"`
		// False if the statistics package is exempted from garbage collection; Otherwise, True.
		// This attribute must be set to False in order to reference the statistics package in a hint or through client API.
		AllowGC bool `spanner:"ALLOW_GC"`
	}
)
