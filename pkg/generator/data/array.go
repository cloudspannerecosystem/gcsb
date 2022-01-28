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
	"errors"
	"math/big"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner/spansql"
)

var (
	// Assert that ArrayGenerator implements Generator
	_ Generator = (*ArrayGenerator)(nil)
)

type (
	ArrayGenerator struct {
		g Generator
		l int
	}
)

func NewArrayGenerator(cfg Config) (Generator, error) {
	if cfg.Generator() == nil {
		return nil, errors.New("array generator requires a generator")
	}

	if cfg.Length() <= 0 {
		return nil, errors.New("array generator length must be <= 0")
	}

	ret := &ArrayGenerator{
		g: cfg.Generator(),
		l: cfg.Length(),
	}

	return ret, nil
}

func (g *ArrayGenerator) Next() interface{} {
	var ret interface{}
	switch g.g.Type() {
	case spansql.Bool:
		return g.nextBool()
	case spansql.String:
		return g.nextString()
	case spansql.Int64:
		return g.nextInt64()
	case spansql.Float64:
		return g.nextFloat64()
	case spansql.Bytes:
		return g.nextBytes()
	case spansql.Timestamp:
		return g.nextTimestamp()
	case spansql.Date:
		return g.nextDate()
	case spansql.Numeric:
		return g.nextNumeric()
	}

	return ret
}

func (g *ArrayGenerator) nextNumeric() []*big.Rat {
	ret := make([]*big.Rat, 0, g.l)
	for i := 0; i < g.l; i++ {
		v, _ := g.g.Next().(*big.Rat)
		ret = append(ret, v)
	}
	return ret
}

func (g *ArrayGenerator) nextDate() []civil.Date {
	ret := make([]civil.Date, 0, g.l)
	for i := 0; i < g.l; i++ {
		v, _ := g.g.Next().(civil.Date)
		ret = append(ret, v)
	}
	return ret
}

func (g *ArrayGenerator) nextTimestamp() []time.Time {
	ret := make([]time.Time, 0, g.l)
	for i := 0; i < g.l; i++ {
		v, _ := g.g.Next().(time.Time)
		ret = append(ret, v)
	}
	return ret
}

func (g *ArrayGenerator) nextBytes() []byte {
	ret := make([]byte, 0, g.l)
	for i := 0; i < g.l; i++ {
		v, _ := g.g.Next().(byte)
		ret = append(ret, v)
	}
	return ret
}

func (g *ArrayGenerator) nextFloat64() []float64 {
	ret := make([]float64, 0, g.l)
	for i := 0; i < g.l; i++ {
		v, _ := g.g.Next().(float64)
		ret = append(ret, v)
	}
	return ret
}

func (g *ArrayGenerator) nextInt64() []int64 {
	ret := make([]int64, 0, g.l)
	for i := 0; i < g.l; i++ {
		v, _ := g.g.Next().(int64)
		ret = append(ret, v)
	}
	return ret
}

func (g *ArrayGenerator) nextString() []string {
	ret := make([]string, 0, g.l)
	for i := 0; i < g.l; i++ {
		v, _ := g.g.Next().(string)
		ret = append(ret, v)
	}
	return ret
}

func (g *ArrayGenerator) nextBool() []bool {
	ret := make([]bool, 0, g.l)
	for i := 0; i < g.l; i++ {
		v, _ := g.g.Next().(bool)
		ret = append(ret, v)
	}
	return ret
}

func (g *ArrayGenerator) Type() spansql.TypeBase {
	return g.g.Type()
}
