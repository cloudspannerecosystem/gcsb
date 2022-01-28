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
	"math/rand"

	"cloud.google.com/go/spanner/spansql"
)

// Assert that RandomByteGenerator implements Generator
var _ Generator = (*RandomByteGenerator)(nil)

type (
	RandomByteGenerator struct {
		len int
		src *rand.Rand
	}
)

func NewRandomByteGenerator(cfg Config) (Generator, error) {
	if cfg.Length() <= 0 {
		return nil, errors.New("string length must be > 0")
	}

	ret := &RandomByteGenerator{
		len: cfg.Length(),
		src: rand.New(cfg.Source()),
	}

	return ret, nil
}

func (g *RandomByteGenerator) Next() interface{} {
	ret := make([]byte, g.len)
	_, _ = g.src.Read(ret)
	return ret
}

func (g *RandomByteGenerator) Type() spansql.TypeBase {
	return spansql.Bytes
}
