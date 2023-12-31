---
sidebar_position: 6
---

# Development

## Prerequisites

For development, you will need:

- [buf](https://github.com/bufbuild/buf)
- [golangci-lint](https://github.com/golangci/golangci-lint)
- [Go](https://go.dev)

The simplest way to install all required tools is to use [asdf](https://asdf-vm.com):

```shell
asdf install
```

See [`.tool-versions`](https://github.com/auditumio/auditum/tree/main/.tool-versions) for currently pinned versions.

## Building and testing

Refer to [Taskfile](https://github.com/auditumio/auditum/tree/main/Taskfile.yml) for available commands to work with the project.

## Integrations

### Tracing

To develop tracing integration, you can deploy Jaeger locally using Docker.

We need the following ports:

- `4317` for gRPC OTLP collector
- `4318` for HTTP OTLP collector
- `16686` for UI

This command starts Jaeger using in-memory storage:

```shell
docker run \
  --name jaeger \
  --rm \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 16686:16686 \
  jaegertracing/all-in-one:1.49
```

For more options, see reference: https://www.jaegertracing.io/docs/1.49/deployment

Example Auditum configuration part to send traces to Jaeger via gRPC OTLP port:

```yaml
tracing:
  enabled: true
  exporter: otlp
  otlp:
    endpoint: "grpc://localhost:4317"
```

Make a few requests to Auditum API and then navigate to http://localhost:16686
to see traces in Jaeger UI.
