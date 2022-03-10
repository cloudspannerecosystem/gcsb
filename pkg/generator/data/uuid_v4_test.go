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
	"reflect"
	"testing"

	"cloud.google.com/go/spanner/spansql"
)

func TestNewUUIDV4Generator(t *testing.T) {
	tests := []struct {
		desc    string
		colType spansql.TypeBase
		colLen  int64
		wantErr bool
	}{
		{
			desc:    "STRING(36) should be valid",
			colType: spansql.String,
			colLen:  36,
			wantErr: false,
		},
		{
			desc:    "STRING(32) should be valid",
			colType: spansql.String,
			colLen:  32,
			wantErr: false,
		},
		{
			desc:    "STRING(40) should be valid",
			colType: spansql.String,
			colLen:  40,
			wantErr: false,
		},
		{
			desc:    "BYTES(16) should be valid",
			colType: spansql.Bytes,
			colLen:  16,
			wantErr: false,
		},
		{
			desc:    "BYTES(32) should be valid",
			colType: spansql.Bytes,
			colLen:  32,
			wantErr: false,
		},
		{
			desc:    "STRING(16) should be invalid",
			colType: spansql.String,
			colLen:  16,
			wantErr: true,
		},
		{
			desc:    "INT64 should be invalid",
			colType: spansql.Int64,
			colLen:  0,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			_, err := NewUUIDV4Generator(test.colType, test.colLen)
			gotErr := err != nil
			if gotErr != test.wantErr {
				t.Errorf("NewUUIDV4Generator(%v, %v) got error = %v, but wantErr = %v", test.colType, test.colLen, err, test.wantErr)
			}
		})
	}
}

func TestUUIDV4GeneratorNext(t *testing.T) {
	tests := []struct {
		desc      string
		generator *UUIDV4Generator
		wantType  string
		wantLen   int
	}{
		{
			desc: "UUID for STRING(36)",
			generator: &UUIDV4Generator{
				colType:   spansql.String,
				colLength: 36,
			},
			wantType: "string",
			wantLen:  36,
		},
		{
			desc: "UUID for STRING(40)",
			generator: &UUIDV4Generator{
				colType:   spansql.String,
				colLength: 40,
			},
			wantType: "string",
			wantLen:  36,
		},
		{
			desc: "UUID for STRING(32)",
			generator: &UUIDV4Generator{
				colType:   spansql.String,
				colLength: 32,
			},
			wantType: "string",
			wantLen:  32,
		},
		{
			desc: "UUID for BYTES(16)",
			generator: &UUIDV4Generator{
				colType:   spansql.Bytes,
				colLength: 16,
			},
			wantType: "[]uint8",
			wantLen:  16,
		},
		{
			desc: "UUID for BYTES(32)",
			generator: &UUIDV4Generator{
				colType:   spansql.Bytes,
				colLength: 32,
			},
			wantType: "[]uint8",
			wantLen:  16,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := test.generator.Next()
			rv := reflect.ValueOf(got)
			gotType := rv.Type().String()
			gotLen := rv.Len()

			if gotType != test.wantType {
				t.Errorf("%v Next() returns different type: got = %v, want = %v", test.generator, gotType, test.wantType)
			}
			if gotLen != test.wantLen {
				t.Errorf("%v Next() returns different length: got = %v, want = %v", test.generator, gotLen, test.wantLen)
			}
		})
	}
}
