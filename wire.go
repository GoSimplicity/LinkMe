//go:build wireinject

package main

import (
	"github.com/GoSimplicity/LinkMe/internal/api"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/email"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/post"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/sms"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/ioc"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"

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
		api.NewPlateHandler,
		api.NewActivityHandler,
		service.NewUserService,
		service.NewPostService,
		service.NewInteractiveService,
		service.NewHistoryService,
		service.NewCheckService,
		service.NewPermissionService,
		service.NewRankingService,
		service.NewPlateService,
		service.NewActivityService,
		repository.NewUserRepository,
		repository.NewPostRepository,
		repository.NewInteractiveRepository,
		repository.NewHistoryRepository,
		repository.NewCheckRepository,
		repository.NewSmsRepository,
		repository.NewPermissionRepository,
		repository.NewRankingCache,
		repository.NewEmailRepository,
		repository.NewPlateRepository,
		repository.NewActivityRepository,
		cache.NewRankingLocalCache,
		cache.NewRankingRedisCache,
		cache.NewUserCache,
		cache.NewPostCache,
		cache.NewInteractiveCache,
		cache.NewHistoryCache,
		cache.NewSMSCache,
		cache.NewEmailCache,
		dao.NewUserDAO,
		dao.NewPostDAO,
		dao.NewInteractiveDAO,
		dao.NewCheckDAO,
		dao.NewSmsDAO,
		dao.NewPermissionDAO,
		dao.NewPlateDAO,
		dao.NewActivityDAO,
		post.NewSaramaSyncProducer,
		post.NewInteractiveReadEventConsumer,
		sms.NewSMSConsumer,
		sms.NewSaramaSyncProducer,
		email.NewEmailConsumer,
		email.NewSaramaSyncProducer,
		// limiter.NewRedisSlidingWindowLimiter,
		wire.Struct(new(Cmd), "*"),
	)
	return new(Cmd)
}
