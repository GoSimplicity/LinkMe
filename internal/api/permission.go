package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/req"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PermissionHandler struct {
	svc service.PermissionService
	l   *zap.Logger
}

func NewPermissionHandler(svc service.PermissionService, l *zap.Logger) *PermissionHandler {
	return &PermissionHandler{
		svc: svc,
		l:   l,
	}
}

func (h *PermissionHandler) RegisterRoutes(server *gin.Engine) {
	permissionGroup := server.Group("/api/permissions")

	permissionGroup.POST("/user/assign", h.AssignUserRole)
	permissionGroup.POST("/users/assign", h.AssignUsersRole)
}

// AssignUserRole 为单个用户分配角色和权限
func (h *PermissionHandler) AssignUserRole(c *gin.Context) {
	var r req.AssignUserRoleRequest
	// 绑定请求参数
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 调用服务层分配角色和权限
	if err := h.svc.AssignRoleToUser(c.Request.Context(), r.UserId, r.RoleIds, r.MenuIds, r.ApiIds); err != nil {
		h.l.Error("分配角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// AssignUsersRole 批量为用户分配角色和权限
func (h *PermissionHandler) AssignUsersRole(c *gin.Context) {
	var r req.AssignUsersRoleRequest
	// 绑定请求参数
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 调用服务层批量分配角色和权限
	if err := h.svc.AssignRoleToUsers(c.Request.Context(), r.UserIds, r.RoleIds, r.MenuIds, r.ApiIds); err != nil {
		h.l.Error("分配角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}
