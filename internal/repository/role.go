package repository

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

type PermissionRepository interface {
	// 菜单管理
	CreateMenu(ctx context.Context, menu *domain.Menu) error
	GetMenuById(ctx context.Context, id int) (*domain.Menu, error)
	UpdateMenu(ctx context.Context, menu *domain.Menu) error
	DeleteMenu(ctx context.Context, id int) error
	ListMenus(ctx context.Context, page, pageSize int) ([]*domain.Menu, int, error)
	GetMenuTree(ctx context.Context) ([]*domain.Menu, error)

	// API接口管理
	CreateApi(ctx context.Context, api *domain.Api) error
	GetApiById(ctx context.Context, id int) (*domain.Api, error)
	UpdateApi(ctx context.Context, api *domain.Api) error
	DeleteApi(ctx context.Context, id int) error
	ListApis(ctx context.Context, page, pageSize int) ([]*domain.Api, int, error)

	// 角色管理
	CreateRole(ctx context.Context, role *domain.Role) error
	GetRoleById(ctx context.Context, id int) (*domain.Role, error)
	UpdateRole(ctx context.Context, role *domain.Role) error
	DeleteRole(ctx context.Context, id int) error
	ListRoles(ctx context.Context, page, pageSize int) ([]*domain.Role, int, error)
	AssignPermissions(ctx context.Context, roleId int, menuIds []int, apiIds []int) error
	AssignMenuPermissionsToUser(ctx context.Context, userId int, menuIds []int) error
	AssignApiPermissionsToUser(ctx context.Context, userId int, apiIds []int) error
	AssignRoleToUser(ctx context.Context, userId int, roleIds []int) error
	RemoveUserMenuPermissions(ctx context.Context, userId int, menuIds []int) error
	RemoveUserApiPermissions(ctx context.Context, userId int, apiIds []int) error
	RemoveRoleFromUser(ctx context.Context, userId int, roleIds []int) error
	RemoveRoleApiPermissions(ctx context.Context, roleIds []int, apiIds []int) error
	RemoveRoleMenuPermissions(ctx context.Context, roleIds []int, menuIds []int) error
}

type permissionRepository struct {
	l   *zap.Logger
	dao dao.PermissionDAO
}

func NewPermissionRepository(l *zap.Logger, dao dao.PermissionDAO) PermissionRepository {
	return &permissionRepository{
		l:   l,
		dao: dao,
	}
}

// 菜单管理
func (r *permissionRepository) CreateMenu(ctx context.Context, menu *domain.Menu) error {
	return r.dao.CreateMenu(ctx, r.menuToDAO(menu))
}

func (r *permissionRepository) GetMenuById(ctx context.Context, id int) (*domain.Menu, error) {
	menu, err := r.dao.GetMenuById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.menuFromDAO(menu), nil
}

func (r *permissionRepository) UpdateMenu(ctx context.Context, menu *domain.Menu) error {
	return r.dao.UpdateMenu(ctx, r.menuToDAO(menu))
}

func (r *permissionRepository) DeleteMenu(ctx context.Context, id int) error {
	return r.dao.DeleteMenu(ctx, id)
}

