package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var (
	Key1 = []byte("sadfkhjlkkljKFJDSLAFUDASLFJKLjfj113d1")
	Key2 = []byte("sadfkhjlkkljKFJDSLAFUDASLFJKLjfj113d2")
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	ExtractToken(ctx *gin.Context) string
	CheckSession(ctx *gin.Context, ssid string) error
	ClearToken(ctx *gin.Context) error
	setRefreshToken(ctx *gin.Context, uid int64, ssid string) error
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
}

func NewJWTHandler(c redis.Cmdable) Handler {
	return &handler{
		client:        c,
		signingMethod: jwt.SigningMethodHS512,
		rcExpiration:  time.Hour * 24 * 7,
	}
}

// SetLoginToken 设置长短Token
func (h *handler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	if err := h.setRefreshToken(ctx, uid, ssid); err != nil {
		return err
	}
	return h.SetJWTToken(ctx, uid, ssid)
}

// SetJWTToken 设置短Token
func (h *handler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := UserClaims{
		Uid:         uid,
		Ssid:        ssid,
		UserAgent:   ctx.GetHeader("User-Agent"),
		ContentType: ctx.GetHeader("Content-Type"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}
	token := jwt.NewWithClaims(h.signingMethod, uc)
	// 进行签名
	signedString, err := token.SignedString(Key1)
	if err != nil {
		return err
	}
	ctx.Header("X-JWT-Token", signedString)
	return nil
}

// setRefreshToken 设置长Token
func (h *handler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	rc := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置刷新时间为一周
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.rcExpiration)),
		},
	}
	t := jwt.NewWithClaims(h.signingMethod, rc)
	signedString, err := t.SignedString(Key2)
	if err != nil {
		return err
	}
	ctx.Header("X-Refresh-Token", signedString)
	return err
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
	ctx.Header("X-Refresh-Token", "")
	ctx.Header("X-JWT-Token", "")
	uc := ctx.MustGet("user").(UserClaims)
	return h.client.Set(ctx, fmt.Sprintf("linkme:user:ssid:%s", uc.Ssid), "", h.rcExpiration).Err()
}
