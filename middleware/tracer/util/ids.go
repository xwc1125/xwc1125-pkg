// Package util
package util

import (
	"math"
	"strconv"
	"time"

	"github.com/xwc1125/xwc1125-pkg/utils/randutil"
	"go.opentelemetry.io/otel/trace"
)

var (
	randomInitSequence = int32(randutil.Intn(math.MaxInt32))
	sequence           = randomInitSequence
)

// NewIDs creates and returns a new trace and span ID.
func NewIDs() (traceID trace.TraceID, spanID trace.SpanID) {
	return NewTraceID(), NewSpanID()
}

// NewTraceID creates and returns a trace ID.
func NewTraceID() (traceID trace.TraceID) {
	var (
		timestampNanoBytes = time.Now().UnixNano()
		sequenceBytes      = sequence + 1
		randomBytes        = randutil.B(4)
	)
	timeStr := ([]byte)(strconv.FormatInt(timestampNanoBytes, 10))
	copy(traceID[:], timeStr[:])
	// copy(traceID[:], timestampNanoBytes)
	seqBytes := ([]byte)(strconv.FormatInt(int64(sequenceBytes), 10))
	copy(traceID[8:], seqBytes)
	copy(traceID[12:], randomBytes)
	return
}

// NewSpanID creates and returns a span ID.
func NewSpanID() (spanID trace.SpanID) {
	timeStr := ([]byte)(strconv.FormatInt(time.Now().UnixNano()/1e3, 10))
	copy(spanID[:], timeStr)
	copy(spanID[4:], randutil.B(4))
	return
}
