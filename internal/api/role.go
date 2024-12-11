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

	// 菜单管理
	permissionGroup.POST("/menus/list", h.ListMenus)
	permissionGroup.POST("/menu/create", h.CreateMenu)
	permissionGroup.POST("/menu/update", h.UpdateMenu)
	permissionGroup.DELETE("/menu/:id", h.DeleteMenu)

	// API接口管理
	permissionGroup.POST("/api/list", h.ListApis)
	permissionGroup.POST("/api/create", h.CreateAPI)
	permissionGroup.POST("/api/update", h.UpdateAPI)
	permissionGroup.DELETE("/api/:id", h.DeleteAPI)

	// 角色管理
	permissionGroup.POST("/role/list", h.ListRoles)
	permissionGroup.POST("/role/create", h.CreateRole)
	permissionGroup.POST("/role/update", h.UpdateRole)
	permissionGroup.DELETE("/role/:id", h.DeleteRole)

	permissionGroup.POST("/assign/role/permissions", h.AssignPermissions)    // 分配权限
	permissionGroup.POST("/assign/user/api", h.AssignApiPermissionsToUser)   // 分配用户API权限
	permissionGroup.POST("/assign/user/menu", h.AssignMenuPermissionsToUser) // 分配用户菜单权限
	permissionGroup.POST("/assign/user/role", h.AssignRoleToUser)            // 分配角色
	permissionGroup.DELETE("/remove/user/role", h.RemoveRoleFromUser)        // 移除用户角色
	permissionGroup.DELETE("/remove/user/api", h.RemoveUserApiPermissions)   // 移除用户API权限
	permissionGroup.DELETE("/remove/user/menu", h.RemoveUserMenuPermissions) // 移除用户菜单权限
	permissionGroup.DELETE("/remove/role/api", h.RemoveRoleApiPermissions)   // 移除角色API权限
	permissionGroup.DELETE("/remove/role/menu", h.RemoveRoleMenuPermissions) // 移除角色菜单权限
}

// 菜单管理
func (h *PermissionHandler) ListMenus(c *gin.Context) {
	var r req.ListMenusRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	menus, total, err := h.svc.GetMenus(c.Request.Context(), r.PageNumber, r.PageSize, r.IsTree)
	if err != nil {
		h.l.Error("获取菜单列表失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.SuccessWithData(c, gin.H{
		"list":  menus,
		"total": total,
	})
}

