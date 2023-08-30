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

package grpcx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/auditumio/auditum/pkg/fragma/grpcx"
)

func TestGetPageSize(t *testing.T) {
	tests := []struct {
		name            string
		defaultPageSize int32
		maxPageSize     int32
		req             grpcx.PageSizeRequest
		pageSize        int32
		err             error
	}{
		{
			name:            "Should fallback to defaults",
			defaultPageSize: 10,
			maxPageSize:     100,
			req:             testPageSizeRequest{},
			pageSize:        10,
			err:             nil,
		},
		{
			name:            "Should use page size",
			defaultPageSize: 10,
			maxPageSize:     100,
			req: testPageSizeRequest{
				pageSize: 20,
			},
			pageSize: 20,
			err:      nil,
		},
		{
			name:            "Should use max page size",
			defaultPageSize: 10,
			maxPageSize:     100,
			req: testPageSizeRequest{
				pageSize: 200,
			},
			pageSize: 100,
			err:      nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			pageSize, err := grpcx.GetPageSize(test.defaultPageSize, test.maxPageSize, test.req)
			assert.Equal(t, test.pageSize, pageSize)
			assert.Equal(t, test.err, err)
		})
	}
}

type testPageSizeRequest struct {
	pageSize int32
}

func (r testPageSizeRequest) GetPageSize() int32 {
	return r.pageSize
}
