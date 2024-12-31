package ioc

import (
	prometheus2 "github.com/GoSimplicity/LinkMe/pkg/cachep/prometheus" // 替换为实际路径
	prometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	// 初始化 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
	})

	// 创建并注册自定义的 RedisMetricsHook 插件
	prometheusHook := prometheus2.NewRedisMetricsHook(prometheus.SummaryOpts{
		Namespace: "linkme",
		Subsystem: "redis",
		Name:      "operation_duration_seconds",
		Help:      "Duration of Redis operations in seconds",
		Objectives: map[float64]float64{
			0.5:  0.01,
			0.75: 0.01,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	// 添加 Hook 插件到 Redis 客户端
	client.AddHook(prometheusHook)

	return client
}
