---
sidebar_position: 5
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Delete Records

:::caution Important
By default, audit log records are considered immutable. This means that once a record is created, it cannot be deleted.
This is done to ensure the integrity of the audit trail.

However, for some projects it may be desirable to delete records, e.g. for compliance reasons.
In this case, you can enable the ability to delete records. See [Configuration](/docs/getting-started/configuration) for details.
:::

To update the record, send `DELETE` request to `/projects/{project_id}/records/{record_id}`.

Example request:

<Tabs>
<TabItem value="shell" label="Shell">

```shell
curl \
  --request DELETE \
  "localhost:8080/api/v1alpha1/projects/01886e86-1963-7f3c-b672-b5d93cec6c6e/records/01886e90-69aa-7f3e-97c4-cdef2436ad20"
```

The response does not include any data.

</TabItem>
</Tabs>
