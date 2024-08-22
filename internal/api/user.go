package api

import (
	"errors"
	"github.com/GoSimplicity/LinkMe/internal/api/req"
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/email"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/sms"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/middleware"
	. "github.com/GoSimplicity/LinkMe/pkg/ginp"
	"github.com/GoSimplicity/LinkMe/utils"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/casbin/casbin/v2"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	ce            *casbin.Enforcer
	smsProducer   sms.Producer
	emailProducer email.Producer
}

func NewUserHandler(svc service.UserService, j ijwt.Handler, smsProducer sms.Producer, emailProducer email.Producer, ce *casbin.Enforcer) *UserHandler {
	return &UserHandler{
		Email:         regexp.MustCompile(emailRegexPattern, regexp.None),
		PassWord:      regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:           svc,
		ijwt:          j,
		ce:            ce,
		smsProducer:   smsProducer,
		emailProducer: emailProducer,
	}
}

// RegisterRoutes 注册用户相关路由
func (uh *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 初始化Casbin中间件
	casbinMiddleware := middleware.NewCasbinMiddleware(uh.ce)
	// 创建用户组路由
	userGroup := server.Group("/api/users")
	// 用户注册
	userGroup.POST("/signup", WrapBody(uh.SignUp))
	// 用户登录
	userGroup.POST("/login", WrapBody(uh.Login))
	// 短信登录
	userGroup.POST("/login_sms", WrapBody(uh.LoginSMS))
	// 发送短信验证码
	userGroup.POST("/send_sms", WrapBody(uh.SendSMS))
	// 发送邮件验证码
	userGroup.POST("/send_email", WrapBody(uh.SendEmail))
	// 用户登出
	userGroup.POST("/logout", WrapNoParam(uh.Logout))
	// 刷新令牌
	userGroup.POST("/refresh_token", WrapBody(uh.RefreshToken))
	// 修改密码
	userGroup.POST("/change_password", WrapBody(uh.ChangePassword))
	// 注销用户
	userGroup.DELETE("/write_off", WrapBody(uh.WriteOff))
	// 获取用户资料
	userGroup.GET("/profile", WrapQuery(uh.GetProfile))
	// 更新用户资料
	userGroup.POST("/update_profile", WrapBody(uh.UpdateProfileByID))
	// 获取用户列表（管理员使用）
	userGroup.POST("/list", casbinMiddleware.CheckCasbin(), WrapBody(uh.ListUser))
	// 获取用户统计（管理员使用）
	userGroup.GET("/stats", casbinMiddleware.CheckCasbin(), WrapQuery(uh.GetUserCount))
}

// SignUp 用户注册
func (uh *UserHandler) SignUp(ctx *gin.Context, req req.SignUpReq) (Result, error) {
	// 验证邮箱格式
	emailValid, err := uh.Email.MatchString(req.Username)
	if err != nil {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserSignUpFailure,
		}, err
	}
	if !emailValid {
		return Result{
			Code: UserEmailFormatErrorCode,
			Msg:  UserEmailFormatError,
		}, nil
	}
	// 验证密码是否一致
	if req.Password != req.ConfirmPassword {
		return Result{
			Code: UserPasswordMismatchErrorCode,
			Msg:  UserPasswordMismatchError,
		}, nil
	}
	// 验证密码格式
	passwordValid, err := uh.PassWord.MatchString(req.Password)
	if err != nil {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserSignUpFailure,
		}, err
	}
	if !passwordValid {
		return Result{
			Code: UserPasswordFormatErrorCode,
			Msg:  UserPasswordFormatError,
		}, nil
	}
	// 尝试注册用户
	err = uh.svc.SignUp(ctx.Request.Context(), domain.User{
		Email:    req.Username,
		Password: req.Password,
	})
	if err != nil {
		// 检查是否为重复邮箱错误
		if errors.Is(err, service.ErrDuplicateEmail) {
			return Result{
				Code: UserEmailConflictErrorCode,
				Msg:  UserEmailConflictError,
			}, nil
		}
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserSignUpFailure,
		}, err
	}
	// 注册成功
	return Result{
		Code: RequestsOK,
		Msg:  UserSignUpSuccess,
	}, nil
}

// Login 登陆
func (uh *UserHandler) Login(ctx *gin.Context, req req.LoginReq) (Result, error) {
	du, err := uh.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			return Result{
				Code: UserInvalidOrPasswordCode,
				Msg:  UserLoginFailure,
			}, nil
		}
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserLoginFailure,
		}, err
	}
	token, er := uh.ijwt.SetLoginToken(ctx, du.ID)
	if er != nil {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserLoginFailure,
		}, er
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserLoginSuccess,
		Data: token,
	}, nil
}

// Logout 登出
func (uh *UserHandler) Logout(ctx *gin.Context) (Result, error) {
	// 清除JWT令牌
	if err := uh.ijwt.ClearToken(ctx); err != nil {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserLogoutFailure,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserLogoutSuccess,
	}, nil
}

