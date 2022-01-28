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
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/sample"
)

var (
	_ Generator = (*CompositeKey)(nil)
	_ Generator = (*SingleKey)(nil)
)

type (
	// Composite key returns a spanner key containing multiple static values
	CompositeKey struct {
		g *sample.SampleGenerator
	}

	// SingleKey is used to return a spanner key containing only 1 static value
	SingleKey struct {
		g Generator
	}
)

func (g *CompositeKey) Next() interface{} {
	return g.g.Next()
}

func (g *CompositeKey) Type() spansql.TypeBase {
	return 0
}

func (g *SingleKey) Next() interface{} {
	return spanner.Key{g.g.Next()}
}

func (g *SingleKey) Type() spansql.TypeBase {
	return 0
}
