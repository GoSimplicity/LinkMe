package user

import (
	"github.com/GoSimplicity/LinkMe/middleware"
	"github.com/GoSimplicity/LinkMe/utils"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, userHdl *UserHandler, ijwt utils.Handler) {
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
	authGroup.Use(middleware.NewJWTMiddleware(ijwt).CheckLogin())
	{
		authGroup.POST("/logout", userHdl.Logout)
		authGroup.GET("/profile", userHdl.GetProfile)
		authGroup.POST("/profile", userHdl.UpdateProfile)
		authGroup.POST("/password", userHdl.ChangePassword)
		authGroup.DELETE("/account", userHdl.DeleteUser)
	}

	// 管理员接口组 - 需要管理员权限
	adminGroup := r.Group("/admin/users")
	adminGroup.Use(middleware.NewJWTMiddleware(ijwt).CheckLogin())
	{
		adminGroup.GET("", userHdl.List)
		adminGroup.GET("/:id", userHdl.GetUser)
		adminGroup.POST("/:id", userHdl.UpdateUser)
		adminGroup.DELETE("/:id", userHdl.DeleteUser)
	}
}
