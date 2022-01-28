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

type (
	IndexIterator interface {
		ResetIterator()
		HasNext() bool
		GetNext() Index
	}

	Indexes interface {
		IndexIterator
		Indexes() []Index
		AddIndex(Index)
	}

	indexes struct {
		iteratorIndex int
		indexes       []Index
	}
)

func NewIndexes() Indexes {
	return &indexes{
		indexes: make([]Index, 0),
	}
}

func (t *indexes) Len() int {
	return len(t.indexes)
}

func (t *indexes) ResetIterator() {
	t.iteratorIndex = 0
}

func (t *indexes) HasNext() bool {
	return t.iteratorIndex < len(t.indexes)
}

func (t *indexes) GetNext() Index {
	if t.HasNext() {
		to := t.indexes[t.iteratorIndex]
		t.iteratorIndex++
		return to
	}

	return nil
}

func (t *indexes) Indexes() []Index {
	return t.indexes
}

func (t *indexes) AddIndex(x Index) {
	t.indexes = append(t.indexes, x)
}
