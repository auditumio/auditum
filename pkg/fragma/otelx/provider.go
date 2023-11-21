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

package otelx

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

type Provider struct {
	tracerProvider trace.TracerProvider
}

func NewProvider(opts ...ProviderOption) (*Provider, error) {
	var options providerOptions
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return nil, err
		}
	}

	if options.exporter == nil {
		return nil, fmt.Errorf("exporter option not provided")
	}

	const serviceName = "auditum"

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(options.exporter),
		sdktrace.WithResource(
			resource.NewSchemaless(
				semconv.ServiceName(serviceName),
			),
		),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return &Provider{
		tracerProvider: provider,
	}, nil
}

type ProviderOption func(options *providerOptions) error

func ProviderWithLogExporter(pretty bool) ProviderOption {
	return func(options *providerOptions) error {
		exporter, err := setupLogExporter(pretty)
		if err != nil {
			return err
		}
		options.exporter = exporter
		return nil
	}
}

func ProviderWithJaegerExporter(endpoint string) ProviderOption {
	return func(options *providerOptions) error {
		exporter, err := setupJaegerExporter(endpoint)
		if err != nil {
			return err
		}
		options.exporter = exporter
		return nil
	}
}

func ProviderWithOTLPExporter(endpoint string) ProviderOption {
	return func(options *providerOptions) error {
		exporter, err := setupOTLPExporter(endpoint)
		if err != nil {
			return err
		}
		options.exporter = exporter
		return nil
	}
}

func (p *Provider) Close(ctx context.Context) error {
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if sp, ok := p.tracerProvider.(*sdktrace.TracerProvider); ok {
		return sp.Shutdown(ctx)
	}

	return nil
}

func NoopProvider() *Provider {
	return &Provider{
		tracerProvider: trace.NewNoopTracerProvider(),
	}
}

type providerOptions struct {
	exporter sdktrace.SpanExporter
}

func setupLogExporter(pretty bool) (sdktrace.SpanExporter, error) {
	opts := []stdouttrace.Option{
		stdouttrace.WithWriter(os.Stderr),
	}
	if pretty {
		opts = append(opts, stdouttrace.WithPrettyPrint())
	}

	return stdouttrace.New(opts...)
}

func setupJaegerExporter(endpoint string) (sdktrace.SpanExporter, error) {
	return jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(endpoint),
			jaeger.WithHTTPClient(cleanhttp.DefaultPooledClient()),
		),
	)
}

func setupOTLPExporter(endpoint string) (sdktrace.SpanExporter, error) {
	// We most probably don't need a real context here, since http implementation
	// does nothing on Start, and grpc performs connection in a non-blocking way.
	ctx := context.TODO()

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint as URL: %v", err)
	}

	ep := path.Join(u.Host, u.Path)

	switch u.Scheme {
	case "grpc":
		return otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithEndpoint(ep),
			otlptracegrpc.WithInsecure(),
		)
	case "grpcs":
		return otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithEndpoint(ep),
		)
	case "http":
		return otlptracehttp.New(
			ctx,
			otlptracehttp.WithEndpoint(ep),
			otlptracehttp.WithInsecure(),
		)
	case "https":
		return otlptracehttp.New(
			ctx,
			otlptracehttp.WithEndpoint(ep),
		)
	default:
		return nil, fmt.Errorf("unknown protocol: %s", u.Scheme)
	}
}
