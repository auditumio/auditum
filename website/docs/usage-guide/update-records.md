---
sidebar_position: 4
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Update Records

:::caution Important
By default, audit log records are considered immutable. This means that once a record is created, it cannot be updated.
This is done to ensure the integrity of the audit trail.

However, for some projects it may be desirable to update records, e.g. to add new information or fix mistakes.
In this case, you can enable the ability to update records. See [Configuration](/docs/getting-started/configuration) for details.
:::

To update the record, send `PATCH` request to `/projects/{project_id}/records/{record_id}`.

Example request:

<Tabs>
<TabItem value="shell" label="Shell">

```shell
curl \
  --request PATCH \
  --header "Content-Type: application/json" \
  --header "Accept: application/json+pretty" \
  --data @- \
  "localhost:8080/api/v1alpha1/projects/01886e86-1963-7f3c-b672-b5d93cec6c6e/records/01886e90-69aa-7f3e-97c4-cdef2436ad20" \
  <<EOM
{
  "record": {
    "labels": {
      "post_id": "post-55",
      "category": "funny"
    }
  },
  "update_mask": "labels"
}
EOM
```

Example response:

```json
{
  "record": {
    "id": "01886e90-69aa-7f3e-97c4-cdef2436ad20",
    "project_id": "01886e86-1963-7f3c-b672-b5d93cec6c6e",
    "create_time": "2023-05-30T21:28:58.026109Z",
    "labels": {
      "post_id": "post-55",
      "category": "funny"
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
}
```

</TabItem>
</Tabs>
