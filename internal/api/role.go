package api

import (
	"LinkMe/internal/service"
	"LinkMe/middleware"
	. "LinkMe/pkg/ginp"
	"LinkMe/utils/jwt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type PermissionHandler struct {
	svc service.PermissionService
	l   *zap.Logger
	ce  *casbin.Enforcer
}

func NewPermissionHandler(svc service.PermissionService, l *zap.Logger, ce *casbin.Enforcer) *PermissionHandler {
	return &PermissionHandler{
		svc: svc,
		l:   l,
		ce:  ce,
	}
}

func (h *PermissionHandler) RegisterRoutes(server *gin.Engine) {
	casbinMiddleware := middleware.NewCasbinMiddleware(h.ce, h.l)
	permissionGroup := server.Group("/permissions")
	permissionGroup.Use(casbinMiddleware.CheckCasbin())
	permissionGroup.GET("/list", WrapBody(h.GetPermissions))        // 获取权限列表
	permissionGroup.POST("/assign", WrapBody(h.AssignPermission))   // 分配权限
	permissionGroup.DELETE("/remove", WrapBody(h.RemovePermission)) // 移除权限
}

// GetPermissions 处理获取权限列表的请求
func (h *PermissionHandler) GetPermissions(ctx *gin.Context, req ListPermissionsReq) (Result, error) {
	uc := ctx.MustGet("user").(jwt.UserClaims)
	if uc.Uid != req.UserID {
		return Result{
			Code: http.StatusBadRequest,
			Msg:  "get permissions failed",
		}, nil
	}
	permissions, err := h.svc.GetPermissions(ctx, uc.Uid)
	if err != nil {
		return Result{
			Code: http.StatusInternalServerError,
			Msg:  "get permissions failed",
		}, err
	}
	return Result{
		Code: http.StatusOK,
		Msg:  "get permissions success",
		Data: permissions,
	}, nil
}

// AssignPermission 处理分配权限的请求
func (h *PermissionHandler) AssignPermission(ctx *gin.Context, req AssignPermissionReq) (Result, error) {
	if err := h.svc.AssignPermission(ctx, req.UserID, req.Path, req.Method); err != nil {
		return Result{
			Code: http.StatusInternalServerError,
			Msg:  "assign permission failed",
		}, err
	}
	return Result{
		Code: http.StatusOK,
		Msg:  "assign permission success",
	}, nil
}

// RemovePermission 处理移除权限的请求
func (h *PermissionHandler) RemovePermission(ctx *gin.Context, req RemovePermissionReq) (Result, error) {
	if err := h.svc.RemovePermission(ctx, req.UserID, req.Path, req.Method); err != nil {
		return Result{
			Code: http.StatusInternalServerError,
			Msg:  "remove permission failed",
		}, err
	}
	return Result{
		Code: http.StatusOK,
		Msg:  "remove permission success",
	}, nil
}
