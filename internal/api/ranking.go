package api

import (
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	"github.com/gin-gonic/gin"
)

type RankingHandler struct {
	svc service.RankingService
}

func NewRakingHandler(svc service.RankingService) *RankingHandler {
	return &RankingHandler{
		svc: svc,
	}
}

func (rh *RankingHandler) RegisterRoutes(server *gin.Engine) {
	postGroup := server.Group("/api/raking")

	postGroup.GET("/topN", rh.GetRanking)
}

// GetRanking 获取排行榜
func (rh *RankingHandler) GetRanking(ctx *gin.Context) {
	dp, err := rh.svc.GetTopN(ctx)
	if err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	apiresponse.SuccessWithData(ctx, dp)
}
