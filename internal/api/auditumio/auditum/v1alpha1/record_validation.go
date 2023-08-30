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
	"regexp"

	"google.golang.org/protobuf/types/known/timestamppb"

	auditumv1alpha1 "github.com/auditumio/auditum/api/gen/go/auditumio/auditum/v1alpha1"
	"github.com/auditumio/auditum/internal/aud"
	"github.com/auditumio/auditum/internal/util/tracecontext"
)

const keyRegexpTemplate = `[a-zA-Z0-9-_]+`

var keyRegexp = regexp.MustCompile(keyRegexpTemplate)

func validateLabelsOrMetadata(src map[string]string, restrictions aud.RestrictionsKeyValue) error {
	var totalSizeBytes int

	for key, value := range src {
		if err := validateLabelOrMetadataKey(key, restrictions); err != nil {
			return fmt.Errorf("key %q is invalid: %v", key, err)
		}
		if err := validateLabelOrMetadataValue(value, restrictions); err != nil {
			return fmt.Errorf("value for key %q is invalid: %v", key, err)
		}

		totalSizeBytes += len(key) + len(value)
	}

	if totalSizeBytes > restrictions.TotalMaxSizeBytes {
		return fmt.Errorf("total size of all keys and values must be at most %d bytes", restrictions.TotalMaxSizeBytes)
	}

	return nil
}

func validateLabelOrMetadataKey(src string, restrictions aud.RestrictionsKeyValue) error {
	if len(src) == 0 {
		return fmt.Errorf("must not be empty")
	}

	if !keyRegexp.MatchString(src) {
		return fmt.Errorf("must match regexp %s", keyRegexpTemplate)
	}

	if len(src) > restrictions.KeyMaxSizeBytes {
		return fmt.Errorf("must be at most %d bytes", restrictions.KeyMaxSizeBytes)
	}

	return nil
}

func validateLabelOrMetadataValue(src string, restrictions aud.RestrictionsKeyValue) error {
	// NOTE: value is optional.
	if len(src) == 0 {
		return nil
	}

	if len(src) > restrictions.ValueMaxSizeBytes {
		return fmt.Errorf("must be at most %d bytes", restrictions.ValueMaxSizeBytes)
	}

	return nil
}

func validateResourceType(src string, restrictions aud.RestrictionsString) error {
	if len(src) == 0 {
		return fmt.Errorf("must not be empty")
	}

	if len(src) > restrictions.MaxSizeBytes {
		return fmt.Errorf("must be at most %d bytes", restrictions.MaxSizeBytes)
	}

	return nil
}

func validateResourceID(src string, restrictions aud.RestrictionsString) error {
	if len(src) == 0 {
		return fmt.Errorf("must not be empty")
	}

	if len(src) > restrictions.MaxSizeBytes {
		return fmt.Errorf("must be at most %d bytes", restrictions.MaxSizeBytes)
	}

	return nil
}

func validateResourceChanges(src []*auditumv1alpha1.ResourceChange, restrictions aud.RecordsRestrictionsResourceChanges) error {
	if len(src) > restrictions.TotalMaxCount {
		return fmt.Errorf("must not exceed the limit of %d changes", restrictions.TotalMaxCount)
	}

	return nil
}

func validateResourceChangeName(src string, restrictions aud.RestrictionsString) error {
	if len(src) == 0 {
		return fmt.Errorf("must not be empty")
	}

	if len(src) > restrictions.MaxSizeBytes {
		return fmt.Errorf("must be at most %d bytes", restrictions.MaxSizeBytes)
	}

	return nil
}

func validateResourceChangeDescription(src string, restrictions aud.RestrictionsString) error {
	// NOTE: description is optional.
	if len(src) == 0 {
		return nil
	}

	if len(src) > restrictions.MaxSizeBytes {
		return fmt.Errorf("must be at most %d bytes", restrictions.MaxSizeBytes)
	}

	return nil
}

func validateResourceChangeValue(b []byte, restrictions aud.RestrictionsBytes) error {
	if len(b) > restrictions.MaxSizeBytes {
		return fmt.Errorf("must be at most %d bytes", restrictions.MaxSizeBytes)
	}

	return nil
}

func validateOperationType(src string, restrictions aud.RestrictionsString) error {
	if len(src) == 0 {
		return fmt.Errorf("must not be empty")
	}

	if len(src) > restrictions.MaxSizeBytes {
		return fmt.Errorf("must be at most %d bytes", restrictions.MaxSizeBytes)
	}

	return nil
}

func validateOperationID(src string, restrictions aud.RestrictionsString) error {
	if len(src) == 0 {
		return fmt.Errorf("must not be empty")
	}

	if len(src) > restrictions.MaxSizeBytes {
		return fmt.Errorf("must be at most %d bytes", restrictions.MaxSizeBytes)
	}

	return nil
}

func validateOperationTime(src *timestamppb.Timestamp) error {
	if src.GetSeconds() == 0 && src.GetNanos() == 0 {
		return fmt.Errorf("must not be empty")
	}

	if !src.IsValid() {
		return fmt.Errorf("must be valid time")
	}

	return nil
}

func validateTraceContext(src *auditumv1alpha1.TraceContext) error {
	if src == nil {
		return nil
	}

	if src.GetTraceparent() == "" && src.GetTracestate() != "" {
		return fmt.Errorf(`"tracestate" can be provided only if "traceparent" is provided`)
	}

	if err := validateTraceContextTraceparent(src.GetTraceparent()); err != nil {
		return fmt.Errorf(`invalid "traceparent": %v`, err)
	}

	if err := validateTraceContextTracestate(src.GetTracestate()); err != nil {
		return fmt.Errorf(`invalid "tracestate": %v`, err)
	}

	return nil
}

func validateTraceContextTraceparent(src string) error {
	if src == "" {
		return nil
	}

	if !tracecontext.TraceparentValid(src) {
		return fmt.Errorf("must be valid W3C traceparent")
	}

	return nil
}

func validateTraceContextTracestate(src string) error {
	if src == "" {
		return nil
	}

	const maxLen = 512
	if len(src) > maxLen {
		return fmt.Errorf("must not exceed maximum length of %d bytes", maxLen)
	}

	if !tracecontext.TracestateValid(src) {
		return fmt.Errorf("must be valid W3C tracestate")
	}

	return nil
}

func validateActorType(src string, restrictions aud.RestrictionsString) error {
	if len(src) == 0 {
		return fmt.Errorf("must not be empty")
	}

	if len(src) > restrictions.MaxSizeBytes {
		return fmt.Errorf("must be at most %d bytes", restrictions.MaxSizeBytes)
	}

	return nil
}

func validateActorID(src string, restrictions aud.RestrictionsString) error {
	if len(src) == 0 {
		return fmt.Errorf("must not be empty")
	}

	if len(src) > restrictions.MaxSizeBytes {
		return fmt.Errorf("must be at most %d bytes", restrictions.MaxSizeBytes)
	}

	return nil
}
