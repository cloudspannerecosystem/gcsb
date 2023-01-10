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

package schema

import (
	"testing"

	"cloud.google.com/go/spanner/spansql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestColumn(t *testing.T) {
	Convey("Column", t, func() {
		Convey("Type", func() {
			Convey("BOOL", func() {
				c := NewColumn()
				c.SetSpannerType("BOOL")

				x := c.Type()
				So(x, ShouldNotBeNil)
				So(x.Array, ShouldBeFalse)
				So(x.Base, ShouldEqual, spansql.Bool)
			})

			Convey("String(1024)", func() {
				c := NewColumn()
				c.SetSpannerType("STRING(1024)")

				x := c.Type()
				So(x, ShouldNotBeNil)
				So(x.Array, ShouldBeFalse)
				So(x.Base, ShouldEqual, spansql.String)
				So(x.Len, ShouldEqual, 1024)
			})

			Convey("String(MAX)", func() {
				c := NewColumn()
				c.SetSpannerType("STRING(MAX)")

				x := c.Type()
				So(x, ShouldNotBeNil)
				So(x.Array, ShouldBeFalse)
				So(x.Base, ShouldEqual, spansql.String)
				So(x.Len, ShouldEqual, maxString)
			})

			Convey("INT64", func() {
				c := NewColumn()
				c.SetSpannerType("INT64")

				x := c.Type()
				So(x, ShouldNotBeNil)
				So(x.Array, ShouldBeFalse)
				So(x.Base, ShouldEqual, spansql.Int64)
			})

			Convey("FLOAT64", func() {
				c := NewColumn()
				c.SetSpannerType("FLOAT64")

				x := c.Type()
				So(x, ShouldNotBeNil)
				So(x.Array, ShouldBeFalse)
				So(x.Base, ShouldEqual, spansql.Float64)
			})

			Convey("BYTES", func() {
				c := NewColumn()
				c.SetSpannerType("BYTES")

				x := c.Type()
				So(x, ShouldNotBeNil)
				So(x.Array, ShouldBeFalse)
				So(x.Base, ShouldEqual, spansql.Bytes)
			})

			Convey("TIMESTAMP", func() {
				c := NewColumn()
				c.SetSpannerType("TIMESTAMP")

				x := c.Type()
				So(x, ShouldNotBeNil)
				So(x.Array, ShouldBeFalse)
				So(x.Base, ShouldEqual, spansql.Timestamp)
			})

			Convey("DATE", func() {
				c := NewColumn()
				c.SetSpannerType("DATE")

				x := c.Type()
				So(x, ShouldNotBeNil)
				So(x.Array, ShouldBeFalse)
				So(x.Base, ShouldEqual, spansql.Date)
			})

			Convey("ARRAY<STRING(1024)>", func() {
				c := NewColumn()
				c.SetSpannerType("ARRAY<STRING(1024)>")

				x := c.Type()
				So(x, ShouldNotBeNil)
				So(x.Array, ShouldBeTrue)
				So(x.Base, ShouldEqual, spansql.String)
			})

			Convey("ARRAY<STRING(MAX)>", func() {
				c := NewColumn()
				c.SetSpannerType("ARRAY<STRING(MAX)>")

				x := c.Type()
				So(x, ShouldNotBeNil)
				So(x.Array, ShouldBeTrue)
				So(x.Base, ShouldEqual, spansql.String)
			})

			Convey("UNKNOWN", func() {
				c := NewColumn()
				c.SetSpannerType("UNKNOWN")

				So(func() { c.Type() }, ShouldPanic)
			})
		})
	})
}
