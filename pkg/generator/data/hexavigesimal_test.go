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

func TestHexavigesimalGenerator(t *testing.T) {
	Convey("HexavigesimalGenerator", t, func() {
		hg, _ := NewHexavigesimalGenerator(HexavigesimalGeneratorConfig{})
		Convey("Encode Number", func() {
			// So(hg, ShouldNotBeNil)
			So(hg.Encode(0, 8), ShouldEqual, "AAAAAAAA")
			So(hg.Encode(10, 8), ShouldEqual, "AAAAAAAK")
			So(hg.Encode(25, 8), ShouldEqual, "AAAAAAAZ")
			So(hg.Encode(26, 8), ShouldEqual, "AAAAAABA")
			So(hg.Encode(40, 8), ShouldEqual, "AAAAAABO")
			So(hg.Encode(55, 8), ShouldEqual, "AAAAAACD")
			So(hg.Encode(1000000000, 8), ShouldEqual, "ADGEHTYM")
			So(hg.Encode(208827064575, 8), ShouldEqual, "ZZZZZZZZ")
		})

		Convey("Decode Number", func() {
			// So(hg, ShouldNotBeNil)
			So(hg.Decode("AAAAAAAA"), ShouldEqual, 0)
			So(hg.Decode("AAAAAAAK"), ShouldEqual, 10)
			So(hg.Decode("AAAAAAAZ"), ShouldEqual, 25)
			So(hg.Decode("AAAAAABA"), ShouldEqual, 26)
			So(hg.Decode("AAAAAABO"), ShouldEqual, 40)
			So(hg.Decode("AAAAAACD"), ShouldEqual, 55)
			So(hg.Decode("ADGEHTYM"), ShouldEqual, 1000000000)
			So(hg.Decode("ZZZZZZZZ"), ShouldEqual, 208827064575)
		})

		Convey("GetMaxValue", func() {
			So(hg.GetMaxValue(2, 2), ShouldEqual, 4)
			So(hg.GetMaxValue(3, 2), ShouldEqual, 8)
			So(hg.GetMaxValue(4, 2), ShouldEqual, 16)
			So(hg.GetMaxValue(2, 26), ShouldEqual, 52)
			So(hg.GetMaxValue(8, 26), ShouldEqual, 8353082608)
		})

		Convey("Generate Stuff", func() {
			hg, _ := NewHexavigesimalGenerator(HexavigesimalGeneratorConfig{
				Length: 8,
				Minimum: 0,
				Maximum: 10000000,
			})

			i := 0
			So(hg.Next(), ShouldEqual, "AAAAAAAA")
			for i < 18278 {
				hg.Next()
				i++
			}
			So(hg.Next(), ShouldEqual, "AAAABBBB")
			i = 0
			for i < 18278 {
				hg.Next()
				i++
			}
			So(hg.Next(), ShouldEqual, "AAAACCCC")
		})
	})
}
