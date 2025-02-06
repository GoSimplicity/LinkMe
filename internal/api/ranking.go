package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/req"
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	. "github.com/GoSimplicity/LinkMe/pkg/ginp"
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
	postGroup.GET("/config", rh.GetRankingConfig)
	postGroup.POST("/reset", WrapBody(rh.ResetRanking))
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

func (rh *RankingHandler) GetRankingConfig(ctx *gin.Context) {
	dp, err := rh.svc.GetRankingConfig(ctx)
	if err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}
	apiresponse.SuccessWithData(ctx, dp)
}

func (rh *RankingHandler) ResetRanking(ctx *gin.Context, req req.RankingParameterReq) (Result, error) {
	rankingParameters := domain.RankingParameter{
		Alpha:  req.Alpha,
		Beta:   req.Beta,
		Gamma:  req.Gamma,
		Lambda: req.Lambda,
	}
	err := rh.svc.ResetRankingConfig(ctx, rankingParameters)
	if err != nil {
		return Result{
			Code: SetRankingErrorCode,
			Msg:  SetRankingErrorMsg,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  SetRankingSuccessMsg,
	}, nil
}
