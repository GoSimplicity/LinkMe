package dto

type SignUpReq struct {
	Username        string `json:"username" binding:"required"`        // 用户名
	Password        string `json:"password" binding:"required"`        // 密码
	ConfirmPassword string `json:"confirmPassword" binding:"required"` // 确认密码
}

type LoginReq struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

type SendSMSReq struct {
	Number string `json:"number" binding:"required"` // 手机号码
}

type LoginSMSReq struct {
	Code string `json:"code" binding:"required"`
}

type ChangePasswordReq struct {
	Username        string `json:"username" binding:"required"`        // 用户名
	Password        string `json:"password" binding:"required"`        // 当前密码
	NewPassword     string `json:"newPassword" binding:"required"`     // 新密码
	ConfirmPassword string `json:"confirmPassword" binding:"required"` // 确认新密码
}

type SendEmailReq struct {
	Email string `json:"email" binding:"required"` // 邮箱
}

type DeleteUserReq struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

type UpdateProfileReq struct {
	RealName string `json:"realName,omitempty"` // 真实姓名
	Avatar   string `json:"avatar,omitempty"`   // 头像URL
	About    string `json:"about,omitempty"`    // 个人简介
	Birthday string `json:"birthday,omitempty"` // 生日
	Phone    string `json:"phone,omitempty"`    // 手机号
}

type UpdateProfileAdminReq struct {
	UserID   int64  `json:"userId" binding:"required"` // 用户ID
	RealName string `json:"realName,omitempty"`        // 真实姓名
	Avatar   string `json:"avatar,omitempty"`          // 头像URL
	About    string `json:"about,omitempty"`           // 个人简介
	Birthday string `json:"birthday,omitempty"`        // 生日
	Phone    string `json:"phone,omitempty"`           // 手机号
}

type ListUserReq struct {
	Page int    `json:"page,omitempty" form:"page"`
	Size *int64 `json:"size,omitempty" form:"size"`
}

type LogoutReq struct {
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"` // 刷新令牌
}

type GetProfileReq struct {
}
