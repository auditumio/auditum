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

//go:build integration

package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"

	"github.com/auditumio/auditum/internal/aud"
	"github.com/auditumio/auditum/internal/aud/types"
	"github.com/auditumio/auditum/internal/sql/sqltest"
)

func TestIntegration_Store_CreateProject(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := sqltest.NewDatabase(ctx, t)

	// Seed

	setCleanupProjects(t, db)

	// Test

	t.Run("Should create project", func(t *testing.T) {
		id := aud.MustNewID()

		proj := aud.Project{
			ID:          id,
			CreateTime:  time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			DisplayName: "My Project",
			UpdateRecordEnabled: types.BoolValue{
				Bool:  true,
				Valid: true,
			},
			DeleteRecordEnabled: types.BoolValue{
				Bool:  false,
				Valid: false,
			},
		}

		store := NewStore(db)

		err := store.CreateProject(ctx, proj)
		assert.NoError(t, err)

		// Check if project was created.

		var model projectModel
		err = db.NewSelect().
			Model(&model).
			Where("id = ?", id).
			Scan(ctx)
		assert.NoError(t, err)

		got := fromProjectModel(model)
		want := proj
		assert.Equal(t, want, got)
	})
}

func TestIntegration_Store_UpdateProject(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := sqltest.NewDatabase(ctx, t)

	// Seed

	p1id := aud.MustNewID()
	p2id := aud.MustNewID()

	seededProjectModels := []projectModel{
		{
			ID:              p1id,
			PartitionNumber: 1,
			CreateTime:      time.Date(2023, 1, 2, 3, 1, 0, 0, time.UTC),
			DisplayName:     "Project 1",
			UpdateRecordEnabled: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			DeleteRecordEnabled: sql.NullBool{
				Bool:  false,
				Valid: false,
			},
		},
		{
			ID:              p2id,
			PartitionNumber: 2,
			CreateTime:      time.Date(2023, 1, 2, 3, 2, 0, 0, time.UTC),
			DisplayName:     "Project 2",
			UpdateRecordEnabled: sql.NullBool{
				Bool:  false,
				Valid: false,
			},
			DeleteRecordEnabled: sql.NullBool{
				Bool:  false,
				Valid: false,
			},
		},
	}

	seedProjects(ctx, t, db, seededProjectModels...)
	setCleanupProjects(t, db)

	// Test

	t.Run("Should update project", func(t *testing.T) {
		store := NewStore(db)

		id := p2id

		update := aud.ProjectUpdate{
			DisplayName:       "Project 2 updated",
			UpdateDisplayName: true,
			UpdateRecordEnabled: types.BoolValue{
				Bool:  false,
				Valid: false,
			},
			UpdateUpdateRecordEnabled: true,
			DeleteRecordEnabled: types.BoolValue{
				Bool:  false,
				Valid: true,
			},
			UpdateDeleteRecordEnabled: true,
		}

		updatedProject, err := store.UpdateProject(ctx, id, update)
		assert.NoError(t, err)

		assert.Equal(t, aud.Project{
			ID:                  p2id,
			CreateTime:          seededProjectModels[1].CreateTime,
			DisplayName:         update.DisplayName,
			UpdateRecordEnabled: update.UpdateRecordEnabled,
			DeleteRecordEnabled: update.DeleteRecordEnabled,
		}, updatedProject)
	})
}

