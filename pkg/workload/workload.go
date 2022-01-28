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

// package workload contains an interface for orchestrating goroutines.
// The idea is that we can experiment with different concurrency models
// without disrupting our call sites.
package workload

import (
	"context"

	"github.com/cloudspannerecosystem/gcsb/pkg/config"
	"github.com/cloudspannerecosystem/gcsb/pkg/schema"

	"github.com/rcrowley/go-metrics"
)

type (
	Constructor func(WorkloadConfig) (Workload, error)

	Workload interface {
		Load([]string) error
		Run(string) error
		Stop() error
	}

	WorkloadConfig struct {
		Context        context.Context
		Config         *config.Config
		Schema         schema.Schema
		MetricRegistry metrics.Registry
	}
)

// GetWorkload adds future support for different concurrency models
func GetWorkloadConstructor(workloadType string) (Constructor, error) {
	switch workloadType {
	default:
		// return NewPoolWorkload, nil
		return NewCoreWorkload, nil
	}
}
