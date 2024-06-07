package api

import (
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
	checkGroup.POST("/submit", WrapBody(ch.SubmitCheck))          // 提交审核
	checkGroup.PUT("/approve", WrapBody(ch.ApproveCheck))         // 审核通过
	checkGroup.PUT("/reject", WrapBody(ch.RejectCheck))           // 审核拒绝
	checkGroup.GET("/list", WrapBody(ch.ListChecks))              // 审核列表
	checkGroup.GET("/detail/:checkId", WrapParam(ch.CheckDetail)) // 审核详情
}

func (ch *CheckHandler) SubmitCheck(ctx *gin.Context, req SubmitCheckReq) (Result, error) {
	// TODO: 实现提交审核逻辑
	return Result{}, nil
}

func (ch *CheckHandler) ApproveCheck(ctx *gin.Context, req ApproveCheckReq) (Result, error) {
	// TODO: 实现审核通过逻辑
	return Result{}, nil
}

func (ch *CheckHandler) RejectCheck(ctx *gin.Context, req RejectCheckReq) (Result, error) {
	// TODO: 实现审核拒绝逻辑
	return Result{}, nil
}

func (ch *CheckHandler) ListChecks(ctx *gin.Context, req ListCheckReq) (Result, error) {
	// TODO: 实现获取审核列表逻辑
	return Result{}, nil
}

func (ch *CheckHandler) CheckDetail(ctx *gin.Context, req CheckDetailReq) (Result, error) {
	// TODO: 实现获取审核详情逻辑
	return Result{}, nil
}
