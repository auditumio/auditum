package aud

import (
	"encoding/json"
	"time"
)

type Record struct {
	ID         ID
	ProjectID  ID
	CreateTime time.Time
	Labels     map[string]string
	Resource   Resource
	Operation  Operation
	Actor      Actor
}

type Resource struct {
	Type     string
	ID       string
	Metadata map[string]string
	Changes  []ResourceChange
}

type ResourceChange struct {
	Name        string
	Description string
	OldValue    json.RawMessage
	NewValue    json.RawMessage
}

type Operation struct {
	Type         string
	ID           string
	Time         time.Time
	Metadata     map[string]string
	TraceContext TraceContext
	Status       OperationStatus
}

type TraceContext struct {
	Traceparent string
	Tracestate  string
}

func (tc TraceContext) IsZero() bool {
	return tc.Traceparent == "" && tc.Tracestate == ""
}

type OperationStatus int

func (s OperationStatus) Int() int {
	return int(s)
}

const (
	OperationStatusUnspecified OperationStatus = iota
	OperationStatusSucceeded
	OperationStatusFailed
)

type Actor struct {
	Type     string
	ID       string
	Metadata map[string]string
}
