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

func TestFloat64Generator(t *testing.T) {
	Convey("Float64Generator", t, func() {
		Convey("Nil Source", func() {
			fg, err := NewFloat64Generator(NewConfig())
			So(err, ShouldBeNil)
			So(fg, ShouldNotBeNil)
		})

		Convey("Random", func() {
			fg, err := NewFloat64Generator(NewConfig())

			So(err, ShouldBeNil)
			So(fg, ShouldNotBeNil)

			for i := 0; i < 20; i++ {
				v, ok := fg.Next().(float64)
				So(v, ShouldNotBeNil)
				So(ok, ShouldBeTrue)
			}
		})

		Convey("Ranged", func() {
			cfg := NewConfig()
			cfg.SetRange(true)
			cfg.SetMinimum(5.0)
			cfg.SetMaximum(-5.0)
			fg, err := NewFloat64Generator(cfg)

			So(err, ShouldBeNil)
			So(fg, ShouldNotBeNil)

			for i := 0; i < 20; i++ {
				v, ok := fg.Next().(float64)
				So(v, ShouldNotBeNil)
				So(ok, ShouldBeTrue)
				So(v, ShouldBeLessThanOrEqualTo, 5.0)
				So(v, ShouldBeGreaterThanOrEqualTo, -5.0)
			}
		})
	})
}

func BenchmarkRandomFloatGenerator(b *testing.B) {
	bg, err := NewFloat64Generator(NewConfig())
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bg.Next()
	}
}

func BenchmarkRangedFloatGenerator(b *testing.B) {
	cfg := NewConfig()
	cfg.SetRange(true)
	cfg.SetMinimum(0)
	cfg.SetMaximum(1000)
	bg, err := NewFloat64Generator(cfg)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bg.Next()
	}
}
