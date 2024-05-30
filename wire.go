//go:build wireinject

package main

import (
	"LinkMe/internal/api"
	"LinkMe/internal/domain/events/post"
	"LinkMe/internal/repository"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/service"
	"LinkMe/ioc"
	ijwt "LinkMe/utils/jwt"
	"github.com/google/wire"
	_ "github.com/google/wire"
)

func InitWebServer() *Cmd {
	wire.Build(
		ioc.InitDB,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
		ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitMongoDB,
		ioc.InitSaramaClient,
		ioc.InitConsumers,
		ioc.InitSyncProducer,
		ijwt.NewJWTHandler,
		api.NewUserHandler,
		api.NewPostHandler,
		service.NewUserService,
		service.NewPostService,
		service.NewInteractiveService,
		repository.NewUserRepository,
		repository.NewPostRepository,
		repository.NewInteractiveRepository,
		cache.NewUserCache,
		cache.NewPostCache,
		dao.NewUserDAO,
		dao.NewPostDAO,
		dao.NewInteractiveDAO,
		post.NewSaramaSyncProducer,
		post.NewInteractiveReadEventConsumer,
		wire.Struct(new(Cmd), "*"),
	)
	return new(Cmd)
}
