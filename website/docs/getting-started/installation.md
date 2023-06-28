---
sidebar_position: 1
---

# Installation

There are multiple ways to install Auditum. Select one that suits your needs.

## Binary

You can download a binary for your platform from the [releases page](https://github.com/infragmo/auditum/releases).

## Docker

The easiest way to run Auditum is to use a Docker image. Images are available on
[Docker Hub](https://hub.docker.com/r/infragmo/auditum).

Example of running Auditum in a Docker container:

```shell
docker run \
  -p 8080:8080 \
  infragmo/auditum:latest
```

## Kubernetes

:::caution
Kubernetes Helm chart is currently in development.
:::

## Source

You can build Auditum from source.

### Prerequisites

- [Go](https://go.dev/doc/install)
- [Task](https://taskfile.dev/installation/)

### Instruction

1. Clone the repository:

  ```shell
  git clone https://github.com/infragmo/auditum.git
  cd auditum
  ```
  
2. Build the binary:

  ```shell
  task build
  ```

## Next Steps

Follow the [Configuration](/docs/getting-started/configuration) guide to learn
how to configure Auditum.
