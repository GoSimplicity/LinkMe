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
	ID int64 `json:"id" binding:"required"`
}

type GetProfileReq struct {
	ID int64 `json:"id" binding:"required"`
}

type UpdateProfileReq struct {
	ID       int64   `json:"id"`
	RealName string  `json:"realName"`
	Avatar   string  `json:"avatar"`
	About    string  `json:"about"`
	Birthday string  `json:"birthday"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
}

type ListUserReq struct {
	Page   int    `json:"page,omitempty" form:"page"`
	Size   int    `json:"size,omitempty" form:"size"`
	Search string `json:"search,omitempty" form:"search"`
}

type LogoutReq struct {
	ID int64 `json:"id" binding:"required"` // 用户ID
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"` // 刷新令牌
}

type GetUserReq struct {
	ID int64 `json:"id" binding:"required"` // 用户ID
}

type UpdateUserReq struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	RealName string  `json:"realName"`
	Avatar   string  `json:"avatar"`
	About    string  `json:"about"`
	Birthday string  `json:"birthday"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
}
