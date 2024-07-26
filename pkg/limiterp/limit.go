package limiterp

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed limit_slide_window.lua
var luaScript string

type Limiter interface {
	// Limit true为触发，false为不触发
	Limit(ctx context.Context, key string) (bool, error)
}

type limiter struct {
	cmd      redis.Cmdable
	interval time.Duration
	// 阈值
	rate int
}

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &limiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (l limiter) Limit(ctx context.Context, key string) (bool, error) {
	return l.cmd.Eval(ctx, luaScript, []string{key},
		l.interval.Milliseconds(), l.rate, time.Now().UnixMilli()).Bool()
}
