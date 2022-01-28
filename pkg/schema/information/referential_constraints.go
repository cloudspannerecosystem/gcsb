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
	// ReferentialConstraints is a collection of ReferentialConstraint
	ReferentialConstraints []*ReferentialConstraint

	// ReferentialConstraint is a row from information_schema.referential_constraints (see: https://cloud.google.com/spanner/docs/information-schema#referential-constraints)
	ReferentialConstraint struct {
		// The name of the FOREIGN KEY's catalog. Always an empty string.
		ConstraintCatalog string `spanner:"CONSTRAINT_CATALOG"`
		// The name of the FOREIGN KEY's schema. An empty string if unnamed.
		ConstraintSchema string `spanner:"CONSTRAINT_SCHEMA"`
		// The name of the FOREIGN KEY.
		ConstraintName string `spanner:"CONSTRAINT_NAME"`
		// The catalog name of the PRIMARY KEY or UNIQUE constraint the FOREIGN KEY references. Always an empty string.
		UniqueConstraintCatalog string `spanner:"UNIQUE_CONSTRAINT_CATALOG"`
		// The schema name of the PRIMARY KEY or UNIQUE constraint the FOREIGN KEY references. An empty string if unnamed.
		UniqueConstraintSchema string `spanner:"UNIQUE_CONSTRAINT_SCHEMA"`
		// The name of the PRIMARY KEY or UNIQUE constraint the FOREIGN KEY references.
		UniqueConstraintName string `spanner:"UNIQUE_CONSTRAINT_NAME"`
		// Always SIMPLE.
		MatchOption string `spanner:"MATCH_OPTION"`
		// Always NO ACTION.
		UpdateRule string `spanner:"UPDATE_RULE"`
		// Always NO ACTION.
		DeleteRule string `spanner:"DELETE_RULE"`
		// The current state of the foreign key. Spanner does not begin enforcing the constraint until the foreign key's backing indexes are created and backfilled. Once the indexes are ready, Spanner begins enforcing the constraint for new transactions while it validates the existing data. Possible values and the states they represent are:
		//   BACKFILLING_INDEXES: indexes are being backfilled.
		//   VALIDATING_DATA: existing data and new writes are being validated.
		//   WAITING_FOR_COMMIT: the foreign key bulk operations have completed successfully, or none were needed, but the foreign key is still pending.
		//   COMMITTED: the schema change was committed.
		SpannerState string `spanner:"SPANNER_STATE"`
	}
)
