# Changelog

## [Unreleased]

## [0.3.0] - 2024-07-15

### Added

- New _Project_ field `external_id`. The value can be provided when creating a project.
    This is an optional field that can be used to store a custom identifier
    for the project, e.g. for multi-tenancy purposes.
- New `ListProjects` field `filter.external_ids` allows to filter projects by their
    `external_id` field value.
- Tracing: added support for `otlp` exporter, to replace deprecated `jaeger` exporter,
    which is deprecated and will be removed in the next release.

### Deprecated

- `jaeger` tracing exporter is deprecated, use `otlp` instead. See documentation
    for details on how to configure the new exporter.

## [0.2.0] - 2023-08-30

The project is moved to community organization [auditumio](https://github.com/auditumio).
We hope that will this improve the governance and the development of the project.

### Breaking Changes

- The project is moved from `infragmo` to `auditumio` organization.
This change is breaking in a few ways, e.g. the protobuf package is changed, 
docker registry is changed, etc.

### New Features

- Added `store.postgres.logQueries` and `store.sqlite.logQueries` configuration
options to log SQL queries.

### Improvements

- Improved HTTP API reference documentation.
- Added search on documentation website.

## [0.1.0] - 2023-07-09

The initial release.
