package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/req"
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	. "github.com/GoSimplicity/LinkMe/pkg/ginp"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/gin-gonic/gin"
)

// LotteryDrawHandler 负责处理抽奖和秒杀相关的请求
type LotteryDrawHandler struct {
	svc service.LotteryDrawService
}

func NewLotteryDrawHandler(svc service.LotteryDrawService) *LotteryDrawHandler {
	return &LotteryDrawHandler{
		svc: svc,
	}
}

// RegisterRoutes 注册抽奖和秒杀的路由
func (lh *LotteryDrawHandler) RegisterRoutes(server *gin.Engine) {
	apiGroup := server.Group("/api")
	{
		// 抽奖相关的路由
		lotteryGroup := apiGroup.Group("/lottery")
		{
			lotteryGroup.POST("/list", WrapBody(lh.ListLotteryDraws))              // 获取所有抽奖活动
			lotteryGroup.POST("/create", WrapBody(lh.CreateLotteryDraw))           // 创建新的抽奖活动
			lotteryGroup.GET("/:id", WrapQuery(lh.GetLotteryDraw))                 // 获取指定ID的抽奖活动
			lotteryGroup.POST("/participate", WrapBody(lh.ParticipateLotteryDraw)) // 参与抽奖活动
		}

		// 秒杀相关的路由
		secondKillGroup := apiGroup.Group("/secondKill")
		{
			secondKillGroup.POST("/list", WrapBody(lh.ListKillEvents))               // 获取所有秒杀活动
			secondKillGroup.POST("/create", WrapBody(lh.CreateSecondKillEvent))      // 创建新的秒杀活动
			secondKillGroup.GET("/:id", WrapQuery(lh.GetSecondKillEvent))            // 获取指定ID的秒杀活动
			secondKillGroup.POST("/participate", WrapBody(lh.ParticipateSecondKill)) // 参与秒杀活动
		}
	}
}

// ListLotteryDraws 获取所有抽奖活动
func (lh *LotteryDrawHandler) ListLotteryDraws(ctx *gin.Context, req req.ListLotteryDrawsReq) (Result, error) {
	pagination := domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	}

	ld, err := lh.svc.ListLotteryDraws(ctx, req.Status, pagination)
	if err != nil {
		return Result{
			Code: ServerRequestError,
			Msg:  ListLotteryDrawsError,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  ListLotteryDrawsSuccess,
		Data: ld,
	}, nil
}

// CreateLotteryDraw 创建新的抽奖活动
func (lh *LotteryDrawHandler) CreateLotteryDraw(ctx *gin.Context, req req.CreateLotteryDrawReq) (Result, error) {
	input := domain.LotteryDraw{
		Name:        req.Name,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	err := lh.svc.CreateLotteryDraw(ctx, domain.LotteryDraw{
		Name:        input.Name,
		Description: input.Description,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Status:      domain.LotteryStatusPending,
	})
	if err != nil {
		return Result{
			Code: ServerRequestError,
			Msg:  CreateLotteryDrawError,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  CreateLotteryDrawSuccess,
	}, nil
}

// GetLotteryDraw 获取指定ID的抽奖活动
func (lh *LotteryDrawHandler) GetLotteryDraw(ctx *gin.Context, req req.GetLotteryDrawReq) (Result, error) {
	ld, err := lh.svc.GetLotteryDrawByID(ctx, req.ID)
	if err != nil {
		return Result{
			Code: ServerRequestError,
			Msg:  GetLotteryDrawError,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  GetLotteryDrawSuccess,
		Data: ld,
	}, nil
}

// ParticipateLotteryDraw 参与抽奖活动
func (lh *LotteryDrawHandler) ParticipateLotteryDraw(ctx *gin.Context, req req.ParticipateLotteryDrawReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	err := lh.svc.ParticipateLotteryDraw(ctx, req.ActivityID, uc.Uid)
	if err != nil {
		return Result{
			Code: ServerRequestError,
			Msg:  ParticipateLotteryDrawError,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  ParticipateLotteryDrawSuccess,
	}, nil
}

// ListKillEvents 获取所有秒杀活动
func (lh *LotteryDrawHandler) ListKillEvents(ctx *gin.Context, req req.GetAllSecondKillEventsReq) (Result, error) {
	pagination := domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	}

	ke, err := lh.svc.ListSecondKillEvents(ctx, req.Status, pagination)
	if err != nil {
		return Result{
			Code: ServerRequestError,
			Msg:  ListSecondKillEventsError,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  ListSecondKillEventsSuccess,
		Data: ke,
	}, nil
}

// CreateSecondKillEvent 创建新的秒杀活动
func (lh *LotteryDrawHandler) CreateSecondKillEvent(ctx *gin.Context, req req.CreateSecondKillEventReq) (Result, error) {
	input := domain.SecondKillEvent{
		Name:        req.Name,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	err := lh.svc.CreateSecondKillEvent(ctx, input)
	if err != nil {
		return Result{
			Code: ServerRequestError,
			Msg:  CreateSecondKillEventError,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  CreateSecondKillEventSuccess,
	}, nil
}

// GetSecondKillEvent 获取指定ID的秒杀活动
func (lh *LotteryDrawHandler) GetSecondKillEvent(ctx *gin.Context, req req.GetSecondKillEventReq) (Result, error) {
	ke, err := lh.svc.GetSecondKillEventByID(ctx, req.ID)
	if err != nil {
		return Result{
			Code: ServerRequestError,
			Msg:  GetSecondKillEventError,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  GetSecondKillEventSuccess,
		Data: ke,
	}, nil
}

// ParticipateSecondKill 参与秒杀活动
func (lh *LotteryDrawHandler) ParticipateSecondKill(ctx *gin.Context, req req.ParticipateSecondKillReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	err := lh.svc.ParticipateSecondKill(ctx, req.ActivityID, uc.Uid)
	if err != nil {
		return Result{
			Code: ServerRequestError,
			Msg:  ParticipateSecondKillError,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  ParticipateSecondKillSuccess,
	}, nil
}
