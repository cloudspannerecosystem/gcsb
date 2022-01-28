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
	"runtime"

	"github.com/spf13/viper"
)

// SetDefaults takes a viper instance and sets default configuration values
func SetDefaults(v *viper.Viper) {
	// Defaults
	v.SetDefault("num_conns", runtime.GOMAXPROCS(0))
	v.SetDefault("threads", 10)
	v.SetDefault("batch", true)
	v.SetDefault("batch_size", 5)
	v.SetDefault("max_execution_time", 0)

	// Operations defualts
	v.SetDefault("operations.total", 10000)
	v.SetDefault("operations.read", 50)
	v.SetDefault("operations.write", 50)
	v.SetDefault("operations.sample_size", 50)
	v.SetDefault("operations.read_stale", false)

	// Pool Defaults
	v.SetDefault("pool.max_opened", 1000)
	v.SetDefault("pool.min_opened", 100)
	v.SetDefault("pool.max_idle", 0)
	v.SetDefault("pool.write_sessions", 0.2)
	v.SetDefault("pool.healthcheck_workers", 10)
	v.SetDefault("pool.healthcheck_interval", "50m")
	v.SetDefault("pool.track_session_handles", false)
}
