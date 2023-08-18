// Package shutdown
package shutdown

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/chain5j/logger"
)

func TestNewHook(t *testing.T) {
	// 初始化 HTTP 服务
	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server startup err", "err", err)
		}
	}()

	// 优雅关闭
	NewHook().Close(
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				logger.Error("server shutdown err", "err", err)
			}
		},
	)
}
