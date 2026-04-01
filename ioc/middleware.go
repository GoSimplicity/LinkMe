package ioc

import (
	"time"

	"github.com/GoSimplicity/LinkMe/middleware"
	"github.com/GoSimplicity/LinkMe/pkg/ginp/prometheus"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitMiddlewares 初始化中间件
func InitMiddlewares(ih ijwt.Handler, l *zap.Logger) []gin.HandlerFunc {
	prom := &prometheus.MetricsPlugin{
		Namespace:  "linkme",
		Subsystem:  "api",
		InstanceID: "instance_1",
	}

	// 注册指标
	prom.RegisterMetrics()

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Refresh-Token", "X-Request-ID"},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		MaxAge:           12 * time.Hour,
	}

	if viper.GetBool("cors.allow_all") {
		corsConfig.AllowAllOrigins = true
	} else {
		allowedOrigins := viper.GetStringSlice("cors.allow_origins")
		corsConfig.AllowOriginFunc = func(origin string) bool {
			if origin == "" {
				return true
			}
			for _, item := range allowedOrigins {
				if item == origin {
					return true
				}
			}
			return false
		}
	}

	return []gin.HandlerFunc{
		cors.New(corsConfig),
		// 统计响应时间
		prom.TrackActiveRequestsMiddleware(),
		// 统计活跃请求数
		prom.TrackResponseTimeMiddleware(),
		middleware.NewJWTMiddleware(ih).CheckLogin(),
		middleware.NewLogMiddleware(l).Log(),
	}
}
