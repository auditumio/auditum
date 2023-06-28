package auditumv1alpha1

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	auditumv1alpha1 "github.com/infragmo/auditum/api/gen/go/infragmo/auditum/v1alpha1"
	"github.com/infragmo/auditum/internal/aud"
)

func decodeRecords(projectID string, src []*auditumv1alpha1.Record, restrictions aud.RecordsRestrictions) ([]aud.Record, error) {
	var err error

	dst := make([]aud.Record, len(src))
	for i := range src {
		src[i].ProjectId = projectID
		dst[i], err = decodeRecord(src[i], restrictions)
		if err != nil {
			return nil, fmt.Errorf(`invalid "records[%d]": %v`, i, err)
		}
	}

	return dst, nil
}

func decodeRecord(src *auditumv1alpha1.Record, restrictions aud.RecordsRestrictions) (dst aud.Record, err error) {
	id, err := decodeIDOptional(src.GetId())
	if err != nil {
		return dst, fmt.Errorf(`invalid "id": %v`, err)
	}

	projectID, err := decodeID(src.GetProjectId())
	if err != nil {
		return dst, fmt.Errorf(`invalid "project_id": %v`, err)
	}

	labels, err := decodeLabels(src.GetLabels(), restrictions.Labels)
	if err != nil {
		return dst, fmt.Errorf(`invalid "labels": %v`, err)
	}

	resource, err := decodeResource(src.GetResource(), restrictions.Resource)
	if err != nil {
		return dst, fmt.Errorf(`invalid "resource": %v`, err)
	}

	operation, err := decodeOperation(src.GetOperation(), restrictions.Operation)
	if err != nil {
		return dst, fmt.Errorf(`invalid "operation": %v`, err)
	}

	actor, err := decodeActor(src.GetActor(), restrictions.Actor)
	if err != nil {
		return dst, fmt.Errorf(`invalid "actor": %v`, err)
	}

	return aud.Record{
		ID:         id,
		ProjectID:  projectID,
		CreateTime: time.Time{}, // Ignored as OUTPUT_ONLY.
		Labels:     labels,
		Resource:   resource,
		Operation:  operation,
		Actor:      actor,
	}, nil
}

func decodeLabels(src map[string]string, restrictions aud.RestrictionsKeyValue) (map[string]string, error) {
	if len(src) == 0 {
		return nil, nil
	}

	if err := validateLabelsOrMetadata(src, restrictions); err != nil {
		return nil, err
	}

	return src, nil
}

func decodeResource(src *auditumv1alpha1.Resource, restrictions aud.RecordsRestrictionsResource) (dst aud.Resource, err error) {
	if src == nil {
		return dst, fmt.Errorf("must not be empty")
	}

	typ, err := decodeResourceType(src.GetType(), restrictions.Type)
	if err != nil {
		return dst, fmt.Errorf(`invalid "type": %v`, err)
	}

	id, err := decodeResourceID(src.GetId(), restrictions.ID)
	if err != nil {
		return dst, fmt.Errorf(`invalid "id": %v`, err)
	}

	meta, err := decodeResourceMetadata(src.GetMetadata(), restrictions.Metadata)
	if err != nil {
		return dst, fmt.Errorf(`invalid "metadata": %v`, err)
	}

	changes, err := decodeResourceChanges(src.GetChanges(), restrictions.Changes)
	if err != nil {
		return dst, fmt.Errorf(`invalid "changes": %v`, err)
	}

	return aud.Resource{
		Type:     typ,
		ID:       id,
		Metadata: meta,
		Changes:  changes,
	}, nil
}

func decodeResourceType(src string, restrictions aud.RestrictionsString) (string, error) {
	if err := validateResourceType(src, restrictions); err != nil {
		return "", err
	}

	return src, nil
}

func decodeResourceID(src string, restrictions aud.RestrictionsString) (string, error) {
	if err := validateResourceID(src, restrictions); err != nil {
		return "", err
	}

	return src, nil
}

