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
