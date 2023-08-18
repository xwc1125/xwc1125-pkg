// Package provider
package provider

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/xwc1125/xwc1125-pkg/utils/iputil/ipv4"
	jaegerprop "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	hostIP, _   = ipv4.GetIntranetIps() // ipv4的地址 可选,是否传递看自己
	hostname, _ = os.Hostname()
)

// NewJaegerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func NewJaegerProvider(serviceName, endpoint string, options ...Option) (*sdkTrace.TracerProvider, error) {
	var endpointOption jaeger.EndpointOption
	if serviceName == "" {
		return nil, errors.New("no service name provided")
	}
	if strings.HasPrefix(endpoint, "http") {
		endpointOption = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint))
	} else {
		endpointOption = jaeger.WithAgentEndpoint(jaeger.WithAgentHost(endpoint))
	}

	// Create the Jaeger exporter
	exporter, err := jaeger.New(endpointOption)
	if err != nil {
		return nil, err
	}

	opts := applyOptions(options...)
	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithIDGenerator(NewIDGenerator()),
		// set sample
		sdkTrace.WithSampler(sdkTrace.TraceIDRatioBased(opts.SamplingRatio)),
		// Always be sure to batch in production.
		sdkTrace.WithBatcher(exporter),
		// Record information about this application in an Resource.
		sdkTrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.NetHostIPKey.String(hostIP),
			semconv.HostNameKey.String(hostname),
		)),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(jaegerprop.Jaeger{})

	return tp, nil
}
