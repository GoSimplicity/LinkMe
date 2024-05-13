//go:build wireinject

package main

import (
	"LinkMe/internal/dao"
	"LinkMe/internal/repository"
	"LinkMe/internal/service"
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
		web.NewUserHandler,
		service.NewUserService,
		repository.NewUserRepository,
		dao.NewUserDAO,
	)
	return gin.Default()
}
