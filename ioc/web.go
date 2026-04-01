package ioc

import (
	"github.com/GoSimplicity/LinkMe/internal/api"
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	"github.com/gin-gonic/gin"
)

// InitWeb 初始化web服务
func InitWeb(userHdl *api.UserHandler,
	postHdl *api.PostHandler,
	historyHdl *api.HistoryHandler,
	checkHdl *api.CheckHandler,
	m []gin.HandlerFunc,
	permHdl *api.PermissionHandler,
	rankingHdl *api.RankingHandler,
	plateHdl *api.PlateHandler,
	activityHdl *api.ActivityHandler,
	commentHdl *api.CommentHandler,
	searchHdl *api.SearchHandler,
	relationHdl *api.RelationHandler,
	lotteryDrawHdl *api.LotteryDrawHandler,
	roleHdl *api.RoleHandler,
	menuHdl *api.MenuHandler,
	apiHdl *api.ApiHandler,
) *gin.Engine {
	server := gin.Default()
	server.Use(m...)
	server.GET("/healthz", func(ctx *gin.Context) {
		apiresponse.SuccessWithData(ctx, gin.H{"status": "ok"})
	})
	server.GET("/readyz", func(ctx *gin.Context) {
		apiresponse.SuccessWithData(ctx, gin.H{"status": "ready"})
	})
	userHdl.RegisterRoutes(server)
	postHdl.RegisterRoutes(server)
	historyHdl.RegisterRoutes(server)
	checkHdl.RegisterRoutes(server)
	permHdl.RegisterRoutes(server)
	rankingHdl.RegisterRoutes(server)
	plateHdl.RegisterRoutes(server)
	activityHdl.RegisterRoutes(server)
	commentHdl.RegisterRoutes(server)
	searchHdl.RegisterRoutes(server)
	relationHdl.RegisterRoutes(server)
	lotteryDrawHdl.RegisterRoutes(server)
	roleHdl.RegisterRoutes(server)
	menuHdl.RegisterRoutes(server)
	apiHdl.RegisterRoutes(server)
	return server
}
