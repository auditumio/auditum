---
sidebar_position: 2
---

# Configuration

When you run Auditum without providing any configuration options, it will use 
the default settings. This is fine for quick testing, but you may need to harden
the configuration for production use.

Auditum can be configured using environment variables and a configuration file.

## Configuration File

You can find the [default configuration file](https://github.com/auditumio/auditum/tree/main/config/auditum.yaml) in the repository. All configuration
options are documented in the file. To customize the configuration, you can
create a copy of the file and modify it as needed.

To supply a configuration file to Auditum, use the `--config` flag:

```shell
auditum --config /path/to/config.yaml
```

## Environment Variables

Configuration can also be set via environment variables. Environment variables
are prefixed with `AUDITUM_` and use configuration keys separated by underscores.
For example, to change `log.format` to `text`, set the environment variable
in the following form:

```shell
export AUDITUM_log_format=text
```

## Next Steps

Follow the [Deployment](/docs/getting-started/deployment) to learn how to manage audit logs
with Auditum.
