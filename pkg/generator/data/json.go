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

	"cloud.google.com/go/spanner"
	"cloud.google.com/go/spanner/spansql"
)

var _ Generator = (*JsonGenerator)(nil)

type (
	JsonGenerator struct {
		len int
		gen Generator
	}
)

func NewJsonGenerator(cfg Config) (Generator, error) {
	cfg.SetLength(5)

	// string generator
	sg, err := NewStringGenerator(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize string generator: %s", err.Error())
	}

	ret := &JsonGenerator{
		len: 5, // TODO: Make length configurable
		gen: sg,
	}

	return ret, nil
}

func (j *JsonGenerator) Next() interface{} {
	ret := make(map[string]interface{}, j.len)

	for i := 0; i <= j.len; i++ {
		ret[j.gen.Next().(string)] = j.gen.Next()
	}

	return spanner.NullJSON{
		Value: ret,
		Valid: true,
	}

	// return mret
}

func (j *JsonGenerator) Type() spansql.TypeBase {
	return spansql.JSON
}
