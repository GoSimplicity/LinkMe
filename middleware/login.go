package middleware

import (
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWTMiddleware struct {
	ijwt.Handler
}

func NewJWTMiddleware(hdl ijwt.Handler) *JWTMiddleware {
	return &JWTMiddleware{
		Handler: hdl,
	}
}

// CheckLogin 校验JWT
func (m *JWTMiddleware) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		// 如果请求的路径是下述路径，则不进行token验证
		if path == "/api/user/signup" ||
			path == "/api/user/login" ||
			path == "/api/user/refresh_token" ||
			path == "/api/user/change_password" ||
			path == "/api/user/send_sms" ||
			path == "/api/user/send_email" {
			return
		}
		// 从请求中提取token
		tokenStr := m.ExtractToken(ctx)
		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("jwt.auth_key")), nil
		})
		if err != nil {
			// token 错误
			ctx.AbortWithStatus(RequestsERROR)
			return
		}
		if token == nil || !token.Valid {
			// token 非法或过期
			ctx.AbortWithStatus(RequestsERROR)
			return
		}
		// 检查是否携带ua头
		if uc.UserAgent == "" {
			ctx.AbortWithStatus(RequestsERROR)
			return
		}
		// 检查会话是否有效
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			ctx.AbortWithStatus(RequestsERROR)
			return
		}
		ctx.Set("user", uc)
	}
}
