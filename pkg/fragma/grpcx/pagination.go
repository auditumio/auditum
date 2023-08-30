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
