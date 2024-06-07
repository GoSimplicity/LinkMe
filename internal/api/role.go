package api

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/service"
	"LinkMe/pkg/ginp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RoleHandler struct {
	svc service.RoleService
	l   *zap.Logger
}

func NewRoleHandler(svc service.RoleService, l *zap.Logger) *RoleHandler {
	return &RoleHandler{
		svc: svc,
		l:   l,
	}
}

func (rh *RoleHandler) RegisterRoutes(server *gin.Engine) {
	roleGroup := server.Group("/roles")
	roleGroup.POST("/", ginp.WrapBody(rh.CreateRole))
	roleGroup.POST("/permission", ginp.WrapBody(rh.CreatePermission))
	roleGroup.POST("/assign", ginp.WrapBody(rh.AssignPermissionToRole))
}

func (rh *RoleHandler) CreateRole(ctx *gin.Context, req CreateRoleReq) (ginp.Result, error) {
	err := rh.svc.CreateRole(ctx.Request.Context(), domain.Role{
		Name: req.Name,
	})
	if err != nil {
		rh.l.Error("create role failed", zap.Error(err))
		return ginp.Result{Code: 500, Msg: "Create role failed"}, err
	}
	return ginp.Result{Code: 200, Msg: "Role created successfully"}, nil
}

func (rh *RoleHandler) CreatePermission(ctx *gin.Context, req CreatePermissionReq) (ginp.Result, error) {
	err := rh.svc.CreatePermission(ctx.Request.Context(), domain.Permission{
		Name: req.Name,
	})
	if err != nil {
		rh.l.Error("create permission failed", zap.Error(err))
		return ginp.Result{Code: 500, Msg: "Create permission failed"}, err
	}
	return ginp.Result{Code: 200, Msg: "Permission created successfully"}, nil
}

func (rh *RoleHandler) AssignPermissionToRole(ctx *gin.Context, req AssignPermissionToRoleReq) (ginp.Result, error) {
	err := rh.svc.AssignPermissionToRole(ctx.Request.Context(), req.RoleID, req.PermissionID)
	if err != nil {
		rh.l.Error("assign permission to role failed", zap.Error(err))
		return ginp.Result{Code: 500, Msg: "Assign permission to role failed"}, err
	}
	return ginp.Result{Code: 200, Msg: "Permission assigned to role successfully"}, nil
}
