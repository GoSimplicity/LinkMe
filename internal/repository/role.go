package repository

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

type RoleRepository interface {
	CreateRole(ctx context.Context, role *domain.Role, menuIds []int, apiIds []int) error
	GetRoleById(ctx context.Context, id int) (*domain.Role, error)
	UpdateRole(ctx context.Context, role *domain.Role) error
	DeleteRole(ctx context.Context, id int) error
	ListRoles(ctx context.Context, page, pageSize int) ([]*domain.Role, int, error)
	GetUserRole(ctx context.Context, userId int) (*domain.Role, error)
	GetRole(ctx context.Context, roleId int) (*domain.Role, error)
}

type roleRepository struct {
	l   *zap.Logger
	dao dao.RoleDAO
}

func NewRoleRepository(l *zap.Logger, dao dao.RoleDAO) RoleRepository {
	return &roleRepository{
		l:   l,
		dao: dao,
	}
}

// CreateRole 创建角色
func (r *roleRepository) CreateRole(ctx context.Context, role *domain.Role, menuIds []int, apiIds []int) error {
	return r.dao.CreateRole(ctx, r.roleToDAO(role), menuIds, apiIds)
}

// GetRoleById 根据ID获取角色
func (r *roleRepository) GetRoleById(ctx context.Context, id int) (*domain.Role, error) {
	role, err := r.dao.GetRoleById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.roleFromDAO(role), nil
}

// UpdateRole 更新角色信息
func (r *roleRepository) UpdateRole(ctx context.Context, role *domain.Role) error {
	return r.dao.UpdateRole(ctx, r.roleToDAO(role))
}

// DeleteRole 删除角色
func (r *roleRepository) DeleteRole(ctx context.Context, id int) error {
	return r.dao.DeleteRole(ctx, id)
}

// ListRoles 分页获取角色列表
func (r *roleRepository) ListRoles(ctx context.Context, page, pageSize int) ([]*domain.Role, int, error) {
	roles, total, err := r.dao.ListRoles(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return r.rolesFromDAO(roles), total, nil
}

// GetRole 根据角色ID获取角色信息
func (r *roleRepository) GetRole(ctx context.Context, roleId int) (*domain.Role, error) {
	role, err := r.dao.GetRole(ctx, roleId)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	return r.roleFromDAO(role), nil
}

// GetUserRole 获取用户的角色信息
func (r *roleRepository) GetUserRole(ctx context.Context, userId int) (*domain.Role, error) {
	role, err := r.dao.GetUserRole(ctx, userId)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	return r.roleFromDAO(role), nil
}

// roleToDAO 将领域模型转换为DAO模型
func (r *roleRepository) roleToDAO(role *domain.Role) *dao.Role {
	if role == nil {
		return nil
	}
	return &dao.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		RoleType:    role.RoleType,
		IsDefault:   role.IsDefault,
		CreateTime:  role.CreateTime,
		UpdateTime:  role.UpdateTime,
		IsDeleted:   role.IsDeleted,
	}
}

// roleFromDAO 将DAO模型转换为领域模型
func (r *roleRepository) roleFromDAO(role *dao.Role) *domain.Role {
	if role == nil {
		return nil
	}
	return &domain.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		RoleType:    role.RoleType,
		IsDefault:   role.IsDefault,
		CreateTime:  role.CreateTime,
		UpdateTime:  role.UpdateTime,
		IsDeleted:   role.IsDeleted,
		Menus:       r.menusFromDAO(role.Menus),
		Apis:        r.apisFromDAO(role.Apis),
	}
}

// rolesFromDAO 批量将DAO模型转换为领域模型
func (r *roleRepository) rolesFromDAO(roles []*dao.Role) []*domain.Role {
	if roles == nil {
		return nil
	}
	result := make([]*domain.Role, 0, len(roles))
	for _, role := range roles {
		if role != nil {
			result = append(result, r.roleFromDAO(role))
		}
	}
	return result
}

// menusFromDAO 批量将菜单DAO模型转换为领域模型
func (r *roleRepository) menusFromDAO(menus []*dao.Menu) []*domain.Menu {
	if menus == nil {
		return nil
	}
	result := make([]*domain.Menu, 0, len(menus))
	for _, menu := range menus {
		if menu != nil {
			result = append(result, r.menuFromDAO(menu))
		}
	}
	return result
}

// apisFromDAO 批量将API DAO模型转换为领域模型
func (r *roleRepository) apisFromDAO(apis []*dao.Api) []*domain.Api {
	if apis == nil {
		return nil
	}
	result := make([]*domain.Api, 0, len(apis))
	for _, api := range apis {
		if api != nil {
			result = append(result, r.apiFromDAO(api))
		}
	}
	return result
}

// apiFromDAO 将API DAO模型转换为领域模型
func (r *roleRepository) apiFromDAO(api *dao.Api) *domain.Api {
	if api == nil {
		return nil
	}
	return &domain.Api{
		ID:          api.ID,
		Name:        api.Name,
		Path:        api.Path,
		Method:      api.Method,
		Description: api.Description,
		Version:     api.Version,
		Category:    api.Category,
		IsPublic:    api.IsPublic,
		CreateTime:  api.CreateTime,
		UpdateTime:  api.UpdateTime,
		IsDeleted:   api.IsDeleted,
	}
}

// menuFromDAO 将菜单DAO模型转换为领域模型
func (r *roleRepository) menuFromDAO(menu *dao.Menu) *domain.Menu {
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
		Children:   r.menusFromDAO(menu.Children),
		RouteName:  menu.RouteName,
		Hidden:     menu.Hidden,
		CreateTime: menu.CreateTime,
		UpdateTime: menu.UpdateTime,
		IsDeleted:  menu.IsDeleted,
	}
}
