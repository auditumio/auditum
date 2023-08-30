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

package tracecontext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceparentValid(t *testing.T) {
	tests := []struct {
		traceparent string
		want        bool
	}{
		{
			traceparent: "",
			want:        false,
		},
		{
			traceparent: "123",
			want:        false,
		},
		{
			traceparent: "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01",
			want:        true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run("", func(t *testing.T) {
			got := TraceparentValid(test.traceparent)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestTracestateValid(t *testing.T) {
	tests := []struct {
		tracestate string
		want       bool
	}{
		{
			tracestate: "",
			want:       true,
		},
		{
			tracestate: "123",
			want:       false,
		},
		{
			tracestate: "congo=t61rcWkgMzE",
			want:       true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run("", func(t *testing.T) {
			got := TracestateValid(test.tracestate)
			assert.Equal(t, test.want, got)
		})
	}
}
