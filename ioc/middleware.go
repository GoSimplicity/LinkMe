package ioc

import (
	"github.com/GoSimplicity/LinkMe/middleware"
	"github.com/GoSimplicity/LinkMe/pkg/prometheusp"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"time"
)

// InitMiddlewares 初始化中间件
func InitMiddlewares(ih ijwt.Handler, l *zap.Logger, enforcer *casbin.Enforcer) []gin.HandlerFunc {
	prom := &prometheusp.MetricsPlugin{
		Namespace:  "linkme",
		Subsystem:  "api",
		InstanceID: "instance_1",
	}

	// 注册指标
	prom.RegisterMetrics()

	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization", "X-Refresh-Token"},
			ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "")
			},
			MaxAge: 12 * time.Hour,
		}),
		// 统计响应时间
		prom.TrackActiveRequestsMiddleware(),
		// 统计活跃请求数
		prom.TrackResponseTimeMiddleware(),
		middleware.NewJWTMiddleware(ih).CheckLogin(),
		middleware.NewLogMiddleware(l).Log(),
	}
}
