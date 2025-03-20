//go:build wireinject

package di

import (
	"github.com/GoSimplicity/LinkMe/internal/app/user/service"
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

var ServiceSet = wire.NewSet(
	service.NewUserService,
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
	wire.Build(Injector, UtilsSet, ServiceSet)
	return &App{}, nil
}
