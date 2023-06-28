package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"

	"github.com/infragmo/auditum/internal/aud"
)

type Store struct {
	db *bun.DB
}

func NewStore(db *bun.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateProject(ctx context.Context, project aud.Project) error {
	model := toProjectModel(project)

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(&model).
			Returning("partition_number").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("insert project into db: %v", err)
		}

		if tx.Dialect().Name() == dialect.PG {
			if err := createTablePartitionForProject(
				ctx,
				tx,
				tableNameRecords,
				model.ID,
				model.PartitionNumber,
			); err != nil {
				return fmt.Errorf("create partition of records for project: %v", err)
			}

			if err := createTablePartitionForProject(
				ctx,
				tx,
				tableNameRecordsResourceChanges,
				model.ID,
				model.PartitionNumber,
			); err != nil {
				return fmt.Errorf("create partition of record resource changes for project: %v", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("run transaction: %w", err)
	}

	return nil
}

func (s *Store) GetProject(ctx context.Context, id aud.ID) (aud.Project, error) {
	return getProject(ctx, s.db, id)
}

func (s *Store) ListProjects(
	ctx context.Context,
	limit int32,
	cursor aud.ProjectCursor,
) ([]aud.Project, error) {
	var models []projectModel

	q := s.db.NewSelect().
		Model(&models)

	if cursor.LastID != nil {
		q.Where("id < ?", cursor.LastID)
	}

	q.Order("id DESC")
	q.Limit(int(limit))

	err := q.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("select projects from db: %v", err)
	}

	projects := fromProjectModels(models)
	return projects, nil
}

func (s *Store) UpdateProject(
	ctx context.Context,
	id aud.ID,
	update aud.ProjectUpdate,
) (aud.Project, error) {
	var columns []string
	if update.UpdateDisplayName {
		columns = append(columns, "display_name")
	}
	if update.UpdateUpdateRecordEnabled {
		columns = append(columns, "update_record_enabled")
	}
	if update.UpdateDeleteRecordEnabled {
		columns = append(columns, "delete_record_enabled")
	}
	if len(columns) == 0 {
		return aud.Project{}, fmt.Errorf("nothing to update")
	}

	proj := aud.Project{
		ID:                  id,
		DisplayName:         update.DisplayName,
		UpdateRecordEnabled: update.UpdateRecordEnabled,
		DeleteRecordEnabled: update.DeleteRecordEnabled,
	}
	model := toProjectModel(proj)

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		result, err := tx.NewUpdate().
			Model(&model).
			Column(columns...).
			Where("id = ?", id).
			Returning("*").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("update project in db: %v", err)
		}

		if rowsAffected(result) == 0 {
			return aud.ErrProjectNotFound
		}

		return nil
	})
	if err != nil {
		return aud.Project{}, fmt.Errorf("run transaction: %w", err)
	}

	project := fromProjectModel(model)
	return project, nil
}

func (s *Store) CreateRecord(ctx context.Context, record aud.Record) error {
	model := toRecordModel(record)

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := projectExists(ctx, tx, record.ProjectID); err != nil {
			return err
		}

		_, err := tx.NewInsert().
			Model(&model).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("insert record into db: %v", err)
		}

		if len(model.ResourceChanges) == 0 {
			return nil
		}

		_, err = tx.NewInsert().
			Model(&model.ResourceChanges).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("insert record resource changes into db: %v", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("run transaction: %w", err)
	}

	return nil
}

func (s *Store) CreateRecords(ctx context.Context, records []aud.Record) error {
	if len(records) == 0 {
		return fmt.Errorf("no records to create")
	}

	projectID := records[0].ProjectID
	for i := 1; i < len(records); i++ {
		if records[i].ProjectID != projectID {
			return fmt.Errorf("records must have the same project id")
		}
	}

	recordMods := toRecordModels(records)

	var changeMods []recordResourceChangeModel
	for _, recordMod := range recordMods {
		changeMods = append(changeMods, recordMod.ResourceChanges...)
	}

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := projectExists(ctx, tx, projectID); err != nil {
			return err
		}

		_, err := tx.NewInsert().
			Model(&recordMods).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("insert records into db: %v", err)
		}

		if len(changeMods) == 0 {
			return nil
		}

		_, err = tx.NewInsert().
			Model(&changeMods).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("insert record resource changes into db: %v", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("run transaction: %w", err)
	}

	return nil
}

func (s *Store) GetRecord(
	ctx context.Context,
	projectID aud.ID,
	id aud.ID,
) (aud.Record, error) {
	var model recordModel

	// Transaction is used since the query contains relation.
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := projectExists(ctx, tx, projectID); err != nil {
			return err
		}

		err := tx.NewSelect().
			Model(&model).
			Relation(relationResourceChanges).
			Where("project_id = ?", projectID).
			Where("id = ?", id).
			Scan(ctx)
		if errors.Is(err, sql.ErrNoRows) {
			return aud.ErrRecordNotFound
		}
		if err != nil {
			return fmt.Errorf("select record from db: %v", err)
		}

		return nil
	})
	if err != nil {
		return aud.Record{}, fmt.Errorf("run transaction: %w", err)
	}

	record := fromRecordModel(model)
	return record, nil
}

