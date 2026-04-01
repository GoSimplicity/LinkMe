package middleware

import (
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
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
		tokenStr := m.ExtractToken(ctx)
		if tokenStr == "" {
			return
		}

		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("jwt.auth_key")), nil
		})
		if err != nil {
			apiresponse.UnauthorizedErrorWithDetails(ctx, nil, "登录态无效")
			ctx.Abort()
			return
		}
		if token == nil || !token.Valid {
			apiresponse.UnauthorizedErrorWithDetails(ctx, nil, "登录态无效")
			ctx.Abort()
			return
		}
		if uc.UserAgent == "" {
			apiresponse.UnauthorizedErrorWithDetails(ctx, nil, "登录态无效")
			ctx.Abort()
			return
		}
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			apiresponse.UnauthorizedErrorWithDetails(ctx, nil, "登录态已失效")
			ctx.Abort()
			return
		}
		ctx.Set("user", uc)
	}
}
