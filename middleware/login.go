/*
 * MIT License
 *
 * Copyright (c) 2024 Bamboo
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

package middleware

import (
	"net/http"

	ijwt "github.com/GoSimplicity/LinkMe/utils"
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
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			// token 非法或过期
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 检查是否携带ua头
		if uc.UserAgent == "" {
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
