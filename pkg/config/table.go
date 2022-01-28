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

// Assert that Table implements Validate
var _ Validate = (*Table)(nil)

type (
	Table struct {
		Name       string           `mapstructure:"name"`
		Operations *TableOperations `mapstructure:"operations" yaml:"operations"`
		Columns    []Column         `mapstructure:"columns"`
	}
)

func (t *Table) Validate() error {
	var result *multierror.Error

	// TODO: Validate table config

	return result.ErrorOrNil()
}

func (t *Table) Column(name string) *Column {
	for _, c := range t.Columns {
		if c.Name == name {
			return &c
		}
	}

	return nil
}
