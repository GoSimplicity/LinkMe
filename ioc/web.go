package ioc

import (
	"LinkMe/internal/api"
	"github.com/gin-gonic/gin"
)

// InitWebServer 初始化web服务
func InitWebServer(userHdl *api.UserHandler, m []gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	server.Use(m...)
	userHdl.RegisterRoutes(server)
	return server
}
