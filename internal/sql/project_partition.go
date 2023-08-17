package sql

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/auditumio/auditum/internal/aud"
)

func partitionForProjectTableName(table string, ppn int32) string {
	// ppn stands for "project partition number".
	return fmt.Sprintf("%s_ppn_%d", table, ppn)
}

func createTablePartitionForProject(
	ctx context.Context,
	idb bun.IDB,
	ofTableName string,
	projectID aud.ID,
	partitionNumber int32,
) error {
	q := fmt.Sprintf(
		`CREATE TABLE %s PARTITION OF %s FOR VALUES IN (?);`,
		partitionForProjectTableName(ofTableName, partitionNumber),
		ofTableName,
	)

	_, err := idb.ExecContext(ctx, q, projectID)
	if err != nil {
		return fmt.Errorf("create table in db: %v", err)
	}

	return nil
}

func dropTablePartitionForProject(
	ctx context.Context,
	idb bun.IDB,
	ofTableName string,
	partitionNumber int32,
) error {
	_, err := idb.NewDropTable().
		TableExpr(partitionForProjectTableName(ofTableName, partitionNumber)).
		Cascade().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("drop table in db: %v", err)
	}

	return nil
}
