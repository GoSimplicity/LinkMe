//go:build wireinject

package main

import (
	"LinkMe/internal/cache"
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
		cache.NewUserCache,
		ioc.InitDB,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
		ioc.InitRedis,
		ioc.InitLogger,
		ijwt.NewJWTHandler,
		web.NewUserHandler,
		service.NewUserService,
		repository.NewUserRepository,
		dao.NewUserDAO,
	)
	return gin.Default()
}