func TestIntegration_Store_CreateRecord(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := sqltest.NewDatabase(ctx, t)

	// Seed

	seedTestProject(ctx, t, db)
	setCleanupTestProject(t, db)

	setCleanupRecords(t, db)

	// Test

	t.Run("Should create record", func(t *testing.T) {
		id := aud.MustNewID()

		rec := aud.Record{
			ID:         id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			Resource: aud.Resource{
				Type: "COMMENT",
				ID:   "comment-7",
				Metadata: map[string]string{
					"status": "published",
				},
				Changes: []aud.ResourceChange{
					{
						Name:        "text",
						Description: "Edit text",
						OldValue:    json.RawMessage(`"Hello world"`),
						NewValue:    json.RawMessage(`"Hello, World!"`),
					},
				},
			},
			Operation: aud.Operation{
				Type: "UPDATE",
				ID:   "example.v1.PostService/UpdatePostComment",
				Time: time.Date(2023, 1, 1, 2, 1, 0, 0, time.UTC),
				Metadata: map[string]string{
					"via": "Moderator UI",
				},
				TraceContext: aud.TraceContext{
					Traceparent: "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01",
					Tracestate:  "congo=t61rcWkgMzE",
				},
				Status: aud.OperationStatusSucceeded,
			},
			Actor: aud.Actor{
				Type: "USER",
				ID:   "user-82",
				Metadata: map[string]string{
					"as": "moderator",
				},
			},
		}

		store := NewStore(db)

		err := store.CreateRecord(ctx, rec)
		assert.NoError(t, err)

		// Check if record was created.

		var model recordModel
		err = db.NewSelect().
			Model(&model).
			Relation(relationResourceChanges).
			Where("id = ?", id).
			Scan(ctx)
		assert.NoError(t, err)

		got := fromRecordModel(model)
		want := rec
		assert.Equal(t, want, got)
	})
}

