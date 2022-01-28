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
	"math/rand"
	"unsafe"

	"cloud.google.com/go/spanner/spansql"
)

// Assert that StringGenerator implements Generator
var _ Generator = (*StringGenerator)(nil)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type (
	// StringGenerator returns randomly generated strings of a fixed length
	StringGenerator struct {
		len int
		src rand.Source
	}
)

func NewStringGenerator(cfg Config) (Generator, error) {
	ret := &StringGenerator{
		src: cfg.Source(),
		len: cfg.Length(),
	}

	return ret, nil
}

/*
 * Next returns the next randomly generated value
 *
 * The random string generation method was borrowed from icza
 * See: https://stackoverflow.com/a/31832326/145479
 */
func (s *StringGenerator) Next() interface{} {
	b := make([]byte, s.len)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := s.len-1, s.src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = s.src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func (s *StringGenerator) Type() spansql.TypeBase {
	return spansql.String
}
