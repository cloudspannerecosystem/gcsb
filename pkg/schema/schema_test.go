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
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSchema(t *testing.T) {
	Convey("Schema", t, func() {
		Convey("Getters/Setters", func() {
			s := NewSchema()
			So(s, ShouldNotBeNil)

			t1 := NewTable()
			t1.SetName("single_table")
			s.SetTable(t1)
			So(s.Table().Name(), ShouldEqual, "single_table")

			for i := 1; i <= 3; i++ {
				t := NewTable()
				t.SetName(fmt.Sprintf("table_%d", i))
				s.AddTable(t)
			}

			So(s.Tables(), ShouldNotBeNil)
			So(s.Tables().Len(), ShouldEqual, 3)
			So(s.GetTable("table_1"), ShouldNotBeNil)
			So(s.GetTable("table_1").Name(), ShouldEqual, "table_1")

			So(s.GetTable("table_1").IsApex(), ShouldBeTrue)

			t1 = s.GetTable("table_1")
			t2 := s.GetTable("table_2")
			t3 := s.GetTable("table_3")
			So(t1, ShouldNotBeNil)
			So(t2, ShouldNotBeNil)
			So(t3, ShouldNotBeNil)

			t3.SetParentName(t2.Name())
			t2.SetParentName(t1.Name())

			err := s.Traverse()
			So(err, ShouldBeNil)
			So(t3.IsApex(), ShouldBeFalse)
			So(t2.IsApex(), ShouldBeFalse)
			So(t1.IsApex(), ShouldBeTrue)

			// Child links should be set as well
			So(t1.HasChild(), ShouldBeTrue)
			So(t1.Child().Name(), ShouldEqual, t2.Name())
			So(t2.HasChild(), ShouldBeTrue)
			So(t2.Child().Name(), ShouldEqual, t3.Name())
			So(t3.HasChild(), ShouldBeFalse)

			So(t3.IsBottom(), ShouldBeTrue)

			// Get apex
			apex := t3.GetApex()
			So(apex, ShouldNotBeNil)
			So(apex.Name(), ShouldEqual, t1.Name())

			// get apex on apex returns itself
			So(t1.GetApex().Name(), ShouldEqual, t1.Name())

			// GetAllRelationNames
			relativeNames := t3.GetAllRelationNames()
			So(relativeNames, ShouldHaveLength, 3)
			So(relativeNames, ShouldContain, t1.Name())
			So(relativeNames, ShouldContain, t2.Name())
			So(relativeNames, ShouldContain, t3.Name())

			// Missing parent tables should return error
			// I actually don't even know if this is possible but it's handled
			t4 := NewTable()
			t4.SetName("table_4")
			t4.SetParentName("NOT_EXIST")
			s.AddTable(t4)

			err = s.Traverse()
			So(err, ShouldNotBeNil)

			// Non-interleaved
			t5 := NewTable()
			t5.SetName("table_5")
			So(t5.IsInterleaved(), ShouldBeFalse)

			// getapex on non-interleaved should be nil
			So(t5.GetApex(), ShouldBeNil)
		})
	})
}
