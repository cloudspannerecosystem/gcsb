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

package generator

import (
	"testing"

	"github.com/cloudspannerecosystem/gcsb/pkg/config"
	"github.com/cloudspannerecosystem/gcsb/pkg/schema"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerator(t *testing.T) {
	Convey("GetDataGeneratorMapForTable", t, func() {
		ten := 10
		cfg := config.Config{
			Tables: []config.Table{
				{
					Name: "foo",
					Columns: []config.Column{
						{
							Name: "bar",
							Generator: &config.Generator{
								Length: &ten,
							},
						},
						{
							Name: "baz",
							Generator: &config.Generator{
								Range: []*config.Range{
									{
										Minimum: func() *interface{} { var ret interface{} = 10; return &ret }(),  // :(
										Maximum: func() *interface{} { var ret interface{} = 100; return &ret }(), // :(
									},
								},
							},
						},
					},
				},
			},
		}

		foo := schema.NewTable()
		foo.SetName("foo")

		bar := schema.NewColumn()
		bar.SetName("bar")
		bar.SetSpannerType("STRING(1024)")

		baz := schema.NewColumn()
		baz.SetName("baz")
		baz.SetSpannerType("INT64")

		foo.AddColumn(bar)
		foo.AddColumn(baz)

		gmap, err := GetDataGeneratorMapForTable(cfg, foo)
		So(err, ShouldBeNil)
		So(gmap, ShouldNotBeNil)

		// Test bar column generator
		barVal := gmap["bar"].Next()
		So(barVal, ShouldHaveSameTypeAs, "IM A STRING")
		stringBarVal, ok := barVal.(string)
		So(ok, ShouldBeTrue)
		So(stringBarVal, ShouldHaveLength, 10)

		// Test baz column generator
		bazVal := gmap["baz"].Next()
		So(bazVal, ShouldHaveSameTypeAs, int64(10))
		intBazVal, ok := bazVal.(int64)
		So(ok, ShouldBeTrue)
		So(intBazVal, ShouldBeBetween, 10, 100)
	})
}
