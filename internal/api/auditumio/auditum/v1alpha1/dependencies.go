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

package auditumv1alpha1

import (
	"context"

	"github.com/auditumio/auditum/internal/aud"
)

type Store interface {
	CreateProject(ctx context.Context, project aud.Project) error

	GetProject(ctx context.Context, id aud.ID) (aud.Project, error)

	ListProjects(
		ctx context.Context,
		limit int32,
		cursor aud.ProjectCursor,
	) ([]aud.Project, error)

	UpdateProject(
		ctx context.Context,
		projectID aud.ID,
		update aud.ProjectUpdate,
	) (aud.Project, error)

	CreateRecord(ctx context.Context, record aud.Record) error

	CreateRecords(ctx context.Context, records []aud.Record) error

	GetRecord(
		ctx context.Context,
		projectID aud.ID,
		id aud.ID,
	) (aud.Record, error)

	ListRecords(
		ctx context.Context,
		projectID aud.ID,
		filter aud.RecordFilter,
		limit int32,
		cursor aud.RecordCursor,
	) ([]aud.Record, error)

	UpdateRecord(
		ctx context.Context,
		projectID aud.ID,
		id aud.ID,
		update aud.RecordUpdate,
	) (aud.Record, error)

	DeleteRecord(ctx context.Context, projectID aud.ID, id aud.ID) error
}
