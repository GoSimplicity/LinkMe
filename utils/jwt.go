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

package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int64) (string, string, error)
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) (string, error)
	ExtractToken(ctx *gin.Context) string
	CheckSession(ctx *gin.Context, ssid string) error
	VerifyRefreshToken(ctx *gin.Context, token string) (bool, *RefreshClaims, error)
	ClearToken(ctx *gin.Context) error
	setRefreshToken(ctx *gin.Context, uid int64, ssid string) (string, error)
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid         int64
	Ssid        string
	UserAgent   string
	ContentType string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}

type handler struct {
	client        redis.Cmdable
	signingMethod jwt.SigningMethod
	rcExpiration  time.Duration
	key1          []byte
	key2          []byte
	issuer        string
}

func NewJWTHandler(c redis.Cmdable) Handler {
	key1 := viper.GetString("jwt.auth_key")
	key2 := viper.GetString("jwt.refresh_key")
	issuer := viper.GetString("jwt.issuer")

	return &handler{
		client:        c,
		signingMethod: jwt.SigningMethodHS512,
		rcExpiration:  time.Hour * 24 * 7,
		key1:          []byte(key1),
		key2:          []byte(key2),
		issuer:        issuer,
	}
}

// SetLoginToken 设置长短Token
func (h *handler) SetLoginToken(ctx *gin.Context, uid int64) (string, string, error) {
	ssid := uuid.New().String()
	refreshToken, err := h.setRefreshToken(ctx, uid, ssid)
	if err != nil {
		return "", "", err
	}

	jwtToken, err := h.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return "", "", err
	}

	return jwtToken, refreshToken, nil
}

// SetJWTToken 设置短Token
func (h *handler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) (string, error) {
	// 从配置文件中获取JWT的过期时间
	expirationMinutes := viper.GetInt64("jwt.expiration")

	// 如果未设置或值无效，设置一个默认值，例如30分钟
	if expirationMinutes <= 0 {
		expirationMinutes = 30
	}

	uc := UserClaims{
		Uid:         uid,
		Ssid:        ssid,
		UserAgent:   ctx.GetHeader("User-Agent"),
		ContentType: ctx.GetHeader("Content-Type"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expirationMinutes))),
			Issuer:    h.issuer,
		},
	}

	token := jwt.NewWithClaims(h.signingMethod, uc)
	// 进行签名
	signedString, err := token.SignedString(h.key1)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

// setRefreshToken 设置长Token
func (h *handler) setRefreshToken(_ *gin.Context, uid int64, ssid string) (string, error) {
	rc := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置刷新时间为一周
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.rcExpiration)),
		},
	}

	t := jwt.NewWithClaims(h.signingMethod, rc)
	signedString, err := t.SignedString(h.key2)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

// ExtractToken 提取 Authorization 头部中的 Token
func (h *handler) ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return ""
	}

	// Authorization 头部格式需为 Bearer string
	s := strings.Split(authCode, " ")
	if len(s) != 2 {
		return ""
	}
	return s[1]
}

// CheckSession 检查会话状态
func (h *handler) CheckSession(ctx *gin.Context, ssid string) error {
	// 判断缓存中是否存在指定键
	c, err := h.client.Exists(ctx, fmt.Sprintf("linkme:user:ssid:%s", ssid)).Result()
	if err != nil {
		return err
	}

	if c != 0 {
		return errors.New("token失效")
	}

	return nil
}

// ClearToken 清空token
func (h *handler) ClearToken(ctx *gin.Context) error {
	// 获取 Authorization 头部中的 token
	authToken, err := h.extractBearerToken(ctx)
	if err != nil {
		return err
	}

	// 提取 token 的 claims 信息
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		return h.key1, nil
	})

	if err != nil || !token.Valid {
		return errors.New("invalid authorization token")
	}

	// 将 token 加入 Redis 黑名单
	if err := h.addToBlacklist(ctx, authToken, claims.ExpiresAt.Time); err != nil {
		return err
	}

	return nil
}

// extractBearerToken 提取 Bearer Token
func (h *handler) extractBearerToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization token")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("invalid authorization token format")
	}

	return authHeader[len(bearerPrefix):], nil
}

// addToBlacklist 将 token 加入 Redis 黑名单
func (h *handler) addToBlacklist(ctx *gin.Context, authToken string, expiresAt time.Time) error {
	remainingTime := time.Until(expiresAt)
	blacklistKey := fmt.Sprintf("blacklist:token:%s", authToken)

	// 将 token 存入 Redis，并设置过期时间
	if err := h.client.Set(ctx, blacklistKey, "invalid", remainingTime).Err(); err != nil {
		return err
	}
	return nil
}

// VerifyRefreshToken 验证refresh token
func (h *handler) VerifyRefreshToken(ctx *gin.Context, token string) (bool, *RefreshClaims, error) {
	// 解析refresh token
	refreshClaims := &RefreshClaims{}
	refreshToken, err := jwt.ParseWithClaims(token, refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return h.key2, nil
	})

	// 检查解析和验证结果
	if err != nil || !refreshToken.Valid {
		return false, nil, errors.New("无效的refresh token")
	}

	// 检查token是否在黑名单中
	exists, err := h.client.Exists(ctx, fmt.Sprintf("linkme:user:ssid:%s", refreshClaims.Ssid)).Result()
	if err != nil {
		return false, nil, fmt.Errorf("检查token状态失败: %v", err)
	}

	// 如果token在黑名单中,说明已经失效
	if exists > 0 {
		return false, nil, errors.New("refresh token已失效")
	}

	return true, refreshClaims, nil
}
