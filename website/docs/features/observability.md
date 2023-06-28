---
sidebar_position: 2
---

# Observability

Auditum has built-in support for structured logging, Prometheus metrics and 
OpenTelemetry tracing.

## Logs

By default, Auditum logs its activity to the standard output in JSON format.

Example log:

```json
{"level":"info","time":"2023-06-28T23:08:26.878684Z","logger":"grpc_server_controller","caller":"grpcx/server_controller.go:33","msg":"Starting gRPC server...","addr":":9090"}
```

You can configure the format and level of logs using the `log.format` and 
`log.level` configuration options. See [Configuration](/docs/getting-started/configuration)
for more details.

## Metrics

Auditum exposes Prometheus metrics on the `/metrics` endpoint of the HTTP server.

Example:

```shell
$ curl 'localhost:8080/metrics'
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0.00021668
go_gc_duration_seconds{quantile="0.25"} 0.00023288
go_gc_duration_seconds{quantile="0.5"} 0.000282453
go_gc_duration_seconds{quantile="0.75"} 0.000995577
go_gc_duration_seconds{quantile="1"} 0.000995577
go_gc_duration_seconds_sum 0.00172759
go_gc_duration_seconds_count 4
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 18
# 
# ... redacted ...
#
```

You can find various metrics related to the gRPC server, HTTP server, and Go 
runtime.

## Tracing

Auditum supports OpenTelemetry tracing. By default, tracing is disabled, because
you need to configure exporter that is specific to your infrastructure.

You can configure the tracing exporter using the `tracing.exporter`
configuration option. See [Configuration](/docs/getting-started/configuration) for more details.

### Jaeger

To export traces to [Jaeger](https://www.jaegertracing.io/), you can use the following configuration:

```yaml
tracing:
  enabled: true
  exporter: jaeger
  jaeger:
    endpoint: http://localhost:14268/api/traces
```

Note that other software like [Grafana Tempo](https://grafana.com/oss/tempo/)
also supports Jaeger trace format, so you can use the same configuration for
them.