func TestIntegration_Store_ListRecords(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := sqltest.NewDatabase(ctx, t)

	// Seed

	seedTestProject(ctx, t, db)
	setCleanupTestProject(t, db)

	r1id := aud.MustNewID()
	r2id := aud.MustNewID()
	r3id := aud.MustNewID()
	r4id := aud.MustNewID()
	r5id := aud.MustNewID()
	r6id := aud.MustNewID()

	seededRecordModels := []recordModel{
		{
			ID:         r1id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			ResourceType: "POST",
			ResourceID:   "post-42",
			ResourceMeta: map[string]string{
				"category": "funny",
			},
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:  r1id,
					ProjectID: testProjectID,
					Name:      "text",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"My windows aren’t dirty, that’s my dog’s nose art."`),
				},
				{
					RecordID:  r1id,
					ProjectID: testProjectID,
					Name:      "status",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"published"`),
				},
			},
			OperationType:   "CREATE",
			OperationID:     "example.v1.PostService/CreatePost",
			OperationTime:   time.Date(2023, 1, 1, 1, 1, 0, 0, time.UTC),
			OperationStatus: aud.OperationStatusSucceeded.Int(),
			ActorType:       "USER",
			ActorID:         "user-82",
		},
		{
			ID:         r2id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			ResourceType: "COMMENT",
			ResourceID:   "comment-79",
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:  r2id,
					ProjectID: testProjectID,
					Name:      "text",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"Show us, my fiend!"`),
				},
				{
					RecordID:  r2id,
					ProjectID: testProjectID,
					Name:      "status",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"published"`),
				},
			},
			OperationType:   "CREATE",
			OperationID:     "example.v1.PostService/CreatePostComment",
			OperationTime:   time.Date(2023, 1, 1, 1, 2, 0, 0, time.UTC),
			OperationStatus: aud.OperationStatusSucceeded.Int(),
			ActorType:       "USER",
			ActorID:         "user-83",
		},
		{
			ID:         r3id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			ResourceType: "COMMENT",
			ResourceID:   "comment-79",
			ResourceMeta: map[string]string{
				"status": "published",
			},
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:    r3id,
					ProjectID:   testProjectID,
					Name:        "text",
					Description: "Edit text",
					OldValue:    json.RawMessage(`"Show us, my fiend!"`),
					NewValue:    json.RawMessage(`"Show us, my friend!"`),
				},
			},
			OperationType:   "UPDATE",
			OperationID:     "example.v1.PostService/UpdatePostComment",
			OperationTime:   time.Date(2023, 1, 1, 1, 3, 0, 0, time.UTC),
			OperationStatus: aud.OperationStatusSucceeded.Int(),
			ActorType:       "USER",
			ActorID:         "user-83",
		},
		{
			ID:         r4id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-55",
			},
			ResourceType: "POST",
			ResourceID:   "post-55",
			ResourceMeta: map[string]string{
				"status": "draft",
			},
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:    r4id,
					ProjectID:   testProjectID,
					Name:        "text",
					Description: "Edit text",
					OldValue:    json.RawMessage(`"The dog knows the best seat in the house."`),
					NewValue:    json.RawMessage(`"For the best seat in the house, you’ll have to move the dog."`),
				},
			},
			OperationType:   "UPDATE",
			OperationID:     "example.v1.PostService/UpdatePost",
			OperationTime:   time.Date(2023, 1, 1, 1, 4, 0, 0, time.UTC),
			OperationStatus: aud.OperationStatusSucceeded.Int(),
			ActorType:       "USER",
			ActorID:         "user-83",
		},
		{
			ID:         r5id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			ResourceType: "POST",
			ResourceID:   "post-42",
			ResourceMeta: map[string]string{
				"category": "funny",
			},
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:    r5id,
					ProjectID:   testProjectID,
					Name:        "status",
					Description: "Unpublish post",
					OldValue:    json.RawMessage(`"published"`),
					NewValue:    json.RawMessage(`"draft"`),
				},
			},
			OperationType: "UPDATE",
			OperationID:   "example.v1.PostService/UpdatePost",
			OperationTime: time.Date(2023, 1, 1, 1, 5, 0, 0, time.UTC),
			OperationMeta: map[string]string{
				"failure_reason": "Permission Denied",
			},
			OperationStatus: aud.OperationStatusFailed.Int(),
			ActorType:       "USER",
			ActorID:         "user-10",
			ActorMeta: map[string]string{
				"as": "reporter",
			},
		},
		{
			ID:         r6id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			ResourceType: "POST",
			ResourceID:   "post-42",
			ResourceMeta: map[string]string{
				"category": "funny",
			},
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:    r6id,
					ProjectID:   testProjectID,
					Name:        "status",
					Description: "Unpublish post",
					OldValue:    json.RawMessage(`"published"`),
					NewValue:    json.RawMessage(`"draft"`),
				},
			},
			OperationType: "UPDATE",
			OperationID:   "example.v1.PostService/UpdatePost",
			OperationTime: time.Date(2023, 1, 1, 1, 6, 0, 0, time.UTC),
			OperationMeta: map[string]string{
				"via":    "Moderator UI",
				"reason": "The post is not fun enough. GIF meme is required!",
			},
			OperationStatus: aud.OperationStatusSucceeded.Int(),
			ActorType:       "USER",
			ActorID:         "user-5",
			ActorMeta: map[string]string{
				"as": "moderator",
			},
		},
	}
	seedRecords(ctx, t, db, seededRecordModels...)
	setCleanupRecords(t, db)

	// Test

	t.Run("Should list records - without filter", func(t *testing.T) {
		store := NewStore(db)

		filter := aud.RecordFilter{}
		limit := int32(10)
		pag := aud.RecordCursor{}

		records, err := store.ListRecords(ctx, testProjectID, filter, limit, pag)
		assert.NoError(t, err)

		assert.Equal(t, []aud.Record{
			fromRecordModel(seededRecordModels[5]),
			fromRecordModel(seededRecordModels[4]),
			fromRecordModel(seededRecordModels[3]),
			fromRecordModel(seededRecordModels[2]),
			fromRecordModel(seededRecordModels[1]),
			fromRecordModel(seededRecordModels[0]),
		}, records)
	})

	t.Run("Should list records - filter by label", func(t *testing.T) {
		store := NewStore(db)

		filter := aud.RecordFilter{
			Labels: map[string]string{
				"post_id": "post-42",
			},
		}
		limit := int32(10)
		pag := aud.RecordCursor{}

		records, err := store.ListRecords(ctx, testProjectID, filter, limit, pag)
		assert.NoError(t, err)

		assert.Equal(t, []aud.Record{
			fromRecordModel(seededRecordModels[5]),
			fromRecordModel(seededRecordModels[4]),
			fromRecordModel(seededRecordModels[2]),
			fromRecordModel(seededRecordModels[1]),
			fromRecordModel(seededRecordModels[0]),
		}, records)
	})

	t.Run("Should list records - filter by resource", func(t *testing.T) {
		store := NewStore(db)

		filter := aud.RecordFilter{
			ResourceType: "POST",
			ResourceID:   "post-42",
		}
		limit := int32(10)
		pag := aud.RecordCursor{}

		records, err := store.ListRecords(ctx, testProjectID, filter, limit, pag)
		assert.NoError(t, err)

		assert.Equal(t, []aud.Record{
			fromRecordModel(seededRecordModels[5]),
			fromRecordModel(seededRecordModels[4]),
			fromRecordModel(seededRecordModels[0]),
		}, records)
	})

	t.Run("Should list records - filter by operation", func(t *testing.T) {
		store := NewStore(db)

		filter := aud.RecordFilter{
			OperationType: "UPDATE",
			OperationID:   "example.v1.PostService/UpdatePost",
		}
		limit := int32(10)
		pag := aud.RecordCursor{}

		records, err := store.ListRecords(ctx, testProjectID, filter, limit, pag)
		assert.NoError(t, err)

		assert.Equal(t, []aud.Record{
			fromRecordModel(seededRecordModels[5]),
			fromRecordModel(seededRecordModels[4]),
			fromRecordModel(seededRecordModels[3]),
		}, records)
	})

	t.Run("Should list records - filter by operation time", func(t *testing.T) {
		store := NewStore(db)

		filter := aud.RecordFilter{
			OperationTimeFrom: time.Date(2023, 1, 1, 1, 2, 0, 0, time.UTC),
			OperationTimeTo:   time.Date(2023, 1, 1, 1, 4, 0, 0, time.UTC),
		}
		limit := int32(10)
		pag := aud.RecordCursor{}

		records, err := store.ListRecords(ctx, testProjectID, filter, limit, pag)
		assert.NoError(t, err)

		assert.Equal(t, []aud.Record{
			fromRecordModel(seededRecordModels[2]),
			fromRecordModel(seededRecordModels[1]),
		}, records)
	})

	t.Run("Should list records - filter by actor", func(t *testing.T) {
		store := NewStore(db)

		filter := aud.RecordFilter{
			ActorType: "USER",
			ActorID:   "user-83",
		}
		limit := int32(10)
		pag := aud.RecordCursor{}

		records, err := store.ListRecords(ctx, testProjectID, filter, limit, pag)
		assert.NoError(t, err)

		assert.Equal(t, []aud.Record{
			fromRecordModel(seededRecordModels[3]),
			fromRecordModel(seededRecordModels[2]),
			fromRecordModel(seededRecordModels[1]),
		}, records)
	})
}

