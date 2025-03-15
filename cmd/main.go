package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoSimplicity/LinkMe/di"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	if err := di.InitViper(); err != nil {
		panic("初始化配置失败: " + err.Error())
	}
	app, err := di.ProvideApp()
	if err != nil {
		panic("初始化应用失败: " + err.Error())
	}

	defer app.Logger.Sync()

	// 替换全局logger
	zap.ReplaceGlobals(app.Logger)

	mode := viper.GetString("server.mode")
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	r := app.Server

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    viper.GetString("server.addr"),
		Handler: r,
	}

	// 创建退出信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		zap.L().Info("启动Web服务器", zap.String("地址", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("启动Web服务器失败", zap.Error(err))
		}
	}()

	// 等待退出信号
	<-quit
	zap.L().Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("服务器关闭异常", zap.Error(err))
	}

	zap.L().Info("服务器已成功关闭")
}
