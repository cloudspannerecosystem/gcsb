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

package config

import (
	"time"

	"github.com/hashicorp/go-multierror"
)

// Assert that Operations implements Validate
var _ Validate = (*Operations)(nil)

type (
	Operations struct {
		Total       int           `mapstructure:"total" yaml:"total"`
		Read        int           `mapstructure:"read" yaml:"read"`
		Write       int           `mapstructure:"write" yaml:"write"`
		SampleSize  float64       `mapstructure:"sample_size" yaml:"sample_size"`
		ReadStale   bool          `mapstructure:"read_stale" yaml:"read_stale"`
		Staleness   time.Duration `mapstructure:"staleness" yaml:"staleness"`
		PartialKeys bool          `mapstructure:"partial_keys" yaml:"partial_keys"`
	}

	TableOperations struct {
		// Read  int `mapstructure:"read"`
		// Write int `mapstructure:"write"`
		Total int `mapstructure:"total"`
	}
)

func (o *Operations) Validate() error {
	var result *multierror.Error

	// TODO: Validate table config

	return result.ErrorOrNil()
}
