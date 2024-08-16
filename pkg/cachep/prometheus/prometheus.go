package prometheus

import (
	"context"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

// RedisMetricsHook 实现了 redis.Hook 接口，用于监控 Redis 操作
type RedisMetricsHook struct {
	operationMetrics *prometheus.SummaryVec
}

// NewRedisMetricsHook 初始化 RedisMetricsHook 实例，并注册 Prometheus 指标
func NewRedisMetricsHook(opts prometheus.SummaryOpts) *RedisMetricsHook {
	operationMetrics := prometheus.NewSummaryVec(opts, []string{"operation"})
	prometheus.MustRegister(operationMetrics)
	return &RedisMetricsHook{
		operationMetrics: operationMetrics,
	}
}

// ProcessHook 实现了 redis.ProcessHook，用于监控单个命令的执行时间
func (h *RedisMetricsHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		// 在命令执行之前记录开始时间
		startTime := time.Now()

		// 调用下一个钩子或最终的命令执行
		err := next(ctx, cmd)

		// 计算命令执行的持续时间，并记录到 Prometheus
		duration := time.Since(startTime).Seconds()
		h.operationMetrics.WithLabelValues(cmd.Name()).Observe(duration)

		return err
	}
}

// DialHook 实现了 redis.DialHook，用于监控连接的执行时间
func (h *RedisMetricsHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 在连接操作之前记录开始时间
		startTime := time.Now()

		// 调用下一个钩子或实际的连接操作
		conn, err := next(ctx, network, addr)

		// 计算连接操作的持续时间，并记录到 Prometheus
		duration := time.Since(startTime).Seconds()
		h.operationMetrics.WithLabelValues("dial").Observe(duration)

		return conn, err
	}
}

// ProcessPipelineHook 实现了 redis.ProcessPipelineHook，用于监控 Pipeline 操作的执行时间
func (h *RedisMetricsHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		// 在 Pipeline 执行之前记录开始时间
		startTime := time.Now()

		// 调用下一个钩子或实际的 Pipeline 操作
		err := next(ctx, cmds)

		// 计算 Pipeline 操作的持续时间，并记录到 Prometheus
		duration := time.Since(startTime).Seconds()
		h.operationMetrics.WithLabelValues("pipeline").Observe(duration)

		return err
	}
}
