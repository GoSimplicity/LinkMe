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

type MenuHandler struct {
	svc service.MenuService
	l   *zap.Logger
}

func NewMenuHandler(svc service.MenuService, l *zap.Logger) *MenuHandler {
	return &MenuHandler{
		svc: svc,
		l:   l,
	}
}

func (h *MenuHandler) RegisterRoutes(server *gin.Engine) {
	menuGroup := server.Group("/api/menus")

	menuGroup.POST("/list", h.ListMenus)
	menuGroup.POST("/create", h.CreateMenu)
	menuGroup.POST("/update", h.UpdateMenu)
	menuGroup.DELETE("/:id", h.DeleteMenu)
}

// ListMenus 获取菜单列表
func (m *MenuHandler) ListMenus(c *gin.Context) {
	var req req.ListMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 调用service层获取菜单列表
	menus, total, err := m.svc.GetMenus(c.Request.Context(), req.PageNumber, req.PageSize, req.IsTree)
	if err != nil {
		m.l.Error("获取菜单列表失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.SuccessWithData(c, gin.H{
		"list":  menus,
		"total": total,
	})
}

// CreateMenu 创建菜单
func (m *MenuHandler) CreateMenu(c *gin.Context) {
	var req req.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 构建菜单对象
	menu := &domain.Menu{
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		SortOrder: req.SortOrder,
		ParentID:  req.ParentId,
		Icon:      req.Icon,
		Hidden:    req.Hidden,
		RouteName: req.RouteName,
	}

	if err := m.svc.CreateMenu(c.Request.Context(), menu); err != nil {
		m.l.Error("创建菜单失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// UpdateMenu 更新菜单
func (m *MenuHandler) UpdateMenu(c *gin.Context) {
	var req req.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 构建更新的菜单对象
	menu := &domain.Menu{
		ID:        req.Id,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		SortOrder: req.SortOrder,
		ParentID:  req.ParentId,
		Icon:      req.Icon,
		Hidden:    req.Hidden,
		RouteName: req.RouteName,
	}

	if err := m.svc.UpdateMenu(c.Request.Context(), menu); err != nil {
		m.l.Error("更新菜单失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// DeleteMenu 删除菜单
func (m *MenuHandler) DeleteMenu(c *gin.Context) {
	// 从URL参数中获取菜单ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		m.l.Error("解析ID失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := m.svc.DeleteMenu(c.Request.Context(), id); err != nil {
		m.l.Error("删除菜单失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}
