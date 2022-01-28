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
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner/spansql"
)

var (
	_ Generator = (*StaticBooleanGenerator)(nil)
	_ Generator = (*StaticStringGenerator)(nil)
	_ Generator = (*StaticBytesGenerator)(nil)
	_ Generator = (*StaticInt64Generator)(nil)
	_ Generator = (*StaticFloat64Generator)(nil)
	_ Generator = (*StaticNumericGenerator)(nil)
	_ Generator = (*StaticDateGenerator)(nil)
	_ Generator = (*StaticTimestampGenerator)(nil)
)

type (
	//
	StaticBooleanGenerator struct {
		src  *rand.Rand
		vals []bool
		l    int
	}

	//
	StaticStringGenerator struct {
		src  *rand.Rand
		vals []string
		l    int
	}

	//
	StaticBytesGenerator struct {
		src  *rand.Rand
		vals [][]byte
		l    int
	}

	//
	StaticInt64Generator struct {
		src  *rand.Rand
		vals []int64
		l    int
	}

	//
	StaticFloat64Generator struct {
		src  *rand.Rand
		vals []float64
		l    int
	}

	//
	StaticNumericGenerator struct {
		src  *rand.Rand
		vals []*big.Rat
		l    int
	}

	//
	StaticDateGenerator struct {
		src  *rand.Rand
		vals []civil.Date
		l    int
	}

	//
	StaticTimestampGenerator struct {
		src  *rand.Rand
		vals []time.Time
		l    int
	}
)

func NewStaticBooleanGenerator(cfg Config, vals []bool) (Generator, error) {
	ret := &StaticBooleanGenerator{
		src:  rand.New(cfg.Source()),
		vals: vals,
		l:    len(vals),
	}

	return ret, nil
}

func NewStaticStringGenerator(cfg Config, vals []string) (Generator, error) {
	ret := &StaticStringGenerator{
		src:  rand.New(cfg.Source()),
		vals: vals,
		l:    len(vals),
	}

	return ret, nil
}

func NewStaticBytesGenerator(cfg Config, vals [][]byte) (Generator, error) {
	ret := &StaticBytesGenerator{
		src:  rand.New(cfg.Source()),
		vals: vals,
		l:    len(vals),
	}

	return ret, nil
}

func NewStaticInt64Generator(cfg Config, vals []int64) (Generator, error) {
	ret := &StaticInt64Generator{
		src:  rand.New(cfg.Source()),
		vals: vals,
		l:    len(vals),
	}

	return ret, nil
}

func NewStaticFloat64Generator(cfg Config, vals []float64) (Generator, error) {
	ret := &StaticFloat64Generator{
		src:  rand.New(cfg.Source()),
		vals: vals,
		l:    len(vals),
	}

	return ret, nil
}

func NewStaticNumericGenerator(cfg Config, vals []*big.Rat) (Generator, error) {
	ret := &StaticNumericGenerator{
		src:  rand.New(cfg.Source()),
		vals: vals,
		l:    len(vals),
	}

	return ret, nil
}

func NewStaticDateGenerator(cfg Config, vals []civil.Date) (Generator, error) {
	ret := &StaticDateGenerator{
		src:  rand.New(cfg.Source()),
		vals: vals,
		l:    len(vals),
	}

	return ret, nil
}

func NewStaticTimestampGenerator(cfg Config, vals []time.Time) (Generator, error) {
	ret := &StaticTimestampGenerator{
		src:  rand.New(cfg.Source()),
		vals: vals,
		l:    len(vals),
	}

	return ret, nil
}

func (g *StaticBooleanGenerator) Type() spansql.TypeBase {
	return spansql.Bool
}

func (g *StaticBooleanGenerator) Next() interface{} {
	return g.vals[rand.Intn(g.l)]
}

func (g *StaticStringGenerator) Type() spansql.TypeBase {
	return spansql.String
}

func (g *StaticStringGenerator) Next() interface{} {
	return g.vals[rand.Intn(g.l)]
}

func (g *StaticBytesGenerator) Type() spansql.TypeBase {
	return spansql.Bytes
}

func (g *StaticBytesGenerator) Next() interface{} {
	return g.vals[rand.Intn(g.l)]
}

func (g *StaticInt64Generator) Type() spansql.TypeBase {
	return spansql.Int64
}

func (g *StaticInt64Generator) Next() interface{} {
	return g.vals[rand.Intn(g.l)]
}

func (g *StaticFloat64Generator) Type() spansql.TypeBase {
	return spansql.Float64
}

func (g *StaticFloat64Generator) Next() interface{} {
	return g.vals[rand.Intn(g.l)]
}

func (g *StaticNumericGenerator) Type() spansql.TypeBase {
	return spansql.Numeric
}

func (g *StaticNumericGenerator) Next() interface{} {
	return g.vals[rand.Intn(g.l)]
}

func (g *StaticDateGenerator) Type() spansql.TypeBase {
	return spansql.Date
}

func (g *StaticDateGenerator) Next() interface{} {
	return g.vals[rand.Intn(g.l)]
}

func (g *StaticTimestampGenerator) Type() spansql.TypeBase {
	return spansql.Timestamp
}

func (g *StaticTimestampGenerator) Next() interface{} {
	return g.vals[rand.Intn(g.l)]
}
