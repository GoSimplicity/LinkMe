package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/req"
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	svc service.HistoryService
}

func NewHistoryHandler(svc service.HistoryService) *HistoryHandler {
	return &HistoryHandler{
		svc: svc,
	}
}

func (h *HistoryHandler) RegisterRoutes(server *gin.Engine) {
	historyGroup := server.Group("/api/history")

	historyGroup.POST("/list", h.GetHistory)
	historyGroup.DELETE("/delete", h.DeleteOneHistory)
	historyGroup.DELETE("/delete/all", h.DeleteAllHistory)
}

// GetHistory 获取历史记录
func (h *HistoryHandler) GetHistory(ctx *gin.Context) {
	var req req.ListHistoryReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)
	history, err := h.svc.GetHistory(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
		Uid:  uc.Uid,
	})
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, "获取历史记录失败")
		return
	}

	apiresponse.SuccessWithData(ctx, history)
}

// DeleteOneHistory 删除一条历史记录
func (h *HistoryHandler) DeleteOneHistory(ctx *gin.Context) {
	var req req.DeleteHistoryReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if err := h.svc.DeleteOneHistory(ctx, req.PostId, uc.Uid); err != nil {
		apiresponse.ErrorWithMessage(ctx, HistoryDeleteError)
		return
	}

	apiresponse.Success(ctx)
}

// DeleteAllHistory 删除所有历史记录
func (h *HistoryHandler) DeleteAllHistory(ctx *gin.Context) {
	var req req.DeleteHistoryAllReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithData(ctx, err)
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if req.IsDeleteAll {
		if err := h.svc.DeleteAllHistory(ctx, uc.Uid); err != nil {
			apiresponse.ErrorWithMessage(ctx, HistoryDeleteError)
			return
		}
	}

	apiresponse.Success(ctx)
}
