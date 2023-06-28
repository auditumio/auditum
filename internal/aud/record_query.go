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
