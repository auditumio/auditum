// Copyright 2023 Igor Zibarev
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aud

import (
	"time"
)

type RecordFilter struct {
	Labels map[string]string

	ResourceType string
	ResourceID   string

	OperationType string
	OperationID   string

	OperationTimeFrom time.Time
	OperationTimeTo   time.Time

	ActorType string
	ActorID   string
}

type RecordCursor struct {
	LastOperationTime *time.Time `json:"lot,omitempty"`
	LastID            *ID        `json:"lid,omitempty"`
}

func (p RecordCursor) Empty() bool {
	return p.LastOperationTime == nil && p.LastID == nil
}

func NewRecordCursor(records []Record, pageSize int32) RecordCursor {
	var cursor RecordCursor

	if len(records) >= int(pageSize) {
		last := records[len(records)-1]
		cursor.LastOperationTime = &last.Operation.Time
		cursor.LastID = &last.ID
	}

	return cursor
}

type RecordUpdate struct {
	Labels       map[string]string
	UpdateLabels bool

	Resource       Resource
	UpdateResource bool

	Operation       Operation
	UpdateOperation bool

	Actor       Actor
	UpdateActor bool
}
