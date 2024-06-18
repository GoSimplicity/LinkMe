package api

import (
	. "LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/domain/events/email"
	"LinkMe/internal/domain/events/sms"
	"LinkMe/internal/service"
	. "LinkMe/pkg/ginp"
	"LinkMe/utils"
	ijwt "LinkMe/utils/jwt"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"strconv"
)

const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	Email         *regexp.Regexp
	PassWord      *regexp.Regexp
	svc           service.UserService
	ijwt          ijwt.Handler
	l             *zap.Logger
	smsProducer   sms.Producer
	emailProducer email.Producer
}

func NewUserHandler(svc service.UserService, j ijwt.Handler, l *zap.Logger, smsProducer sms.Producer, emailProducer email.Producer) *UserHandler {
	return &UserHandler{
		Email:         regexp.MustCompile(emailRegexPattern, regexp.None),
		PassWord:      regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:           svc,
		ijwt:          j,
		l:             l,
		smsProducer:   smsProducer,
		emailProducer: emailProducer,
	}
}

func (uh *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/users")
	// 使用插件中的泛型函数
	//userGroup.POST("/signup", WrapBody[SignUpReq](uh.SignUp))
	userGroup.POST("/signup", WrapBody(uh.SignUp))
	userGroup.POST("/login", WrapBody(uh.Login))
	userGroup.POST("/send_sms", WrapBody(uh.SendSMS))
	userGroup.POST("/send_email", WrapBody(uh.SendEmail))
	userGroup.POST("/logout", uh.Logout)
	userGroup.PUT("/refresh_token", uh.RefreshToken)
	userGroup.POST("/change_password", WrapBody(uh.ChangePassword))
	userGroup.DELETE("/write_off", WrapBody(uh.WriteOff))
	userGroup.GET("/profile/:UserID", WrapBody(uh.GetProfile))
	userGroup.PUT("/profile/:UserID", WrapBody(uh.UpdateProfileByID))
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
	uh.l.Error("signup failed", zap.Error(err))
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
	uh.l.Error("login failed", zap.Error(err))
	return Result{
		Code: UserInternalServerError,
	}, err
}

// Logout 登出
func (uh *UserHandler) Logout(ctx *gin.Context) {
	// 清除JWT令牌
	if err := uh.ijwt.ClearToken(ctx); err != nil {
		uh.l.Error("logout failed", zap.Error(err))
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

func (uh *UserHandler) SendSMS(ctx *gin.Context, req SMSReq) (Result, error) {
	valid := utils.IsValidNumber(req.Number)
	if !valid {
		uh.l.Error("电话号码无效", zap.String("number: ", req.Number))
		return Result{
			Code: SMSNumberErr,
			Msg:  InvalidNumber,
		}, nil
	}
	if err := uh.smsProducer.ProduceSMSCode(ctx, sms.SMSCodeEvent{Number: req.Number}); err != nil {
		uh.l.Error("kafka produce sms failed", zap.Error(err))
		return Result{}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserSendSMSCodeSuccess,
	}, nil
}

func (uh *UserHandler) ChangePassword(ctx *gin.Context, req ChangeReq) (Result, error) {
	// 检查新密码和确认密码是否匹配
	if req.NewPassword != req.ConfirmPassword {
		return Result{
			Code: UserInvalidInput,
			Msg:  "新密码和确认密码不一致",
		}, nil
	}
	err := uh.svc.ChangePassword(ctx.Request.Context(), req.Email, req.Password, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			return Result{
				Code: UserInvalidOrPassword,
				Msg:  "旧密码错误或用户不存在",
			}, nil
		}
		uh.l.Error("change password failed", zap.Error(err))
		return Result{
			Code: UserInternalServerError,
			Msg:  "更改密码失败",
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  "密码更改成功",
	}, nil
}

func (uh *UserHandler) SendEmail(ctx *gin.Context, req EmailReq) (Result, error) {
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
	if err = uh.emailProducer.ProduceEmail(ctx, email.EmailEvent{Email: req.Email}); err != nil {
		return Result{}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserSendEmailCodeSuccess,
	}, nil
}

func (uh *UserHandler) WriteOff(ctx *gin.Context, req DeleteUserReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	err := uh.svc.DeleteUser(ctx, req.Email, req.Password, uc.Uid)
	if err != nil {
		return Result{
			Code: ServerERROR,
			Msg:  UserDeletedFailure,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserDeletedSuccess,
	}, nil
}
func (uh *UserHandler) GetProfile(ctx *gin.Context, req ProfileReq) (Result, error) {
	userIDStr := ctx.Param("UserID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		return Result{
			Code: UserInvalidOrProfileError,
			Msg:  "Invalid userID",
		}, err
	}
	profile, err := uh.svc.GetProfileByUserID(ctx, userID)
	if err != nil {
		return Result{
			Code: UserInvalidOrProfileError,
			Msg:  UserProfileGetFailure,
		}, err
	}
	return Result{
		Code: UserValidProfile,
		Msg:  UserProfileGetSuccess,
		Data: profile,
	}, nil
}

func (uh *UserHandler) UpdateProfileByID(ctx *gin.Context, req ProfileReq) (Result, error) {
	userIDStr := ctx.Param("UserID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		return Result{
			Code: UserInvalidOrProfileError,
			Msg:  "Invalid userID",
		}, err
	}
	req.UserID = userID

	profile, err := uh.svc.GetProfileByUserID(ctx, req.UserID)
	if err != nil {
		return Result{
			Code: UserInvalidOrProfileError,
			Msg:  UserProfileGetFailure,
			Data: profile,
		}, err
	}

	profile.Bio = req.Bio
	profile.NickName = req.NickName
	profile.Avatar = req.Avatar

	err = uh.svc.UpdateProfile(ctx, profile)
	if err != nil {
		return Result{
			Code: UserInvalidOrProfileError,
			Msg:  UserProfileUpdateFailure,
			Data: profile,
		}, err
	}
	return Result{
		Code: UserValidProfile,
		Msg:  UserProfileUpdateSuccess,
		Data: profile,
	}, nil
}
