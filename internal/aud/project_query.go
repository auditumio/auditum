package aud

import "github.com/infragmo/auditum/internal/aud/types"

type ProjectCursor struct {
	LastID *ID `json:"lid,omitempty"`
}

func (p ProjectCursor) Empty() bool {
	return p.LastID == nil
}

func NewProjectCursor(projects []Project, pageSize int32) ProjectCursor {
	var cursor ProjectCursor

	if len(projects) >= int(pageSize) {
		last := projects[len(projects)-1]
		cursor.LastID = &last.ID
	}

	return cursor
}

type ProjectUpdate struct {
	DisplayName       string
	UpdateDisplayName bool

	UpdateRecordEnabled       types.BoolValue
	UpdateUpdateRecordEnabled bool

	DeleteRecordEnabled       types.BoolValue
	UpdateDeleteRecordEnabled bool
}
