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

// package sample is used to sample rows from a table in order to generate point reads
package sample

import (
	"fmt"
	"math/rand"
	"reflect"

	"cloud.google.com/go/spanner"
)

type (
	SampleGenerator struct {
		src  *rand.Rand
		l    int                    // sample length
		s    map[string]interface{} // sample map
		cols []string               // columns and order
	}
)

func NewSampleGenerator(src *rand.Rand, samples map[string]interface{}, cols []string) (*SampleGenerator, error) {
	if len(samples) <= 0 {
		return nil, fmt.Errorf("can not use zero length table samples to generate reads (is there data loaded?)")
	}

	if len(samples) != len(cols) {
		return nil, fmt.Errorf("sample cols and sample map must be of equal len (%d != %d)", len(samples), len(cols))
	}

	ret := &SampleGenerator{
		src:  src,
		s:    samples,
		cols: cols,
	}

	i := 0
	for k, v := range samples {
		if reflect.TypeOf(v).Kind() != reflect.Slice {
			return nil, fmt.Errorf("sample for column '%s' is not a slice", k)
		}

		vv := reflect.ValueOf(v)
		if i == 0 {
			ret.l = vv.Len()
		} else {
			if vv.Len() != ret.l {
				return nil, fmt.Errorf("samples for composite primary keys must be of equal length (%s column mismatch)", k)
			}
		}

		i++
	}

	if ret.l == 0 {
		// return nil, fmt.Errorf("can not calculate maximum sample index. this is a bug")
		return nil, fmt.Errorf("can not use zero length table samples to generate reads (is there data loaded?)")
	}

	return ret, nil
}

func (s *SampleGenerator) Next() interface{} {
	ret := spanner.Key{}

	idx := s.src.Intn(s.l)
	for _, col := range s.cols {
		// This is terrible, inefficient, and unsafe
		s := reflect.ValueOf(s.s[col])
		iv := s.Index(idx)
		ret = append(ret, iv.Interface())
	}

	return ret
}
