//go:build wireinject

package ioc

import (
	"github.com/GoSimplicity/LinkMe/internal/api"
	cache2 "github.com/GoSimplicity/LinkMe/internal/domain/events/cache"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/check"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/email"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/es"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/post"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/publish"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/sms"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/sync"
	"github.com/GoSimplicity/LinkMe/internal/mock"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/pkg/cachep/bloom"
	"github.com/GoSimplicity/LinkMe/pkg/cachep/local"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/google/wire"
	_ "github.com/google/wire"
)

func InitWebServer() *Cmd {
	wire.Build(
		InitDB,
		InitWeb,
		InitMiddlewares,
		InitRedis,
		InitLogger,
		InitMongoDB,
		InitSaramaClient,
		InitConsumers,
		InitSyncProducer,
		InitializeSnowflakeNode,
		InitCasbin,
		InitSms,
		InitRanking,
		InitES,
		ijwt.NewJWTHandler,
		api.NewUserHandler,
		api.NewPostHandler,
		api.NewHistoryHandler,
		api.NewCheckHandler,
		api.NewPermissionHandler,
		api.NewRakingHandler,
		api.NewPlateHandler,
		api.NewActivityHandler,
		api.NewCommentHandler,
		api.NewSearchHandler,
		api.NewRelationHandler,
		api.NewLotteryDrawHandler,
		service.NewUserService,
		service.NewPostService,
		service.NewInteractiveService,
		service.NewHistoryService,
		service.NewCheckService,
		service.NewPermissionService,
		service.NewRankingService,
		service.NewPlateService,
		service.NewActivityService,
		service.NewCommentService,
		service.NewSearchService,
		service.NewRelationService,
		service.NewLotteryDrawService,
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
		repository.NewCommentRepository,
		repository.NewSearchRepository,
		repository.NewRelationRepository,
		repository.NewLotteryDrawRepository,
		cache.NewRankingLocalCache,
		cache.NewRankingRedisCache,
		cache.NewUserCache,
		cache.NewInteractiveCache,
		cache.NewHistoryCache,
		cache.NewSMSCache,
		cache.NewEmailCache,
		cache.NewRelationCache,
		cache.NewCheckCache,
		dao.NewUserDAO,
		dao.NewPostDAO,
		dao.NewInteractiveDAO,
		dao.NewCheckDAO,
		dao.NewSmsDAO,
		dao.NewPermissionDAO,
		dao.NewPlateDAO,
		dao.NewActivityDAO,
		dao.NewCommentService,
		dao.NewSearchDAO,
		dao.NewRelationDAO,
		dao.NewLotteryDrawDAO,
		post.NewSaramaSyncProducer,
		post.NewReadEventConsumer,
		sms.NewSMSConsumer,
		sms.NewSaramaSyncProducer,
		email.NewEmailConsumer,
		email.NewSaramaSyncProducer,
		bloom.NewCacheBloom,
		local.NewLocalCacheManager,
		sync.NewSyncConsumer,
		cache2.NewCacheConsumer,
		publish.NewPublishPostEventConsumer,
		publish.NewSaramaSyncProducer,
		check.NewCheckConsumer,
		es.NewEsConsumer,
		mock.NewMockUserRepository,
		// limiter.NewRedisSlidingWindowLimiter,
		wire.Struct(new(Cmd), "*"),
	)
	return new(Cmd)
}
