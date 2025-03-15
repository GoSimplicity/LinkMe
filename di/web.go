package di

import (
	"github.com/gin-gonic/gin"
)

// InitWebServer 初始化web服务
func InitWebServer(m []gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	server.Use(m...)

	return server
}