func (s *Store) ListRecords(
	ctx context.Context,
	projectID aud.ID,
	filter aud.RecordFilter,
	limit int32,
	cursor aud.RecordCursor,
) ([]aud.Record, error) {
	var models []recordModel

	// Transaction is used since the query contains relation.
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := projectExists(ctx, tx, projectID); err != nil {
			return err
		}

		q := tx.NewSelect().
			Model(&models).
			Relation(relationResourceChanges)

		q.Where("project_id = ?", projectID)

		if len(filter.Labels) > 0 {
			switch tx.Dialect().Name() {
			case dialect.PG:
				q.Where("labels @> ?", filter.Labels)
			case dialect.SQLite:
				for k, v := range filter.Labels {
					q.Where("json_extract(labels, ?) = ?", "$."+k, v)
				}
			default:
				return fmt.Errorf("unsupported dialect: %s", tx.Dialect().Name().String())
			}
		}

		if filter.ResourceType != "" {
			q.Where("resource_type = ?", filter.ResourceType)
		}
		if filter.ResourceID != "" {
			q.Where("resource_id = ?", filter.ResourceID)
		}

		if filter.OperationType != "" {
			q.Where("operation_type = ?", filter.OperationType)
		}
		if filter.OperationID != "" {
			q.Where("operation_id = ?", filter.OperationID)
		}
		if !filter.OperationTimeFrom.IsZero() {
			q.Where("operation_time >= ?", filter.OperationTimeFrom)
		}
		if !filter.OperationTimeTo.IsZero() {
			q.Where("operation_time < ?", filter.OperationTimeTo)
		}

		if filter.ActorType != "" {
			q.Where("actor_type = ?", filter.ActorType)
		}
		if filter.ActorID != "" {
			q.Where("actor_id = ?", filter.ActorID)
		}

		if !cursor.Empty() {
			q.Where("operation_time < ?", cursor.LastOperationTime)
			q.Where("id < ?", cursor.LastID)
		}

		q.Order("operation_time DESC", "id DESC")
		q.Limit(int(limit))

		err := q.Scan(ctx)
		if err != nil {
			return fmt.Errorf("select records from db: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("run transaction: %w", err)
	}

	records := fromRecordModels(models)
	return records, nil
}

func (s *Store) UpdateRecord(
	ctx context.Context,
	projectID aud.ID,
	id aud.ID,
	update aud.RecordUpdate,
) (aud.Record, error) {
	var columns []string
	if update.UpdateLabels {
		columns = append(
			columns,
			"labels",
		)
	}
	if update.UpdateResource {
		columns = append(
			columns,
			"resource_type",
			"resource_id",
			"resource_metadata",
		)
	}
	if update.UpdateOperation {
		columns = append(
			columns,
			"operation_type",
			"operation_id",
			"operation_time",
			"operation_metadata",
			"operation_traceparent",
			"operation_tracestate",
			"operation_status",
		)
	}
	if update.UpdateActor {
		columns = append(
			columns,
			"actor_type",
			"actor_id",
			"actor_metadata",
		)
	}
	if len(columns) == 0 {
		return aud.Record{}, fmt.Errorf("nothing to update")
	}

	rec := aud.Record{
		ID:        id,
		ProjectID: projectID,
		Labels:    update.Labels,
		Resource:  update.Resource,
		Operation: update.Operation,
		Actor:     update.Actor,
	}
	model := toRecordModel(rec)

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		proj, err := getProject(ctx, tx, projectID)
		if err != nil {
			return err
		}

		if proj.UpdateRecordEnabled.False() {
			return aud.ErrDisabled
		}

		result, err := tx.NewUpdate().
			Model(&model).
			Column(columns...).
			Where("project_id = ?", projectID).
			Where("id = ?", id).
			Returning("*").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("update record in db: %v", err)
		}

		if rowsAffected(result) == 0 {
			return aud.ErrRecordNotFound
		}

		if !update.UpdateResource {
			return nil
		}

		_, err = tx.NewDelete().
			Model(&model.ResourceChanges).
			Where("project_id = ?", projectID).
			Where("record_id = ?", id).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete resource changes from db: %v", err)
		}

		_, err = tx.NewInsert().
			Model(&model.ResourceChanges).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("insert record resource changes into db: %v", err)
		}

		return nil
	})
	if err != nil {
		return aud.Record{}, fmt.Errorf("run transaction: %w", err)
	}

	record := fromRecordModel(model)
	return record, nil
}

func (s *Store) DeleteRecord(ctx context.Context, projectID aud.ID, id aud.ID) error {
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		proj, err := getProject(ctx, tx, projectID)
		if err != nil {
			return err
		}

		if proj.DeleteRecordEnabled.False() {
			return aud.ErrDisabled
		}

		_, err = tx.NewDelete().
			Model((*recordModel)(nil)).
			Where("project_id = ?", projectID).
			Where("id = ?", id).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete record into db: %v", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("run transaction: %w", err)
	}

	return nil
}

func getProject(ctx context.Context, idb bun.IDB, id aud.ID) (aud.Project, error) {
	var model projectModel

	err := idb.NewSelect().
		Model(&model).
		Where("id = ?", id).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return aud.Project{}, aud.ErrProjectNotFound
	}
	if err != nil {
		return aud.Project{}, fmt.Errorf("select project from db: %v", err)
	}

	project := fromProjectModel(model)
	return project, nil
}

func projectExists(ctx context.Context, idb bun.IDB, id aud.ID) error {
	exists, err := idb.NewSelect().
		Model((*projectModel)(nil)).
		Where("id = ?", id).
		Exists(ctx)
	if err != nil {
		return fmt.Errorf("select project from db: %v", err)
	}

	if !exists {
		return aud.ErrProjectNotFound
	}

	return nil
}

// rowsAffected returns the number of rows affected for the result.
// It panics on error. It assumes that driver implements RowsAffected,
// making it easier to use.
func rowsAffected(result sql.Result) int64 {
	n, err := result.RowsAffected()
	if err != nil {
		panic("rows affected: " + err.Error())
	}
	return n
}
