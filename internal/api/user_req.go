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
