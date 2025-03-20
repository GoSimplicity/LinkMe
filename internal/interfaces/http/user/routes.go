package user

import (
	"github.com/GoSimplicity/LinkMe/internal/app/user/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, userSvc service.UserService) {
	h := NewHandler(userSvc)
	r := server.Group("/api/v1")
	// 公共接口组 - 无需认证
	publicGroup := r.Group("/users")
	{
		// 注册登录相关
		publicGroup.POST("/register", h.Register)
		publicGroup.POST("/login", h.Login)
		publicGroup.POST("/login-sms", h.LoginSMS)
		publicGroup.POST("/send-sms", h.SendSMS)
		publicGroup.POST("/send-email", h.SendEmail)
		publicGroup.POST("/refresh-token", h.RefreshToken)
	}

	// 需要认证的接口组 - 用户自身操作
	authGroup := r.Group("/users")
	// TODO: 添加中间件进行认证
	// authGroup.Use(middleware.JWTMiddleware())
	{
		authGroup.POST("/logout", h.Logout)
		authGroup.GET("/profile", h.GetProfile)
		authGroup.PUT("/profile", h.UpdateProfile)
		authGroup.PUT("/password", h.ChangePassword)
		authGroup.DELETE("/account", h.DeleteUser)
	}

	// 管理员接口组 - 需要管理员权限
	adminGroup := r.Group("/admin/users")
	// TODO: 添加管理员权限验证中间件
	// adminGroup.Use(middleware.AdminAuth())
	{
		adminGroup.GET("", h.List)
		adminGroup.GET("/:id", h.GetUserById)
		adminGroup.PUT("/:id", h.UpdateUserById)
		adminGroup.DELETE("/:id", h.DeleteUserById)
	}
}
