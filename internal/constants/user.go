package constants

const (
	RequestsOK              = 200
	RequestsERROR           = 401
	ServerERROR             = 501
	UserInvalidInput        = 401001 // 输入错误
	UserInvalidOrPassword   = 401002 // 用户名或密码不对
	UserDuplicateEmail      = 401003 // 邮箱冲突
	UserInternalServerError = 501001 // 系统错误
)
