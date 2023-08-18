package provider

import (
	"context"

	"github.com/xwc1125/xwc1125-pkg/middleware/tracer/util"
	"go.opentelemetry.io/otel/trace"
)

type IDGenerator struct{}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// NewIDs creates and returns a new trace and span ID.
func (id *IDGenerator) NewIDs(ctx context.Context) (traceID trace.TraceID, spanID trace.SpanID) {
	return util.NewIDs()
}

// NewSpanID returns an ID for a new span in the trace with traceID.
func (id *IDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) (spanID trace.SpanID) {
	return util.NewSpanID()
}
