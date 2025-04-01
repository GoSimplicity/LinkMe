package di

import (
	"github.com/GoSimplicity/LinkMe/internal/interfaces/http/user"
	"github.com/gin-gonic/gin"
)

// InitWebServer 初始化web服务
func InitWebServer(m []gin.HandlerFunc, userHdl *user.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(m...)

	user.RegisterRoutes(server, userHdl)

	return server
}
