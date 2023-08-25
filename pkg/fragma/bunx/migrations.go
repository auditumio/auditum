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
