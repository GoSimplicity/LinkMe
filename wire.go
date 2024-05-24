//go:build wireinject

package main

import (
	"LinkMe/internal/api"
	"LinkMe/internal/repository"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/service"
	"LinkMe/ioc"
	ijwt "LinkMe/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	_ "github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
		ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitMongoDB,
		ioc.InitSaramaClient,
		ijwt.NewJWTHandler,
		api.NewUserHandler,
		api.NewPostHandler,
		service.NewUserService,
		service.NewPostService,
		repository.NewUserRepository,
		repository.NewPostRepository,
		cache.NewUserCache,
		cache.NewPostCache,
		dao.NewUserDAO,
		dao.NewPostDAO,
	)
	return gin.Default()
}
