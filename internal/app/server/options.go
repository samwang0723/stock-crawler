// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"net/http"

	config "github.com/samwang0723/stock-crawler/configs"
	"github.com/samwang0723/stock-crawler/internal/app/handlers"
	"github.com/samwang0723/stock-crawler/internal/concurrent"
	structuredlog "github.com/samwang0723/stock-crawler/internal/logger/structured"
)

type Options struct {
	Name    string
	Logger  structuredlog.ILogger
	Handler handlers.IHandler

	Config      *config.Config
	Dispatcher  *concurrent.Dispatcher
	HealthCheck *http.Server

	// Before funcs
	BeforeStart []func() error
	BeforeStop  []func() error

	ProfilingEnabled bool
}

type Option func(o *Options)

func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

func Handler(handler handlers.IHandler) Option {
	return func(o *Options) {
		o.Handler = handler
	}
}

func Logger(logger structuredlog.ILogger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func Config(cfg *config.Config) Option {
	return func(o *Options) {
		o.Config = cfg
	}
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func Dispatcher(dispatcher *concurrent.Dispatcher) Option {
	return func(o *Options) {
		o.Dispatcher = dispatcher
	}
}

func HealthCheck(healthCheck *http.Server) Option {
	return func(o *Options) {
		o.HealthCheck = healthCheck
	}
}