func decodeResourceMetadata(src map[string]string, restrictions aud.RestrictionsKeyValue) (map[string]string, error) {
	if err := validateLabelsOrMetadata(src, restrictions); err != nil {
		return nil, err
	}

	return src, nil
}

func decodeResourceChanges(src []*auditumv1alpha1.ResourceChange, restrictions aud.RecordsRestrictionsResourceChanges) (dst []aud.ResourceChange, err error) {
	if len(src) == 0 {
		return nil, nil
	}

	if err := validateResourceChanges(src, restrictions); err != nil {
		return nil, err
	}

	dst = make([]aud.ResourceChange, len(src))
	for i := range src {
		dst[i], err = decodeResourceChange(src[i], restrictions)
		if err != nil {
			return nil, fmt.Errorf(`invalid "changes[%d]": %v`, i, err)
		}
	}
	return dst, nil
}

func decodeResourceChange(src *auditumv1alpha1.ResourceChange, restrictions aud.RecordsRestrictionsResourceChanges) (dst aud.ResourceChange, err error) {
	if src == nil {
		return dst, fmt.Errorf("must not be empty")
	}

	name, err := decodeResourceChangeName(src.GetName(), restrictions.Name)
	if err != nil {
		return dst, fmt.Errorf(`invalid "name": %v`, err)
	}

	desc, err := decodeResourceChangeDescription(src.GetDescription(), restrictions.Description)
	if err != nil {
		return dst, fmt.Errorf(`invalid "description": %v`, err)
	}

	ovb, err := decodeResourceChangeValue(src.GetOldValue(), restrictions.OldValue)
	if err != nil {
		return dst, fmt.Errorf(`invalid "old_value": %v`, err)
	}

	nvb, err := decodeResourceChangeValue(src.GetNewValue(), restrictions.NewValue)
	if err != nil {
		return dst, fmt.Errorf(`invalid "new_value": %v`, err)
	}

	return aud.ResourceChange{
		Name:        name,
		Description: desc,
		OldValue:    ovb,
		NewValue:    nvb,
	}, nil
}

func decodeResourceChangeName(src string, restrictions aud.RestrictionsString) (string, error) {
	if err := validateResourceChangeName(src, restrictions); err != nil {
		return "", err
	}

	return src, nil
}

func decodeResourceChangeDescription(src string, restrictions aud.RestrictionsString) (string, error) {
	if err := validateResourceChangeDescription(src, restrictions); err != nil {
		return "", err
	}

	return src, nil
}

func decodeResourceChangeValue(src *structpb.Value, restrictions aud.RestrictionsBytes) ([]byte, error) {
	// NOTE: value is optional.
	if src == nil {
		return nil, nil
	}

	b, err := src.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("invalid json")
	}

	if err := validateResourceChangeValue(b, restrictions); err != nil {
		return nil, err
	}

	return b, nil
}

func decodeOperation(src *auditumv1alpha1.Operation, restrictions aud.RecordsRestrictionsOperation) (dst aud.Operation, err error) {
	if src == nil {
		return dst, fmt.Errorf("must not be empty")
	}

	typ, err := decodeOperationType(src.GetType(), restrictions.Type)
	if err != nil {
		return dst, fmt.Errorf(`invalid "type": %v`, err)
	}

	id, err := decodeOperationID(src.GetId(), restrictions.ID)
	if err != nil {
		return dst, fmt.Errorf(`invalid "id": %v`, err)
	}

	ts, err := decodeOperationTime(src.GetTime())
	if err != nil {
		return dst, fmt.Errorf(`invalid "time": %v`, err)
	}

	meta, err := decodeOperationMetadata(src.GetMetadata(), restrictions.Metadata)
	if err != nil {
		return dst, fmt.Errorf(`invalid "metadata": %v`, err)
	}

	traceContext, err := decodeTraceContext(src.GetTraceContext())
	if err != nil {
		return dst, fmt.Errorf(`invalid "trace_context": %v`, err)
	}

	status := decodeOperationStatus(src.GetStatus())

	return aud.Operation{
		Type:         typ,
		ID:           id,
		Time:         ts,
		Metadata:     meta,
		TraceContext: traceContext,
		Status:       status,
	}, nil
}

