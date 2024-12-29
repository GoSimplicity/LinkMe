package api

import (
	"errors"

	"github.com/GoSimplicity/LinkMe/internal/api/req"
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/email"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/sms"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	"github.com/GoSimplicity/LinkMe/utils"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
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
	//casbinMiddleware := middleware.NewCasbinMiddleware(uh.ce)
	userGroup := server.Group("/api/user")

	userGroup.POST("/signup", uh.SignUp)                    // 用户注册
	userGroup.POST("/login", uh.Login)                      // 用户登录
	userGroup.POST("/login_sms", uh.LoginSMS)               // 短信登录
	userGroup.POST("/send_sms", uh.SendSMS)                 // 发送短信验证码
	userGroup.POST("/send_email", uh.SendEmail)             // 发送邮件验证码
	userGroup.POST("/logout", uh.Logout)                    // 用户登出
	userGroup.POST("/refresh_token", uh.RefreshToken)       // 刷新令牌
	userGroup.POST("/change_password", uh.ChangePassword)   // 修改密码
	userGroup.DELETE("/write_off", uh.WriteOff)             // 注销用户
	userGroup.GET("/profile", uh.GetProfile)                // 获取用户资料
	userGroup.POST("/profile/update", uh.UpdateProfile)     // 更新用户资料(管理员)
	userGroup.POST("/update_profile", uh.UpdateProfileByID) // 更新用户资料
	//userGroup.POST("/list", casbinMiddleware.CheckCasbin(), uh.ListUser) // 获取用户列表（管理员使用）
	userGroup.POST("/list", uh.ListUser) // 获取用户列表（管理员使用）
	userGroup.GET("/codes", uh.GetCodes) // 获取权限码
}

// SignUp 用户注册
func (uh *UserHandler) SignUp(ctx *gin.Context) {
	var req req.SignUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	if req.Password != req.ConfirmPassword {
		apiresponse.ErrorWithMessage(ctx, UserPasswordMismatchError)
		return
	}

	// 尝试注册用户
	err := uh.svc.SignUp(ctx.Request.Context(), domain.User{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	apiresponse.Success(ctx)
}

// Login 用户登录
func (uh *UserHandler) Login(ctx *gin.Context) {
	var req req.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	// 登录验证
	du, err := uh.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			apiresponse.ErrorWithMessage(ctx, UserLoginFailure)
			return
		}
		apiresponse.ErrorWithMessage(ctx, UserLoginFailure)
		return
	}

	// 生成令牌
	jwtToken, refreshToken, err := uh.ijwt.SetLoginToken(ctx, du.ID)
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, UserLoginFailure)
		return
	}

	apiresponse.SuccessWithData(ctx, map[string]string{
		"accessToken":  jwtToken,
		"refreshToken": refreshToken,
	})
}

// Logout 用户登出
func (uh *UserHandler) Logout(ctx *gin.Context) {
	if err := uh.ijwt.ClearToken(ctx); err != nil {
		apiresponse.ErrorWithMessage(ctx, UserLogoutFailure)
		return
	}

	apiresponse.Success(ctx)
}

// RefreshToken 刷新令牌
func (uh *UserHandler) RefreshToken(ctx *gin.Context) {
	var req req.RefreshTokenReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	var rc ijwt.RefreshClaims

	// 验证refresh token
	if ok, claims, err := uh.ijwt.VerifyRefreshToken(ctx, req.RefreshToken); !ok || err != nil {
		apiresponse.ErrorWithMessage(ctx, UserRefreshTokenFailure)
		return
	} else {
		rc = *claims
	}

	// 刷新令牌
	tokenStr, err := uh.ijwt.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, UserRefreshTokenFailure)
		return
	}

	apiresponse.SuccessWithData(ctx, tokenStr)
}

