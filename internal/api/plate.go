package api

import (
	. "LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/service"
	. "LinkMe/pkg/ginp"
	ijwt "LinkMe/utils/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PlateHandler struct {
	svc service.PlateService
	l   *zap.Logger
}

func NewPlateHandler(svc service.PlateService, l *zap.Logger) *PlateHandler {
	return &PlateHandler{
		svc: svc,
		l:   l,
	}
}

func (h *PlateHandler) RegisterRoutes(server *gin.Engine) {
	permissionGroup := server.Group("/plate")
	permissionGroup.POST("/create", WrapBody(h.CreatePlate))
	permissionGroup.POST("/update", WrapBody(h.UpdatePlate))
	permissionGroup.DELETE("/delete", WrapBody(h.DeletePlate))
	permissionGroup.DELETE("/list", WrapBody(h.ListPlate))
}

func (h *PlateHandler) CreatePlate(ctx *gin.Context, req CreatePlateReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if err := h.svc.CreatePlate(ctx, domain.Plate{
		Name:        req.Name,
		Description: req.Description,
		Uid:         uc.Uid,
	}); err != nil {
		return Result{
			Code: RequestsERROR,
			Msg:  PlateCreateError,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  PlateCreateSuccess,
	}, nil
}

func (h *PlateHandler) UpdatePlate(ctx *gin.Context, req CreatePlateReq) (Result, error) {
	return Result{}, nil
}

func (h *PlateHandler) DeletePlate(ctx *gin.Context, req CreatePlateReq) (Result, error) {
	return Result{}, nil
}

func (h *PlateHandler) ListPlate(ctx *gin.Context, req CreatePlateReq) (Result, error) {
	return Result{}, nil
}
