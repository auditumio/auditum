---
sidebar_position: 2
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Quickstart

This guide will show you how to try out Auditum in 5 minutes.

The easiest way to get started is to run Auditum as a Docker container.

First, run the container:

<Tabs>
<TabItem value="shell" label="Shell">

```shell
docker run \
  -p 8080:8080 \
  auditumio/auditum:latest
```

</TabItem>
</Tabs>

This will start Auditum on port 8080. You can now send audit log records using the
Auditum API. We will use **curl** in the following examples to send requests.

To create records, you need to set up a project. A project is a logical grouping of records.

<Tabs>
<TabItem value="shell" label="Shell" default>

```shell
curl \
  --request POST \
  --header "Content-Type: application/json" \
  --header "Accept: application/json+pretty" \
  --data @- \
  "localhost:8080/api/v1alpha1/projects" \
  <<EOM
{
  "project": {
    "display_name": "My Project"
  }
}
EOM
```

</TabItem>
</Tabs>

We get the created project ID from the response:

```json
{
  "project": {
    "id": "01886e86-1963-7f3c-b672-b5d93cec6c6e",
    "create_time": "2023-05-30T21:17:42.115752Z",
    "display_name": "My Project"
  }
}
```

Now let's create a record:

<Tabs>
<TabItem value="shell" label="Shell">

```shell
curl \
  --request POST \
  --header "Content-Type: application/json" \
  --header "Accept: application/json+pretty" \
  --data @- \
  "localhost:8080/api/v1alpha1/projects/01886e86-1963-7f3c-b672-b5d93cec6c6e/records" \
  <<EOM
{
  "record": {
    "labels": {
      "post_id": "post-42"
    },
    "resource": {
      "type": "POST",
      "id": "post-42",
      "metadata": {
        "category": "funny"
      },
      "changes": [
        {
          "name": "text",
          "old_value": null,
          "new_value": "My windows aren’t dirty, that’s my dog’s nose art."
        },
        {
          "name": "status",
          "old_value": null,
          "new_value": "published"
        }
      ]
    },
    "operation": {
      "type": "CREATE",
      "id": "example.v1.PostService/CreatePost",
      "time": "2023-01-02T03:01:00Z",
      "trace_context": {
        "traceparent": "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01",
        "tracestate": "congo=t61rcWkgMzE"
      },
      "status": "SUCCEEDED"
    },
    "actor": {
      "type": "USER",
      "id": "user-82"
    }
  }
}
EOM
```

</TabItem>
</Tabs>

Now let's query the records we just sent:

<Tabs>
<TabItem value="shell" label="Shell">

```shell
curl \
  --request GET \
  --header "Accept: application/json+pretty" \
  "localhost:8080/api/v1alpha1/projects/01886e86-1963-7f3c-b672-b5d93cec6c6e/records"
```

</TabItem>
</Tabs>

You should see the record we just sent in the response:

```json
{
  "records": [
    {
      "id": "01886e87-8f8f-7f3c-b987-2680f63120b8",
      "project_id": "01886e86-1963-7f3c-b672-b5d93cec6c6e",
      "create_time": "2023-05-30T21:19:17.903783Z",
      "labels": {
        "post_id": "post-42"
      },
      "resource": {
        "type": "POST",
        "id": "post-42",
        "metadata": {
          "category": "funny"
        },
        "changes": [
          {
            "name": "text",
            "old_value": null,
            "new_value": "My windows aren’t dirty, that’s my dog’s nose art."
          },
          {
            "name": "status",
            "old_value": null,
            "new_value": "published"
          }
        ]
      },
      "operation": {
        "type": "CREATE",
        "id": "example.v1.PostService/CreatePost",
        "time": "2023-01-02T03:01:00Z",
        "trace_context": {
          "traceparent": "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01",
          "tracestate": "congo=t61rcWkgMzE"
        },
        "status": "SUCCEEDED"
      },
      "actor": {
        "type": "USER",
        "id": "user-82"
      }
    }
  ],
  "next_page_token": ""
}
```

That's it!

You can now configure and deploy Auditum to your environment and start sending
audit log records from your applications.

## Next Steps

Follow the [Getting Started](/docs/getting-started) guide to learn ways to
install and configure Auditum.
