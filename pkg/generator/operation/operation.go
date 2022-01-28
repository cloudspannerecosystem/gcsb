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

package operation

import (
	"math/rand"
	"time"

	"github.com/cloudspannerecosystem/gcsb/pkg/config"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/selector"
)

type (
	Operation int
)

const (
	READ Operation = 1 + iota
	WRITE
)

func NewOperationSelector(cfg *config.Config) (selector.Selector, error) {
	return selector.NewWeightedRandomSelector(
		rand.New(rand.NewSource(time.Now().UnixNano())),
		selector.NewWeightedChoice(READ, uint(cfg.Operations.Read)),
		selector.NewWeightedChoice(WRITE, uint(cfg.Operations.Write)),
	)
}
