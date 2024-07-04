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
	UserGetCountErrorCode         = 401009
	UserSignUpSuccess             = "User registration successful" // 用户注册成功
	UserGetCountError             = "User get count error"
	UserGetCountSuccess           = "User get count success"
	UserListError                 = "User list error"                                                                                   // 用户获取失败
	UserListSuccess               = "User get success"                                                                                  // 用户获取成功
	UserSignUpFailure             = "User registration failed"                                                                          // 用户注册失败
	UserLoginSuccess              = "User login successful"                                                                             // 用户登陆成功
	UserLoginFailure              = "User login failed"                                                                                 // 用户登陆失败
	UserLogoutSuccess             = "User logout successful"                                                                            // 用户登出成功
	UserLogoutFailure             = "User logout failed"                                                                                // 用户登出失败
	UserRefreshTokenSuccess       = "Token refresh successful"                                                                          // 令牌刷新成功
	UserRefreshTokenFailure       = "Token refresh failed"                                                                              // 令牌刷新失败
	UserProfileGetSuccess         = "Profile get successful"                                                                            // 获取用户资料成功
	UserProfileGetFailure         = "Profile get failed"                                                                                // 获取用户资料失败
	UserProfileUpdateSuccess      = "Profile update successful"                                                                         // 更新用户资料成功
	UserProfileUpdateFailure      = "Profile update failed"                                                                             // 更新用户资料失败
	UserPasswordChangeSuccess     = "Password change successful"                                                                        // 密码更改成功
	UserPasswordChangeFailure     = "Password change failed"                                                                            // 密码更改失败
	UserDeletedSuccess            = "User deleted successfully"                                                                         // 用户删除成功
	UserDeletedFailure            = "User deletion failed"                                                                              // 用户删除失败
	UserSendSMSCodeSuccess        = "Code sent successfully"                                                                            // 短信验证码发送成功
	UserSendEmailCodeSuccess      = "Code sent successfully"                                                                            // 邮箱验证码发送成功
	UserEmailFormatError          = "Invalid email format, please check"                                                                // 邮箱格式错误
	UserPasswordMismatchError     = "The two passwords entered are different, please re-enter"                                          // 两次输入的密码不一致
	UserPasswordFormatError       = "Password must contain letters, numbers, and special characters, and be at least 8 characters long" // 密码格式错误
	UserEmailConflictError        = "Email already exists"                                                                              // 邮箱冲突
)
