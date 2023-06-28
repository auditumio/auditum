---
sidebar_position: 1
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Create Project

To create a project, send `POST` request to `/projects`.

Example request:

<Tabs>
<TabItem value="shell" label="Shell">

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

:::tip Tip
In examples, we use `Accept: application/json+pretty` header to get a pretty-printed
JSON response. This is useful for demonstration purposes. In production requests, use 
`Accept: application/json` or omit the header altogether to reduce response size.
:::

Example response:

```json
{
  "project": {
    "id": "01886e86-1963-7f3c-b672-b5d93cec6c6e",
    "create_time": "2023-05-30T21:17:42.115752Z",
    "display_name": "My Project"
  }
}
```
