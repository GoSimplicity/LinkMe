package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	Uid       int64
	Ssid      string
	UserAgent string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}

type handler struct {
	signingMethod jwt.SigningMethod
	rcExpiration  time.Duration
}

func NewJWTHandler() Handler {
	return &handler{
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
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.GetHeader("User-Agent"),
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
			//
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
	// 约定Authorization 头部格式为 Bearer tokenString
	s := strings.Split(authCode, " ")
	if len(s) != 2 {
		return ""
	}
	return s[1]
}

func (h *handler) CheckSession(ctx *gin.Context, ssid string) error {
	return nil
}

func (h *handler) ClearToken(ctx *gin.Context) error {
	return nil
}
