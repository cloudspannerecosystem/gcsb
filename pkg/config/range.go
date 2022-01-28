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

import "github.com/hashicorp/go-multierror"

// Assert that Generator implements Validate
var _ Validate = (*Range)(nil)

type (
	Range struct {
		Begin   *interface{} `mapstructure:"begin"`   // Begin for ranges like ranged string & date
		End     *interface{} `mapstructure:"end"`     // End of ranges like ranged string & date
		Length  *int         `mapstructure:"length"`  // Length for generators like string or bytes
		Static  *bool        `mapstructure:"static"`  // Static value indicator for bool generator
		Value   *interface{} `mapstructure:"value"`   // Value for static generation
		Minimum *interface{} `mapstructure:"minimum"` // Minimum for numeric generators
		Maximum *interface{} `mapstructure:"maximum"` // Maximum for numeric generators
	}
)

func (r *Range) Validate() error {
	var result *multierror.Error

	// TODO: Validate Range config

	return result.ErrorOrNil()
}
