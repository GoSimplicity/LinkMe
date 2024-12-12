package service

import (
	"context"
	"errors"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"go.uber.org/zap"
)

var (
	ErrInvalidRole   = errors.New("角色不能为空")
	ErrInvalidID     = errors.New("无效的ID")
	ErrInvalidParams = errors.New("无效的参数")
)

type RoleService interface {
	CreateRole(ctx context.Context, role *domain.Role, menuIds []int, apiIds []int) error
	GetRoleById(ctx context.Context, id int) (*domain.Role, error)
	UpdateRole(ctx context.Context, role *domain.Role) error
	DeleteRole(ctx context.Context, id int) error
	ListRoles(ctx context.Context, page, pageSize int) ([]*domain.Role, int, error)
	GetUserRole(ctx context.Context, userId int) (*domain.Role, error)
	GetRole(ctx context.Context, roleId int) (*domain.Role, error)
}

type roleService struct {
	repo           repository.RoleRepository
	permissionRepo repository.PermissionRepository
	l              *zap.Logger
}

func NewRoleService(repo repository.RoleRepository, permissionRepo repository.PermissionRepository, l *zap.Logger) RoleService {
	return &roleService{
		repo:           repo,
		permissionRepo: permissionRepo,
		l:              l,
	}
}

// CreateRole 创建新角色
func (r *roleService) CreateRole(ctx context.Context, role *domain.Role, menuIds []int, apiIds []int) error {
	if role == nil {
		r.l.Warn("角色不能为空")
		return ErrInvalidRole
	}

	return r.repo.CreateRole(ctx, role, menuIds, apiIds)
}

// GetRoleById 根据ID获取角色信息
func (r *roleService) GetRoleById(ctx context.Context, id int) (*domain.Role, error) {
	if id <= 0 {
		r.l.Warn("角色ID无效", zap.Int("ID", id))
		return nil, ErrInvalidID
	}

	return r.repo.GetRoleById(ctx, id)
}

// UpdateRole 更新角色信息
func (r *roleService) UpdateRole(ctx context.Context, role *domain.Role) error {
	if role == nil {
		r.l.Warn("角色不能为空")
		return ErrInvalidRole
	}

	return r.repo.UpdateRole(ctx, role)
}

// DeleteRole 删除角色及其相关权限
func (r *roleService) DeleteRole(ctx context.Context, id int) error {
	if id <= 0 {
		r.l.Warn("角色ID无效", zap.Int("ID", id))
		return ErrInvalidID
	}

	// 删除角色前先删除相关权限
	if err := r.permissionRepo.RemoveRolePermissions(ctx, id); err != nil {
		r.l.Error("删除角色API权限失败", zap.Error(err))
		return err
	}

	return r.repo.DeleteRole(ctx, id)
}

// ListRoles 分页获取角色列表
func (r *roleService) ListRoles(ctx context.Context, page, pageSize int) ([]*domain.Role, int, error) {
	if page < 1 || pageSize < 1 {
		r.l.Warn("分页参数无效", zap.Int("页码", page), zap.Int("每页数量", pageSize))
		return nil, 0, ErrInvalidParams
	}

	return r.repo.ListRoles(ctx, page, pageSize)
}

// GetRole 根据角色ID获取角色信息
func (r *roleService) GetRole(ctx context.Context, roleId int) (*domain.Role, error) {
	if roleId <= 0 {
		r.l.Warn("角色ID无效", zap.Int("roleId", roleId))
		return nil, ErrInvalidID
	}

	return r.repo.GetRole(ctx, roleId)
}

// GetUserRole 获取用户的角色信息
func (r *roleService) GetUserRole(ctx context.Context, userId int) (*domain.Role, error) {
	if userId <= 0 {
		r.l.Warn("用户ID无效", zap.Int("userId", userId))
		return nil, ErrInvalidID
	}

	return r.repo.GetUserRole(ctx, userId)
}
