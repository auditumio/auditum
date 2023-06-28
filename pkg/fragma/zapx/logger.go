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
