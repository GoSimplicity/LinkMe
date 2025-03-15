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

	logger := di.InitZap()
	defer logger.Sync()

	// 替换全局logger
	zap.ReplaceGlobals(logger)

	mode := viper.GetString("server.mode")
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	r := gin.New()
	// 使用自定义的日志中间件和恢复中间件
	r.Use(ginLogger(), gin.Recovery())

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

// ginLogger 自定义 Gin 的日志中间件
func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 构建基础日志字段
		fields := []zap.Field{
			zap.Int("状态码", statusCode),
			zap.String("请求方法", method),
			zap.String("请求路径", path),
			zap.String("客户端IP", clientIP),
			zap.Duration("响应时间", latency),
		}

		// 只有在有用户代理信息时才添加该字段
		if userAgent != "" {
			fields = append(fields, zap.String("用户代理", userAgent))
		}

		// 只有在有错误信息时才添加该字段
		if errorMessage != "" {
			fields = append(fields, zap.String("错误信息", errorMessage))
		}

		// 添加响应时间的人性化显示
		if latency > time.Second {
			fields = append(fields, zap.String("耗时", latency.Round(time.Millisecond).String()))
		} else {
			fields = append(fields, zap.String("耗时", latency.Round(time.Microsecond).String()))
		}

		// 根据状态码选择日志级别和添加状态描述
		statusDesc := http.StatusText(statusCode)
		if statusDesc != "" {
			fields = append(fields, zap.String("状态描述", statusDesc))
		}

		if statusCode >= 500 {
			zap.L().Error("HTTP请求失败", fields...)
		} else if statusCode >= 400 {
			zap.L().Warn("HTTP请求异常", fields...)
		} else {
			zap.L().Info("HTTP请求成功", fields...)
		}
	}
}
