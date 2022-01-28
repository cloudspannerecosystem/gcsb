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

func TestArrayGenerator(t *testing.T) {
	Convey("ArrayGenerator", t, func() {

		Convey("Missing  Generator", func() {
			bg, err := NewBooleanGenerator(NewConfig())
			So(err, ShouldBeNil)
			So(bg, ShouldNotBeNil)

			cfg := NewConfig()
			cfg.SetLength(10)
			ag, err := NewArrayGenerator(cfg)

			So(err, ShouldNotBeNil)
			So(ag, ShouldBeNil)
		})

		Convey("Invalid Length", func() {
			bg, err := NewBooleanGenerator(NewConfig())
			So(err, ShouldBeNil)
			So(bg, ShouldNotBeNil)

			cfg := NewConfig()
			cfg.SetGenerator(bg)
			cfg.SetLength(-5)
			ag, err := NewArrayGenerator(cfg)

			So(err, ShouldNotBeNil)
			So(ag, ShouldBeNil)
		})

		Convey("Next", func() {
			bcfg := NewConfig()
			bcfg.SetStatic(true)
			bcfg.SetValue(true)
			bg, err := NewBooleanGenerator(bcfg)
			So(err, ShouldBeNil)
			So(bg, ShouldNotBeNil)

			cfg := NewConfig()
			cfg.SetLength(10)
			cfg.SetGenerator(bg)

			ag, err := NewArrayGenerator(cfg)

			So(err, ShouldBeNil)
			So(ag, ShouldNotBeNil)

			v, ok := ag.Next().([]bool)
			So(ok, ShouldBeTrue)
			So(v, ShouldNotBeNil)

			tagl, ok := ag.(*ArrayGenerator)
			So(ok, ShouldBeTrue)
			So(tagl, ShouldNotBeNil)

			So(v, ShouldHaveLength, tagl.l)

			for _, e := range v {
				So(e, ShouldBeTrue)
			}
		})
	})
}

func BenchmarkBooleanArrayGenerator(b *testing.B) {
	bg, err := NewBooleanGenerator(NewConfig())
	if err != nil {
		panic(err)
	}

	cfg := NewConfig()
	cfg.SetLength(10)
	cfg.SetGenerator(bg)
	ag, err := NewArrayGenerator(cfg)

	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ag.Next()
	}
}

func BenchmarkInt64ArrayGenerator(b *testing.B) {
	ig, err := NewInt64Generator(NewConfig())
	if err != nil {
		panic(err)
	}

	cfg := NewConfig()
	cfg.SetLength(10)
	cfg.SetGenerator(ig)
	ag, err := NewArrayGenerator(cfg)

	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ag.Next()
	}
}
