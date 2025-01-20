package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/req"
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/middleware"
	. "github.com/GoSimplicity/LinkMe/pkg/ginp"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type CheckHandler struct {
	svc service.CheckService
	ce  *casbin.Enforcer
	biz string
}

func NewCheckHandler(svc service.CheckService, ce *casbin.Enforcer) *CheckHandler {
	return &CheckHandler{
		svc: svc,
		ce:  ce,
		biz: "check",
	}
}

func (ch *CheckHandler) RegisterRoutes(server *gin.Engine) {
	casbinMiddleware := middleware.NewCasbinMiddleware(ch.ce)
	checkGroup := server.Group("/api/checks")
	checkGroup.Use(casbinMiddleware.CheckCasbin())
	checkGroup.POST("/approve", WrapBody(ch.ApproveCheck))
	checkGroup.POST("/reject", WrapBody(ch.RejectCheck))
	checkGroup.POST("/list", WrapBody(ch.ListChecks))
	checkGroup.GET("/detail", WrapBody(ch.CheckDetail))
}

// ApproveCheck 审核通过
func (ch *CheckHandler) ApproveCheck(ctx *gin.Context, req req.ApproveCheckReq) (Result, error) {
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

// RejectCheck 审核拒绝
func (ch *CheckHandler) RejectCheck(ctx *gin.Context, req req.RejectCheckReq) (Result, error) {
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

// ListChecks 获取审核列表
func (ch *CheckHandler) ListChecks(ctx *gin.Context, req req.ListCheckReq) (Result, error) {
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

// CheckDetail 获取审核详情
func (ch *CheckHandler) CheckDetail(ctx *gin.Context, req req.CheckDetailReq) (Result, error) {
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
