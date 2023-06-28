package sql

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun"

	"github.com/infragmo/auditum/internal/aud"
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
