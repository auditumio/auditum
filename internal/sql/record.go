package sql

import (
	"time"

	"github.com/uptrace/bun"

	"github.com/infragmo/auditum/internal/aud"
)

const tableNameRecords = "records"

type recordModel struct {
	bun.BaseModel `bun:"table:records,alias:records"`

	ID                   aud.ID                      `bun:"id,pk"`
	ProjectID            aud.ID                      `bun:"project_id,notnull"`
	CreateTime           time.Time                   `bun:"create_time,notnull"`
	Labels               map[string]string           `bun:"labels,type:jsonb"`
	ResourceType         string                      `bun:"resource_type,notnull,nullzero"`
	ResourceID           string                      `bun:"resource_id,notnull,nullzero"`
	ResourceMeta         map[string]string           `bun:"resource_metadata,type:jsonb"`
	ResourceChanges      []recordResourceChangeModel `bun:"rel:has-many,join:id=record_id"`
	OperationType        string                      `bun:"operation_type,notnull,nullzero"`
	OperationID          string                      `bun:"operation_id,notnull,nullzero"`
	OperationTime        time.Time                   `bun:"operation_time,notnull"`
	OperationMeta        map[string]string           `bun:"operation_metadata,type:jsonb"`
	OperationTraceparent string                      `bun:"operation_traceparent,nullzero"`
	OperationTracestate  string                      `bun:"operation_tracestate,nullzero"`
	OperationStatus      int                         `bun:"operation_status,nullzero"`
	ActorType            string                      `bun:"actor_type,notnull,nullzero"`
	ActorID              string                      `bun:"actor_id,notnull,nullzero"`
	ActorMeta            map[string]string           `bun:"actor_metadata,type:jsonb"`
}

func normalizeRecordModel(model *recordModel) {
	model.CreateTime = model.CreateTime.UTC()
	model.OperationTime = model.OperationTime.UTC()
}

func toRecordModel(record aud.Record) recordModel {
	changes := toRecordResourceChangeModels(
		record.ID,
		record.ProjectID,
		record.Resource.Changes,
	)

	return recordModel{
		ID:                   record.ID,
		ProjectID:            record.ProjectID,
		CreateTime:           record.CreateTime,
		Labels:               record.Labels,
		ResourceType:         record.Resource.Type,
		ResourceID:           record.Resource.ID,
		ResourceMeta:         record.Resource.Metadata,
		ResourceChanges:      changes,
		OperationType:        record.Operation.Type,
		OperationID:          record.Operation.ID,
		OperationTime:        record.Operation.Time,
		OperationMeta:        record.Operation.Metadata,
		OperationTraceparent: record.Operation.TraceContext.Traceparent,
		OperationTracestate:  record.Operation.TraceContext.Tracestate,
		OperationStatus:      record.Operation.Status.Int(),
		ActorType:            record.Actor.Type,
		ActorID:              record.Actor.ID,
		ActorMeta:            record.Actor.Metadata,
	}
}

func toRecordModels(records []aud.Record) []recordModel {
	models := make([]recordModel, len(records))
	for i, record := range records {
		models[i] = toRecordModel(record)
	}
	return models
}

func fromRecordModel(model recordModel) aud.Record {
	normalizeRecordModel(&model)

	return aud.Record{
		ID:         model.ID,
		ProjectID:  model.ProjectID,
		CreateTime: model.CreateTime,
		Labels:     model.Labels,
		Resource: aud.Resource{
			Type:     model.ResourceType,
			ID:       model.ResourceID,
			Metadata: model.ResourceMeta,
			Changes:  fromRecordResourceChangeModels(model.ResourceChanges),
		},
		Operation: aud.Operation{
			Type:     model.OperationType,
			ID:       model.OperationID,
			Time:     model.OperationTime,
			Metadata: model.OperationMeta,
			TraceContext: aud.TraceContext{
				Traceparent: model.OperationTraceparent,
				Tracestate:  model.OperationTracestate,
			},
			Status: aud.OperationStatus(model.OperationStatus),
		},
		Actor: aud.Actor{
			Type:     model.ActorType,
			ID:       model.ActorID,
			Metadata: model.ActorMeta,
		},
	}
}

func fromRecordModels(models []recordModel) []aud.Record {
	records := make([]aud.Record, len(models))
	for i, model := range models {
		records[i] = fromRecordModel(model)
	}
	return records
}
