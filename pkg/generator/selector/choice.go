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

var (
	_ Choice         = (*choice)(nil)         // Assert that choice implements Choice
	_ WeightedChoice = (*weightedChoice)(nil) // Assert that weightedChoice implements WeightedChoice
)

type (
	// Choice wraps a selection target
	Choice interface {
		Item() interface{}
	}

	// WeightedChoice is an item selection target with a non negative weight
	WeightedChoice interface {
		Choice
		Weight() uint
	}

	choice struct {
		item interface{}
	}

	weightedChoice struct {
		choice
		weight uint
	}
)

func NewChoice(item interface{}) Choice {
	return &choice{
		item: item,
	}
}

func (c *choice) Item() interface{} {
	return c.item
}

// NewWeightedChoice returns a choice that has a probability weight
func NewWeightedChoice(item interface{}, weight uint) WeightedChoice {
	return &weightedChoice{
		choice: choice{
			item: item,
		},
		weight: weight,
	}
}

func (c *weightedChoice) Weight() uint {
	return c.weight
}
