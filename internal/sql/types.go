package sql

import (
	"database/sql"

	"github.com/auditumio/auditum/internal/aud/types"
)

func toBoolValueModel(src types.BoolValue) sql.NullBool {
	return sql.NullBool{
		Bool:  src.Bool,
		Valid: src.Valid,
	}
}

func fromBoolValueModel(src sql.NullBool) types.BoolValue {
	return types.BoolValue{
		Bool:  src.Bool,
		Valid: src.Valid,
	}
}
