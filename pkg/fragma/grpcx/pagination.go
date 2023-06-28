package grpcx

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PageSizeRequest describes a request with page size parameter.
type PageSizeRequest interface {
	GetPageSize() int32
}

func GetPageSize(
	defaultPageSize int32,
	maxPageSize int32,
	req PageSizeRequest,
) (pageSize int32, err error) {
	pageSize = defaultPageSize
	if req.GetPageSize() != 0 {
		pageSize = req.GetPageSize()
		if pageSize < 0 {
			return pageSize, status.Error(
				codes.InvalidArgument,
				`Invalid "page_size": value is less than zero.`,
			)
		}
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return pageSize, nil
}
