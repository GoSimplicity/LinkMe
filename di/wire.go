//go:build wireinject

package di

import (
	"github.com/GoSimplicity/LinkMe/internal/app/user/repository"
	"github.com/GoSimplicity/LinkMe/internal/app/user/service"
	"github.com/GoSimplicity/LinkMe/internal/core/cache"
	"github.com/GoSimplicity/LinkMe/internal/interfaces/http/user"
	"github.com/GoSimplicity/LinkMe/internal/pkg/infra/database/dao"
	ijwt "github.com/GoSimplicity/LinkMe/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Logger *zap.Logger
	Server *gin.Engine
}

var HandlerSet = wire.NewSet(
	user.NewUserHandler,
)

var ServiceSet = wire.NewSet(
	service.NewUserService,
)

var RepositorySet = wire.NewSet(
	repository.NewUserRepository,
)

var DatabaseSet = wire.NewSet(
	dao.NewUserDao,
)

var CacheSet = wire.NewSet(
	cache.NewCoreCache,
)

var UtilsSet = wire.NewSet(
	ijwt.NewJWTHandler,
)

var Injector = wire.NewSet(
	InitLogger,
	InitDB,
	InitMiddlewares,
	InitWebServer,
	InitRedis,
	wire.Struct(new(App), "*"),
)

func ProvideApp() (*App, error) {
	wire.Build(Injector, UtilsSet, ServiceSet, RepositorySet, DatabaseSet, HandlerSet)
	return &App{}, nil
}
