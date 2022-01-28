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
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

var (
	// Assert that config implements Validate
	_ Validate = (*Config)(nil)
)

const (
	DefaultTableOperations = 5
)

type (
	Validate interface {
		Validate() error
	}

	Config struct {
		Project          string        `mapstructure:"project" yaml:"project"`
		Instance         string        `mapstructure:"instance" yaml:"instance"`
		Database         string        `mapstructure:"database" yaml:"database"`
		Threads          int           `mapstructure:"threads" yaml:"threads"`
		NumConns         int           `mapstructure:"num_conns" yaml:"num_cons"`
		MaxExecutionTime time.Duration `mapstructure:"max_execution_time" yaml:"max_execution_time"`
		Operations       Operations    `mapstructure:"operations" yaml:"operations"`
		Pool             Pool          `mapstructure:"pool" yaml:"pool"`
		Tables           []Table       `mapstructure:"tables" yaml:"tables"`
		Batch            bool          `mapstructure:"batch"`
		BatchSize        int           `mapstructure:"batch_size"`
		clientOnce       sync.Once
		client           *spanner.Client
		contextOnce      sync.Once
		ctx              context.Context
		context          context.Context
		contextCancel    context.CancelFunc
	}
)

// NewConfig will unmarshal a viper instance into *Config and validate it
func NewConfig(v *viper.Viper) (*Config, error) {
	// Bind env vars
	Bind(v)

	// Set Default Values
	SetDefaults(v)

	// Unmarshal the config
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Validate will ensure the configuration is valid for attempting to establish a connection
func (c *Config) Validate() error {
	var result *multierror.Error

	if c.Project == "" {
		result = multierror.Append(result, errors.New("project can not be empty"))
	}

	if c.Instance == "" {
		result = multierror.Append(result, errors.New("instance can not be empty"))
	}

	if c.Database == "" {
		result = multierror.Append(result, errors.New("database can not be empty"))
	}

	// Validate pool block
	errs := c.Pool.Validate()
	if errs != nil {
		result = multierror.Append(result, errs)
	}

	return result.ErrorOrNil()
}

// Client returns a configured spanner client
func (c *Config) Client(ctx context.Context) (*spanner.Client, error) {
	var err error
	c.clientOnce.Do(func() {
		c.client, err = spanner.NewClientWithConfig(ctx, c.DB(), spanner.ClientConfig{
			SessionPoolConfig: spanner.SessionPoolConfig{
				MaxOpened:           uint64(c.Pool.MaxOpened),
				MinOpened:           uint64(c.Pool.MinOpened),
				MaxIdle:             uint64(c.Pool.MaxIdle),
				WriteSessions:       c.Pool.WriteSessions,
				HealthCheckWorkers:  c.Pool.HealthcheckWorkers,
				HealthCheckInterval: c.Pool.HealthcheckInterval,
				TrackSessionHandles: c.Pool.TrackSessionHandles,
			},
		},
			option.WithGRPCConnectionPool(c.NumConns),

			// TODO(grpc/grpc-go#1388) using connection pool without WithBlock
			// can cause RPCs to fail randomly. We can delete this after the issue is fixed.
			option.WithGRPCDialOption(grpc.WithBlock()),
		)

	})

	return c.client, err
}

// DB returns the database DSN
func (c *Config) DB() string {
	return fmt.Sprintf("projects/%s/instances/%s/databases/%s", c.Project, c.Instance, c.Database)
}

func (c *Config) Context() (context.Context, context.CancelFunc) {
	c.contextOnce.Do(func() {
		c.ctx = context.Background()
		c.context, c.contextCancel = context.WithCancel(c.ctx)
	})

	return c.context, c.contextCancel
}

func (c *Config) Table(name string) *Table {
	for _, t := range c.Tables {
		if t.Name == name {
			return &t
		}
	}

	return nil
}
