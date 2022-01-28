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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBooleanGenerator(t *testing.T) {
	Convey("BooleanGenerator", t, func() {

		Convey("Missing  Source", func() {
			bg, err := NewBooleanGenerator(NewConfig())

			So(err, ShouldBeNil)
			So(bg, ShouldNotBeNil)
		})

		Convey("Next", func() {
			bg, err := NewBooleanGenerator(NewConfig())

			So(bg, ShouldNotBeNil)
			So(err, ShouldBeNil)

			var falseSeen, trueSeen bool
			for i := 0; i <= 200; i++ {
				if bg.Next().(bool) {
					trueSeen = true
					continue
				}
				falseSeen = true
			}

			So(falseSeen, ShouldBeTrue)
			So(trueSeen, ShouldBeTrue)
		})

		Convey("Static value", func() {
			cfg := NewConfig()
			cfg.SetStatic(true)
			cfg.SetValue(false)
			bg, err := NewBooleanGenerator(cfg)

			So(err, ShouldBeNil)
			So(bg, ShouldNotBeNil)

			for i := 0; i < 20; i++ {
				So(bg.Next(), ShouldBeFalse)
			}
		})
	})
}

func BenchmarkBooleanGenerator(b *testing.B) {
	bg, err := NewBooleanGenerator(NewConfig())
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bg.Next()
	}
}

func BenchmarkStaticBooleanGenerator(b *testing.B) {
	cfg := NewConfig()
	cfg.SetStatic(true)
	cfg.SetValue(false)
	bg, err := NewBooleanGenerator(cfg)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bg.Next()
	}
}
