package api

import (
	. "LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/service"
	"LinkMe/middleware"
	. "LinkMe/pkg/ginp"
	ijwt "LinkMe/utils/jwt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CheckHandler struct {
	svc service.CheckService
	l   *zap.Logger
	ce  *casbin.Enforcer
	biz string
}

func NewCheckHandler(svc service.CheckService, l *zap.Logger, ce *casbin.Enforcer) *CheckHandler {
	return &CheckHandler{
		svc: svc,
		l:   l,
		ce:  ce,
		biz: "check",
	}
}

func (ch *CheckHandler) RegisterRoutes(server *gin.Engine) {
	casbinMiddleware := middleware.NewCasbinMiddleware(ch.ce, ch.l)
	checkGroup := server.Group("/checks")
	checkGroup.Use(casbinMiddleware.CheckCasbin())
	checkGroup.POST("/approve", WrapBody(ch.ApproveCheck)) // 审核通过
	checkGroup.POST("/reject", WrapBody(ch.RejectCheck))   // 审核拒绝
	checkGroup.GET("/list", WrapBody(ch.ListChecks))       // 审核列表
	checkGroup.GET("/detail", WrapBody(ch.CheckDetail))    // 审核详情
	checkGroup.GET("/stats", WrapQuery(ch.GetCheckCount))  // 管理员使用
}

func (ch *CheckHandler) ApproveCheck(ctx *gin.Context, req ApproveCheckReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	err := ch.svc.ApproveCheck(ctx, req.CheckID, req.Remark, uc.Uid)
	if err != nil {
		return Result{
			Code: RequestsERROR,
			Msg:  "failed to approve check",
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  "success to approve check",
	}, nil
}

func (ch *CheckHandler) RejectCheck(ctx *gin.Context, req RejectCheckReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	err := ch.svc.RejectCheck(ctx, req.CheckID, req.Remark, uc.Uid)
	if err != nil {
		return Result{
			Code: RequestsERROR,
			Msg:  "failed to reject check",
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  "success to reject check",
	}, nil
}

func (ch *CheckHandler) ListChecks(ctx *gin.Context, req ListCheckReq) (Result, error) {
	checks, err := ch.svc.ListChecks(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return Result{
			Code: RequestsERROR,
			Msg:  "failed to list checks",
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  "success to list checks",
		Data: checks,
	}, nil
}

func (ch *CheckHandler) CheckDetail(ctx *gin.Context, req CheckDetailReq) (Result, error) {
	check, err := ch.svc.CheckDetail(ctx, req.CheckID)
	if err != nil {
		return Result{
			Code: RequestsERROR,
			Msg:  "failed to get check detail",
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  "success to get check detail",
		Data: check,
	}, nil
}

func (ch *CheckHandler) GetCheckCount(ctx *gin.Context, _ GetCheckCount) (Result, error) {
	count, err := ch.svc.GetCheckCount(ctx)
	if err != nil {
		return Result{
			Code: GetCheckERRORCode,
			Msg:  GetCheckERROR,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  GetCheckSuccess,
		Data: count,
	}, nil
}
