package ioc

import (
	"LinkMe/internal/web"
	"github.com/gin-gonic/gin"
)

func InitWebServer(userHdl *web.UserHandler, m []gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	server.Use(m...)
	userHdl.RegisterRoutes(server)
	return server
}
