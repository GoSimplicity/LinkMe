package api

import (
	"strconv"

	"github.com/GoSimplicity/LinkMe/internal/api/req"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RoleHandler struct {
	svc           service.RoleService
	menuSvc       service.MenuService
	apiSvc        service.ApiService
	permissionSvc service.PermissionService
	l             *zap.Logger
}

func NewRoleHandler(svc service.RoleService, menuSvc service.MenuService, apiSvc service.ApiService, permissionSvc service.PermissionService, l *zap.Logger) *RoleHandler {
	return &RoleHandler{
		svc:           svc,
		menuSvc:       menuSvc,
		apiSvc:        apiSvc,
		permissionSvc: permissionSvc,
		l:             l,
	}
}

func (r *RoleHandler) RegisterRoutes(server *gin.Engine) {
	roleGroup := server.Group("/api/role")

	roleGroup.POST("/list", r.ListRoles)
	roleGroup.POST("/create", r.CreateRole)
	roleGroup.POST("/update", r.UpdateRole)
	roleGroup.DELETE("/:id", r.DeleteRole)
	roleGroup.GET("/user/:id", r.GetUserRoles)
	roleGroup.GET("/:id", r.GetRoles)
}

// ListRoles 获取角色列表
func (r *RoleHandler) ListRoles(c *gin.Context) {
	var req req.ListRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 调用service获取角色列表
	roles, total, err := r.svc.ListRoles(c.Request.Context(), req.PageNumber, req.PageSize)
	if err != nil {
		r.l.Error("获取角色列表失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.SuccessWithData(c, gin.H{
		"list":  roles,
		"total": total,
	})
}

// CreateRole 创建角色
func (r *RoleHandler) CreateRole(c *gin.Context) {
	var req req.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 构建角色对象
	role := &domain.Role{
		Name:        req.Name,
		Description: req.Description,
		RoleType:    req.RoleType,
		IsDefault:   req.IsDefault,
	}

	// 创建角色并分配权限
	if err := r.svc.CreateRole(c.Request.Context(), role, req.MenuIds, req.ApiIds); err != nil {
		r.l.Error("创建角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// UpdateRole 更新角色
func (r *RoleHandler) UpdateRole(c *gin.Context) {
	var req req.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 构建角色对象
	role := &domain.Role{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		RoleType:    req.RoleType,
		IsDefault:   req.IsDefault,
	}

	// 更新角色基本信息
	if err := r.svc.UpdateRole(c.Request.Context(), role); err != nil {
		r.l.Error("更新角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 更新角色权限
	if err := r.permissionSvc.AssignRole(c.Request.Context(), role.ID, req.MenuIds, req.ApiIds); err != nil {
		r.l.Error("更新权限失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// DeleteRole 删除角色
func (r *RoleHandler) DeleteRole(c *gin.Context) {
	// 从URL参数中获取角色ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.l.Error("解析ID失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := r.svc.DeleteRole(c.Request.Context(), id); err != nil {
		r.l.Error("删除角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// UpdateUserRole 更新用户角色
func (r *RoleHandler) UpdateUserRole(c *gin.Context) {
	var req req.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 分配用户角色和权限
	if err := r.permissionSvc.AssignRoleToUser(c.Request.Context(), req.UserId, req.RoleIds, req.MenuIds, req.ApiIds); err != nil {
		r.l.Error("分配API权限失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// GetUserRoles 获取用户角色
func (r *RoleHandler) GetUserRoles(c *gin.Context) {
	// 从URL参数中获取用户ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.l.Error("解析ID失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	role, err := r.svc.GetUserRole(c.Request.Context(), id)
	if err != nil {
		r.l.Error("获取用户角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.SuccessWithData(c, role)
}

// GetRoles 获取角色详情
func (r *RoleHandler) GetRoles(c *gin.Context) {
	// 从URL参数中获取角色ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.l.Error("解析ID失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	role, err := r.svc.GetRole(c.Request.Context(), id)
	if err != nil {
		r.l.Error("获取角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.SuccessWithData(c, role)
}
