package ioc

import (
	. "github.com/GoSimplicity/LinkMe/pkg/limiterp"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitLimiter(redis redis.Cmdable) Limiter {
	return NewRedisSlidingWindowLimiter(redis, time.Second, 100)
}
