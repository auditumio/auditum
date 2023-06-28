---
sidebar_position: 1
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# gRPC API

Auditum offers classic HTTP/JSON API and modern gRPC API. Both types of API 
provide the same features, but gRPC is more efficient, easier to integrate with,
and it is recommended for better performance and developer experience.

This page describes how to use Auditum gRPC API.

## Protobuf

Auditum gRPC API is defined in [Protobuf](https://protobuf.dev/). The definitions
are well-documented and available [in the same repository](https://github.com/infragmo/auditum/tree/main/api/proto/infragmo/auditum).
The generated code is [available for Go](https://github.com/infragmo/auditum/tree/main/api/gen/go/infragmo/auditum),
and you can generate it for any other language. See the official [guide](https://protobuf.dev/getting-started/)
for your language.

## Try it out

To try out Auditum gRPC API, you can use [grpcurl](https://github.com/fullstorydev/grpcurl).

Let's demonstrate how to create a project.

Example request:

<Tabs>
<TabItem value="shell" label="Shell">

```shell
grpcurl \
  -plaintext \
  -d @ \
  localhost:9090 \
  infragmo.auditum.v1alpha1.ProjectService/CreateProject \
  <<EOM
{
  "project": {
    "display_name": "My Great Project"
  }
}
EOM
```

</TabItem>
</Tabs>

Example response:

```json
{
  "project": {
    "id": "01886e86-1963-7f3c-b672-b5d93cec6c6e",
    "createTime": "2023-05-30T21:17:42.115752Z",
    "displayName": "My Great Project"
  }
}
```

Let's check that the project was created.

Example request:

<Tabs>
<TabItem value="shell" label="Shell">

```shell
grpcurl \
  -plaintext \
  -d @ \
  localhost:9090 \
  infragmo.auditum.v1alpha1.ProjectService/ListProjects \
  <<EOM
{}
EOM
```

</TabItem>
</Tabs>

Example response:

```json
{
  "projects": [
    {
      "id": "01886e86-1963-7f3c-b672-b5d93cec6c6e",
      "createTime": "2023-05-30T21:17:42.115752Z",
      "displayName": "My Great Project"
    }
  ]
}
```
