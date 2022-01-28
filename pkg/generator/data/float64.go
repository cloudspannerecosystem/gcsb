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

// Assert that Float64Generator implements Generator
var _ Generator = (*Float64Generator)(nil)

type (
	Float64Generator struct {
		src *rand.Rand
		f   func() interface{}
		r   bool
		min float64
		max float64
	}
)

func NewFloat64Generator(cfg Config) (Generator, error) {
	ret := &Float64Generator{
		src: rand.New(cfg.Source()),
	}

	ret.f = ret.nextRandom
	if cfg.Range() {
		ret.f = ret.nextRanged
		ret.r = true

		switch min := cfg.Minimum().(type) {
		case float64:
			ret.min = min
		default:
			return nil, fmt.Errorf("minimum '%s' of type '%T' invalid for float64 generator", min, min)
		}

		switch max := cfg.Maximum().(type) {
		case float64:
			ret.max = max
		default:
			return nil, fmt.Errorf("maximum '%s' of type '%T' invalid for float64 generator", max, max)
		}
	}

	return ret, nil
}

func (g *Float64Generator) Next() interface{} {
	return g.f()
}

func (g *Float64Generator) nextRandom() interface{} {
	return g.src.Float64()
}

func (g *Float64Generator) nextRanged() interface{} {
	return g.min + g.src.Float64()*(g.max-g.min)
}

func (g *Float64Generator) Type() spansql.TypeBase {
	return spansql.Float64
}
