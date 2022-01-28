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

package selector

import (
	"errors"
	"math"
	"math/rand"
	"sort"
)

var (
	// Assert that WeightedRandomSelector implements Selector
	_                         Selector = (*WeightedRandomSelector)(nil)
	errWeightOverflow                  = errors.New("sum of WeightedChoice Weights exceeds max int")
	errNoValidWeightedChoices          = errors.New("zero WeightedChoices with Weight() >= 1")
	errNoRandomSource                  = errors.New("missing random source")
)

type (
	WeightedRandomSelector struct {
		source *rand.Rand
		data   []WeightedChoice
		totals []int
		max    int
	}
)

func NewWeightedRandomSelector(rs *rand.Rand, choices ...WeightedChoice) (Selector, error) {
	if rs == nil {
		return nil, errNoRandomSource
	}

	sort.Slice(choices, func(i, j int) bool {
		return choices[i].Weight() < choices[j].Weight()
	})

	totals := make([]int, len(choices))
	runningTotal := 0
	for i, c := range choices {
		weight := int(c.Weight())
		if (math.MaxInt64 - runningTotal) <= weight {
			return nil, errWeightOverflow
		}
		runningTotal += weight
		totals[i] = runningTotal
	}

	if runningTotal < 1 {
		return nil, errNoValidWeightedChoices
	}

	return &WeightedRandomSelector{
		source: rs,
		data:   choices,
		totals: totals,
		max:    runningTotal,
	}, nil
}

// Select returns a choice
func (s *WeightedRandomSelector) Select() Choice {
	r := s.source.Intn(s.max) + 1
	i := searchInts(s.totals, r)
	return s.data[i]
}

// Borrowed from https://github.com/mroth/weightedrand
//
// The standard library sort.SearchInts() just wraps the generic sort.Search()
// function, which takes a function closure to determine truthfulness. However,
// since this function is utilized within a for loop, it cannot currently be
// properly inlined by the compiler, resulting in non-trivial performance
// overhead.
func searchInts(a []int, x int) int {
	i, j := 0, len(a)
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
