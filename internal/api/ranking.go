package api

import (
	"LinkMe/internal/constants"
	"LinkMe/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RankingHandler struct {
	svc service.RankingService
	l   *zap.Logger
	biz string
}

func NewRakingHandler(svc service.RankingService, l *zap.Logger) *RankingHandler {
	return &RankingHandler{
		svc: svc,
		l:   l,
		biz: "raking",
	}
}

func (rh *RankingHandler) RegisterRoutes(server *gin.Engine) {
	postGroup := server.Group("/raking")
	postGroup.GET("/topN", rh.GetRanking)
}

func (rh *RankingHandler) GetRanking(ctx *gin.Context) {
	dp, err := rh.svc.GetTopN(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(constants.RequestsOK, gin.H{
		"data": dp,
	})
}
