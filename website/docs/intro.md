---
sidebar_position: 1
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Introduction

Welcome to Auditum documentation 👋

Let's discover what Auditum is and how it fits into your architecture.

## What is Auditum?

Auditum is an audit trail management software. It allows you to collect, store,
and query audit records (also known as audit logs) for multiple applications
in one place and query them using a simple API. It is designed to integrate well
with business applications, and it can be used with any software that produces
audit logs.

## What are audit logs?

Audit logs answer the question "_who did what and when?_". They are useful for
business features (think "change history"), security, compliance, incident 
management and debugging.

As an example, let's imagine we are building a Reddit-like social website. We
have users, posts and comments. Users can create posts (and comments). Some of
the users are moderators and can modify or delete posts of other users. We want
to keep track of all the change history of our posts and comments. Auditum can
help us with that.

## How can Auditum help?

Every time a user creates, modifies or deletes a post or a comment, we can send
an audit log record to Auditum. Auditum will store the record and allow us to
query it on demand. For example, we can query the chronological history of:

- the changes made to a particular post and its comments;
- the changes made by a particular user;
- the posts created yesterday;
- the changes made by moderators in the last 24 hours;
- the calls to a particular operation in your API;
- the operations that failed in the last hour;
- the operations that match the trace ID in your distributed tracing system;
- and many more.

## Why Auditum?

It is a good idea to have a separate service for keeping audit logs. One of the
many reasons is standardization, and Auditum helps make records consistent 
across applications. Another reason is that, in practice, audit logs typically
take up more storage space than all the other data in the database combined. 
Auditum is designed to handle large volumes of audit logs efficiently and 
offload the main database of your application.

Auditum a standalone service that is easy to configure, deploy and use. It
values developer experience and simplicity over complexity and feature bloat,
and it is focused on minimizing operational overhead.

Auditum is open source and free to use. It is licensed under the Apache 2.0 License.
You can find the source code on [GitHub](https://github.com/auditumio/auditum) and
contribute to the project.

## Next Steps

Follow the [Quickstart](/docs/quickstart) guide to try out Auditum in 5 minutes.
