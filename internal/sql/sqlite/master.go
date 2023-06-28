package sqlite

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type MasterModel struct {
	bun.BaseModel `bun:"table:sqlite_master,alias:sqlite_master"`

	Name string `bun:"name"`
}

func ListTables(ctx context.Context, db *bun.DB) ([]MasterModel, error) {
	var models []MasterModel

	err := db.NewSelect().
		Model(&models).
		Where("type = 'table'").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("select tables from db: %v", err)
	}

	return models, nil
}
