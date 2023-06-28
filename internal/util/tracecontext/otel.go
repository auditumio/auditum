package tracecontext

// Code in this file is taken and adapted from otel-go.
// See: https://github.com/open-telemetry/opentelemetry-go/blob/7c552505a9cb0462d6859cec18bfe9cbac6e38c3/propagation/trace_context.go

import (
	"encoding/hex"
	"regexp"

	"go.opentelemetry.io/otel/trace"
)

const (
	maxVersion = 254
)

var traceCtxRegExp = regexp.MustCompile("^(?P<version>[0-9a-f]{2})-(?P<traceID>[a-f0-9]{32})-(?P<spanID>[a-f0-9]{16})-(?P<traceFlags>[a-f0-9]{2})(?:-.*)?$")

func extractTraceparent(src string) trace.SpanContext {
	matches := traceCtxRegExp.FindStringSubmatch(src)

	if len(matches) == 0 {
		return trace.SpanContext{}
	}

	if len(matches) < 5 { // four subgroups plus the overall match
		return trace.SpanContext{}
	}

	if len(matches[1]) != 2 {
		return trace.SpanContext{}
	}
	ver, err := hex.DecodeString(matches[1])
	if err != nil {
		return trace.SpanContext{}
	}
	version := int(ver[0])
	if version > maxVersion {
		return trace.SpanContext{}
	}

	if version == 0 && len(matches) != 5 { // four subgroups plus the overall match
		return trace.SpanContext{}
	}

	if len(matches[2]) != 32 {
		return trace.SpanContext{}
	}

	var scc trace.SpanContextConfig

	scc.TraceID, err = trace.TraceIDFromHex(matches[2][:32])
	if err != nil {
		return trace.SpanContext{}
	}

	if len(matches[3]) != 16 {
		return trace.SpanContext{}
	}
	scc.SpanID, err = trace.SpanIDFromHex(matches[3])
	if err != nil {
		return trace.SpanContext{}
	}

	if len(matches[4]) != 2 {
		return trace.SpanContext{}
	}
	opts, err := hex.DecodeString(matches[4])
	if err != nil || len(opts) < 1 || (version == 0 && opts[0] > 2) {
		return trace.SpanContext{}
	}
	// Clear all flags other than the trace-context supported sampling bit.
	scc.TraceFlags = trace.TraceFlags(opts[0]) & trace.FlagsSampled

	sc := trace.NewSpanContext(scc)
	if !sc.IsValid() {
		return trace.SpanContext{}
	}

	return sc
}

func extractTracestate(src string) (trace.TraceState, error) {
	return trace.ParseTraceState(src)
}
