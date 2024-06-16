//go:build wireinject

package main

import (
	"LinkMe/internal/api"
	"LinkMe/internal/domain/events/post"
	"LinkMe/internal/domain/events/sms"
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
		ioc.InitializeSnowflakeNode,
		ioc.InitCasbin,
		ioc.InitSms,
		ioc.InitRanking,
		ijwt.NewJWTHandler,
		api.NewUserHandler,
		api.NewPostHandler,
		api.NewHistoryHandler,
		api.NewCheckHandler,
		api.NewPermissionHandler,
		api.NewRakingHandler,
		service.NewUserService,
		service.NewPostService,
		service.NewInteractiveService,
		service.NewHistoryService,
		service.NewCheckService,
		service.NewPermissionService,
		service.NewSmsService,
		service.NewRankingService,
		repository.NewUserRepository,
		repository.NewPostRepository,
		repository.NewInteractiveRepository,
		repository.NewHistoryRepository,
		repository.NewCheckRepository,
		repository.NewSmsRepository,
		repository.NewPermissionRepository,
		repository.NewRankingCache,
		cache.NewRankingLocalCache,
		cache.NewRankingRedisCache,
		cache.NewUserCache,
		cache.NewPostCache,
		cache.NewInteractiveCache,
		cache.NewHistoryCache,
		cache.NewSMSCache,
		dao.NewUserDAO,
		dao.NewPostDAO,
		dao.NewInteractiveDAO,
		dao.NewCheckDAO,
		dao.NewSmsDAO,
		dao.NewPermissionDAO,
		post.NewSaramaSyncProducer,
		post.NewInteractiveReadEventConsumer,
		sms.NewSMSConsumer,
		sms.NewSaramaSyncProducer,
		wire.Struct(new(Cmd), "*"),
	)
	return new(Cmd)
}
