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
	"database/sql"
	"time"

	"github.com/uptrace/bun"

	"github.com/auditumio/auditum/internal/aud"
)

type projectModel struct {
	bun.BaseModel `bun:"table:projects,alias:projects"`

	ID                  aud.ID       `bun:"id,pk"`
	PartitionNumber     int32        `bun:"partition_number,autoincrement"`
	CreateTime          time.Time    `bun:"create_time,notnull"`
	DisplayName         string       `bun:"display_name,notnull,nullzero"`
	UpdateRecordEnabled sql.NullBool `bun:"update_record_enabled"`
	DeleteRecordEnabled sql.NullBool `bun:"delete_record_enabled"`
}

func normalizeProjectModel(model *projectModel) {
	model.CreateTime = model.CreateTime.UTC()
}

func toProjectModel(project aud.Project) projectModel {
	return projectModel{
		ID:                  project.ID,
		CreateTime:          project.CreateTime,
		DisplayName:         project.DisplayName,
		UpdateRecordEnabled: toBoolValueModel(project.UpdateRecordEnabled),
		DeleteRecordEnabled: toBoolValueModel(project.DeleteRecordEnabled),
	}
}

func fromProjectModel(model projectModel) aud.Project {
	normalizeProjectModel(&model)

	return aud.Project{
		ID:                  model.ID,
		CreateTime:          model.CreateTime,
		DisplayName:         model.DisplayName,
		UpdateRecordEnabled: fromBoolValueModel(model.UpdateRecordEnabled),
		DeleteRecordEnabled: fromBoolValueModel(model.DeleteRecordEnabled),
	}
}

func fromProjectModels(models []projectModel) []aud.Project {
	projects := make([]aud.Project, len(models))
	for i, model := range models {
		projects[i] = fromProjectModel(model)
	}
	return projects
}
