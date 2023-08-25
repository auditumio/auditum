package aud

import (
	"time"

	"github.com/auditumio/auditum/internal/aud/types"
)

type Project struct {
	ID                  ID
	CreateTime          time.Time
	DisplayName         string
	UpdateRecordEnabled types.BoolValue
	DeleteRecordEnabled types.BoolValue
}
