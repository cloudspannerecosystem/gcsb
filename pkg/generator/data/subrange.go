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
	"math/rand"

	"cloud.google.com/go/spanner/spansql"
)

// Assert that SubRangeGenerator implements Generator
var _ Generator = (*SubRangeGenerator)(nil)

type (
	// SubRangeGenerator randomly chooses a generator and returns it's Next() value
	SubRangeGenerator struct {
		src        *rand.Rand
		generators []Generator
	}
)

func NewSubRangeGenerator(cfg Config) (*SubRangeGenerator, error) {
	return &SubRangeGenerator{
		src:        rand.New(cfg.Source()),
		generators: make([]Generator, 0),
	}, nil
}

func (g *SubRangeGenerator) Next() interface{} {
	return g.generators[g.src.Intn(len(g.generators))].Next()
}

func (g *SubRangeGenerator) AddGenerator(x Generator) {
	g.generators = append(g.generators, x)
}

func (g *SubRangeGenerator) Type() spansql.TypeBase {
	if len(g.generators) < 1 {
		panic("Can not determine type base for subrange generator that has no generators")
	}

	return g.generators[0].Type()
}
