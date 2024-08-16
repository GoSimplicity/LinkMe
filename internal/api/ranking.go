package api

import (
	"github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/gin-gonic/gin"
)

type RankingHandler struct {
	svc service.RankingService
	biz string
}

func NewRakingHandler(svc service.RankingService) *RankingHandler {
	return &RankingHandler{
		svc: svc,
		biz: "raking",
	}
}

func (rh *RankingHandler) RegisterRoutes(server *gin.Engine) {
	postGroup := server.Group("/api/raking")
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
