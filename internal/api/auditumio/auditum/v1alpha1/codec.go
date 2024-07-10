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

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/auditumio/auditum/internal/aud"
	"github.com/auditumio/auditum/internal/aud/types"
)

// Common codec functions.

func decodeID(src string) (dst aud.ID, err error) {
	if src == "" {
		return dst, fmt.Errorf("must not be empty")
	}

	id, err := aud.ParseID(src)
	if err != nil {
		return dst, fmt.Errorf("must be a valid id")
	}

	return id, nil
}

func decodeIDOptional(src string) (dst aud.ID, err error) {
	if src == "" {
		return dst, nil
	}

	id, err := aud.ParseID(src)
	if err != nil {
		return dst, fmt.Errorf("must be a valid id")
	}

	return id, nil
}

func decodeBoolValue(src *wrapperspb.BoolValue) types.BoolValue {
	return types.BoolValue{
		Bool:  src.GetValue(),
		Valid: src != nil,
	}
}

func encodeBoolValue(src types.BoolValue) *wrapperspb.BoolValue {
	if !src.Valid {
		return nil
	}

	return wrapperspb.Bool(src.Bool)
}

func encodeOptionalString(src string) *string {
	if src == "" {
		return nil
	}

	return &src
}