// RefreshToken 刷新令牌
func (uh *UserHandler) RefreshToken(ctx *gin.Context, _ req.RefreshTokenReq) (Result, error) {
	var rc ijwt.RefreshClaims
	// 从前端的Authorization中取出token
	tokenString := uh.ijwt.ExtractToken(ctx)
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.Key2, nil
	})
	if err != nil || token == nil || !token.Valid {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserRefreshTokenFailure,
		}, err
	}
	// 检查会话状态是否异常
	if err = uh.ijwt.CheckSession(ctx, rc.Ssid); err != nil {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserRefreshTokenFailure,
		}, err
	}
	// 刷新短token
	tokenStr, er := uh.ijwt.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if er != nil {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserRefreshTokenFailure,
		}, er
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserRefreshTokenSuccess,
		Data: tokenStr,
	}, nil
}

// SendSMS 发送短信验证码
func (uh *UserHandler) SendSMS(ctx *gin.Context, req req.SMSReq) (Result, error) {
	// 验证手机号码格式
	if !utils.IsValidNumber(req.Number) {
		return Result{
			Code: SMSNumberErr,
			Msg:  InvalidNumber,
		}, nil
	}
	// 发送短信验证码
	if err := uh.smsProducer.ProduceSMSCode(ctx, sms.SMSCodeEvent{Number: req.Number}); err != nil {
		return Result{}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserSendSMSCodeSuccess,
	}, nil
}

// ChangePassword 修改密码
func (uh *UserHandler) ChangePassword(ctx *gin.Context, req req.ChangeReq) (Result, error) {
	// 检查新密码和确认密码是否匹配
	if req.NewPassword != req.ConfirmPassword {
		return Result{
			Code: UserInvalidInputCode,
			Msg:  UserPasswordMismatchError,
		}, nil
	}
	// 修改密码
	if err := uh.svc.ChangePassword(ctx.Request.Context(), req.Username, req.Password, req.NewPassword, req.ConfirmPassword); err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			return Result{
				Code: UserInvalidOrPasswordCode,
				Msg:  UserLoginFailure,
			}, nil
		}
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserPasswordChangeFailure,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserPasswordChangeSuccess,
	}, nil
}

// SendEmail 发送邮件
func (uh *UserHandler) SendEmail(ctx *gin.Context, req req.UsernameReq) (Result, error) {
	// 验证邮箱格式
	if emailValid, err := uh.Email.MatchString(req.Username); err != nil {
		return Result{}, err
	} else if !emailValid {
		return Result{
			Code: UserInvalidInputCode,
			Msg:  UserEmailFormatError,
		}, nil
	}
	// 发送邮件
	if err := uh.emailProducer.ProduceEmail(ctx, email.EmailEvent{Email: req.Username}); err != nil {
		return Result{}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserSendEmailCodeSuccess,
	}, nil
}

// WriteOff 注销用户
func (uh *UserHandler) WriteOff(ctx *gin.Context, req req.DeleteUserReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	// 删除用户
	if err := uh.svc.DeleteUser(ctx, req.Username, req.Password, uc.Uid); err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			return Result{
				Code: UserInvalidOrPasswordCode,
				Msg:  UserLoginFailure,
			}, nil
		}
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserDeletedFailure,
		}, err
	}
	// 清除JWT令牌
	if err := uh.ijwt.ClearToken(ctx); err != nil {
		return Result{}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserDeletedSuccess,
	}, nil
}

// GetProfile 获取用户资料
func (uh *UserHandler) GetProfile(ctx *gin.Context, _ req.GetProfileReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	// 获取用户资料
	profile, err := uh.svc.GetProfileByUserID(ctx, uc.Uid)
	if err != nil {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserProfileGetFailure,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserProfileGetSuccess,
		Data: profile,
	}, nil
}

// UpdateProfileByID 更新用户资料
func (uh *UserHandler) UpdateProfileByID(ctx *gin.Context, req req.UpdateProfileReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	// 更新用户资料
	err := uh.svc.UpdateProfile(ctx, domain.Profile{
		NickName: req.NickName,
		Avatar:   req.Avatar,
		About:    req.About,
		Birthday: req.Birthday,
		UserID:   uc.Uid,
	})
	if err != nil {
		return Result{
			Code: UserInvalidOrProfileErrorCode,
			Msg:  UserProfileUpdateFailure,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  UserProfileUpdateSuccess,
	}, nil
}

// LoginSMS 短信登录
func (uh *UserHandler) LoginSMS(ctx *gin.Context, req req.LoginSMSReq) (Result, error) {
	// 此处可以添加短信登录逻辑
	return Result{}, nil
}

// ListUser 获取用户列表（管理员使用）
func (uh *UserHandler) ListUser(ctx *gin.Context, req req.ListUserReq) (Result, error) {
	// 分页查询用户列表
	users, err := uh.svc.ListUser(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return Result{
			Code: UserListErrorCode,
			Msg:  UserListError,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserListSuccess,
		Data: users,
	}, nil
}

// GetUserCount 获取用户统计（管理员使用）
func (uh *UserHandler) GetUserCount(ctx *gin.Context, _ req.GetUserCountReq) (Result, error) {
	// 获取用户总数
	count, err := uh.svc.GetUserCount(ctx)
	if err != nil {
		return Result{
			Code: UserGetCountErrorCode,
			Msg:  UserGetCountError,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  UserGetCountSuccess,
		Data: count,
	}, nil
}