func decodeOperationType(src string, restrictions aud.RestrictionsString) (string, error) {
	if err := validateOperationType(src, restrictions); err != nil {
		return "", err
	}

	return src, nil
}

func decodeOperationID(src string, restrictions aud.RestrictionsString) (string, error) {
	if err := validateOperationID(src, restrictions); err != nil {
		return "", err
	}

	return src, nil
}

func decodeOperationTime(src *timestamppb.Timestamp) (time.Time, error) {
	if err := validateOperationTime(src); err != nil {
		return time.Time{}, err
	}

	return src.AsTime(), nil
}

func decodeOperationMetadata(src map[string]string, restrictions aud.RestrictionsKeyValue) (map[string]string, error) {
	if err := validateLabelsOrMetadata(src, restrictions); err != nil {
		return nil, err
	}

	return src, nil
}

func decodeTraceContext(src *auditumv1alpha1.TraceContext) (dst aud.TraceContext, err error) {
	if err := validateTraceContext(src); err != nil {
		return dst, err
	}

	return aud.TraceContext{
		Traceparent: src.GetTraceparent(),
		Tracestate:  src.GetTracestate(),
	}, nil
}

func decodeOperationStatus(src auditumv1alpha1.OperationStatus_Enum) aud.OperationStatus {
	// NOTE: no validation for status, as at this layer it must be always valid
	// (either existing enum value or UNSPECIFIED).

	switch src {
	case auditumv1alpha1.OperationStatus_SUCCEEDED:
		return aud.OperationStatusSucceeded
	case auditumv1alpha1.OperationStatus_FAILED:
		return aud.OperationStatusFailed
	default:
		return aud.OperationStatusUnspecified
	}
}

func decodeActor(src *auditumv1alpha1.Actor, restrictions aud.RecordsRestrictionsActor) (dst aud.Actor, err error) {
	if src == nil {
		return dst, fmt.Errorf("must not be empty")
	}

	typ, err := decodeActorType(src.GetType(), restrictions.Type)
	if err != nil {
		return dst, fmt.Errorf(`invalid "type": %v`, err)
	}

	id, err := decodeActorID(src.GetId(), restrictions.ID)
	if err != nil {
		return dst, fmt.Errorf(`invalid "id": %v`, err)
	}

	meta, err := decodeActorMetadata(src.GetMetadata(), restrictions.Metadata)
	if err != nil {
		return dst, fmt.Errorf(`invalid "metadata": %v`, err)
	}

	return aud.Actor{
		Type:     typ,
		ID:       id,
		Metadata: meta,
	}, nil
}

func decodeActorType(src string, restrictions aud.RestrictionsString) (string, error) {
	if err := validateActorType(src, restrictions); err != nil {
		return "", err
	}

	return src, nil
}

func decodeActorID(src string, restrictions aud.RestrictionsString) (string, error) {
	if err := validateActorID(src, restrictions); err != nil {
		return "", err
	}

	return src, nil
}

func decodeActorMetadata(src map[string]string, restrictions aud.RestrictionsKeyValue) (map[string]string, error) {
	if err := validateLabelsOrMetadata(src, restrictions); err != nil {
		return nil, err
	}

	return src, nil
}

func encodeRecords(src []aud.Record) []*auditumv1alpha1.Record {
	dst := make([]*auditumv1alpha1.Record, len(src))
	for i := range src {
		dst[i] = encodeRecord(src[i])
	}
	return dst
}

func encodeRecord(src aud.Record) *auditumv1alpha1.Record {
	return &auditumv1alpha1.Record{
		Id:         src.ID.String(),
		ProjectId:  src.ProjectID.String(),
		CreateTime: timestamppb.New(src.CreateTime),
		Labels:     src.Labels,
		Resource:   encodeResource(src.Resource),
		Operation:  encodeOperation(src.Operation),
		Actor:      encodeActor(src.Actor),
	}
}

