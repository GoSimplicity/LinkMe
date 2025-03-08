package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GoSimplicity/LinkMe/internal/domain/events"
	"github.com/GoSimplicity/LinkMe/ioc"
	"github.com/spf13/viper"

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
	cmd := ioc.InitWebServer()

	server := cmd.Server
	server.GET("/headers", printHeaders)

	// 创建一个用于接收系统信号的通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动 Prometheus 监控
	go func() {

		if err := startMetricsServer(); err != nil {
			zap.L().Fatal("启动监控服务器失败", zap.Error(err))
		}
	}()

	// 启动定时任务和worker
	go func() {
		if err := cmd.Scheduler.RegisterTimedTasks(); err != nil {
			zap.L().Fatal("注册定时任务失败", zap.Error(err))
		}

		if err := cmd.Scheduler.Run(); err != nil {
			zap.L().Fatal("启动定时任务失败", zap.Error(err))
		}
	}()

	// 启动消费者
	for _, s := range cmd.Consumer {
		go func(consumer events.Consumer) {
			if err := consumer.Start(context.Background()); err != nil {
				zap.L().Fatal("启动消费者失败", zap.Error(err))
			}
		}(s)
	}

	// 注册任务处理器并启动异步任务服务器
	go func() {
		mux := cmd.Routes.RegisterHandlers()
		if err := cmd.Asynq.Run(mux); err != nil {
			zap.L().Fatal("启动异步任务服务器失败", zap.Error(err))
		}
	}()

	// 在新的goroutine中启动服务器
	go func() {
		if err := server.Run(viper.GetString("server.addr")); err != nil {
			zap.L().Fatal("启动Web服务器失败", zap.Error(err))
		}
	}()

	// 等待中断信号
	<-quit
	zap.L().Info("正在关闭服务器...")

	// 关闭异步任务服务器
	cmd.Asynq.Shutdown()

	cmd.Scheduler.Stop()

	zap.L().Info("服务器已成功关闭")
	os.Exit(0)
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
	// 启动 HTTP 服务器并捕获可能的错误
	if err := http.ListenAndServe(":9091", nil); err != nil {
		log.Fatalf("Prometheus 启动失败: %v", err)
		return err
	}
	return nil
}