func TestIntegration_Store_UpdateRecord(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := sqltest.NewDatabase(ctx, t)

	// Seed

	seedTestProject(ctx, t, db)
	setCleanupTestProject(t, db)

	r1id := aud.MustNewID()
	r2id := aud.MustNewID()
	r3id := aud.MustNewID()

	seededRecordModels := []recordModel{
		{
			ID:         r1id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			ResourceType: "POST",
			ResourceID:   "post-42",
			ResourceMeta: map[string]string{
				"category": "funny",
			},
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:  r1id,
					ProjectID: testProjectID,
					Name:      "text",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"My windows aren’t dirty, that’s my dog’s nose art."`),
				},
				{
					RecordID:  r1id,
					ProjectID: testProjectID,
					Name:      "status",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"published"`),
				},
			},
			OperationType:   "CREATE",
			OperationID:     "example.v1.PostService/CreatePost",
			OperationTime:   time.Date(2023, 1, 1, 1, 1, 0, 0, time.UTC),
			OperationStatus: aud.OperationStatusSucceeded.Int(),
			ActorType:       "USER",
			ActorID:         "user-82",
		},
		{
			ID:         r2id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			ResourceType: "COMMENT",
			ResourceID:   "comment-79",
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:  r2id,
					ProjectID: testProjectID,
					Name:      "text",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"Show us, my fiend!"`),
				},
				{
					RecordID:  r2id,
					ProjectID: testProjectID,
					Name:      "status",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"published"`),
				},
			},
			OperationType:   "CREATE",
			OperationID:     "example.v1.PostService/CreatePostComment",
			OperationTime:   time.Date(2023, 1, 1, 1, 2, 0, 0, time.UTC),
			OperationStatus: aud.OperationStatusSucceeded.Int(),
			ActorType:       "USER",
			ActorID:         "user-83",
		},
		{
			ID:         r3id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			ResourceType: "COMMENT",
			ResourceID:   "comment-79",
			ResourceMeta: map[string]string{
				"status": "published",
			},
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:    r3id,
					ProjectID:   testProjectID,
					Name:        "text",
					Description: "Edit text",
					OldValue:    json.RawMessage(`"Show us, my fiend!"`),
					NewValue:    json.RawMessage(`"Show us, my friend!"`),
				},
			},
			OperationType:   "UPDATE",
			OperationID:     "example.v1.PostService/UpdatePostComment",
			OperationTime:   time.Date(2023, 1, 1, 1, 3, 0, 0, time.UTC),
			OperationStatus: aud.OperationStatusSucceeded.Int(),
			ActorType:       "USER",
			ActorID:         "user-83",
		},
	}
	seedRecords(ctx, t, db, seededRecordModels...)
	setCleanupRecords(t, db)

	// Test

	t.Run("Should update record", func(t *testing.T) {
		store := NewStore(db)

		id := r2id

		update := aud.RecordUpdate{
			Labels: map[string]string{
				"post_id": "post-43",
			},
			UpdateLabels: true,
			Resource: aud.Resource{
				Type: "COMMENT",
				ID:   "COMMENT-80",
				Metadata: map[string]string{
					"status": "published",
				},
				Changes: []aud.ResourceChange{
					{
						Name:     "text",
						OldValue: json.RawMessage(`null`),
						NewValue: json.RawMessage(`"Please show us, my friend!"`),
					},
				},
			},
			UpdateResource: true,
			Operation: aud.Operation{
				Type: "CREATE",
				ID:   "example.v2.PostService/CreatePostComment",
				Time: time.Date(2023, 2, 1, 1, 2, 0, 0, time.UTC),
				Metadata: map[string]string{
					"authorized": "true",
				},
				TraceContext: aud.TraceContext{
					Traceparent: "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01",
					Tracestate:  "congo=t61rcWkgMzE",
				},
				Status: aud.OperationStatusSucceeded,
			},
			UpdateOperation: true,
			Actor: aud.Actor{
				Type: "USER",
				ID:   "user-84",
				Metadata: map[string]string{
					"device_type": "mobile",
				},
			},
			UpdateActor: true,
		}

		updatedRecord, err := store.UpdateRecord(ctx, testProjectID, id, update)
		assert.NoError(t, err)

		assert.Equal(t, aud.Record{
			ID:         r2id,
			ProjectID:  testProjectID,
			CreateTime: seededRecordModels[1].CreateTime,
			Labels:     update.Labels,
			Resource:   update.Resource,
			Operation:  update.Operation,
			Actor:      update.Actor,
		}, updatedRecord)
	})

	t.Run("Should return error when record does not exist", func(t *testing.T) {
		store := NewStore(db)

		id := aud.MustNewID()

		update := aud.RecordUpdate{
			Labels: map[string]string{
				"post_id": "post-43",
			},
			UpdateLabels: true,
		}

		_, err := store.UpdateRecord(ctx, testProjectID, id, update)
		assert.ErrorIs(t, err, aud.ErrRecordNotFound)
	})
}

