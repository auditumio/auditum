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
