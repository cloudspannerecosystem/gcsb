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

import (
	"math/rand"
	"time"

	"cloud.google.com/go/spanner/spansql"
)

type (
	Config interface {
		Copy() Config
		SetBegin(interface{})
		Begin() interface{}
		SetEnd(interface{})
		End() interface{}
		SetLength(int)
		Length() int
		SetStatic(bool)
		Static() bool
		SetValue(interface{})
		Value() interface{}
		SetMinimum(interface{})
		Minimum() interface{}
		SetMaximum(interface{})
		Maximum() interface{}
		SetSource(rand.Source)
		Source() rand.Source
		SetRange(bool)
		Range() bool
		SetGenerator(Generator)
		Generator() Generator
		SetSpannerType(spansql.Type)
		SpannerType() spansql.Type
	}

	generatorConfig struct {
		source      rand.Source
		begin       interface{}
		end         interface{}
		length      int
		static      bool
		value       interface{}
		minimum     interface{}
		maximum     interface{}
		ranged      bool
		generator   Generator
		spannerType spansql.Type
	}
)

func NewConfig() Config {
	return &generatorConfig{}
}

func (c *generatorConfig) Copy() Config {
	cp := *c
	return &cp
}

func (c *generatorConfig) SetSource(x rand.Source) {
	c.source = x
}

func (c *generatorConfig) Source() rand.Source {
	if c.source == nil {
		c.source = rand.NewSource(time.Now().UnixNano())
	}

	return c.source
}

func (c *generatorConfig) SetBegin(x interface{}) {
	c.begin = x
}

func (c *generatorConfig) Begin() interface{} {
	return c.begin
}

func (c *generatorConfig) SetEnd(x interface{}) {
	c.end = x
}

func (c *generatorConfig) End() interface{} {
	return c.end
}

func (c *generatorConfig) SetLength(x int) {
	c.length = x
}

func (c *generatorConfig) Length() int {
	return c.length
}

func (c *generatorConfig) SetStatic(x bool) {
	c.static = x
}

func (c *generatorConfig) Static() bool {
	return c.static
}

func (c *generatorConfig) SetValue(x interface{}) {
	c.value = x
}

func (c *generatorConfig) Value() interface{} {
	return c.value
}

func (c *generatorConfig) SetMinimum(x interface{}) {
	c.minimum = x
}

func (c *generatorConfig) Minimum() interface{} {
	return c.minimum
}

func (c *generatorConfig) SetMaximum(x interface{}) {
	c.maximum = x
}

func (c *generatorConfig) Maximum() interface{} {
	return c.maximum
}

func (c *generatorConfig) SetRange(x bool) {
	c.ranged = x
}

func (c *generatorConfig) Range() bool {
	return c.ranged
}

func (c *generatorConfig) SetGenerator(x Generator) {
	c.generator = x
}

func (c *generatorConfig) Generator() Generator {
	return c.generator
}

func (c *generatorConfig) SetSpannerType(x spansql.Type) {
	c.spannerType = x
}

func (c *generatorConfig) SpannerType() spansql.Type {
	return c.spannerType
}
