package ginp

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type TokenResult struct {
	Code         int    `json:"code"`
	Msg          string `json:"msg"`
	JWTToken     string `json:"jwt_token"`
	RefreshToken string `json:"refresh_token"`
}
