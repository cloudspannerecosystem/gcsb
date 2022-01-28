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
	"math"
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func PercentOf(part int, total int) float64 {
	return (float64(part) * float64(100)) / float64(total)
}

func TestWeightedRandomSelector(t *testing.T) {
	Convey("WeightedRandomSelector", t, func() {
		Convey("NewWeightedRandomSelector", func() {
			Convey("Missing Random Source", func() {
				wrs, err := NewWeightedRandomSelector(
					nil,
					NewWeightedChoice("FOO", 1),
					NewWeightedChoice("BAR", 2),
				)
				So(wrs, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, errNoRandomSource)
			})

			Convey("Weight int overflow", func() {
				wrs, err := NewWeightedRandomSelector(
					rand.New(rand.NewSource(time.Now().UnixNano())),
					NewWeightedChoice("FOO", math.MaxInt64),
					NewWeightedChoice("BAR", 1),
				)

				So(wrs, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, errWeightOverflow)
			})

			Convey("No valid weights", func() {
				wrs, err := NewWeightedRandomSelector(
					rand.New(rand.NewSource(time.Now().UnixNano())),
					NewWeightedChoice("FOO", 0),
					NewWeightedChoice("BAR", 0),
				)

				So(wrs, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, errNoValidWeightedChoices)
			})

			Convey("Select", func() {
				iterations := 10000
				wrs, err := NewWeightedRandomSelector(
					rand.New(rand.NewSource(time.Now().UnixNano())),
					NewWeightedChoice("ONE", 25),
					NewWeightedChoice("TWO", 25),
					NewWeightedChoice("THREE", 25),
					NewWeightedChoice("FOUR", 25),
					NewWeightedChoice("NOPICK", 0),
				)

				So(wrs, ShouldNotBeNil)
				So(err, ShouldBeNil)

				seen := make(map[string]int, 2)
				for i := 0; i < iterations; i++ {
					c := wrs.Select()
					v := c.Item().(string)

					if _, e := seen[v]; !e {
						seen[v] = 1
					} else {
						seen[v]++
					}
				}

				So(seen, ShouldNotContainKey, "NOPICK") // 0 weight never returns
				for _, v := range []string{"ONE", "TWO", "THREE", "FOUR"} {
					So(seen, ShouldContainKey, v)
				}

				margin := 5.0
				for _, v := range seen {
					// I don't know if this is a wise test but essentially,
					// if the sum of all weights is 100, I am treating the
					// amount as a percentage of iterations. Here, I'm asserting
					// that each thing was seen ~25% of the time within a margin.
					// Currently I'm guessing that margin is 5%, but hopefully
					// lower would be better.
					So(PercentOf(v, iterations), ShouldAlmostEqual, 25.0, margin)
				}
			})
		})
	})
}

func BenchmarkWeightedRandomSelector(b *testing.B) {
	num := 100
	choices := make([]WeightedChoice, 0, num)
	for i := 1; i < num; i++ {
		choices = append(choices, NewWeightedChoice(i, uint(i)))
	}

	wrs, err := NewWeightedRandomSelector(
		rand.New(rand.NewSource(time.Now().UnixNano())),
		choices...,
	)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wrs.Select()
	}
}
