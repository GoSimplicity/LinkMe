package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/req"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/gin-gonic/gin"
)

type CheckHandler struct {
	svc service.CheckService
}

func NewCheckHandler(svc service.CheckService) *CheckHandler {
	return &CheckHandler{
		svc: svc,
	}
}

func (ch *CheckHandler) RegisterRoutes(server *gin.Engine) {
	checkGroup := server.Group("/api/checks")

	checkGroup.POST("/approve", ch.ApproveCheck)
	checkGroup.POST("/reject", ch.RejectCheck)
	checkGroup.GET("/list", ch.ListChecks)
	checkGroup.GET("/detail", ch.CheckDetail)
}

// ApproveCheck 审核通过
func (ch *CheckHandler) ApproveCheck(ctx *gin.Context) {
	var req req.ApproveCheckReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	err := ch.svc.ApproveCheck(ctx, req.CheckID, req.Remark, uc.Uid)
	if err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	apiresponse.Success(ctx)
}

// RejectCheck 审核拒绝
func (ch *CheckHandler) RejectCheck(ctx *gin.Context) {
	var req req.RejectCheckReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	err := ch.svc.RejectCheck(ctx, req.CheckID, req.Remark, uc.Uid)
	if err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	apiresponse.Success(ctx)
}

// ListChecks 获取审核列表
func (ch *CheckHandler) ListChecks(ctx *gin.Context) {
	var req req.ListCheckReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	checks, err := ch.svc.ListChecks(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	apiresponse.SuccessWithData(ctx, checks)
}

// CheckDetail 获取审核详情
func (ch *CheckHandler) CheckDetail(ctx *gin.Context) {
	var req req.CheckDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	check, err := ch.svc.CheckDetail(ctx, req.CheckID)
	if err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	apiresponse.SuccessWithData(ctx, check)
}
