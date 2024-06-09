package api

import (
	. "LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/service"
	. "LinkMe/pkg/ginp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CheckHandler struct {
	svc service.CheckService
	l   *zap.Logger
	biz string
}

func NewCheckHandler(svc service.CheckService, l *zap.Logger) *CheckHandler {
	return &CheckHandler{
		svc: svc,
		l:   l,
		biz: "check",
	}
}

func (ch *CheckHandler) RegisterRoutes(server *gin.Engine) {
	checkGroup := server.Group("/checks")
	//checkGroup.POST("/submit", WrapBody(ch.SubmitCheck))        // 提交审核
	checkGroup.POST("/approve", WrapBody(ch.ApproveCheck))        // 审核通过
	checkGroup.POST("/reject", WrapBody(ch.RejectCheck))          // 审核拒绝
	checkGroup.GET("/list", WrapBody(ch.ListChecks))              // 审核列表
	checkGroup.GET("/detail/:checkId", WrapParam(ch.CheckDetail)) // 审核详情
}

//func (ch *CheckHandler) SubmitCheck(ctx *gin.Context, req SubmitCheckReq) (Result, error) {
//	uc := ctx.MustGet("user").(ijwt.UserClaims)
//	check, err := ch.svc.SubmitCheck(ctx, domain.Check{
//		PostID:  req.PostID,
//		Content: req.Content,
//		Title:   req.Title,
//		UserID:  uc.Uid,
//	})
//	if err != nil {
//		ch.l.Error("failed to submit check", zap.Error(err))
//		return Result{
//			Code: RequestsERROR,
//			Msg:  "failed to submit check",
//		}, err
//	}
//	return Result{
//		Code: RequestsOK,
//		Msg:  "success to submit check",
//		Data: check,
//	}, nil
//}

func (ch *CheckHandler) ApproveCheck(ctx *gin.Context, req ApproveCheckReq) (Result, error) {
	err := ch.svc.ApproveCheck(ctx, req.CheckID, req.Remark)
	if err != nil {
		ch.l.Error("failed to approve check", zap.Error(err))
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
	err := ch.svc.RejectCheck(ctx, req.CheckID, req.Remark)
	if err != nil {
		ch.l.Error("failed to reject check", zap.Error(err))
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
		ch.l.Error("failed to list checks", zap.Error(err))
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
		ch.l.Error("failed to get check detail", zap.Error(err))
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
