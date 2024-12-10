package service

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"go.uber.org/zap"
)

type PermissionService interface {
	// 菜单管理
	GetMenus(ctx context.Context, pageNum, pageSize int) ([]*domain.Menu, int, error)
	CreateMenu(ctx context.Context, menu *domain.Menu) error
	GetMenuById(ctx context.Context, id int) (*domain.Menu, error)
	UpdateMenu(ctx context.Context, menu *domain.Menu) error
	DeleteMenu(ctx context.Context, id int) error
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
	AssignRoleToUser(ctx context.Context, userId int, roleIds []int) error
	RemoveUserPermissions(ctx context.Context, userId int) error
	RemoveRoleFromUser(ctx context.Context, userId int, roleIds []int) error
}

type permissionService struct {
	repo repository.PermissionRepository
	l    *zap.Logger
}

func NewPermissionService(repo repository.PermissionRepository, l *zap.Logger) PermissionService {
	return &permissionService{
		repo: repo,
		l:    l,
	}
}

// 菜单管理
func (s *permissionService) GetMenus(ctx context.Context, pageNum, pageSize int) ([]*domain.Menu, int, error) {
	if pageNum < 1 || pageSize < 1 {
		s.l.Warn("分页参数无效", zap.Int("页码", pageNum), zap.Int("每页数量", pageSize))
		return nil, 0, nil
	}

	return s.repo.ListMenus(ctx, pageNum, pageSize)
}

func (s *permissionService) CreateMenu(ctx context.Context, menu *domain.Menu) error {
	if menu == nil {
		s.l.Warn("菜单不能为空")
		return nil
	}

	return s.repo.CreateMenu(ctx, menu)
}

func (s *permissionService) GetMenuById(ctx context.Context, id int) (*domain.Menu, error) {
	if id <= 0 {
		s.l.Warn("菜单ID无效", zap.Int("ID", id))
		return nil, nil
	}

	return s.repo.GetMenuById(ctx, id)
}

func (s *permissionService) UpdateMenu(ctx context.Context, menu *domain.Menu) error {
	if menu == nil {
		s.l.Warn("菜单不能为空")
		return nil
	}

	return s.repo.UpdateMenu(ctx, menu)
}

func (s *permissionService) DeleteMenu(ctx context.Context, id int) error {
	if id <= 0 {
		s.l.Warn("菜单ID无效", zap.Int("ID", id))
		return nil
	}

	return s.repo.DeleteMenu(ctx, id)
}

func (s *permissionService) GetMenuTree(ctx context.Context) ([]*domain.Menu, error) {
	return s.repo.GetMenuTree(ctx)
}

// API接口管理
func (s *permissionService) CreateApi(ctx context.Context, api *domain.Api) error {
	if api == nil {
		s.l.Warn("API不能为空")
		return nil
	}

	return s.repo.CreateApi(ctx, api)
}

func (s *permissionService) GetApiById(ctx context.Context, id int) (*domain.Api, error) {
	if id <= 0 {
		s.l.Warn("API ID无效", zap.Int("ID", id))
		return nil, nil
	}

	return s.repo.GetApiById(ctx, id)
}

func (s *permissionService) UpdateApi(ctx context.Context, api *domain.Api) error {
	if api == nil {
		s.l.Warn("API不能为空")
		return nil
	}

	return s.repo.UpdateApi(ctx, api)
}

func (s *permissionService) DeleteApi(ctx context.Context, id int) error {
	if id <= 0 {
		s.l.Warn("API ID无效", zap.Int("ID", id))
		return nil
	}

	return s.repo.DeleteApi(ctx, id)
}

func (s *permissionService) ListApis(ctx context.Context, page, pageSize int) ([]*domain.Api, int, error) {
	if page < 1 || pageSize < 1 {
		s.l.Warn("分页参数无效", zap.Int("页码", page), zap.Int("每页数量", pageSize))
		return nil, 0, nil
	}

	return s.repo.ListApis(ctx, page, pageSize)
}

// 角色管理
func (s *permissionService) CreateRole(ctx context.Context, role *domain.Role) error {
	if role == nil {
		s.l.Warn("角色不能为空")
		return nil
	}

	return s.repo.CreateRole(ctx, role)
}

func (s *permissionService) GetRoleById(ctx context.Context, id int) (*domain.Role, error) {
	if id <= 0 {
		s.l.Warn("角色ID无效", zap.Int("ID", id))
		return nil, nil
	}

	return s.repo.GetRoleById(ctx, id)
}

func (s *permissionService) UpdateRole(ctx context.Context, role *domain.Role) error {
	if role == nil {
		s.l.Warn("角色不能为空")
		return nil
	}

	return s.repo.UpdateRole(ctx, role)
}

func (s *permissionService) DeleteRole(ctx context.Context, id int) error {
	if id <= 0 {
		s.l.Warn("角色ID无效", zap.Int("ID", id))
		return nil
	}

	return s.repo.DeleteRole(ctx, id)
}

func (s *permissionService) ListRoles(ctx context.Context, page, pageSize int) ([]*domain.Role, int, error) {
	if page < 1 || pageSize < 1 {
		s.l.Warn("分页参数无效", zap.Int("页码", page), zap.Int("每页数量", pageSize))
		return nil, 0, nil
	}

	return s.repo.ListRoles(ctx, page, pageSize)
}

func (s *permissionService) AssignPermissions(ctx context.Context, roleId int, menuIds []int, apiIds []int) error {
	if roleId <= 0 {
		s.l.Warn("角色ID无效", zap.Int("roleId", roleId))
		return nil
	}

	return s.repo.AssignPermissions(ctx, roleId, menuIds, apiIds)
}

func (s *permissionService) AssignRoleToUser(ctx context.Context, userId int, roleIds []int) error {
	if userId <= 0 {
		s.l.Warn("用户ID无效", zap.Int("userId", userId))
		return nil
	}

	return s.repo.AssignRoleToUser(ctx, userId, roleIds)
}

func (s *permissionService) RemoveUserPermissions(ctx context.Context, userId int) error {
	if userId <= 0 {
		s.l.Warn("用户ID无效", zap.Int("userId", userId))
		return nil
	}

	return s.repo.RemoveUserPermissions(ctx, userId)
}

func (s *permissionService) RemoveRoleFromUser(ctx context.Context, userId int, roleIds []int) error {
	if userId <= 0 {
		s.l.Warn("用户ID无效", zap.Int("userId", userId))
		return nil
	}

	return s.repo.RemoveRoleFromUser(ctx, userId, roleIds)
}
