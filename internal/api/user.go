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
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	svc           service.UserService
	ijwt          ijwt.Handler
	ce            *casbin.Enforcer
	smsProducer   sms.Producer
	emailProducer email.Producer
}

func NewUserHandler(svc service.UserService, j ijwt.Handler, smsProducer sms.Producer, emailProducer email.Producer, ce *casbin.Enforcer) *UserHandler {
	return &UserHandler{
		svc:           svc,
		ijwt:          j,
		ce:            ce,
		smsProducer:   smsProducer,
		emailProducer: emailProducer,
	}
}

func (uh *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 初始化Casbin中间件
	casbinMiddleware := middleware.NewCasbinMiddleware(uh.ce)
	userGroup := server.Group("/api/user")

	userGroup.POST("/signup", WrapBody(uh.SignUp))                                      // 用户注册
	userGroup.POST("/login", WrapBody(uh.Login))                                        // 用户登录
	userGroup.POST("/login_sms", WrapBody(uh.LoginSMS))                                 // 短信登录
	userGroup.POST("/send_sms", WrapBody(uh.SendSMS))                                   // 发送短信验证码
	userGroup.POST("/send_email", WrapBody(uh.SendEmail))                               // 发送邮件验证码
	userGroup.POST("/logout", WrapNoParam(uh.Logout))                                   // 用户登出
	userGroup.POST("/refresh_token", WrapBody(uh.RefreshToken))                         // 刷新令牌
	userGroup.POST("/change_password", WrapBody(uh.ChangePassword))                     // 修改密码
	userGroup.DELETE("/write_off", WrapBody(uh.WriteOff))                               // 注销用户
	userGroup.GET("/profile", WrapQuery(uh.GetProfile))                                 // 获取用户资料
	userGroup.POST("/update_profile", WrapBody(uh.UpdateProfileByID))                   // 更新用户资料
	userGroup.POST("/list", casbinMiddleware.CheckCasbin(), WrapBody(uh.ListUser))      // 获取用户列表（管理员使用）
	userGroup.GET("/stats", casbinMiddleware.CheckCasbin(), WrapQuery(uh.GetUserCount)) // 获取用户统计（管理员使用）
}

// SignUp 用户注册
func (uh *UserHandler) SignUp(ctx *gin.Context, req req.SignUpReq) (Result, error) {
	// 验证密码是否一致
	if req.Password != req.ConfirmPassword {
		return Result{
			Code: UserPasswordMismatchErrorCode,
			Msg:  UserPasswordMismatchError,
		}, nil
	}

	// 尝试注册用户
	err := uh.svc.SignUp(ctx.Request.Context(), domain.User{
		Email:    req.Username,
		Password: req.Password,
	})

	// 根据错误类型返回不同的响应
	switch {
	case err == nil:
		return Result{
			Code: RequestsOK,
			Msg:  UserSignUpSuccess,
		}, nil
	case errors.Is(err, service.ErrDuplicateUsername):
		return Result{
			Code: UserEmailConflictErrorCode,
			Msg:  UserEmailConflictError,
		}, nil
	case err.Error() == "invalid email format":
		return Result{
			Code: UserEmailFormatErrorCode,
			Msg:  UserEmailFormatError,
		}, nil
	case err.Error() == "invalid password format":
		return Result{
			Code: UserPasswordFormatErrorCode,
			Msg:  UserPasswordFormatError,
		}, nil
	default:
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserSignUpFailure,
		}, err
	}
}

// Login 用户登录
func (uh *UserHandler) Login(ctx *gin.Context, req req.LoginReq) (Result, error) {
	// 登录验证
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

	// 生成令牌
	jwtToken, refreshToken, err := uh.ijwt.SetLoginToken(ctx, du.ID)
	if err != nil {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserLoginFailure,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  UserLoginSuccess,
		Data: map[string]string{
			"jwt_token":     jwtToken,
			"refresh_token": refreshToken,
		},
	}, nil
}

// Logout 用户登出
func (uh *UserHandler) Logout(ctx *gin.Context) (Result, error) {
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
	tokenString := uh.ijwt.ExtractToken(ctx)

	// 解析并验证令牌
	token, err := jwt.ParseWithClaims(tokenString, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.Key2, nil
	})
	if err != nil || token == nil || !token.Valid {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserRefreshTokenFailure,
		}, err
	}

	// 验证会话
	if err = uh.ijwt.CheckSession(ctx, rc.Ssid); err != nil {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserRefreshTokenFailure,
		}, err
	}

	// 刷新令牌
	tokenStr, err := uh.ijwt.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserRefreshTokenFailure,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  UserRefreshTokenSuccess,
		Data: tokenStr,
	}, nil
}

// SendSMS 发送短信验证码
func (uh *UserHandler) SendSMS(ctx *gin.Context, req req.SMSReq) (Result, error) {
	if !utils.IsValidNumber(req.Number) {
		return Result{
			Code: SMSNumberErr,
			Msg:  "无效的手机号码",
		}, nil
	}

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
	if req.NewPassword != req.ConfirmPassword {
		return Result{
			Code: UserInvalidInputCode,
			Msg:  UserPasswordMismatchError,
		}, nil
	}

	err := uh.svc.ChangePassword(ctx.Request.Context(), req.Username, req.Password, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			return Result{
				Code: UserInvalidOrPasswordCode,
				Msg:  UserPasswordChangeFailure,
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

// SendEmail 发送邮件验证码
func (uh *UserHandler) SendEmail(ctx *gin.Context, req req.UsernameReq) (Result, error) {
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
				Msg:  UserDeletedFailure,
			}, nil
		}
		return Result{
			Code: UserServerErrorCode,
			Msg:  UserDeletedFailure,
		}, err
	}

	// 清除令牌
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
	// TODO: 实现短信登录逻辑
	return Result{}, nil
}

// ListUser 获取用户列表（管理员使用）
func (uh *UserHandler) ListUser(ctx *gin.Context, req req.ListUserReq) (Result, error) {
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
