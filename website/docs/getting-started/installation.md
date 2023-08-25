---
sidebar_position: 1
---

# Installation

There are multiple ways to install Auditum. Select one that suits your needs.

## Binary

You can download a binary for your platform from the [releases page](https://github.com/auditumio/auditum/releases).

## Docker

The easiest way to run Auditum is to use a Docker image. Images are available on
these registries:

- [Docker Hub](https://hub.docker.com/r/auditumio/auditum) - `auditumio/auditum`
- [GitHub Container Registry](https://github.com/auditumio/auditum/pkgs/container/auditum) - `ghcr.io/auditumio/auditum`

Example of running Auditum in a Docker container:

```shell
docker run \
  -p 8080:8080 \
  auditumio/auditum:latest
```

Or, using a specific version from GitHub Container Registry:

```shell
docker run \
  -p 8080:8080 \
  ghcr.io/auditumio/auditum:0.1.0
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
  git clone https://github.com/auditumio/auditum.git
  cd auditum
  ```
  
2. Build the binary:

  ```shell
  task build
  ```

## Next Steps

Follow the [Configuration](/docs/getting-started/configuration) guide to learn
how to configure Auditum.