func (h *PermissionHandler) CreateMenu(c *gin.Context) {
	var r req.CreateMenuRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	menu := &domain.Menu{
		Name:       r.Name,
		Path:       r.Path,
		Component:  r.Component,
		SortOrder:  r.SortOrder,
		ParentID:   int64(r.ParentId),
		Icon:       r.Icon,
		Hidden:     r.Hidden,
		RouteName:  r.RouteName,
		CreateTime: 0,
		UpdateTime: 0,
		IsDeleted:  0,
	}

	if err := h.svc.CreateMenu(c.Request.Context(), menu); err != nil {
		h.l.Error("创建菜单失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) UpdateMenu(c *gin.Context) {
	var r req.UpdateMenuRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	menu := &domain.Menu{
		ID:         int64(r.Id),
		Name:       r.Name,
		Path:       r.Path,
		Component:  r.Component,
		SortOrder:  r.SortOrder,
		ParentID:   int64(r.ParentId),
		Icon:       r.Icon,
		Hidden:     r.Hidden,
		RouteName:  r.RouteName,
		UpdateTime: 0,
	}

	if err := h.svc.UpdateMenu(c.Request.Context(), menu); err != nil {
		h.l.Error("更新菜单失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.l.Error("解析ID失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.DeleteMenu(c.Request.Context(), id); err != nil {
		h.l.Error("删除菜单失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// API接口管理
func (h *PermissionHandler) ListApis(c *gin.Context) {
	var r req.ListApisRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apis, total, err := h.svc.ListApis(c.Request.Context(), r.PageNumber, r.PageSize)
	if err != nil {
		h.l.Error("获取API列表失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.SuccessWithData(c, gin.H{
		"list":  apis,
		"total": total,
	})
}

func (h *PermissionHandler) CreateAPI(c *gin.Context) {
	var r req.CreateApiRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	api := &domain.Api{
		Name:        r.Name,
		Path:        r.Path,
		Method:      r.Method,
		Description: r.Description,
		Version:     r.Version,
		Category:    r.Category,
		IsPublic:    r.IsPublic,
		CreateTime:  0,
		UpdateTime:  0,
		IsDeleted:   0,
	}

	if err := h.svc.CreateApi(c.Request.Context(), api); err != nil {
		h.l.Error("创建API失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) UpdateAPI(c *gin.Context) {
	var r req.UpdateApiRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	api := &domain.Api{
		ID:          r.Id,
		Name:        r.Name,
		Path:        r.Path,
		Method:      r.Method,
		Description: r.Description,
		Version:     r.Version,
		Category:    r.Category,
		IsPublic:    r.IsPublic,
		UpdateTime:  0,
	}

	if err := h.svc.UpdateApi(c.Request.Context(), api); err != nil {
		h.l.Error("更新API失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) DeleteAPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.l.Error("解析ID失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.DeleteApi(c.Request.Context(), id); err != nil {
		h.l.Error("删除API失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// 角色管理
func (h *PermissionHandler) ListRoles(c *gin.Context) {
	var r req.ListRolesRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	roles, total, err := h.svc.ListRoles(c.Request.Context(), r.PageNumber, r.PageSize)
	if err != nil {
		h.l.Error("获取角色列表失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.SuccessWithData(c, gin.H{
		"list":  roles,
		"total": total,
	})
}

func (h *PermissionHandler) CreateRole(c *gin.Context) {
	var r req.CreateRoleRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	role := &domain.Role{
		Name:        r.Name,
		Description: r.Description,
		RoleType:    r.RoleType,
		IsDefault:   r.IsDefault,
		CreateTime:  0,
		UpdateTime:  0,
		IsDeleted:   0,
	}

	if err := h.svc.CreateRole(c.Request.Context(), role); err != nil {
		h.l.Error("创建角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) UpdateRole(c *gin.Context) {
	var r req.UpdateRoleRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	role := &domain.Role{
		ID:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		RoleType:    r.RoleType,
		IsDefault:   r.IsDefault,
		UpdateTime:  0,
	}

	if err := h.svc.UpdateRole(c.Request.Context(), role); err != nil {
		h.l.Error("更新角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.l.Error("解析ID失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.DeleteRole(c.Request.Context(), id); err != nil {
		h.l.Error("删除角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// 权限分配
func (h *PermissionHandler) AssignPermissions(c *gin.Context) {
	var r req.AssignPermissionsRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.AssignPermissions(c.Request.Context(), r.RoleId, r.MenuIds, r.ApiIds); err != nil {
		h.l.Error("分配权限失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) AssignApiPermissionsToUser(c *gin.Context) {
	var r req.AssignApiPermissionsToUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.AssignApiPermissionsToUser(c.Request.Context(), r.UserId, r.ApiIds); err != nil {
		h.l.Error("分配API权限失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) AssignMenuPermissionsToUser(c *gin.Context) {
	var r req.AssignMenuPermissionsToUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.AssignMenuPermissionsToUser(c.Request.Context(), r.UserId, r.MenuIds); err != nil {
		h.l.Error("分配菜单权限失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) AssignRoleToUser(c *gin.Context) {
	var r req.AssignRoleToUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.AssignRoleToUser(c.Request.Context(), r.UserId, r.RoleIds); err != nil {
		h.l.Error("分配角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) RemoveRoleFromUser(c *gin.Context) {
	var r req.RemoveRoleFromUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.RemoveRoleFromUser(c.Request.Context(), r.UserId, r.RoleIds); err != nil {
		h.l.Error("移除用户角色失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) RemoveUserApiPermissions(c *gin.Context) {
	var r req.RemoveUserApiPermissionsRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.RemoveUserApiPermissions(c.Request.Context(), r.UserId, r.ApiIds); err != nil {
		h.l.Error("移除用户API权限失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) RemoveUserMenuPermissions(c *gin.Context) {
	var r req.RemoveUserMenuPermissionsRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.RemoveUserMenuPermissions(c.Request.Context(), r.UserId, r.MenuIds); err != nil {
		h.l.Error("移除用户菜单权限失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) RemoveRoleApiPermissions(c *gin.Context) {
	var r req.RemoveRoleApiPermissionsRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.RemoveRoleApiPermissions(c.Request.Context(), r.RoleIds, r.ApiIds); err != nil {
		h.l.Error("移除角色API权限失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

func (h *PermissionHandler) RemoveRoleMenuPermissions(c *gin.Context) {
	var r req.RemoveRoleMenuPermissionsRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := h.svc.RemoveRoleMenuPermissions(c.Request.Context(), r.RoleIds, r.MenuIds); err != nil {
		h.l.Error("移除角色菜单权限失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}
