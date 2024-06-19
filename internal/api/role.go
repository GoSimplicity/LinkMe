package api

import (
	"LinkMe/internal/service"
	"LinkMe/middleware"
	"LinkMe/pkg/ginp"
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
	permissionGroup.GET("/list", ginp.WrapBody(h.GetPermissions))        // 获取权限列表
	permissionGroup.POST("/assign", ginp.WrapBody(h.AssignPermission))   // 分配权限
	permissionGroup.DELETE("/remove", ginp.WrapBody(h.RemovePermission)) // 移除权限
}

// GetPermissions 处理获取权限列表的请求
func (h *PermissionHandler) GetPermissions(ctx *gin.Context, req ListPermissionsReq) (ginp.Result, error) {
	uc := ctx.MustGet("user").(jwt.UserClaims)
	if uc.Uid != req.UserID {
		return ginp.Result{
			Code: http.StatusBadRequest,
			Msg:  "传入的userId不对",
		}, nil
	}
	permissions, err := h.svc.GetPermissions(ctx, uc.Uid)
	if err != nil {
		h.l.Error("获取权限失败", zap.Error(err))
		return ginp.Result{
			Code: http.StatusInternalServerError,
			Msg:  "获取权限失败",
		}, err
	}
	return ginp.Result{
		Code: http.StatusOK,
		Msg:  "成功获取权限",
		Data: permissions,
	}, nil
}

// AssignPermission 处理分配权限的请求
func (h *PermissionHandler) AssignPermission(ctx *gin.Context, req AssignPermissionReq) (ginp.Result, error) {
	if err := h.svc.AssignPermission(ctx, req.UserID, req.Path, req.Method); err != nil {
		h.l.Error("分配权限失败", zap.Error(err))
		return ginp.Result{
			Code: http.StatusInternalServerError,
			Msg:  "分配权限失败",
		}, err
	}
	return ginp.Result{
		Code: http.StatusOK,
		Msg:  "成功分配权限",
	}, nil
}

// RemovePermission 处理移除权限的请求
func (h *PermissionHandler) RemovePermission(ctx *gin.Context, req RemovePermissionReq) (ginp.Result, error) {
	if err := h.svc.RemovePermission(ctx, req.UserID, req.Path, req.Method); err != nil {
		h.l.Error("移除权限失败", zap.Error(err))
		return ginp.Result{
			Code: http.StatusInternalServerError,
			Msg:  "移除权限失败",
		}, err
	}
	return ginp.Result{
		Code: http.StatusOK,
		Msg:  "成功移除权限",
	}, nil
}
