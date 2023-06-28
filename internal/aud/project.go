package aud

import (
	"time"

	"github.com/infragmo/auditum/internal/aud/types"
)

type Project struct {
	ID                  ID
	CreateTime          time.Time
	DisplayName         string
	UpdateRecordEnabled types.BoolValue
	DeleteRecordEnabled types.BoolValue
}
