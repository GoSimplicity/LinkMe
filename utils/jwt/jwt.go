package jwt

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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
	client         redis.Cmdable
	signingMethod  jwt.SigningMethod
	rcExpiration   time.Duration
	AuthExpiration time.Duration
	Issuer         string
	AuthKey        []byte
	RefreshKey     []byte
}

func NewJWTHandler(c redis.Cmdable) Handler {
	return &handler{
		client:         c,
		signingMethod:  jwt.SigningMethodHS512,
		rcExpiration:   time.Hour * viper.GetDuration("jwt.refresh_expire"),
		AuthExpiration: time.Minute * viper.GetDuration("jwt.auth_expire"),
		AuthKey:        []byte(viper.GetString("jwt.auth_key")),
		RefreshKey:     []byte(viper.GetString("jwt.refresh_key")),
		Issuer:         viper.GetString("jwt.issuer"),
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
	uc := UserClaims{
		Uid:         uid,
		Ssid:        ssid,
		UserAgent:   ctx.GetHeader("User-Agent"),
		ContentType: ctx.GetHeader("Content-Type"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.AuthExpiration)),
			Issuer:    h.Issuer,
		},
	}
	token := jwt.NewWithClaims(h.signingMethod, uc)
	// 进行签名
	signedString, err := token.SignedString(h.AuthKey)
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
	signedString, err := t.SignedString(h.RefreshKey)
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
	uc := ctx.MustGet("user").(UserClaims)
	// 获取 refresh token
	refreshTokenString := ctx.GetHeader("X-Refresh-Token")
	if refreshTokenString == "" {
		return errors.New("missing refresh token")
	}

	ctx.Header("X-Refresh-Token", "")

	// 解析 refresh token
	refreshClaims := &RefreshClaims{}
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return h.RefreshKey, nil
	})
	if err != nil || !refreshToken.Valid {
		return errors.New("invalid refresh token")
	}
	// 设置redis中的会话ID键的过期时间为refresh token的剩余过期时间
	remainingTime := refreshClaims.ExpiresAt.Time.Sub(time.Now())
	if er := h.client.Set(ctx, fmt.Sprintf("linkme:user:ssid:%s", uc.Ssid), "", remainingTime).Err(); er != nil {
		return er
	}
	return nil
}

// VerifyRefreshToken 验证refresh token
func (h *handler) VerifyRefreshToken(ctx *gin.Context, token string) (bool, *RefreshClaims, error) {
	// 解析refresh token
	refreshClaims := &RefreshClaims{}
	refreshToken, err := jwt.ParseWithClaims(token, refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return h.RefreshKey, nil
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
