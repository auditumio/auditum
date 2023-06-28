---
sidebar_position: 4
---

import DocCardList from '@theme/DocCardList';

# Usage Guide

Learn the fundamentals of audit log record management.

You can jump straight to a specific section of the guide, or follow the 
description below to get started.

## Basic Workflow

It is quite simple to start managing audit log records with Auditum.

First, you need to [create a project](./create-project). A project is a logical grouping of records.
For example, you can create a project for each of your applications.

Then, you can start [creating audit log records](./create-records) in the project.
This is the most important part, as you will need to map your application data to
audit log records.

After you have created records, you can [query them](./search-records.md) on demand.
The API for searching records is quite flexible to meet your needs. And Auditum
is designed to scale with your data, so you should expect fast responses even with
tons of records.

Optionally, you can configure Auditum to enable [updating](./update-records.md)
and [deleting](./delete-records.md) records, if your requirements allow it.
However, we recommend that you keep the records immutable to ensure data integrity.

And that's it! Now you are able to manage audit log records with Auditum.

<DocCardList />
