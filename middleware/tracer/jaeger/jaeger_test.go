// Package jaeger
package jaeger

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// 启动jaeger：
// docker run \
// --rm \
// --name jaeger \
// -p6831:6831/udp \
// -p5775:5775/udp \
// -p16686:16686 \
// -p14250:14250 \
// -p14268:14268 \
// jaegertracing/all-in-one:latest
func TestJaegerRaw(t *testing.T) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1, // 全部采样
		},
		Reporter: &jaegercfg.ReporterConfig{
			// 当span发送到服务器时要不要打日志
			LogSpans: true,
			// IP:PORT
			// LocalAgentHostPort: "127.0.0.1:6831",
			CollectorEndpoint: "http://127.0.0.1:14268/api/traces",
		},
		ServiceName: "jaeger-test1",
	}
	// 生成链路
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		return
	}
	defer closer.Close()
	// 创建父span时的名称
	parentSpan := tracer.StartSpan("main")
	// opentracing.ChildOf：父span   传输使用的Context 很重要！
	span := tracer.StartSpan("funcA", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 500) // 业务逻辑
	span.Finish()
	// 嵌套
	span2 := tracer.StartSpan("funcB", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 1000) // 业务逻辑
	span2.Finish()
	parentSpan.Finish()
}

func TestNewJaeger(t *testing.T) {
	var config = new(Config)
	config.SetDefaults()
	tracer, closer, err := config.New("jaeger-test2")
	if err != nil {
		return
	}
	defer closer.Close()
	// 创建父span时的名称
	parentSpan := tracer.StartSpan("root")
	defer parentSpan.Finish()
	parentSpan.SetTag("type", "demo")
	parentSpan.LogFields(
		log.String("demo.log", "this is tracing demo"),
	)

	// 将 span 传递给 demoFun
	ctx := opentracing.ContextWithSpan(context.Background(), parentSpan)
	demoFun(ctx)

	// opentracing.ChildOf：父span 传输使用的Context 很重要！
	span := tracer.StartSpan("funcA2", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 500) // 业务逻辑
	span.Finish()
	// 嵌套
	span2 := tracer.StartSpan("funcB2", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 1000) // 业务逻辑
	span2.Finish()
}

func demoFun(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "demoFun")
	defer span.Finish()

	// 假设出错
	err := errors.New("do something erro")
	span.SetTag("error", true)
	span.LogFields(
		log.String("event", "error"),
		log.String("message", err.Error()),
	)

	// 将  ctx 传递给 demoFoo
	demoFoo(ctx)

	//  模拟耗时
	time.Sleep(time.Second * 1)
}

func demoFoo(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "demoFoo")
	defer span.Finish()

	//  模拟耗时
	time.Sleep(time.Second * 2)
}
