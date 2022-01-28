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

func TestThreadGenerator(t *testing.T) {
	Convey("ThreadGenerator", t, func() {
		Convey("GetThreadGenerators", func() {
			cg, _ := NewThreadDataGenerator(ThreadDataGeneratorConfig{
				PrefixLength: 8,
				StringLength: 64,
				RowCount:     10,
				ThreadCount:  2,
			})
			cg.GetThreadGenerators()
			Convey("Should Return 2 Generators", func() {
				So(len(cg.threadGenerators), ShouldEqual, 2)
			})
			Convey("First Generator Next", func() {
				So(cg.threadGenerators[0].Next(), ShouldStartWith, "AAAAAAAA")
				So(cg.threadGenerators[0].Next(), ShouldStartWith, "AAAAAAAB")
				So(cg.threadGenerators[0].Next(), ShouldStartWith, "AAAAAAAC")
				So(cg.threadGenerators[0].Next(), ShouldStartWith, "AAAAAAAD")
				So(cg.threadGenerators[0].Next(), ShouldStartWith, "AAAAAAAE")
				So(cg.threadGenerators[0].Next(), ShouldStartWith, "AAAAAAAA")
				So(cg.threadGenerators[0].Next(), ShouldStartWith, "AAAAAAAB")
			})
			Convey("Second Generator Next", func() {
				So(cg.threadGenerators[1].Next(), ShouldStartWith, "AAAAAAAF")
				So(cg.threadGenerators[1].Next(), ShouldStartWith, "AAAAAAAG")
				So(cg.threadGenerators[1].Next(), ShouldStartWith, "AAAAAAAH")
				So(cg.threadGenerators[1].Next(), ShouldStartWith, "AAAAAAAI")
				So(cg.threadGenerators[1].Next(), ShouldStartWith, "AAAAAAAJ")
				So(cg.threadGenerators[1].Next(), ShouldStartWith, "AAAAAAAF")
				So(cg.threadGenerators[1].Next(), ShouldStartWith, "AAAAAAAG")
			})
		})
		// Convey("Second Generator Next11", func() {
		// 	cg, _ := NewThreadDataGenerator(ThreadDataGeneratorConfig{
		// 		PrefixLength: 8,
		// 		StringLength: 64,
		// 		RowCount:     1000,
		// 		ThreadCount:  10,
		// 	})
		// 	cg.GetThreadGenerators()
		// 	for i := range cg.threadGenerators {
		// 		c := 0
		// 		fmt.Println("----", "Generator", i, "-----------------------------------------------")
		// 		for c < 200 {
		// 			fmt.Println(cg.threadGenerators[i].Next())
		// 			c++
		// 		}
		// 	}
		// })

	})
}
