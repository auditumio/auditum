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

func toNullString(src string) sql.NullString {
	return sql.NullString{
		String: src,
		Valid:  src != "",
	}
}

func fromNullString(src sql.NullString) string {
	if !src.Valid {
		return ""
	}
	return src.String
}