func encodeResource(src aud.Resource) *auditumv1alpha1.Resource {
	return &auditumv1alpha1.Resource{
		Type:     src.Type,
		Id:       src.ID,
		Metadata: src.Metadata,
		Changes:  encodeResourceChanges(src.Changes),
	}
}

func encodeResourceChanges(src []aud.ResourceChange) []*auditumv1alpha1.ResourceChange {
	dst := make([]*auditumv1alpha1.ResourceChange, len(src))
	for i := range src {
		dst[i] = encodeResourceChange(src[i])
	}
	return dst
}

func encodeResourceChange(src aud.ResourceChange) *auditumv1alpha1.ResourceChange {
	var oldValue structpb.Value
	if err := oldValue.UnmarshalJSON(src.OldValue); err != nil {
		// This is exceptional.
		panic(fmt.Errorf("unmarshal old value from json: %v", err))
	}
	var newValue structpb.Value
	if err := newValue.UnmarshalJSON(src.NewValue); err != nil {
		// This is exceptional.
		panic(fmt.Errorf("unmarshal new value from json: %v", err))
	}

	return &auditumv1alpha1.ResourceChange{
		Name:        src.Name,
		Description: src.Description,
		OldValue:    &oldValue,
		NewValue:    &newValue,
	}
}

func encodeOperation(src aud.Operation) *auditumv1alpha1.Operation {
	var traceContext *auditumv1alpha1.TraceContext
	if !src.TraceContext.IsZero() {
		traceContext = &auditumv1alpha1.TraceContext{
			Traceparent: src.TraceContext.Traceparent,
			Tracestate:  src.TraceContext.Tracestate,
		}
	}

	return &auditumv1alpha1.Operation{
		Type:         src.Type,
		Id:           src.ID,
		Time:         timestamppb.New(src.Time),
		Metadata:     src.Metadata,
		TraceContext: traceContext,
		Status:       encodeOperationStatus(src.Status),
	}
}

func encodeOperationStatus(src aud.OperationStatus) auditumv1alpha1.OperationStatus_Enum {
	switch src {
	case aud.OperationStatusSucceeded:
		return auditumv1alpha1.OperationStatus_SUCCEEDED
	case aud.OperationStatusFailed:
		return auditumv1alpha1.OperationStatus_FAILED
	default:
		return auditumv1alpha1.OperationStatus_UNSPECIFIED
	}
}

func encodeActor(src aud.Actor) *auditumv1alpha1.Actor {
	return &auditumv1alpha1.Actor{
		Type:     src.Type,
		Id:       src.ID,
		Metadata: src.Metadata,
	}
}

func decodeRecordFilter(src *auditumv1alpha1.ListRecordsRequest_Filter) (dst aud.RecordFilter, err error) {
	var operationTimeFrom time.Time
	if v := src.GetOperationTimeFrom(); v != nil {
		if !v.IsValid() {
			return dst, fmt.Errorf(`invalid "operation_time_from" time value`)
		}
		operationTimeFrom = v.AsTime()
	}

	var operationTimeTo time.Time
	if v := src.GetOperationTimeTo(); v != nil {
		if !v.IsValid() {
			return dst, fmt.Errorf(`invalid "operation_time_to" time value`)
		}
		operationTimeTo = v.AsTime()
	}

	return aud.RecordFilter{
		Labels:            src.GetLabels(),
		ResourceType:      src.GetResourceType(),
		ResourceID:        src.GetResourceId(),
		OperationType:     src.GetOperationType(),
		OperationID:       src.GetOperationId(),
		OperationTimeFrom: operationTimeFrom,
		OperationTimeTo:   operationTimeTo,
		ActorType:         src.GetActorType(),
		ActorID:           src.GetActorId(),
	}, nil
}
