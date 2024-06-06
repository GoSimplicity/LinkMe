package api

import (
	. "LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/service"
	. "LinkMe/pkg/ginp"
	ijwt "LinkMe/utils/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type HistoryHandler struct {
	svc service.HistoryService
	l   *zap.Logger
}

func NewHistoryHandler(svc service.HistoryService, l *zap.Logger) *HistoryHandler {
	return &HistoryHandler{
		svc: svc,
		l:   l,
	}
}

func (h *HistoryHandler) RegisterRoutes(server *gin.Engine) {
	historyGroup := server.Group("/history")
	historyGroup.GET("/list", WrapBody(h.GetHistory))
	historyGroup.DELETE("/delete", WrapBody(h.DeleteOneHistory))
	historyGroup.DELETE("/delete/all", WrapBody(h.DeleteAllHistory))
}

func (h *HistoryHandler) GetHistory(ctx *gin.Context, req ListHistoryReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	history, err := h.svc.GetHistory(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
		Uid:  uc.Uid,
	})
	if err != nil {
		h.l.Error("get history failed", zap.Error(err))
		return Result{
			Code: 500,
			Msg:  HistoryListError,
		}, err
	}
	return Result{
		Code: http.StatusOK,
		Msg:  HistoryListSuccess,
		Data: history,
	}, err
}

func (h *HistoryHandler) DeleteOneHistory(ctx *gin.Context, req DeleteHistoryReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if err := h.svc.DeleteOneHistory(ctx, req.PostId, uc.Uid); err != nil {
		h.l.Error("delete history failed", zap.Error(err))
		return Result{
			Code: 500,
			Msg:  HistoryDeleteError,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  HistoryDeleteSuccess,
	}, nil
}

func (h *HistoryHandler) DeleteAllHistory(ctx *gin.Context, req DeleteHistoryAllReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if req.IsDeleteAll == true {
		if err := h.svc.DeleteAllHistory(ctx, uc.Uid); err != nil {
			h.l.Error("delete history failed", zap.Error(err))
			return Result{
				Code: 500,
				Msg:  HistoryDeleteError,
			}, err
		}
	}
	return Result{
		Code: RequestsOK,
		Msg:  HistoryDeleteSuccess,
	}, nil
}
