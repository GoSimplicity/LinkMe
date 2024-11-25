package constants

const (
	UserInvalidInputCode          = 401001 // 用户输入错误
	UserInvalidOrPasswordCode     = 401002 // 用户名或密码错误
	UserInvalidOrProfileErrorCode = 401003 // 用户资料无效
	UserEmailFormatErrorCode      = 401004 // 邮箱格式错误
	UserPasswordMismatchErrorCode = 401005 // 两次输入的密码不一致
	UserPasswordFormatErrorCode   = 401006 // 密码格式错误
	UserEmailConflictErrorCode    = 401007 // 邮箱冲突
	UserListErrorCode             = 401008 // 用户获取失败
	UserServerErrorCode           = 500001 // 用户服务内部错误
	UserGetCountErrorCode         = 401009 // 获取用户数量失败
	UserSignUpSuccess            = "用户注册成功"
	UserGetCountError            = "获取用户数量失败"
	UserGetCountSuccess          = "获取用户数量成功"
	UserListError                = "用户获取失败"
	UserListSuccess              = "用户获取成功"
	UserSignUpFailure            = "用户注册失败"
	UserLoginSuccess             = "用户登录成功"
	UserLoginFailure             = "用户登录失败"
	UserLogoutSuccess            = "用户登出成功"
	UserLogoutFailure            = "用户登出失败"
	UserRefreshTokenSuccess      = "令牌刷新成功"
	UserRefreshTokenFailure      = "令牌刷新失败"
	UserProfileGetSuccess        = "获取用户资料成功"
	UserProfileGetFailure        = "获取用户资料失败"
	UserProfileUpdateSuccess     = "更新用户资料成功"
	UserProfileUpdateFailure     = "更新用户资料失败"
	UserPasswordChangeSuccess    = "密码修改成功"
	UserPasswordChangeFailure    = "密码修改失败"
	UserDeletedSuccess           = "用户删除成功"
	UserDeletedFailure           = "用户删除失败"
	UserSendSMSCodeSuccess       = "短信验证码发送成功"
	UserSendEmailCodeSuccess     = "邮箱验证码发送成功"
	UserEmailFormatError         = "邮箱格式错误，请检查"
	UserPasswordMismatchError    = "两次输入的密码不一致，请重新输入"
	UserPasswordFormatError      = "密码必须包含字母、数字和特殊字符，且长度不少于8位"
	UserEmailConflictError       = "该邮箱已被注册"
)
