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

	auditumv1alpha1 "github.com/infragmo/auditum/api/gen/go/infragmo/auditum/v1alpha1"
	"github.com/infragmo/auditum/internal/aud"
	"github.com/infragmo/auditum/pkg/fragma/grpcx"
)

type ProjectServiceServer struct {
	auditumv1alpha1.UnimplementedProjectServiceServer

	store Store
	log   *zap.Logger

	id  func() aud.ID
	now func() time.Time
}

func NewProjectServiceServer(
	store Store,
	log *zap.Logger,
) *ProjectServiceServer {
	return &ProjectServiceServer{
		store: store,
		log:   log.Named("project_service_server"),
		id:    aud.MustNewID,
		now:   time.Now,
	}
}

func (s *ProjectServiceServer) CreateProject(
	ctx context.Context,
	req *auditumv1alpha1.CreateProjectRequest,
) (*auditumv1alpha1.CreateProjectResponse, error) {
	project, err := decodeProject(req.GetProject())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "project": %v.`,
			err.Error(),
		)
	}

	project.ID = s.id()
	project.CreateTime = s.now().UTC()

	err = s.store.CreateProject(ctx, project)
	if err != nil {
		s.log.Error("Create project in store", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}

	return &auditumv1alpha1.CreateProjectResponse{
		Project: encodeProject(project),
	}, nil
}

func (s *ProjectServiceServer) GetProject(
	ctx context.Context,
	req *auditumv1alpha1.GetProjectRequest,
) (*auditumv1alpha1.GetProjectResponse, error) {
	id, err := decodeID(req.GetProjectId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid "project_id": %v.`,
			err.Error(),
		)
	}

	project, err := s.store.GetProject(ctx, id)
	if errors.Is(err, aud.ErrProjectNotFound) {
		return nil, status.Errorf(codes.NotFound, "")
	}
	if err != nil {
		s.log.Error("Get project from store", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}

	return &auditumv1alpha1.GetProjectResponse{
		Project: encodeProject(project),
	}, nil
}

func (s *ProjectServiceServer) ListProjects(
	ctx context.Context,
	req *auditumv1alpha1.ListProjectsRequest,
) (*auditumv1alpha1.ListProjectsResponse, error) {
	const (
		defaultPageSize = 10
		maxPageSize     = 100
	)
	pageSize, err := grpcx.GetPageSize(defaultPageSize, maxPageSize, req)
	if err != nil {
		return nil, err
	}

	var cursor aud.ProjectCursor
	if err := aud.DecodePageToken(req.GetPageToken(), &cursor); err != nil {
		s.log.Warn("Decode page token", zap.Error(err))
		return nil, status.Error(
			codes.InvalidArgument,
			`Request is invalid. Invalid "page_token".`,
		)
	}

	projects, err := s.store.ListProjects(
		ctx,
		pageSize,
		cursor,
	)
	if err != nil {
		s.log.Error("List projects in store", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}

	cursor = aud.NewProjectCursor(projects, pageSize)
	nextPageToken, err := aud.EncodePageToken(cursor)
	if err != nil {
		s.log.Error("Encode page token", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	return &auditumv1alpha1.ListProjectsResponse{
		Projects:      encodeProjects(projects),
		NextPageToken: nextPageToken,
	}, nil
}

func (s *ProjectServiceServer) UpdateProject(
	ctx context.Context,
	req *auditumv1alpha1.UpdateProjectRequest,
) (*auditumv1alpha1.UpdateProjectResponse, error) {
	projectID, err := decodeID(req.GetProject().GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			`Request is invalid. Invalid project id: %v.`,
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

	var update aud.ProjectUpdate

	for _, path := range paths {
		switch path {
		case "display_name":
			update.DisplayName = req.GetProject().GetDisplayName()
			update.UpdateDisplayName = true
		case "update_record_enabled":
			update.UpdateRecordEnabled = decodeBoolValue(req.GetProject().GetUpdateRecordEnabled())
			update.UpdateUpdateRecordEnabled = true
		case "delete_record_enabled":
			update.DeleteRecordEnabled = decodeBoolValue(req.GetProject().GetDeleteRecordEnabled())
			update.UpdateDeleteRecordEnabled = true
		default:
			return nil, status.Errorf(
				codes.InvalidArgument,
				`Request is invalid. Invalid "update_mask": path %v is not supported.`,
				path,
			)
		}
	}

	updatedProject, err := s.store.UpdateProject(
		ctx,
		projectID,
		update,
	)
	if errors.Is(err, aud.ErrProjectNotFound) {
		return nil, status.Errorf(codes.NotFound, "")
	}
	if err != nil {
		s.log.Error("Update project in store",
			zap.String("project_id", projectID.String()),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "")
	}

	return &auditumv1alpha1.UpdateProjectResponse{
		Project: encodeProject(updatedProject),
	}, nil
}

func (s *ProjectServiceServer) RegisterServer(srv *grpc.Server) {
	auditumv1alpha1.RegisterProjectServiceServer(srv, s)
}

func (s *ProjectServiceServer) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return auditumv1alpha1.RegisterProjectServiceHandler(ctx, mux, conn)
}
