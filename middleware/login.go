package middleware

import (
	ijwt "LinkMe/internal/tools/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
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
		if path == "/users/signup" ||
			path == "/users/login" {
			return
		}
		// 从请求中提取token
		tokenStr := m.ExtractToken(ctx)
		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return ijwt.Key1, nil
		})
		if err != nil {
			// token 错误
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			// token 非法或过期
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 检查会话是否有效
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("user", uc)
	}
}
