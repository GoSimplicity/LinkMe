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

func (h *HistoryHandler) HistoryRoutes(server *gin.Engine) {
	historyGroup := server.Group("/history")
	historyGroup.GET("/list", WrapBody(h.GetHistory))
	historyGroup.DELETE("/:id", WrapParam(h.DeleteHistory))
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

func (h *HistoryHandler) DeleteHistory(ctx *gin.Context, req DeleteHistoryReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if err := h.svc.DeleteHistory(ctx, req.ID, uc.Uid); err != nil {
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
