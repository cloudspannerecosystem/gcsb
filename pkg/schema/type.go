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
	"regexp"
	"strconv"
	"strings"

	"cloud.google.com/go/spanner/spansql"
)

const (
	maxString = 1024 // Maximum string size
	maxByte   = 1024 // Maximum bytes size
)

var lengthRegexp = regexp.MustCompile(`\(([0-9]+|MAX)\)$`)

// TODO: Return an error rather than panic
func ParseSpannerType(spannerType string) spansql.Type {
	ret := spansql.Type{}

	dt := spannerType

	if strings.HasPrefix(dt, "ARRAY<") {
		ret.Array = true
		dt = strings.TrimSuffix(strings.TrimPrefix(dt, "ARRAY<"), ">")
	}

	// separate type and length from dt with length such as STRING(32) or BYTES(256)
	m := lengthRegexp.FindStringSubmatchIndex(dt)
	if m != nil {
		lengthStr := dt[m[2]:m[3]]
		if lengthStr == "MAX" {
			ret.Len = spansql.MaxLen
		} else {
			l, err := strconv.Atoi(lengthStr)
			if err != nil {
				panic("could not convert precision")
			}
			ret.Len = int64(l)
		}

		// trim length from dt
		dt = dt[:m[0]] + dt[m[1]:]
	}

	ret.Base = parseType(dt)

	// Clip length for certain types
	switch ret.Base {
	case spansql.String:
		if ret.Len > maxString {
			ret.Len = maxString
		}
	case spansql.Bytes:
		if ret.Len > maxByte {
			ret.Len = maxByte
		}
	}

	return ret
}

func parseType(dt string) spansql.TypeBase {
	var ret spansql.TypeBase
	switch dt {
	case "BOOL":
		ret = spansql.Bool
	case "STRING":
		ret = spansql.String
	case "INT64":
		ret = spansql.Int64
	case "FLOAT64":
		ret = spansql.Float64
	case "BYTES":
		ret = spansql.Bytes
	case "TIMESTAMP":
		ret = spansql.Timestamp
	case "DATE":
		ret = spansql.Date
	case "NUMERIC":
		ret = spansql.Numeric
	case "JSON":
		ret = spansql.JSON
	default:
		panic(fmt.Sprintf("unknown spanner type '%s'", dt)) // TODO: return error. dont panic
	}

	return ret
}