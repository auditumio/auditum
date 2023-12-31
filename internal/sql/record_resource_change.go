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

package sql

import (
	"encoding/json"

	"github.com/uptrace/bun"

	"github.com/auditumio/auditum/internal/aud"
)

const relationResourceChanges = "ResourceChanges"

const tableNameRecordsResourceChanges = "records_resource_changes"

type recordResourceChangeModel struct {
	bun.BaseModel `bun:"table:records_resource_changes,alias:records_resource_changes"`

	RecordID    aud.ID          `bun:"record_id,notnull"`
	ProjectID   aud.ID          `bun:"project_id,notnull"`
	Name        string          `bun:"name,notnull,nullzero"`
	Description string          `bun:"description,nullzero"`
	OldValue    json.RawMessage `bun:"old_value,nullzero"`
	NewValue    json.RawMessage `bun:"new_value,nullzero"`
}

func toRecordResourceChangeModel(recordID, projectID aud.ID, src aud.ResourceChange) recordResourceChangeModel {
	return recordResourceChangeModel{
		RecordID:    recordID,
		ProjectID:   projectID,
		Name:        src.Name,
		Description: src.Description,
		OldValue:    src.OldValue,
		NewValue:    src.NewValue,
	}
}

func toRecordResourceChangeModels(recordID, projectID aud.ID, src []aud.ResourceChange) []recordResourceChangeModel {
	dst := make([]recordResourceChangeModel, len(src))
	for i, r := range src {
		dst[i] = toRecordResourceChangeModel(recordID, projectID, r)
	}
	return dst
}

func fromRecordResourceChangeModel(src recordResourceChangeModel) aud.ResourceChange {
	return aud.ResourceChange{
		Name:        src.Name,
		Description: src.Description,
		OldValue:    src.OldValue,
		NewValue:    src.NewValue,
	}
}

func fromRecordResourceChangeModels(src []recordResourceChangeModel) []aud.ResourceChange {
	dst := make([]aud.ResourceChange, len(src))
	for i, r := range src {
		dst[i] = fromRecordResourceChangeModel(r)
	}
	return dst
}
