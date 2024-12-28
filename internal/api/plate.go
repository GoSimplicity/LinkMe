package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/req"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/middleware"
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
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
	permissionGroup.POST("/create", h.CreatePlate)
	permissionGroup.POST("/update", h.UpdatePlate)
	permissionGroup.DELETE("/delete/:plateId", h.DeletePlate)
	permissionGroup.POST("/list", h.ListPlate)
}

func (h *PlateHandler) CreatePlate(ctx *gin.Context) {
	var req req.CreatePlateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if err := h.svc.CreatePlate(ctx, domain.Plate{
		Name:        req.Name,
		Description: req.Description,
		Uid:         uc.Uid,
	}); err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}
	apiresponse.Success(ctx)
}

func (h *PlateHandler) UpdatePlate(ctx *gin.Context) {
	var req req.UpdatePlateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if err := h.svc.UpdatePlate(ctx, domain.Plate{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Uid:         uc.Uid,
	}); err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}
	apiresponse.Success(ctx)
}

func (h *PlateHandler) DeletePlate(ctx *gin.Context) {
	var req req.DeletePlateReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if err := h.svc.DeletePlate(ctx, req.PlateID, uc.Uid); err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}
	apiresponse.Success(ctx)
}

func (h *PlateHandler) ListPlate(ctx *gin.Context) {
	var req req.ListPlateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	plates, err := h.svc.ListPlate(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}
	apiresponse.SuccessWithData(ctx, plates)
}
