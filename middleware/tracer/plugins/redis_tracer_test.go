// Package plugins
package plugins

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"testing"
// 	"time"
//
// 	"github.com/xwc1125/xwc1125-pkg/database/redis"
// 	"github.com/xwc1125/xwc1125-pkg/middleware/tracer/provider"
// 	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
// )
//
// var (
// 	jaegerUrl = "http://127.0.0.1:14268/api/traces"
// 	ctx       context.Context
// 	cancel    context.CancelFunc
// 	tp        *sdkTrace.TracerProvider
// )
//
// func init() {
// 	tp1, err := provider.NewJaegerProvider("test-redis", jaegerUrl)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// 创建context
// 	ctx, cancel = context.WithCancel(context.Background())
// 	tp = tp1
// }
//
// func stop() {
// 	// 优雅退出
// 	defer func(ctx context.Context) {
// 		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
// 		defer cancel()
// 		if err := tp.Shutdown(ctx); err != nil {
// 			panic(err)
// 		}
// 	}(ctx)
// }
//
// func TestNewTracerForRedis(t *testing.T) {
// 	db, err := redis.New(redis.RedisConfig{
// 		Mode:     0,
// 		Addr:     []string{"127.0.0.1:6379"},
// 		Username: "",
// 		Password: "Xwc123~",
// 		DB:       0,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	db.AddHook(NewTracerForRedis())
// 	statusCmd := db.Set(context.Background(), "test-1", "1", 0)
// 	if statusCmd.Err() != nil {
// 		t.Fatal(statusCmd.Err())
// 	}
// 	stringCmd := db.Get(context.Background(), "test-1")
// 	if stringCmd.Err() != nil {
// 		t.Fatal(stringCmd.Err())
// 	}
// 	fmt.Println(stringCmd.String())
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	wg.Wait()
// 	stop()
// }
