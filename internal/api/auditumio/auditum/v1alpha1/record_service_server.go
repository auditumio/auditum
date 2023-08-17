package auditumv1alpha1

import (
	"context"
	"errors"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	auditumv1alpha1 "github.com/infragmo/auditum/api/gen/go/auditumio/auditum/v1alpha1"
	"github.com/infragmo/auditum/internal/aud"
	"github.com/infragmo/auditum/pkg/fragma/grpcx"
)

type RecordServiceServer struct {
	auditumv1alpha1.UnimplementedRecordServiceServer

	store    Store
	log      *zap.Logger
	settings aud.Settings

	id  func() aud.ID
	now func() time.Time
}

func NewRecordServiceServer(
	store Store,
	log *zap.Logger,
	settings aud.Settings,
) *RecordServiceServer {
	return &RecordServiceServer{
		store:    store,
		log:      log.Named("record_service_server"),
		settings: settings,
		id:       aud.MustNewID,
		now:      time.Now,
	}
}

func (s *RecordServiceServer) CreateRecord(
	ctx context.Context,
	req *auditumv1alpha1.CreateRecordRequest,
) (*auditumv1alpha1.CreateRecordResponse, error) {
	record, err := decodeRecord(req.GetRecord(), s.settings.Records.Restrictions)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "record": %v.`,
			err.Error(),
		)
	}

	record.ID = s.id()
	record.CreateTime = s.now().UTC()

	err = s.store.CreateRecord(ctx, record)
	if errors.Is(err, aud.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "Project not found.")
	}
	if err != nil {
		s.log.Error("Create record in store", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}

	return &auditumv1alpha1.CreateRecordResponse{
		Record: encodeRecord(record),
	}, nil
}

func (s *RecordServiceServer) BatchCreateRecords(
	ctx context.Context,
	req *auditumv1alpha1.BatchCreateRecordsRequest,
) (*auditumv1alpha1.BatchCreateRecordsResponse, error) {
	records, err := decodeRecords(
		req.GetProjectId(),
		req.GetRecords(),
		s.settings.Records.Restrictions,
	)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "records": %v.`,
			err.Error(),
		)
	}

	now := s.now().UTC()

	for i := range records {
		records[i].ID = s.id()
		records[i].CreateTime = now
	}

	err = s.store.CreateRecords(ctx, records)
	if errors.Is(err, aud.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "Project not found.")
	}
	if err != nil {
		s.log.Error("Create records in store", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}

	return &auditumv1alpha1.BatchCreateRecordsResponse{
		Records: encodeRecords(records),
	}, nil
}

func (s *RecordServiceServer) GetRecord(
	ctx context.Context,
	req *auditumv1alpha1.GetRecordRequest,
) (*auditumv1alpha1.GetRecordResponse, error) {
	projectID, err := decodeID(req.GetProjectId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "project_id": %v.`,
			err.Error(),
		)
	}

	recordID, err := decodeID(req.GetRecordId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "record_id": %v.`,
			err.Error(),
		)
	}

	record, err := s.store.GetRecord(ctx, projectID, recordID)
	if errors.Is(err, aud.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "Project not found.")
	}
	if errors.Is(err, aud.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "Record not found.")
	}
	if err != nil {
		s.log.Error("Get record from store",
			zap.String("project_id", projectID.String()),
			zap.String("record_id", recordID.String()),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "")
	}

	return &auditumv1alpha1.GetRecordResponse{
		Record: encodeRecord(record),
	}, nil
}

