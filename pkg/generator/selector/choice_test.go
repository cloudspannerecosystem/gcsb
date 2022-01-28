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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChoice(t *testing.T) {
	Convey("Choice", t, func() {
		Convey("NewChoice", func() {
			c := NewChoice("test")
			So(c, ShouldNotBeNil)

			_, ok := c.(*choice)
			So(ok, ShouldBeTrue)
		})

		Convey("Item", func() {
			c := NewChoice("test")
			So(c, ShouldNotBeNil)

			v := c.Item()
			vv, ok := v.(string)
			So(ok, ShouldBeTrue)
			So(vv, ShouldEqual, "test")
		})
	})

	Convey("Weighted Choice", t, func() {
		Convey("NewWeightedChoice", func() {
			c := NewWeightedChoice("test", 100)
			So(c, ShouldNotBeNil)

			_, ok := c.(*weightedChoice)
			So(ok, ShouldBeTrue)
		})

		Convey("Item", func() {
			c := NewWeightedChoice("test", 100)
			So(c, ShouldNotBeNil)

			v := c.Item()
			vv, ok := v.(string)
			So(ok, ShouldBeTrue)
			So(vv, ShouldEqual, "test")
		})

		Convey("Weight", func() {
			c := NewWeightedChoice("test", 100)
			So(c, ShouldNotBeNil)

			v := c.Weight()
			So(v, ShouldNotBeNil)
			So(v, ShouldHaveSameTypeAs, uint(1))
		})
	})
}
