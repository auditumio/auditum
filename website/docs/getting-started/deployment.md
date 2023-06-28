---
sidebar_position: 3
---

# Deployment

Auditum can be deployed in a variety of ways. This guide will help you prepare
your environment for deployment and operate deployed instances.

## Database

### PostgreSQL

If you have chosen to use PostgreSQL as your database, you will need to run
migrations before starting Auditum. You can do this by running the following
command:

```shell
auditum migrate --config /path/to/config.yaml
```

This command will create the necessary data schema in your database.

:::caution
Beware that you will need to run migrations every time you upgrade Auditum to a
new version.
:::

Then you can start Auditum as usual, e.g. with the following command:

```shell
auditum serve --config /path/to/config.yaml
```

## Scaling

You can run multiple instances of Auditum behind a load balancer to scale the
application. Auditum itself is stateless, it uses a database to store audit logs,
so you can run as many instances as you need.

## Next Steps

Follow the [Usage Guide](/docs/usage-guide) to learn how to manage audit logs
with Auditum.
