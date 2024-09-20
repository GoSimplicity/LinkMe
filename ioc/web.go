package ioc

import (
	"github.com/GoSimplicity/LinkMe/internal/api"
	"github.com/gin-gonic/gin"
)

// InitWebServer 初始化web服务
func InitWebServer(userHdl *api.UserHandler, postHdl *api.PostHandler, historyHdl *api.HistoryHandler, checkHdl *api.CheckHandler, m []gin.HandlerFunc, permHdl *api.PermissionHandler, rankingHdl *api.RankingHandler, plateHdl *api.PlateHandler, activityHdl *api.ActivityHandler, commentHdl *api.CommentHandler, searchHdl *api.SearchHandler, relationHdl *api.RelationHandler, lotteryDrawHdl *api.LotteryDrawHandler) *gin.Engine {
	server := gin.Default()
	server.Use(m...)
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
	return server
}
