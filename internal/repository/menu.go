package repository

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

type MenuRepository interface {
	CreateMenu(ctx context.Context, menu *domain.Menu) error
	GetMenuById(ctx context.Context, id int) (*domain.Menu, error)
	UpdateMenu(ctx context.Context, menu *domain.Menu) error
	DeleteMenu(ctx context.Context, id int) error
	ListMenus(ctx context.Context, page, pageSize int) ([]*domain.Menu, int, error)
	GetMenuTree(ctx context.Context) ([]*domain.Menu, error)
}

type menuRepository struct {
	l   *zap.Logger
	dao dao.MenuDAO
}

func NewMenuRepository(l *zap.Logger, dao dao.MenuDAO) MenuRepository {
	return &menuRepository{
		l:   l,
		dao: dao,
	}
}

// CreateMenu 创建新的菜单
func (m *menuRepository) CreateMenu(ctx context.Context, menu *domain.Menu) error {
	return m.dao.CreateMenu(ctx, m.menuToDAO(menu))
}

// GetMenuById 根据ID获取菜单
func (m *menuRepository) GetMenuById(ctx context.Context, id int) (*domain.Menu, error) {
	menu, err := m.dao.GetMenuById(ctx, id)
	if err != nil {
		return nil, err
	}
	return m.menuFromDAO(menu), nil
}

// UpdateMenu 更新菜单信息
func (m *menuRepository) UpdateMenu(ctx context.Context, menu *domain.Menu) error {
	return m.dao.UpdateMenu(ctx, m.menuToDAO(menu))
}

// DeleteMenu 删除菜单
func (m *menuRepository) DeleteMenu(ctx context.Context, id int) error {
	return m.dao.DeleteMenu(ctx, id)
}

// ListMenus 分页获取菜单列表
func (m *menuRepository) ListMenus(ctx context.Context, page, pageSize int) ([]*domain.Menu, int, error) {
	menus, total, err := m.dao.ListMenus(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return m.menusFromDAO(menus), total, nil
}

// GetMenuTree 获取菜单树
func (m *menuRepository) GetMenuTree(ctx context.Context) ([]*domain.Menu, error) {
	menus, err := m.dao.GetMenuTree(ctx)
	if err != nil {
		return nil, err
	}
	return m.menusFromDAO(menus), nil
}

// menuToDAO 将领域模型转换为数据访问对象
func (m *menuRepository) menuToDAO(menu *domain.Menu) *dao.Menu {
	if menu == nil {
		return nil
	}
	return &dao.Menu{
		ID:         menu.ID,
		Name:       menu.Name,
		ParentID:   menu.ParentID,
		Path:       menu.Path,
		Component:  menu.Component,
		Icon:       menu.Icon,
		SortOrder:  menu.SortOrder,
		RouteName:  menu.RouteName,
		Hidden:     menu.Hidden,
		CreateTime: menu.CreateTime,
		UpdateTime: menu.UpdateTime,
		IsDeleted:  menu.IsDeleted,
	}
}

// menuFromDAO 将数据访问对象转换为领域模型
func (m *menuRepository) menuFromDAO(menu *dao.Menu) *domain.Menu {
	if menu == nil {
		return nil
	}
	return &domain.Menu{
		ID:         menu.ID,
		Name:       menu.Name,
		ParentID:   menu.ParentID,
		Path:       menu.Path,
		Component:  menu.Component,
		Icon:       menu.Icon,
		SortOrder:  menu.SortOrder,
		Children:   m.menusFromDAO(menu.Children), // 递归转换子菜单
		RouteName:  menu.RouteName,
		Hidden:     menu.Hidden,
		CreateTime: menu.CreateTime,
		UpdateTime: menu.UpdateTime,
		IsDeleted:  menu.IsDeleted,
	}
}

// menusFromDAO 批量转换菜单列表
func (m *menuRepository) menusFromDAO(menus []*dao.Menu) []*domain.Menu {
	if menus == nil {
		return nil
	}
	result := make([]*domain.Menu, 0, len(menus))
	for _, menu := range menus {
		if menu != nil {
			result = append(result, m.menuFromDAO(menu))
		}
	}
	return result
}
