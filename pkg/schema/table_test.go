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

	. "github.com/smartystreets/goconvey/convey"
)

func TestTable(t *testing.T) {
	Convey("Table", t, func() {
		// TODO: Test single column and no column error return
		Convey("PointInsertStatement", func() {
			t := NewTable()
			t.SetName("test")

			c1 := NewColumn()
			c1.SetName("foo")
			c2 := NewColumn()
			c2.SetName("bar")

			t.AddColumn(c1)
			t.AddColumn(c2)

			stmt, err := t.PointInsertStatement()
			So(err, ShouldBeNil)
			So(stmt, ShouldEqual, "INSERT INTO test(foo, bar) VALUES(@foo, @bar)")
		})

		// TODO: Test single predicates
		// TODO: Test no predicates returns error
		// TODO: Test that predicates that are not valid column names return error
		Convey("PointReadStatement", func() {
			t := NewTable()
			t.SetName("test")

			c1 := NewColumn()
			c1.SetName("foo")
			c2 := NewColumn()
			c2.SetName("bar")
			c3 := NewColumn()
			c3.SetName("baz")

			t.AddColumn(c1)
			t.AddColumn(c2)
			t.AddColumn(c3)

			stmt, err := t.PointReadStatement("foo", "bar")
			So(err, ShouldBeNil)
			So(stmt, ShouldEqual, "SELECT foo, bar, baz FROM test WHERE foo = @foo AND bar = @bar")
		})
	})
}
