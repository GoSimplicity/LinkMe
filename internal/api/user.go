package api

import (
	. "LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/service"
	ijwt "LinkMe/internal/utils/jwt"
	. "LinkMe/pkg/ginp"
	"LinkMe/pkg/logger"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	Email    *regexp.Regexp
	PassWord *regexp.Regexp
	svc      service.UserService
	ijwt     ijwt.Handler
	l        logger.Logger
}

func NewUserHandler(svc service.UserService, j ijwt.Handler, l logger.Logger) *UserHandler {
	return &UserHandler{
		Email:    regexp.MustCompile(emailRegexPattern, regexp.None),
		PassWord: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:      svc,
		ijwt:     j,
		l:        l,
	}
}

func (uh *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/users")
	userGroup.POST("/signup", WrapBody(uh.SignUp))
	userGroup.POST("/login", WrapBody(uh.Login))
	userGroup.POST("/logout", uh.Logout)
	userGroup.POST("/refresh_token", uh.RefreshToken)
	userGroup.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, "hello world!")
	})
}

// SignUp 注册
func (uh *UserHandler) SignUp(ctx *gin.Context, req SignUpReq) (Result, error) {
	emailBool, err := uh.Email.MatchString(req.Email)
	if err != nil {
		return Result{}, err
	}
	if !emailBool {
		return Result{
			Code: UserInvalidInput,
			Msg:  "账号注册失败,请检查邮箱格式",
		}, nil
	}
	if req.Password != req.ConfirmPassword {
		return Result{
			Code: UserInvalidInput,
			Msg:  "输入的两次密码不同,请重新输入",
		}, nil
	}
	passwordBool, err := uh.PassWord.MatchString(req.Password)
	if err != nil {
		return Result{}, err
	}
	if !passwordBool {
		return Result{
			Code: UserInvalidInput,
			Msg:  "密码必须包含字母、数字、特殊字符",
		}, nil
	}
	err = uh.svc.SignUp(ctx.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == nil {
		return Result{
			Code: RequestsOK,
			Msg:  "注册成功",
		}, nil
	} else if errors.Is(err, service.ErrDuplicateEmail) {
		return Result{
			Code: UserDuplicateEmail,
			Msg:  "邮箱冲突",
		}, nil
	}
	uh.l.Error("注册失败", logger.Error(err))
	return Result{
		Code: UserInternalServerError,
		Msg:  "系统异常",
	}, err
}

// Login 登陆
func (uh *UserHandler) Login(ctx *gin.Context, req LoginReq) (Result, error) {
	du, err := uh.svc.Login(ctx, req.Email, req.Password)
	if err == nil {
		err = uh.ijwt.SetLoginToken(ctx, du.ID)
		return Result{
			Code: RequestsOK,
			Msg:  "登陆成功",
			Data: du,
		}, nil
	} else if errors.Is(err, service.ErrInvalidUserOrPassword) {
		return Result{
			Code: UserInvalidOrPassword,
			Msg:  "用户名或密码不对",
		}, nil
	}
	uh.l.Error("登陆失败", logger.Error(err))
	return Result{
		Code: UserInternalServerError,
		Msg:  "系统错误",
	}, err
}

// Logout 登出
func (uh *UserHandler) Logout(ctx *gin.Context) {
	// 清除JWT令牌
	if err := uh.ijwt.ClearToken(ctx); err != nil {
		uh.l.Error("登出失败", logger.Error(err))
		ctx.JSON(ServerERROR, gin.H{"error": "系统异常"})
		return
	}
	ctx.JSON(RequestsOK, gin.H{"message": "登出成功"})
}

// RefreshToken 刷新令牌
func (uh *UserHandler) RefreshToken(ctx *gin.Context) {
	// 该方法需配合前端使用，前端在Authorization中携带长token
	// 长token只用于刷新短token，短token用于身份验证
	var rc ijwt.RefreshClaims
	// 从前端的Authorization中取出token
	tokenString := uh.ijwt.ExtractToken(ctx)
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.Key2, nil
	})
	if err != nil {
		ctx.AbortWithStatus(ServerERROR)
		return
	}
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(ServerERROR)
		return
	}
	// 检查会话状态是否异常
	if err = uh.ijwt.CheckSession(ctx, rc.Ssid); err != nil {
		ctx.AbortWithStatus(ServerERROR)
		return
	}
	// 刷新短token
	if err = uh.ijwt.SetJWTToken(ctx, rc.Uid, rc.Ssid); err != nil {
		ctx.AbortWithStatus(ServerERROR)
		return
	}
	ctx.JSON(RequestsOK, gin.H{
		"message": "令牌刷新成功",
	})
}
