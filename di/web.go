package di

import (
	"github.com/GoSimplicity/LinkMe/internal/app/user/service"
	"github.com/GoSimplicity/LinkMe/internal/interfaces/http/user"
	"github.com/gin-gonic/gin"
)

// InitWebServer 初始化web服务
func InitWebServer(m []gin.HandlerFunc, userSvc service.UserService) *gin.Engine {
	server := gin.Default()
	server.Use(m...)

	user.RegisterRoutes(server, userSvc)

	return server
}
