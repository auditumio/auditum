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

package bunx

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"

	"github.com/auditumio/auditum/pkg/fragma/zapx/zapxmigrate"
)

func RunMigrations(mig *migrate.Migrate, log any) error {
	switch l := log.(type) {
	case nil:
		// Do nothing.
	case *zap.Logger:
		mig.Log = zapxmigrate.NewLogger(l)
	case migrate.Logger:
		mig.Log = l
	default:
		return fmt.Errorf("unsupported log type: %T", log)
	}

	err := mig.Up()
	if err == migrate.ErrNoChange {
		if mig.Log != nil {
			mig.Log.Printf("No migrations to apply")
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("migrate up: %v", err)
	}

	return nil
}
