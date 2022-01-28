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

type (
	ThreadDataGenerator struct {
		prefixLength     int
		stringLength     int
		rowCount         int
		threadCount      int
		threadGenerators []*CombinedGenerator
	}

	ThreadDataGeneratorConfig struct {
		PrefixLength int
		StringLength int
		RowCount     int
		ThreadCount  int
	}
)

func NewThreadDataGenerator(cfg ThreadDataGeneratorConfig) (*ThreadDataGenerator, error) {
	ret := &ThreadDataGenerator{
		prefixLength: cfg.PrefixLength,
		stringLength: cfg.StringLength,
		threadCount:  cfg.ThreadCount,
		rowCount:     cfg.RowCount,
	}

	return ret, nil
}

func (s *ThreadDataGenerator) GetThreadGenerators() {
	rowsPerThread := s.rowCount / s.threadCount
	i := 0
	threads := make([]*CombinedGenerator, s.threadCount)
	for i < s.threadCount {
		gen, _ := NewCombinedGenerator(CombinedGeneratorConfig{
			Min:          i * rowsPerThread,
			Max:          (i+1)*rowsPerThread - 1,
			PrefixLength: s.prefixLength,
			StringLength: s.stringLength,
		})
		threads[i] = gen
		i++
	}
	s.threadGenerators = threads
}
