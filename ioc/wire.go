//go:build wireinject

package ioc

import (
	"github.com/GoSimplicity/LinkMe/internal/api"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/check"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/email"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/es"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/post"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/publish"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/sms"
	"github.com/GoSimplicity/LinkMe/internal/job"
	"github.com/GoSimplicity/LinkMe/internal/mock"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/GoSimplicity/LinkMe/internal/service"
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
		InitSaramaClient,
		InitConsumers,
		InitSyncProducer,
		// InitializeSnowflakeNode,
		InitCasbin,
		InitSms,
		InitES,
		InitAsynqServer,
		InitAsynqClient,
		InitScheduler,
		InitRankingService,
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
		api.NewRoleHandler,
		api.NewMenuHandler,
		api.NewApiHandler,
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
		service.NewRoleService,
		service.NewMenuService,
		service.NewApiService,
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
		repository.NewRoleRepository,
		repository.NewMenuRepository,
		repository.NewApiRepository,
		cache.NewRankingLocalCache,
		cache.NewRankingRedisCache,
		cache.NewUserCache,
		cache.NewHistoryCache,
		cache.NewSMSCache,
		cache.NewEmailCache,
		cache.NewRelationCache,
		cache.NewPostCache,
		dao.NewUserDAO,
		dao.NewPostDAO,
		dao.NewInteractiveDAO,
		dao.NewCheckDAO,
		dao.NewSmsDAO,
		dao.NewPermissionDAO,
		dao.NewPlateDAO,
		dao.NewActivityDAO,
		dao.NewCommentDAO,
		dao.NewSearchDAO,
		dao.NewRelationDAO,
		dao.NewLotteryDrawDAO,
		dao.NewRoleDAO,
		dao.NewMenuDAO,
		dao.NewApiDAO,
		post.NewSaramaSyncProducer,
		post.NewEventConsumer,
		post.NewPostDeadLetterConsumer,
		sms.NewSMSConsumer,
		sms.NewSaramaSyncProducer,
		email.NewEmailConsumer,
		email.NewSaramaSyncProducer,
		publish.NewPublishPostEventConsumer,
		publish.NewSaramaSyncProducer,
		publish.NewPublishDeadLetterConsumer,
		check.NewCheckEventConsumer,
		check.NewSaramaCheckProducer,
		check.NewCheckDeadLetterConsumer,
		es.NewEsConsumer,
		mock.NewMockUserRepository,
		job.NewRoutes,
		job.NewRefreshCacheTask,
		job.NewTimedTask,
		job.NewTimedScheduler,
		// limiter.NewRedisSlidingWindowLimiter,
		wire.Struct(new(Cmd), "*"),
	)
	return new(Cmd)
}
