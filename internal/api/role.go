package api

import (
	"LinkMe/internal/service"
	"LinkMe/middleware"
	. "LinkMe/pkg/ginp"
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
	permissionGroup.GET("/list", WrapQuery(h.GetPermissions))                // 获取权限列表
	permissionGroup.POST("/assign", WrapBody(h.AssignPermission))            // 分配权限
	permissionGroup.POST("/assign_role", WrapBody(h.AssignRole))             // 分配角色
	permissionGroup.DELETE("/remove", WrapBody(h.RemovePermission))          // 移除权限
	permissionGroup.DELETE("/remove_role", WrapBody(h.RemovePermissionRole)) // 移除角色
}

// GetPermissions 处理获取权限列表的请求
func (h *PermissionHandler) GetPermissions(ctx *gin.Context, req ListPermissionsReq) (Result, error) {
	permissions, err := h.svc.GetPermissions(ctx)
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
	if err := h.svc.AssignPermission(ctx, req.UserName, req.Path, req.Method); err != nil {
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
	if err := h.svc.RemovePermission(ctx, req.UserName, req.Path, req.Method); err != nil {
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

func (h *PermissionHandler) AssignRole(ctx *gin.Context, req AssignPermissionRoleReq) (Result, error) {
	if err := h.svc.AssignRoleToUser(ctx, req.UserName, req.RoleName); err != nil {
		return Result{
			Code: http.StatusInternalServerError,
			Msg:  "assign role failed",
		}, err
	}
	return Result{
		Code: http.StatusOK,
		Msg:  "assign role success",
	}, nil
}

func (h *PermissionHandler) RemovePermissionRole(ctx *gin.Context, req RemovePermissionRoleReq) (Result, error) {
	if err := h.svc.RemoveRoleFromUser(ctx, req.UserName, req.RoleName); err != nil {
		return Result{
			Code: http.StatusInternalServerError,
			Msg:  "remove role failed",
		}, err
	}
	return Result{
		Code: http.StatusOK,
		Msg:  "remove role success",
	}, nil
}
