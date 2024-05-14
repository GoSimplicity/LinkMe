//go:build wireinject

package main

import (
	"LinkMe/internal/dao"
	"LinkMe/internal/repository"
	"LinkMe/internal/service"
	ijwt "LinkMe/internal/tools/jwt"
	"LinkMe/internal/web"
	"LinkMe/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	_ "github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
		ijwt.NewJWTHandler,
		web.NewUserHandler,
		service.NewUserService,
		repository.NewUserRepository,
		dao.NewUserDAO,
	)
	return gin.Default()
}
