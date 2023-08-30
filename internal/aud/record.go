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
