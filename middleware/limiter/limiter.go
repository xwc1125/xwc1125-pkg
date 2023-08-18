package limiter

import (
	"context"
	"time"
)

type Resource struct {
	Name   string        // 限流器名称
	Limit  int           // 每秒事件数
	Burst  int           // 突发令牌数
	Window time.Duration // 滑动窗口时间
}

type Entry interface {
	Allow() bool  // 是否允许
	Finish()      // 结束
	Error() error // 错误信息
}

type Limiter interface {
	// Check 获取一个entry
	Check(ctx context.Context, r Resource) Entry
	// SetLimit 动态地改变 Token 桶大小
	SetLimit(ctx context.Context, r Resource)
	// SetBurst 动态地改变 生成速率
	SetBurst(ctx context.Context, r Resource)
	// SetWindow 动态地改变 滑动窗口
	SetWindow(ctx context.Context, r Resource)
}
