package otelx

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
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
