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
	"math/big"
	"math/rand"

	"cloud.google.com/go/spanner/spansql"
)

// Assert that NumericGenerator implements Generator
var _ Generator = (*NumericGenerator)(nil)

type (
	NumericGenerator struct {
		src *rand.Rand
	}
)

func NewNumericGenerator(cfg Config) (Generator, error) {
	return &NumericGenerator{
		src: rand.New(cfg.Source()),
	}, nil
}

func (g *NumericGenerator) Next() interface{} {
	return big.NewRat(g.src.Int63(), g.src.Int63())
}

func (g *NumericGenerator) Type() spansql.TypeBase {
	return spansql.Numeric
}
