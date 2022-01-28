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

package sample

import "cloud.google.com/go/spanner/spansql"

var (
	singleTableSingleKey = spansql.CreateTable{
		Name: "Singers",
		Columns: []spansql.ColumnDef{
			{
				Name: "SingerId",
				Type: spansql.Type{
					Base: spansql.Int64,
				},
				NotNull: true,
			},
			{
				Name: "FirstName",
				Type: spansql.Type{
					Base: spansql.String,
					Len:  1024,
				},
			},
			{
				Name: "LastName",
				Type: spansql.Type{
					Base: spansql.String,
					Len:  1024,
				},
			},
			{
				Name: "BirthDate",
				Type: spansql.Type{
					Base: spansql.Date,
				},
			},
			{
				Name: "ByteField",
				Type: spansql.Type{
					Base: spansql.Bytes,
					Len:  1024,
				},
			},
			{
				Name: "FloatField",
				Type: spansql.Type{
					Base: spansql.Float64,
				},
			},
			{
				Name: "ArrayField",
				Type: spansql.Type{
					Array: true,
					Base:  spansql.Int64,
				},
			},
			{
				Name: "TSField",
				Type: spansql.Type{
					Base: spansql.Timestamp,
				},
			},
			{
				Name: "NumericField",
				Type: spansql.Type{
					Base: spansql.Numeric,
				},
			},
		},
		PrimaryKey: []spansql.KeyPart{
			{Column: "SingerId"},
		},
	}
)
