// Package plugins
package plugins

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xwc1125/xwc1125-pkg/middleware/tracer"
	"github.com/xwc1125/xwc1125-pkg/utils/iputil/ipv4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerNameForRedis       = "redis.client/otelgo"
	kvClientHostForRedis     = "redis.client.host"
	kvClientHostNameForRedis = "redis.client.hostname"
	tracingTraceForRedis     = "redis.trace"
)

type redisTracingHook struct {
	tracer *tracer.Tracer
	attrs  []attribute.KeyValue
}

func NewTracerForRedis() redis.Hook {
	hostName, _ := os.Hostname()
	host, _ := ipv4.GetIntranetIps()
	return &redisTracingHook{
		tracer: tracer.NewTracer(tracerNameForRedis),
		attrs: []attribute.KeyValue{
			semconv.DBSystemRedis,
			attribute.Key(kvClientHostForRedis).String(host),
			attribute.Key(kvClientHostNameForRedis).String(hostName),
		},
	}
}

func (t *redisTracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	// span := trace.SpanFromContext(ctx)
	// if !span.IsRecording() {
	// 	return ctx, nil
	// }
	ctx, _ = tracer.NewTracer(tracerNameForRedis).Start(ctx, tracingTraceForRedis, trace.WithSpanKind(trace.SpanKindClient))
	return ctx, nil
}

func (t *redisTracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	var (
		span = trace.SpanFromContext(ctx)
		tn   = time.Now()
	)
	if !span.IsRecording() {
		return nil
	}
	defer span.End()
	if err := cmd.Err(); err != nil && err != redis.Nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.SetAttributes(t.attrs...)
	attrs := make([]attribute.KeyValue, 0)
	argsB, _ := json.Marshal(cmd.Args())
	attrs = append(attrs,
		attribute.String("name", cmd.FullName()),
		attribute.String("args", string(argsB)),
		attribute.String("cmd", rediscmd.CmdString(cmd)),
		attribute.String("redis.driver", "go-redis"),
	)
	if cmd.Err() != nil {
		attrs = append(attrs, attribute.String("err", cmd.Err().Error()))
	}
	span.AddEvent(tracingTraceForRedis, trace.WithAttributes(attrs...))
	span.End(trace.WithStackTrace(true), trace.WithTimestamp(tn))
	return nil
}

func (t *redisTracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}

	summary, _ := rediscmd.CmdsString(cmds)
	ctx, _ = tracer.NewTracer(tracerNameForRedis).Start(ctx, "pipeline."+tracingTraceForRedis+"."+summary, trace.WithSpanKind(trace.SpanKindClient))
	return ctx, nil
}

func (t *redisTracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	var (
		span = trace.SpanFromContext(ctx)
		tn   = time.Now()
	)
	span.SetAttributes(t.attrs...)
	summary, cmdsString := rediscmd.CmdsString(cmds)
	attrs := make([]attribute.KeyValue, 0)
	attrs = append(attrs,
		attribute.String("summary", summary),
		attribute.String("cmds", cmdsString),
		attribute.Int("numCmds", len(cmds)))

	if err := cmds[0].Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		attrs = append(attrs, attribute.String("error", err.Error()))
	}

	span.AddEvent(tracingTraceForRedis, trace.WithAttributes(attrs...))
	span.End(trace.WithStackTrace(true), trace.WithTimestamp(tn))
	return nil
}
