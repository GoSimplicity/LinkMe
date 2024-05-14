package web

import (
	. "LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/service"
	tools "LinkMe/internal/tools/jwt"
	. "LinkMe/pkg/gin-plug"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	Email    *regexp.Regexp
	PassWord *regexp.Regexp
	svc      service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		Email:    regexp.MustCompile(emailRegexPattern, regexp.None),
		PassWord: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:      svc,
	}
}

func (uh *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/users")
	userGroup.POST("/signup", WrapBody(uh.SignUp))
	userGroup.POST("/login", WrapBody(uh.Login))
	userGroup.POST("/logout", uh.Logout)
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
	return Result{
		Code: UserInternalServerError,
		Msg:  "系统异常",
	}, err
}

// Login 登陆
func (uh *UserHandler) Login(ctx *gin.Context, req LoginReq) (Result, error) {
	u, err := uh.svc.Login(ctx, req.Email, req.Password)
	t := tools.NewJWTHandler()
	if err == nil {
		err = t.SetLoginToken(ctx, u.ID)
		return Result{
			Code: RequestsOK,
			Msg:  "登陆成功",
			Data: u,
		}, nil
	} else if errors.Is(err, service.ErrInvalidUserOrPassword) {
		return Result{
			Code: UserInvalidOrPassword,
			Msg:  "用户名或密码不对",
		}, nil
	}
	return Result{
		Code: UserInternalServerError,
		Msg:  "系统错误",
	}, err
}

// Logout 登出
func (uh *UserHandler) Logout(ctx *gin.Context) {
	fmt.Println("退出登陆")
}
