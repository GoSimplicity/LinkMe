package main

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain/events"
	"net/http"

	"github.com/GoSimplicity/LinkMe/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	Init()
}

func Init() {
	// 初始化配置
	config.InitViper()

	// 初始化 Web 服务器和其他组件
	cmd := InitWebServer()

	// 设置请求头打印路由
	server := cmd.server
	server.GET("/headers", printHeaders)

	// 启动 Prometheus 监控
	go func() {
		if err := startMetricsServer(); err != nil {
			zap.L().Error("Failed to start metrics server", zap.Error(err))
		}
	}()

	// 启动消费者
	for _, s := range cmd.consumer {
		go func(consumer events.Consumer) { // 将每个消费者启动放入goroutine中并发执行
			if err := consumer.Start(context.Background()); err != nil {
				zap.L().Error("Failed to start consumer", zap.Error(err))
			}
		}(s)
	}

	// 启动定时任务
	cmd.Cron.Start()

	// 启动 Web 服务器
	if err := server.Run(":9999"); err != nil {
		zap.L().Fatal("Failed to start web server", zap.Error(err))
	}
}

// printHeaders 打印请求头信息
func printHeaders(c *gin.Context) {
	headers := c.Request.Header
	for key, values := range headers {
		for _, value := range values {
			c.String(http.StatusOK, "%s: %s\n", key, value)
		}
	}
}

// startMetricsServer 启动 Prometheus 监控服务器
func startMetricsServer() error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":9090", nil)
}
