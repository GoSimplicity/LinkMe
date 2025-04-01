package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, userHdl *UserHandler) {
	r := server.Group("/api/v1")
	// 公共接口组 - 无需认证
	publicGroup := r.Group("/users")
	{
		// 注册登录相关
		publicGroup.POST("/register", userHdl.Register)
		publicGroup.POST("/login", userHdl.Login)
		publicGroup.POST("/login-sms", userHdl.LoginSMS)
		publicGroup.POST("/send-sms", userHdl.SendSMS)
		publicGroup.POST("/send-email", userHdl.SendEmail)
		publicGroup.POST("/refresh-token", userHdl.RefreshToken)
	}

	// 需要认证的接口组 - 用户自身操作
	authGroup := r.Group("/users")
	// TODO: 添加中间件进行认证
	// authGroup.Use(middleware.JWTMiddleware())
	{
		authGroup.POST("/logout", userHdl.Logout)
		authGroup.GET("/profile", userHdl.GetProfile)
		authGroup.PUT("/profile", userHdl.UpdateProfile)
		authGroup.PUT("/password", userHdl.ChangePassword)
		authGroup.DELETE("/account", userHdl.DeleteUser)
	}

	// 管理员接口组 - 需要管理员权限
	adminGroup := r.Group("/admin/users")
	// TODO: 添加管理员权限验证中间件
	// adminGroup.Use(middleware.AdminAuth())
	{
		adminGroup.GET("", userHdl.List)
		adminGroup.GET("/:id", userHdl.GetUserById)
		adminGroup.PUT("/:id", userHdl.UpdateUserById)
		adminGroup.DELETE("/:id", userHdl.DeleteUserById)
	}
}
