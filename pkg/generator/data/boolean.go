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

// Assert that BooleanGenerator implements Generator
var _ Generator = (*BooleanGenerator)(nil)

type (
	BooleanGenerator struct {
		src       rand.Source
		cache     int64
		remaining int
		f         func() interface{}
		v         bool
	}
)

func NewBooleanGenerator(cfg Config) (Generator, error) {
	if cfg == nil {
		cfg = NewConfig()
	}

	ret := &BooleanGenerator{
		src: cfg.Source(),
	}

	// Generate  values by defualt

	// If the config is for a static value, return it instead
	ret.f = ret.next
	if cfg.Static() {
		v, ok := cfg.Value().(bool)
		if !ok {
			return nil, fmt.Errorf("value '%s' of type '%T' is invalid for static bool generation", cfg.Value(), cfg.Value())
		}
		ret.v = v
		ret.f = ret.nextStatic
	}

	return ret, nil
}

func (b *BooleanGenerator) Next() interface{} {
	return b.f()
}

// nextStatic returns a static bool value
func (b *BooleanGenerator) nextStatic() interface{} {
	return b.v
}

// next returns a ly generated bool value
func (b *BooleanGenerator) next() interface{} {

	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 63
	}

	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--

	return result
}

func (b *BooleanGenerator) Type() spansql.TypeBase {
	return spansql.Bool
}
