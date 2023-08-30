// Copyright 2023 Igor Zibarev
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

package zapx

import (
	"go.uber.org/zap"
)

func New(opts ...Option) (*zap.Logger, error) {
	o := options{
		format: FormatJSON,
		level:  zap.InfoLevel.String(),
	}
	for _, opt := range opts {
		opt(&o)
	}

	conf, err := getConfig(o.format, o.level)
	if err != nil {
		return nil, err
	}

	return conf.Build()
}

type Option func(*options)

func WithFormat(format string) Option {
	return func(opts *options) {
		opts.format = format
	}
}

func WithLevel(level string) Option {
	return func(opts *options) {
		opts.level = level
	}
}

type options struct {
	format string
	level  string
}