// SendSMS 发送短信验证码
func (uh *UserHandler) SendSMS(ctx *gin.Context) {
	var req req.SMSReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	if !utils.IsValidNumber(req.Number) {
		apiresponse.ErrorWithMessage(ctx, "无效的手机号码")
		return
	}

	if err := uh.smsProducer.ProduceSMSCode(ctx, sms.SMSCodeEvent{Number: req.Number}); err != nil {
		apiresponse.ErrorWithMessage(ctx, "发送短信验证码失败")
		return
	}

	apiresponse.Success(ctx)
}

// ChangePassword 修改密码
func (uh *UserHandler) ChangePassword(ctx *gin.Context) {
	var req req.ChangeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		apiresponse.ErrorWithMessage(ctx, UserPasswordMismatchError)
		return
	}

	err := uh.svc.ChangePassword(ctx.Request.Context(), req.Username, req.Password, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			apiresponse.ErrorWithMessage(ctx, UserPasswordChangeFailure)
			return
		}
		apiresponse.ErrorWithMessage(ctx, UserPasswordChangeFailure)
		return
	}

	apiresponse.Success(ctx)
}

// SendEmail 发送邮件验证码
func (uh *UserHandler) SendEmail(ctx *gin.Context) {
	var req req.UsernameReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	if err := uh.emailProducer.ProduceEmail(ctx, email.EmailEvent{Email: req.Username}); err != nil {
		apiresponse.ErrorWithMessage(ctx, "发送邮件验证码失败")
		return
	}

	apiresponse.Success(ctx)
}

// WriteOff 注销用户
func (uh *UserHandler) WriteOff(ctx *gin.Context) {
	var req req.DeleteUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	// 删除用户
	if err := uh.svc.DeleteUser(ctx, req.Username, req.Password, uc.Uid); err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			apiresponse.ErrorWithMessage(ctx, UserDeletedFailure)
			return
		}
		apiresponse.ErrorWithMessage(ctx, UserDeletedFailure)
		return
	}

	// 清除令牌
	if err := uh.ijwt.ClearToken(ctx); err != nil {
		apiresponse.ErrorWithMessage(ctx, "清除令牌失败")
		return
	}

	apiresponse.Success(ctx)
}

// GetProfile 获取用户资料
func (uh *UserHandler) GetProfile(ctx *gin.Context) {
	var req req.GetProfileReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	profile, err := uh.svc.GetProfileByUserID(ctx, uc.Uid)
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, UserProfileGetFailure)
		return
	}

	apiresponse.SuccessWithData(ctx, profile)
}

// UpdateProfileByID 更新用户资料
func (uh *UserHandler) UpdateProfileByID(ctx *gin.Context) {
	var req req.UpdateProfileReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	err := uh.svc.UpdateProfile(ctx, domain.Profile{
		RealName: req.RealName,
		Avatar:   req.Avatar,
		About:    req.About,
		Birthday: req.Birthday,
		Phone:    &req.Phone,
		UserID:   uc.Uid,
	})
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, UserProfileUpdateFailure)
		return
	}

	apiresponse.Success(ctx)
}

// LoginSMS 短信登录
func (uh *UserHandler) LoginSMS(ctx *gin.Context) {
	var req req.LoginSMSReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}
	// TODO: 实现短信登录逻辑
	apiresponse.Success(ctx)
}

// ListUser 获取用户列表（管理员使用）
func (uh *UserHandler) ListUser(ctx *gin.Context) {
	var req req.ListUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	users, err := uh.svc.ListUser(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, UserListError)
		return
	}

	apiresponse.SuccessWithData(ctx, users)
}

// GetCodes 获取权限码
func (uh *UserHandler) GetCodes(ctx *gin.Context) {
	apiresponse.Success(ctx)
}

// UpdateProfile 更新用户资料(管理员)
func (uh *UserHandler) UpdateProfile(ctx *gin.Context) {
	var req req.UpdateProfileAdminReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	err := uh.svc.UpdateProfileAdmin(ctx, domain.Profile{
		RealName: req.RealName,
		Avatar:   req.Avatar,
		About:    req.About,
		Birthday: req.Birthday,
		Phone:    &req.Phone,
		UserID:   req.UserID,
	})
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, UserProfileUpdateFailure)
		return
	}

	apiresponse.Success(ctx)
}
