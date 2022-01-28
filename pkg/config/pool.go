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

// Assert that Connection implements Validate
var _ Validate = (*Pool)(nil)

type (
	Pool struct {
		MaxOpened           int           `mapstructure:"max_opened" yaml:"max_opened"`
		MinOpened           int           `mapstructure:"min_opened" yaml:"min_opened"`
		MaxIdle             int           `mapstructure:"max_idle" yaml:"max_idle"`
		WriteSessions       float64       `mapstructure:"write_sessions" yaml:"write_sessions"`
		HealthcheckWorkers  int           `mapstructure:"healthcheck_workers" yaml:"healthcheck_workers"`
		HealthcheckInterval time.Duration `mapstructure:"healthcheck_interval" yaml:"healthcheck_interval"`
		TrackSessionHandles bool          `mapstructure:"track_session_handles" yaml:"track_session_handles"`
	}
)

func (c *Pool) Validate() error {
	var result *multierror.Error

	// TODO: Validate pool config

	return result.ErrorOrNil()
}
