package api

type SignUpReq struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SMSReq struct {
	Number string `json:"number"`
}

type ChangeReq struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

type EmailReq struct {
	Email string `json:"email"`
}

type DeleteUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UpdateProfileReq struct {
	NickName string `json:"nickName"` // 昵称
	Avatar   string `json:"avatar"`   // 头像URL
	About    string `json:"about"`    // 个人简介
	Birthday string `json:"birthday"` // 生日
}

type LoginSMSReq struct {
	Code string `json:"code"`
}

type ListUserReq struct {
	Page int    `json:"page,omitempty"` // 当前页码
	Size *int64 `json:"size,omitempty"` // 每页数据量
}

type GetUserCountReq struct {
}