func TestIntegration_Store_DeleteRecord(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := sqltest.NewDatabase(ctx, t)

	// Seed

	seedTestProject(ctx, t, db)
	setCleanupTestProject(t, db)

	r1id := aud.MustNewID()
	r2id := aud.MustNewID()
	r3id := aud.MustNewID()

	seededRecordModels := []recordModel{
		{
			ID:         r1id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			ResourceType: "POST",
			ResourceID:   "post-42",
			ResourceMeta: map[string]string{
				"category": "funny",
			},
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:  r1id,
					ProjectID: testProjectID,
					Name:      "text",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"My windows aren’t dirty, that’s my dog’s nose art."`),
				},
				{
					RecordID:  r1id,
					ProjectID: testProjectID,
					Name:      "status",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"published"`),
				},
			},
			OperationType:   "CREATE",
			OperationID:     "example.v1.PostService/CreatePost",
			OperationTime:   time.Date(2023, 1, 1, 1, 1, 0, 0, time.UTC),
			OperationStatus: aud.OperationStatusSucceeded.Int(),
			ActorType:       "USER",
			ActorID:         "user-82",
		},
		{
			ID:         r2id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			ResourceType: "COMMENT",
			ResourceID:   "comment-79",
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:  r2id,
					ProjectID: testProjectID,
					Name:      "text",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"Show us, my fiend!"`),
				},
				{
					RecordID:  r2id,
					ProjectID: testProjectID,
					Name:      "status",
					OldValue:  json.RawMessage(`null`),
					NewValue:  json.RawMessage(`"published"`),
				},
			},
			OperationType:   "CREATE",
			OperationID:     "example.v1.PostService/CreatePostComment",
			OperationTime:   time.Date(2023, 1, 1, 1, 2, 0, 0, time.UTC),
			OperationStatus: aud.OperationStatusSucceeded.Int(),
			ActorType:       "USER",
			ActorID:         "user-83",
		},
		{
			ID:         r3id,
			ProjectID:  testProjectID,
			CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 0, time.UTC),
			Labels: map[string]string{
				"post_id": "post-42",
			},
			ResourceType: "COMMENT",
			ResourceID:   "comment-79",
			ResourceMeta: map[string]string{
				"status": "published",
			},
			ResourceChanges: []recordResourceChangeModel{
				{
					RecordID:    r3id,
					ProjectID:   testProjectID,
					Name:        "text",
					Description: "Edit text",
					OldValue:    json.RawMessage(`"Show us, my fiend!"`),
					NewValue:    json.RawMessage(`"Show us, my friend!"`),
				},
			},
			OperationType:   "UPDATE",
			OperationID:     "example.v1.PostService/UpdatePostComment",
			OperationTime:   time.Date(2023, 1, 1, 1, 3, 0, 0, time.UTC),
			OperationStatus: aud.OperationStatusSucceeded.Int(),
			ActorType:       "USER",
			ActorID:         "user-83",
		},
	}
	seedRecords(ctx, t, db, seededRecordModels...)
	setCleanupRecords(t, db)

	// Test

	t.Run("Should delete record", func(t *testing.T) {
		store := NewStore(db)

		id := r2id

		err := store.DeleteRecord(ctx, testProjectID, id)
		assert.NoError(t, err)

		// Check if record was deleted.

		var model recordModel
		err = db.NewSelect().
			Model(&model).
			Relation(relationResourceChanges).
			Where("id = ?", id).
			Scan(ctx)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestIntegration_Store_recordsIDs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := sqltest.NewDatabase(ctx, t)

	// Seed

	seedTestProject(ctx, t, db)
	setCleanupTestProject(t, db)

	setCleanupRecords(t, db)

	// Test

	t.Run("Created records should maintain order of ids", func(t *testing.T) {
		store := NewStore(db)

		const size = 5

		records := make([]aud.Record, size)
		for i := 0; i < size; i++ {
			rec := aud.Record{
				ID:         aud.MustNewID(),
				ProjectID:  testProjectID,
				CreateTime: time.Date(2023, 1, 1, 2, 3, 4, 5, time.UTC),
				Labels: map[string]string{
					"i": fmt.Sprintf("%d", i),
				},
				Resource: aud.Resource{
					Type: "COMMENT",
					ID:   "comment-1",
				},
				Operation: aud.Operation{
					Type: "UPDATE",
					ID:   "example.v1.PostService/UpdatePostComment",
					Time: time.Date(2023, 1, 1, 2, 1, 0, 0, time.UTC),
				},
				Actor: aud.Actor{
					Type: "USER",
					ID:   "user-1",
				},
			}

			records[i] = rec

			err := store.CreateRecord(ctx, rec)
			assert.NoError(t, err)
		}

		var models []recordModel
		err := db.NewSelect().
			Model(&models).
			Where("project_id = ?", testProjectID).
			Order("id ASC").
			Scan(ctx)
		assert.NoError(t, err)

		if assert.Len(t, models, 5) {
			for i, model := range models {
				assert.Equal(t, records[i].ID, model.ID)
				assert.Equal(t, fmt.Sprintf("%d", i), model.Labels["i"])
			}
		}
	})
}

var (
	testProjectID              = aud.MustNewID()
	testProjectPartitionNumber = int32(1)
)

func seedTestProject(ctx context.Context, t *testing.T, db *bun.DB) {
	t.Helper()

	projectMod := projectModel{
		ID:              testProjectID,
		CreateTime:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		PartitionNumber: testProjectPartitionNumber,
		DisplayName:     "Test Project",
	}

	seedProjects(ctx, t, db, projectMod)
}

func setCleanupTestProject(t *testing.T, db *bun.DB) {
	t.Helper()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		deleteProject(ctx, t, db, projectModel{
			ID:              testProjectID,
			PartitionNumber: testProjectPartitionNumber,
		})
	})
}

func seedProjects(ctx context.Context, t *testing.T, db *bun.DB, models ...projectModel) {
	t.Helper()

	for _, model := range models {
		model := model

		_, err := db.NewInsert().
			Model(&model).
			Exec(ctx)
		require.NoError(t, err)

		if db.Dialect().Name() == dialect.PG {
			err = createTablePartitionForProject(
				ctx,
				db,
				tableNameRecords,
				model.ID,
				model.PartitionNumber,
			)
			require.NoError(t, err)

			err = createTablePartitionForProject(
				ctx,
				db,
				tableNameRecordsResourceChanges,
				model.ID,
				model.PartitionNumber,
			)
			require.NoError(t, err)
		}
	}
}

func setCleanupProjects(t *testing.T, db *bun.DB) {
	t.Helper()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cleanupProjects(ctx, t, db)
	})
}

func cleanupProjects(ctx context.Context, t *testing.T, db *bun.DB) {
	t.Helper()

	var models []projectModel

	err := db.NewSelect().
		Model(&models).
		Scan(ctx)
	require.NoError(t, err)

	for _, model := range models {
		deleteProject(ctx, t, db, model)
	}
}

func deleteProject(ctx context.Context, t *testing.T, db *bun.DB, model projectModel) {
	t.Helper()

	if db.Dialect().Name() == dialect.PG {
		err := dropTablePartitionForProject(
			ctx,
			db,
			tableNameRecordsResourceChanges,
			model.PartitionNumber,
		)
		require.NoError(t, err)

		err = dropTablePartitionForProject(
			ctx,
			db,
			tableNameRecords,
			model.PartitionNumber,
		)
		require.NoError(t, err)
	}

	_, err := db.NewDelete().
		Model(&projectModel{}).
		Where("id = ?", model.ID).
		Exec(ctx)
	require.NoError(t, err)
}

func seedRecords(ctx context.Context, t *testing.T, db *bun.DB, recordModels ...recordModel) {
	t.Helper()

	for _, recordMod := range recordModels {
		recordMod := recordMod

		_, err := db.NewInsert().
			Model(&recordMod).
			Exec(ctx)
		require.NoError(t, err)

		_, err = db.NewInsert().
			Model(&recordMod.ResourceChanges).
			Exec(ctx)
		require.NoError(t, err)
	}
}

func setCleanupRecords(t *testing.T, db *bun.DB) {
	t.Helper()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cleanupRecords(ctx, t, db)
	})
}

func cleanupRecords(ctx context.Context, t *testing.T, db *bun.DB) {
	t.Helper()

	_, err := db.NewTruncateTable().
		Model((*recordModel)(nil)).
		Cascade().
		Exec(ctx)
	require.NoError(t, err)
}
