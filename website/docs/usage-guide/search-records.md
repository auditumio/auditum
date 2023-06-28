---
sidebar_position: 3
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Search Records

To search records, send `GET` request to `/projects/{project_id}/records`.

Example request:

<Tabs>
<TabItem value="shell" label="Shell">

```shell
curl \
  --request GET \
  --header "Accept: application/json+pretty" \
  --get \
  --data "page_size=2" \
  "localhost:8080/api/v1alpha1/projects/01886e86-1963-7f3c-b672-b5d93cec6c6e/records"
```

Example response:

```json
{
  "records": [
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
    }
  ],
  "next_page_token": "eyJsb3QiOiIyMDIzLTAxLTAyVDAzOjAzOjAwWiIsImxpZCI6WzEsMTM2LDExMCwxNDQsMTA1LDE3MCwxMjcsNjEsMTg0LDE2MSwyMjMsMjMsNjYsMTUwLDYxLDE1MF19"
}
```

</TabItem>
</Tabs>
