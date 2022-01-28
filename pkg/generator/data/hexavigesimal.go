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
	"math"
	"math/big"

	"cloud.google.com/go/spanner/spansql"
	"github.com/cloudspannerecosystem/gcsb/pkg/config"
)

var (
	// Assert that HexavigesimalGenerator implements Generator
	_ Generator = (*HexavigesimalGenerator)(nil)

	base26 = [26]byte{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
		'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
		'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
		'Y', 'Z'}

	index = map[byte]int{
		'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4, 'F': 5,
		'G': 6, 'H': 7, 'I': 8, 'J': 9, 'K': 10, 'L': 11,
		'M': 12, 'N': 13, 'O': 14, 'P': 15, 'Q': 16,
		'R': 17, 'S': 18, 'T': 19, 'U': 20, 'V': 21,
		'W': 22, 'X': 23, 'Y': 24, 'Z': 25,
		'a': 0, 'b': 1, 'c': 2, 'd': 3, 'e': 4, 'f': 5,
		'g': 6, 'h': 7, 'i': 8, 'j': 9, 'k': 10, 'l': 11,
		'm': 12, 'n': 13, 'o': 14, 'p': 15, 'q': 16,
		'r': 17, 's': 18, 't': 19, 'u': 20, 'v': 21,
		'w': 22, 'x': 23, 'y': 24, 'z': 25,
	}
)

type (
	HexavigesimalGenerator struct {
		min    int
		max    int
		cur    int
		length int
	}

	HexavigesimalGeneratorConfig struct {
		Minimum  int
		Maximum  int
		Length   int
		KeyRange *config.TableConfigGeneratorRange
	}
)

func NewHexavigesimalGenerator(cfg HexavigesimalGeneratorConfig) (*HexavigesimalGenerator, error) {
	ret := &HexavigesimalGenerator{
		min:    cfg.Minimum,
		max:    cfg.Maximum,
		cur:    cfg.Minimum,
		length: cfg.Length,
	}
	if cfg.KeyRange != nil {
		ret.min = int(ret.Decode(cfg.KeyRange.Start))
		ret.max = int(ret.Decode(cfg.KeyRange.End))
	}

	return ret, nil
}

func (g *HexavigesimalGenerator) Next() interface{} {
	// If Next() is called more than max - min, wrap back around to min
	if g.cur > g.max {
		g.cur = g.min
	}
	ret := g.Encode(uint64(g.cur), g.length)
	g.cur++
	return ret
}

// Get the maximum number for encoding string of specified length by base of specified length
func (*HexavigesimalGenerator) GetMaxValue(stringLength int, baseLength int) int {
	i := 1
	length := baseLength
	for i < stringLength {
		length += int(math.Pow(float64(baseLength), float64(i)))
		i++
	}
	return length
}

// Encode encodes a uint64 value to string in base26 format
func (g *HexavigesimalGenerator) Encode(value uint64, stringLength int) string {
	chars := make([]uint64, stringLength)
	baseLength := uint64(len(base26))
	current := value
	counter := stringLength - 1
	for current >= baseLength {
		code := current % baseLength
		chars[counter] = code
		counter--
		current = (current - code) / baseLength
	}
	chars[counter] = current
	encoded := ""
	for _, i := range chars {
		encoded = encoded + string(base26[i])
	}
	return encoded
}

// Decode decodes a base26-encoded string back to uint64
func (g *HexavigesimalGenerator) Decode(s string) uint64 {
	res := uint64(0)
	l := len(s) - 1
	b26 := big.NewInt(26)
	bidx := big.NewInt(0)
	bpow := big.NewInt(0)
	for idx := range s {
		c := s[l-idx]
		byteOffset := index[c]
		bidx.SetUint64(uint64(idx))
		res += uint64(byteOffset) * bpow.Exp(b26, bidx, nil).Uint64()
	}
	return res
}

func (g *HexavigesimalGenerator) Type() spansql.TypeBase {
	return spansql.String
}
