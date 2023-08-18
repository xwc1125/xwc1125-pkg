// Package provider
package provider

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// NewTracerProvider
// serviceName 服务名称，一般是apollo上面的服务名称，注意冲突
// http:
//
// client := otlptracehttp.NewClient(
//
//	otlptracehttp.WithEndpoint(endpoint),
//	otlptracehttp.WithInsecure(),
//	otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
//		Enabled:         true,
//		InitialInterval: 1 * time.Second,
//		MaxInterval:     1 * time.Second,
//		MaxElapsedTime:  0,
//	}))
//
// grpc:
// client := otlptracegrpc.NewClient(
//
//	otlptracegrpc.WithEndpoint(endpoint),
//	otlptracegrpc.WithInsecure(),
//	otlptracegrpc.WithReconnectionPeriod(50*time.Millisecond),
//	otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
//		Enabled:         true,
//		InitialInterval: 1 * time.Second,
//		MaxInterval:     1 * time.Second,
//		MaxElapsedTime:  0,
//	}))
func NewTracerProvider(ctx context.Context, client otlptrace.Client, serviceName string) (*sdkTrace.TracerProvider, error) {
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithIDGenerator(NewIDGenerator()),
		sdkTrace.WithBatcher(exp),
		sdkTrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.NetHostIPKey.String(hostIP),
			semconv.HostNameKey.String(hostname),
		)),
	)
	// registers `tp` as the global trace provider.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return tp, nil
}
