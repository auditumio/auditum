package auditumv1alpha1

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	auditumv1alpha1 "github.com/infragmo/auditum/api/gen/go/auditumio/auditum/v1alpha1"
	"github.com/infragmo/auditum/internal/aud"
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

	return aud.Project{
		ID:                  id,
		CreateTime:          time.Time{}, // Ignored as OUTPUT_ONLY.
		DisplayName:         displayName,
		UpdateRecordEnabled: decodeBoolValue(src.GetUpdateRecordEnabled()),
		DeleteRecordEnabled: decodeBoolValue(src.GetDeleteRecordEnabled()),
	}, nil
}

func decodeProjectDisplayName(src string) (string, error) {
	if err := validateProjectDisplayName(src); err != nil {
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
	}
}

func encodeProjects(src []aud.Project) []*auditumv1alpha1.Project {
	dst := make([]*auditumv1alpha1.Project, len(src))
	for i := range src {
		dst[i] = encodeProject(src[i])
	}
	return dst
}
