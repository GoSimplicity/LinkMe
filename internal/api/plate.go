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

type PlateHandler struct {
	svc service.PlateService
	l   *zap.Logger
	ce  *casbin.Enforcer
}

func NewPlateHandler(svc service.PlateService, l *zap.Logger, ce *casbin.Enforcer) *PlateHandler {
	return &PlateHandler{
		svc: svc,
		l:   l,
		ce:  ce,
	}
}

func (h *PlateHandler) RegisterRoutes(server *gin.Engine) {
	casbinMiddleware := middleware.NewCasbinMiddleware(h.ce, h.l)
	permissionGroup := server.Group("/api/plate")
	permissionGroup.Use(casbinMiddleware.CheckCasbin())
	permissionGroup.POST("/create", WrapBody(h.CreatePlate))
	permissionGroup.PUT("/update", WrapBody(h.UpdatePlate))
	permissionGroup.DELETE("/delete", WrapBody(h.DeletePlate))
	permissionGroup.GET("/list", WrapBody(h.ListPlate))
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

func (h *PlateHandler) UpdatePlate(ctx *gin.Context, req UpdatePlateReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if err := h.svc.UpdatePlate(ctx, domain.Plate{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Uid:         uc.Uid,
	}); err != nil {
		return Result{
			Code: RequestsERROR,
			Msg:  PlateUpdateError,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  PlateUpdateSuccess,
	}, nil
}

func (h *PlateHandler) DeletePlate(ctx *gin.Context, req DeletePlateReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	err := h.svc.DeletePlate(ctx, req.ID, uc.Uid)
	if err != nil {
		return Result{
			Code: RequestsERROR,
			Msg:  PlateDeleteError,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  PlateDeleteSuccess,
	}, nil
}

func (h *PlateHandler) ListPlate(ctx *gin.Context, req ListPlateReq) (Result, error) {
	plates, err := h.svc.ListPlate(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return Result{
			Code: RequestsERROR,
			Msg:  PlateListError,
			Data: "",
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  PlateListSuccess,
		Data: plates,
	}, nil
}
