package ioc

import (
	"LinkMe/internal/web"
	"github.com/gin-gonic/gin"
)

func InitWebServer(userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	//server.Use(mdwl...)
	userHdl.RegisterRoutes(server)
	return server
}
