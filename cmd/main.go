package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain/events"
	"github.com/GoSimplicity/LinkMe/ioc"
	"github.com/spf13/viper"

	"github.com/GoSimplicity/LinkMe/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	Init()
}

func Init() {
	config.InitViper()
	cmd := ioc.InitWebServer()
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	metricsServer := startMetricsServer(viper.GetString("metrics.addr"))

	go func() {
		if err := cmd.Scheduler.RegisterTimedTasks(); err != nil {
			zap.L().Fatal("注册定时任务失败", zap.Error(err))
		}

		if err := cmd.Scheduler.Run(); err != nil {
			zap.L().Fatal("启动定时任务失败", zap.Error(err))
		}
	}()

	for _, s := range cmd.Consumer {
		go func(consumer events.Consumer) {
			if err := consumer.Start(rootCtx); err != nil && !errors.Is(err, context.Canceled) {
				zap.L().Fatal("启动消费者失败", zap.Error(err))
			}
		}(s)
	}

	go func() {
		mux := cmd.Routes.RegisterHandlers()
		if err := cmd.Asynq.Run(mux); err != nil {
			zap.L().Fatal("启动异步任务服务器失败", zap.Error(err))
		}
	}()

	appServer := &http.Server{
		Addr:    viper.GetString("server.addr"),
		Handler: cmd.Server,
	}

	go func() {
		if err := appServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("启动Web服务器失败", zap.Error(err))
		}
	}()

	<-rootCtx.Done()
	zap.L().Info("正在关闭服务器...")

	cmd.Asynq.Shutdown()
	cmd.Scheduler.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		zap.L().Error("关闭监控服务器失败", zap.Error(err))
	}
	if err := appServer.Shutdown(shutdownCtx); err != nil {
		zap.L().Error("关闭Web服务器失败", zap.Error(err))
	}

	zap.L().Info("服务器已成功关闭")
}

// startMetricsServer 启动 Prometheus 监控服务器
func startMetricsServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("Prometheus 启动失败", zap.Error(err))
		}
	}()
	return server
}
