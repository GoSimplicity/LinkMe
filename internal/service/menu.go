package service

import (
	"context"
	"errors"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"go.uber.org/zap"
)

type MenuService interface {
	GetMenus(ctx context.Context, pageNum, pageSize int, isTree bool) ([]*domain.Menu, int, error)
	CreateMenu(ctx context.Context, menu *domain.Menu) error
	GetMenuById(ctx context.Context, id int) (*domain.Menu, error)
	UpdateMenu(ctx context.Context, menu *domain.Menu) error
	DeleteMenu(ctx context.Context, id int) error
	GetMenuTree(ctx context.Context) ([]*domain.Menu, error)
}

type menuService struct {
	l    *zap.Logger
	repo repository.MenuRepository
}

func NewMenuService(l *zap.Logger, repo repository.MenuRepository) MenuService {
	return &menuService{
		l:    l,
		repo: repo,
	}
}

// GetMenus 获取菜单列表,支持分页和树形结构
func (m *menuService) GetMenus(ctx context.Context, pageNum, pageSize int, isTree bool) ([]*domain.Menu, int, error) {
	if pageNum < 1 || pageSize < 1 {
		m.l.Warn("分页参数无效", zap.Int("页码", pageNum), zap.Int("每页数量", pageSize))
		return nil, 0, errors.New("分页参数无效")
	}

	// 如果需要树形结构,则调用GetMenuTree
	if isTree {
		menus, err := m.repo.GetMenuTree(ctx)
		if err != nil {
			m.l.Error("获取菜单树失败", zap.Error(err))
			return nil, 0, err
		}
		return menus, len(menus), nil
	}

	return m.repo.ListMenus(ctx, pageNum, pageSize)
}

// CreateMenu 创建新菜单
func (m *menuService) CreateMenu(ctx context.Context, menu *domain.Menu) error {
	if menu == nil {
		m.l.Warn("菜单不能为空")
		return errors.New("菜单不能为空")
	}

	return m.repo.CreateMenu(ctx, menu)
}

// GetMenuById 根据ID获取菜单
func (m *menuService) GetMenuById(ctx context.Context, id int) (*domain.Menu, error) {
	if id <= 0 {
		m.l.Warn("菜单ID无效", zap.Int("ID", id))
		return nil, errors.New("菜单ID无效")
	}

	return m.repo.GetMenuById(ctx, id)
}

// UpdateMenu 更新菜单信息
func (m *menuService) UpdateMenu(ctx context.Context, menu *domain.Menu) error {
	if menu == nil {
		m.l.Warn("菜单不能为空")
		return errors.New("菜单不能为空")
	}

	return m.repo.UpdateMenu(ctx, menu)
}

// DeleteMenu 删除指定ID的菜单
func (m *menuService) DeleteMenu(ctx context.Context, id int) error {
	if id <= 0 {
		m.l.Warn("菜单ID无效", zap.Int("ID", id))
		return errors.New("菜单ID无效")
	}

	return m.repo.DeleteMenu(ctx, id)
}

// GetMenuTree 获取菜单树形结构
func (m *menuService) GetMenuTree(ctx context.Context) ([]*domain.Menu, error) {
	return m.repo.GetMenuTree(ctx)
}
