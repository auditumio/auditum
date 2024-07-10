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
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	auditumv1alpha1 "github.com/auditumio/auditum/api/gen/go/auditumio/auditum/v1alpha1"
	"github.com/auditumio/auditum/internal/aud"
)

func decodeProject(src *auditumv1alpha1.Project) (dst aud.Project, err error) {
	id, err := decodeIDOptional(src.GetId())
	if err != nil {
		return dst, fmt.Errorf(`invalid "id": %v`, err)
	}

	displayName, err := decodeProjectDisplayName(src.GetDisplayName())
	if err != nil {
		return dst, fmt.Errorf(`invalid "display_name": %v`, err)
	}

	externalID, err := decodeExternalID(src.GetExternalId())
	if err != nil {
		return dst, fmt.Errorf(`invalid "external_id": %v`, err)
	}

	return aud.Project{
		ID:                  id,
		CreateTime:          time.Time{}, // Ignored as OUTPUT_ONLY.
		DisplayName:         displayName,
		UpdateRecordEnabled: decodeBoolValue(src.GetUpdateRecordEnabled()),
		DeleteRecordEnabled: decodeBoolValue(src.GetDeleteRecordEnabled()),
		ExternalID:          externalID,
	}, nil
}

func decodeProjectDisplayName(src string) (string, error) {
	if err := validateProjectDisplayName(src); err != nil {
		return "", err
	}

	return src, nil
}

func decodeExternalID(src string) (string, error) {
	if err := validateProjectExternalID(src); err != nil {
		return "", err
	}

	return src, nil
}

func encodeProject(src aud.Project) *auditumv1alpha1.Project {
	return &auditumv1alpha1.Project{
		Id:                  src.ID.String(),
		CreateTime:          timestamppb.New(src.CreateTime),
		DisplayName:         src.DisplayName,
		UpdateRecordEnabled: encodeBoolValue(src.UpdateRecordEnabled),
		DeleteRecordEnabled: encodeBoolValue(src.DeleteRecordEnabled),
		ExternalId:          encodeOptionalString(src.ExternalID),
	}
}

func encodeProjects(src []aud.Project) []*auditumv1alpha1.Project {
	dst := make([]*auditumv1alpha1.Project, len(src))
	for i := range src {
		dst[i] = encodeProject(src[i])
	}
	return dst
}

func decodeProjectFilter(src *auditumv1alpha1.ListProjectsRequest_Filter) (dst aud.ProjectFilter) {
	return aud.ProjectFilter{
		ExternalIDs: src.GetExternalIds(),
	}
}
