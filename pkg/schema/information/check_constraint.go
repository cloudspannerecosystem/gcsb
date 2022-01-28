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
	// CheckConstraints is a collection of CheckConstraint
	CheckConstraints []*CheckConstraint

	// CheckConstraint is a row from information_schema.check_constraints (see: https://cloud.google.com/spanner/docs/information-schema#information_schemacheck_constraints)
	CheckConstraint struct {
		// The name of the constraint's catalog. This column is never null, but always an empty string.
		ConstraintCatalog string `spanner:"CONSTRAINT_CATALOG"`
		// The name of the constraint's schema. An empty string if unnamed.
		ConstraintSchema string `spanner:"CONSTRAINT_SCHEMA"`
		// The name of the constraint. This column is never null. If not explicitly specified in the schema definition, a system-defined name is assigned.
		ConstraintName string `spanner:"CONSTRAINT_NAME"`
		// The expressions of the CHECK constraint. This column is never null.
		CheckClause string `spanner:"CHECK_CLAUSE"`
		// The current state of the CHECK constraint. This column is never null. The possible states are as follows:
		//   VALIDATING: Cloud Spanner is validating the existing data.
		//   COMMITTED: There is no active schema change for this constraint.
		SpannerState string `spanner:"SPANNER_STATE"`
	}
)