func (s *RecordServiceServer) ListRecords(
	ctx context.Context,
	req *auditumv1alpha1.ListRecordsRequest,
) (*auditumv1alpha1.ListRecordsResponse, error) {
	projectID, err := decodeID(req.GetProjectId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "project_id": %v.`,
			err.Error(),
		)
	}

	filter, err := decodeRecordFilter(req.GetFilter())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "filter": %v.`,
			err.Error(),
		)
	}

	const (
		defaultPageSize = 10
		maxPageSize     = 100
	)
	pageSize, err := grpcx.GetPageSize(defaultPageSize, maxPageSize, req)
	if err != nil {
		return nil, err
	}

	var cursor aud.RecordCursor
	if err := aud.DecodePageToken(req.GetPageToken(), &cursor); err != nil {
		s.log.Warn("Decode page token", zap.Error(err))
		return nil, status.Error(
			codes.InvalidArgument,
			`Request is invalid. Invalid "page_token".`,
		)
	}

	records, err := s.store.ListRecords(
		ctx,
		projectID,
		filter,
		pageSize,
		cursor,
	)
	if errors.Is(err, aud.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "Project not found.")
	}
	if err != nil {
		s.log.Error("List records in store", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}

	cursor = aud.NewRecordCursor(records, pageSize)
	nextPageToken, err := aud.EncodePageToken(cursor)
	if err != nil {
		s.log.Error("Encode page token", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	return &auditumv1alpha1.ListRecordsResponse{
		Records:       encodeRecords(records),
		NextPageToken: nextPageToken,
	}, nil
}

func (s *RecordServiceServer) UpdateRecord(ctx context.Context, req *auditumv1alpha1.UpdateRecordRequest) (*auditumv1alpha1.UpdateRecordResponse, error) {
	if !s.settings.Records.UpdateEnabled {
		return nil, status.Error(codes.Unimplemented, "UpdateRecord is disabled.")
	}

	projectID, err := decodeID(req.GetRecord().GetProjectId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "record": invalid "project_id": %v.`,
			err.Error(),
		)
	}

	recordID, err := decodeID(req.GetRecord().GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "record": invalid "id": %v.`,
			err.Error(),
		)
	}

	paths := req.GetUpdateMask().GetPaths()
	if len(paths) == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "update_mask": must not be empty.`,
		)
	}

	var update aud.RecordUpdate

	for _, path := range paths {
		switch path {
		case "labels":
			labels, err := decodeLabels(
				req.GetRecord().GetLabels(),
				s.settings.Records.Restrictions.Labels,
			)
			if err != nil {
				return nil, status.Errorf(
					codes.InvalidArgument,
					`Request is invalid. Invalid "record": invalid "labels": %v.`,
					err.Error(),
				)
			}
			update.Labels = labels
			update.UpdateLabels = true
		case "resource":
			resource, err := decodeResource(
				req.GetRecord().GetResource(),
				s.settings.Records.Restrictions.Resource,
			)
			if err != nil {
				return nil, status.Errorf(
					codes.InvalidArgument,
					`Request is invalid. Invalid "record": invalid "resource": %v.`,
					err.Error(),
				)
			}
			update.Resource = resource
			update.UpdateResource = true
		case "operation":
			operation, err := decodeOperation(
				req.GetRecord().GetOperation(),
				s.settings.Records.Restrictions.Operation,
			)
			if err != nil {
				return nil, status.Errorf(
					codes.InvalidArgument,
					`Request is invalid. Invalid "record": invalid "operation": %v.`,
					err.Error(),
				)
			}
			update.Operation = operation
			update.UpdateOperation = true
		case "actor":
			actor, err := decodeActor(
				req.GetRecord().GetActor(),
				s.settings.Records.Restrictions.Actor,
			)
			if err != nil {
				return nil, status.Errorf(
					codes.InvalidArgument,
					`Request is invalid. Invalid "record": invalid "actor": %v.`,
					err.Error(),
				)
			}
			update.Actor = actor
			update.UpdateActor = true
		default:
			return nil, status.Errorf(
				codes.InvalidArgument,
				`Request is invalid. Invalid "update_mask": path %v is not supported.`,
				path,
			)
		}
	}

	updatedRecord, err := s.store.UpdateRecord(
		ctx,
		projectID,
		recordID,
		update,
	)
	if errors.Is(err, aud.ErrProjectNotFound) {
		return nil, status.Errorf(codes.NotFound, "Project not found.")
	}
	if errors.Is(err, aud.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "Record not found.")
	}
	if err != nil {
		s.log.Error("Update record in store",
			zap.String("project_id", projectID.String()),
			zap.String("record_id", recordID.String()),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "")
	}

	return &auditumv1alpha1.UpdateRecordResponse{
		Record: encodeRecord(updatedRecord),
	}, nil
}

func (s *RecordServiceServer) DeleteRecord(ctx context.Context, req *auditumv1alpha1.DeleteRecordRequest) (*auditumv1alpha1.DeleteRecordResponse, error) {
	if !s.settings.Records.DeleteEnabled {
		return nil, status.Error(codes.Unimplemented, "DeleteRecord is disabled.")
	}

	projectID, err := decodeID(req.GetProjectId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "project_id": %v.`,
			err.Error(),
		)
	}

	recordID, err := decodeID(req.GetRecordId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "record_id": %v.`,
			err.Error(),
		)
	}

	err = s.store.DeleteRecord(ctx, projectID, recordID)
	if errors.Is(err, aud.ErrProjectNotFound) {
		return nil, status.Errorf(codes.NotFound, "Project not found.")
	}
	if err != nil {
		s.log.Error("Delete record from store",
			zap.String("project_id", projectID.String()),
			zap.String("record_id", recordID.String()),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "")
	}

	return &auditumv1alpha1.DeleteRecordResponse{}, nil
}

func (s *RecordServiceServer) RegisterServer(srv *grpc.Server) {
	auditumv1alpha1.RegisterRecordServiceServer(srv, s)
}

func (s *RecordServiceServer) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return auditumv1alpha1.RegisterRecordServiceHandler(ctx, mux, conn)
}
