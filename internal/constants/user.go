package constants

const (
	UserInvalidInput          = 401001                                                                                              // 用户输入错误
	UserInvalidOrPassword     = 401002                                                                                              // 用户名或密码错误
	UserDuplicateEmail        = 401003                                                                                              // 用户邮箱重复
	UserInternalServerError   = 500001                                                                                              // 用户服务内部错误
	UserSignUpSuccess         = "User registration successful"                                                                      // 用户注册成功
	UserSignUpFailure         = "User registration failed"                                                                          // 用户注册失败
	UserLoginSuccess          = "User login successful"                                                                             // 用户登陆成功
	UserLoginFailure          = "User login failed"                                                                                 // 用户登陆失败
	UserLogoutSuccess         = "User logout successful"                                                                            // 用户登出成功
	UserLogoutFailure         = "User logout failed"                                                                                // 用户登出失败
	UserRefreshTokenSuccess   = "Token refresh successful"                                                                          // 令牌刷新成功
	UserRefreshTokenFailure   = "Token refresh failed"                                                                              // 令牌刷新失败
	UserEmailFormatError      = "Invalid email format, please check"                                                                // 邮箱格式错误
	UserPasswordMismatchError = "The two passwords entered are different, please re-enter"                                          // 两次输入的密码不一致
	UserPasswordFormatError   = "Password must contain letters, numbers, and special characters, and be at least 8 characters long" // 密码格式错误
	UserEmailConflictError    = "Email already exists"                                                                              // 邮箱冲突
)
