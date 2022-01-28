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

package data

import (
	"cloud.google.com/go/spanner"
	"cloud.google.com/go/spanner/spansql"
)

// Assert that CommitTimestampGenerator implements Generator
var _ Generator = (*CommitTimestampGenerator)(nil)

type (
	CommitTimestampGenerator struct{}
)

func NewCommitTimestampGenerator(cfg Config) (Generator, error) {
	return &CommitTimestampGenerator{}, nil
}

func (g *CommitTimestampGenerator) Next() interface{} {
	return spanner.CommitTimestamp
}

func (g *CommitTimestampGenerator) Type() spansql.TypeBase {
	return spansql.Timestamp
}