func (r *permissionRepository) ListMenus(ctx context.Context, page, pageSize int) ([]*domain.Menu, int, error) {
	menus, total, err := r.dao.ListMenus(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return r.menusFromDAO(menus), total, nil
}

func (r *permissionRepository) GetMenuTree(ctx context.Context) ([]*domain.Menu, error) {
	menus, err := r.dao.GetMenuTree(ctx)
	if err != nil {
		return nil, err
	}
	return r.menusFromDAO(menus), nil
}

// API接口管理
func (r *permissionRepository) CreateApi(ctx context.Context, api *domain.Api) error {
	return r.dao.CreateApi(ctx, r.apiToDAO(api))
}

func (r *permissionRepository) GetApiById(ctx context.Context, id int) (*domain.Api, error) {
	api, err := r.dao.GetApiById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.apiFromDAO(api), nil
}

func (r *permissionRepository) UpdateApi(ctx context.Context, api *domain.Api) error {
	return r.dao.UpdateApi(ctx, r.apiToDAO(api))
}

func (r *permissionRepository) DeleteApi(ctx context.Context, id int) error {
	return r.dao.DeleteApi(ctx, id)
}

func (r *permissionRepository) ListApis(ctx context.Context, page, pageSize int) ([]*domain.Api, int, error) {
	apis, total, err := r.dao.ListApis(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return r.apisFromDAO(apis), total, nil
}

// 角色管理
func (r *permissionRepository) CreateRole(ctx context.Context, role *domain.Role) error {
	return r.dao.CreateRole(ctx, r.roleToDAO(role))
}

func (r *permissionRepository) GetRoleById(ctx context.Context, id int) (*domain.Role, error) {
	role, err := r.dao.GetRoleById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.roleFromDAO(role), nil
}

func (r *permissionRepository) UpdateRole(ctx context.Context, role *domain.Role) error {
	return r.dao.UpdateRole(ctx, r.roleToDAO(role))
}

func (r *permissionRepository) DeleteRole(ctx context.Context, id int) error {
	return r.dao.DeleteRole(ctx, id)
}

func (r *permissionRepository) ListRoles(ctx context.Context, page, pageSize int) ([]*domain.Role, int, error) {
	roles, total, err := r.dao.ListRoles(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return r.rolesFromDAO(roles), total, nil
}

func (r *permissionRepository) AssignPermissions(ctx context.Context, roleId int, menuIds []int, apiIds []int) error {
	return r.dao.AssignPermissions(ctx, roleId, menuIds, apiIds)
}

func (r *permissionRepository) AssignApiPermissionsToUser(ctx context.Context, userId int, apiIds []int) error {
	return r.dao.AssignApiPermissionsToUser(ctx, userId, apiIds)
}

func (r *permissionRepository) AssignMenuPermissionsToUser(ctx context.Context, userId int, menuIds []int) error {
	return r.dao.AssignMenuPermissionsToUser(ctx, userId, menuIds)
}

func (r *permissionRepository) AssignRoleToUser(ctx context.Context, userId int, roleIds []int) error {
	return r.dao.AssignRoleToUser(ctx, userId, roleIds)
}

func (r *permissionRepository) RemoveRoleFromUser(ctx context.Context, userId int, roleIds []int) error {
	return r.dao.RemoveRoleFromUser(ctx, userId, roleIds)
}

func (r *permissionRepository) RemoveUserApiPermissions(ctx context.Context, userId int, apiIds []int) error {
	return r.dao.RemoveUserApiPermissions(ctx, userId, apiIds)
}

func (r *permissionRepository) RemoveUserMenuPermissions(ctx context.Context, userId int, menuIds []int) error {
	return r.dao.RemoveUserMenuPermissions(ctx, userId, menuIds)
}

func (r *permissionRepository) RemoveRoleApiPermissions(ctx context.Context, roleIds []int, apiIds []int) error {
	return r.dao.RemoveRoleApiPermissions(ctx, roleIds, apiIds)
}

func (r *permissionRepository) RemoveRoleMenuPermissions(ctx context.Context, roleIds []int, menuIds []int) error {
	return r.dao.RemoveRoleMenuPermissions(ctx, roleIds, menuIds)
}

// Menu转换方法
func (r *permissionRepository) menuToDAO(menu *domain.Menu) *dao.Menu {
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

func (r *permissionRepository) menuFromDAO(menu *dao.Menu) *domain.Menu {
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

func (r *permissionRepository) menusFromDAO(menus []*dao.Menu) []*domain.Menu {
	var result []*domain.Menu
	for _, m := range menus {
		result = append(result, r.menuFromDAO(m))
	}
	return result
}

// Api转换方法
func (r *permissionRepository) apiToDAO(api *domain.Api) *dao.Api {
	return &dao.Api{
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

func (r *permissionRepository) apiFromDAO(api *dao.Api) *domain.Api {
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

func (r *permissionRepository) apisFromDAO(apis []*dao.Api) []*domain.Api {
	var result []*domain.Api
	for _, a := range apis {
		result = append(result, r.apiFromDAO(a))
	}
	return result
}

// Role转换方法
func (r *permissionRepository) roleToDAO(role *domain.Role) *dao.Role {
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

func (r *permissionRepository) roleFromDAO(role *dao.Role) *domain.Role {
	return &domain.Role{
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

func (r *permissionRepository) rolesFromDAO(roles []*dao.Role) []*domain.Role {
	var result []*domain.Role
	for _, role := range roles {
		result = append(result, r.roleFromDAO(role))
	}
	return result
}
