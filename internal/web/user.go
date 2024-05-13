package web

import (
	"LinkMe/internal/constants"
	"LinkMe/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

const (
	usernameRegexPattern = `^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	UserName *regexp.Regexp
	PassWord *regexp.Regexp
	svc      service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		UserName: regexp.MustCompile(usernameRegexPattern, regexp.None),
		PassWord: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:      svc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/users")
	userGroup.POST("/signup", u.SignUp)
	userGroup.POST("/login", u.Login)
	userGroup.POST("/logout", u.Logout)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	ctx.JSON(constants.RequestsOK, "注册成功")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	ctx.JSON(constants.RequestsOK, "登陆成功")

}

func (u *UserHandler) Logout(ctx *gin.Context) {
	ctx.JSON(constants.RequestsOK, "退出登陆")

}
