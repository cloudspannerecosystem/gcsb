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
	"fmt"
	"math/rand"

	"cloud.google.com/go/spanner/spansql"
)

// Assert that Int64Generator implements Generator
var _ Generator = (*Int64Generator)(nil)

type (
	Int64Generator struct {
		src *rand.Rand
		f   func() interface{}
		r   bool
		min int64
		max int64
	}
)

func NewInt64Generator(cfg Config) (Generator, error) {
	ret := &Int64Generator{
		src: rand.New(cfg.Source()),
	}

	ret.f = ret.nextRandom
	if cfg.Range() {
		ret.r = true

		switch min := cfg.Minimum().(type) {
		case int:
			ret.min = int64(min)
		case int64:
			ret.min = min
		default:
			return nil, fmt.Errorf("minimum '%s' of type '%T' invalid for int64 generator", min, min)
		}

		switch max := cfg.Maximum().(type) {
		case int:
			ret.max = int64(max)
		case int64:
			ret.max = max
		default:
			return nil, fmt.Errorf("maximum '%s' of type '%T' invalid for int64 generator", max, max)
		}

		ret.f = ret.nextRanged
	}

	return ret, nil
}

func (g *Int64Generator) Next() interface{} {
	return g.f()
}

func (g *Int64Generator) nextRandom() interface{} {
	return g.src.Int63()
}

func (g *Int64Generator) nextRanged() interface{} {
	return g.src.Int63n(g.max-g.min) + g.min
}

func (g *Int64Generator) Type() spansql.TypeBase {
	return spansql.Int64
}
