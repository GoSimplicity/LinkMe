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

type ProfileReq struct {
	UserId   int64  `json:"userId"`
	Avatar   string `json:"avatar"`
	ID       int64  `json:"id"`
	NickName string `json:"nickName"`
	Bio      string `json:"bio"`
}
