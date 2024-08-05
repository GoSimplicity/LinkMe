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

type PlateHandler struct {
	svc service.PlateService
	ce  *casbin.Enforcer
}

func NewPlateHandler(svc service.PlateService, ce *casbin.Enforcer) *PlateHandler {
	return &PlateHandler{
		svc: svc,
		ce:  ce,
	}
}

func (h *PlateHandler) RegisterRoutes(server *gin.Engine) {
	casbinMiddleware := middleware.NewCasbinMiddleware(h.ce)
	permissionGroup := server.Group("/api/plate")
	permissionGroup.Use(casbinMiddleware.CheckCasbin())
	permissionGroup.POST("/create", WrapBody(h.CreatePlate))
	permissionGroup.POST("/update", WrapBody(h.UpdatePlate))
	permissionGroup.DELETE("/delete/:plateId", WrapParam(h.DeletePlate))
	permissionGroup.POST("/list", WrapBody(h.ListPlate))
}

func (h *PlateHandler) CreatePlate(ctx *gin.Context, req req.CreatePlateReq) (Result, error) {
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

func (h *PlateHandler) UpdatePlate(ctx *gin.Context, req req.UpdatePlateReq) (Result, error) {
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

func (h *PlateHandler) DeletePlate(ctx *gin.Context, req req.DeletePlateReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	err := h.svc.DeletePlate(ctx, req.PlateID, uc.Uid)
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

func (h *PlateHandler) ListPlate(ctx *gin.Context, req req.ListPlateReq) (Result, error) {
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
