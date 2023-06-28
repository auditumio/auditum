package auditumv1alpha1

import (
	"fmt"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/infragmo/auditum/internal/aud"
	"github.com/infragmo/auditum/internal/aud/types"
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
