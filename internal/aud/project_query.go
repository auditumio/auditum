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

import "github.com/auditumio/auditum/internal/aud/types"

type ProjectFilter struct {
	ExternalIDs []string
}

type ProjectCursor struct {
	LastID *ID `json:"lid,omitempty"`
}

func (p ProjectCursor) Empty() bool {
	return p.LastID == nil
}

func NewProjectCursor(projects []Project, pageSize int32) ProjectCursor {
	var cursor ProjectCursor

	if len(projects) >= int(pageSize) {
		last := projects[len(projects)-1]
		cursor.LastID = &last.ID
	}

	return cursor
}

type ProjectUpdate struct {
	DisplayName       string
	UpdateDisplayName bool

	UpdateRecordEnabled       types.BoolValue
	UpdateUpdateRecordEnabled bool

	DeleteRecordEnabled       types.BoolValue
	UpdateDeleteRecordEnabled bool
}
