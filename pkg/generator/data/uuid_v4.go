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
	"fmt"
	"strings"

	"cloud.google.com/go/spanner/spansql"
	"github.com/google/uuid"
)

const (
	uuidV4LenString                 = 36
	uuidV4LenStringWithoutSeparator = 32
	uuidV4LenBytes                  = 16
)

// Assert that UUIDV4Generator implements Generator
var _ Generator = (*UUIDV4Generator)(nil)

// UUIDV4Generator is a generator for UUID v4.
type UUIDV4Generator struct {
	colType   spansql.TypeBase
	colLength int64
}

// NewUUIDV4Generator returns a generator for UUID v4.
func NewUUIDV4Generator(colType spansql.TypeBase, colLength int64) (Generator, error) {
	// Validate column.
	switch colType {
	case spansql.String:
		if colLength != uuidV4LenString && colLength != uuidV4LenStringWithoutSeparator {
			return nil, fmt.Errorf("invalid column length for string UUID: %d", colLength)
		}
	case spansql.Bytes:
		if colLength != uuidV4LenBytes {
			return nil, fmt.Errorf("invalid column length for string UUID: %d", colLength)
		}
	default:
		return nil, fmt.Errorf("invalid column type for UUID: %v", colType.SQL())
	}

	return &UUIDV4Generator{
		colType:   colType,
		colLength: colLength,
	}, nil
}

// Next returns the random UUID v4 value.
func (g *UUIDV4Generator) Next() interface{} {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("unexpected UUID v4 generation error: %v", err))
	}

	if g.colType == spansql.Bytes {
		// NOTE: uuid.MarshalBinary always returns nil error.
		b, _ := id.MarshalBinary()
		return b
	}

	if g.colLength == uuidV4LenStringWithoutSeparator {
		// Returns UUID without separators "-".
		return strings.Join(strings.Split(id.String(), "-"), "")
	}

	return id.String()
}

// Type returns the type for UUID v4 generator.
func (g *UUIDV4Generator) Type() spansql.TypeBase {
	return g.colType
}
