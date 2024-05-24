package api

import (
	. "LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/service"
	. "LinkMe/pkg/ginp"
	ijwt "LinkMe/utils/jwt"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
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
	l        *zap.Logger
}

func NewUserHandler(svc service.UserService, j ijwt.Handler, l *zap.Logger) *UserHandler {
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
	userGroup.PUT("/refresh_token", uh.RefreshToken)
	// 测试接口
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
			Msg:  UserEmailFormatError,
		}, nil
	}
	if req.Password != req.ConfirmPassword {
		return Result{
			Code: UserInvalidInput,
			Msg:  UserPasswordMismatchError,
		}, nil
	}
	passwordBool, err := uh.PassWord.MatchString(req.Password)
	if err != nil {
		return Result{}, err
	}
	if !passwordBool {
		return Result{
			Code: UserInvalidInput,
			Msg:  UserPasswordFormatError,
		}, nil
	}
	err = uh.svc.SignUp(ctx.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == nil {
		return Result{
			Code: RequestsOK,
			Msg:  UserSignUpSuccess,
		}, nil
	} else if errors.Is(err, service.ErrDuplicateEmail) {
		return Result{
			Code: UserDuplicateEmail,
			Msg:  UserEmailConflictError,
		}, nil
	}
	uh.l.Error("注册失败", zap.Error(err))
	return Result{
		Code: UserInternalServerError,
		Msg:  UserSignUpFailure,
	}, err
}

// Login 登陆
func (uh *UserHandler) Login(ctx *gin.Context, req LoginReq) (Result, error) {
	du, err := uh.svc.Login(ctx, req.Email, req.Password)
	if err == nil {
		err = uh.ijwt.SetLoginToken(ctx, du.ID)
		return Result{
			Code: RequestsOK,
			Msg:  UserLoginSuccess,
			Data: du,
		}, nil
	} else if errors.Is(err, service.ErrInvalidUserOrPassword) {
		return Result{
			Code: UserInvalidOrPassword,
			Msg:  UserLoginFailure,
		}, nil
	}
	uh.l.Error("登陆失败", zap.Error(err))
	return Result{
		Code: UserInternalServerError,
	}, err
}

// Logout 登出
func (uh *UserHandler) Logout(ctx *gin.Context) {
	// 清除JWT令牌
	if err := uh.ijwt.ClearToken(ctx); err != nil {
		uh.l.Error("登出失败", zap.Error(err))
		ctx.JSON(ServerERROR, gin.H{"error": UserLogoutFailure})
		return
	}
	ctx.JSON(RequestsOK, gin.H{"message": UserLogoutSuccess})
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
		"message": UserRefreshTokenSuccess,
	})
}
