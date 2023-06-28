---
sidebar_position: 2
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Create Records

Records can be created in two ways: creating one record or a batch of records.

## Create a Record

To create one record, send `POST` request to `/projects/{project_id}/records`.

Example request:

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

Example response:

```json
{
  "record": {
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
}
```

## Create a Batch of Records

To create a batch of records, send `POST` request to `/projects/{project_id}/records:batchCreate`.

Example request:

<Tabs>
<TabItem value="shell" label="Shell">

```shell
curl \
  --request POST \
  --header "Content-Type: application/json" \
  --header "Accept: application/json+pretty" \
  --data @- \
  "localhost:8080/api/v1alpha1/projects/01886e86-1963-7f3c-b672-b5d93cec6c6e/records:batchCreate" \
  <<EOM
{
  "records": [
    {
      "labels": {
        "post_id": "post-42"
      },
      "resource": {
        "type": "COMMENT",
        "id": "comment-79",
        "changes": [
          {
            "name": "title",
            "old_value": null,
            "new_value": "Show us, my fiend!"
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
        "id": "example.v1.PostService/CreatePostComment",
        "time": "2023-01-02T03:02:00Z",
        "status": "SUCCEEDED"
      },
      "actor": {
        "type": "USER",
        "id": "user-83"
      }
    },
    {
      "labels": {
        "post_id": "post-42"
      },
      "resource": {
        "type": "COMMENT",
        "id": "comment-79",
        "metadata": {
          "status": "published"
        },
        "changes": [
          {
            "name": "text",
            "description": "Edit text",
            "old_value": "Show us, my fiend!",
            "new_value": "Show us, my friend!"
          }
        ]
      },
      "operation": {
        "type": "UPDATE",
        "id": "example.v1.PostService/UpdatePostComment",
        "time": "2023-01-02T03:03:00Z",
        "status": "SUCCEEDED"
      },
      "actor": {
        "type": "USER",
        "id": "user-83"
      }
    },
    {
      "labels": {
        "post_id": "post-55"
      },
      "resource": {
        "type": "POST",
        "id": "post-55",
        "metadata": {
          "status": "draft"
        },
        "changes": [
          {
            "name": "text",
            "description": "Edit text",
            "old_value": "The dog knows the best seat in the house.",
            "new_value": "For the best seat in the house, you’ll have to move the dog."
          }
        ]
      },
      "operation": {
        "type": "UPDATE",
        "id": "example.v1.PostService/Post",
        "time": "2023-01-02T03:04:00Z",
        "status": "SUCCEEDED"
      },
      "actor": {
        "type": "USER",
        "id": "user-83"
      }
    }
  ]
}
EOM
```

Example response:

```json
{
  "records": [
    {
      "id": "01886e90-69aa-7f3c-8be4-ecc65845e380",
      "project_id": "01886e86-1963-7f3c-b672-b5d93cec6c6e",
      "create_time": "2023-05-30T21:28:58.026109Z",
      "labels": {
        "post_id": "post-42"
      },
      "resource": {
        "type": "COMMENT",
        "id": "comment-79",
        "changes": [
          {
            "name": "title",
            "old_value": null,
            "new_value": "Show us, my fiend!"
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
        "id": "example.v1.PostService/CreatePostComment",
        "time": "2023-01-02T03:02:00Z",
        "status": "SUCCEEDED"
      },
      "actor": {
        "type": "USER",
        "id": "user-83"
      }
    },
    {
      "id": "01886e90-69aa-7f3d-b8a1-df1742963d96",
      "project_id": "01886e86-1963-7f3c-b672-b5d93cec6c6e",
      "create_time": "2023-05-30T21:28:58.026109Z",
      "labels": {
        "post_id": "post-42"
      },
      "resource": {
        "type": "COMMENT",
        "id": "comment-79",
        "metadata": {
          "status": "published"
        },
        "changes": [
          {
            "name": "text",
            "description": "Edit text",
            "old_value": "Show us, my fiend!",
            "new_value": "Show us, my friend!"
          }
        ]
      },
      "operation": {
        "type": "UPDATE",
        "id": "example.v1.PostService/UpdatePostComment",
        "time": "2023-01-02T03:03:00Z",
        "status": "SUCCEEDED"
      },
      "actor": {
        "type": "USER",
        "id": "user-83"
      }
    },
    {
      "id": "01886e90-69aa-7f3e-97c4-cdef2436ad20",
      "project_id": "01886e86-1963-7f3c-b672-b5d93cec6c6e",
      "create_time": "2023-05-30T21:28:58.026109Z",
      "labels": {
        "post_id": "post-55"
      },
      "resource": {
        "type": "POST",
        "id": "post-55",
        "metadata": {
          "status": "draft"
        },
        "changes": [
          {
            "name": "text",
            "description": "Edit text",
            "old_value": "The dog knows the best seat in the house.",
            "new_value": "For the best seat in the house, you’ll have to move the dog."
          }
        ]
      },
      "operation": {
        "type": "UPDATE",
        "id": "example.v1.PostService/Post",
        "time": "2023-01-02T03:04:00Z",
        "status": "SUCCEEDED"
      },
      "actor": {
        "type": "USER",
        "id": "user-83"
      }
    }
  ]
}
```

</TabItem>
</Tabs>
