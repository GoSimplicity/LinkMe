package repository

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

// PermissionRepository 定义了权限仓库接口
type PermissionRepository interface {
	GetPermissions(ctx context.Context) ([]domain.Permission, error)                         // 获取权限列表
	AssignPermission(ctx context.Context, userName string, path string, method string) error // 分配权限
	AssignRoleToUser(ctx context.Context, userName, roleName string) error
	RemovePermission(ctx context.Context, userName string, path string, method string) error // 移除权限
	RemoveRoleFromUser(ctx context.Context, userName, roleName string) error
}

// permissionRepository 是 PermissionRepository 的实现
type permissionRepository struct {
	l   *zap.Logger       // 日志记录器
	dao dao.PermissionDAO // 数据库访问对象
}

// NewPermissionRepository 创建一个新的 PermissionRepository
func NewPermissionRepository(l *zap.Logger, dao dao.PermissionDAO) PermissionRepository {
	return &permissionRepository{
		l:   l,
		dao: dao,
	}
}

// GetPermissions 获取指定用户的权限列表
func (r *permissionRepository) GetPermissions(ctx context.Context) ([]domain.Permission, error) {
	permissions, err := r.dao.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// AssignPermission 分配权限给指定用户
func (r *permissionRepository) AssignPermission(ctx context.Context, userName string, path string, method string) error {
	if err := r.dao.AssignPermission(ctx, userName, path, method); err != nil {
		return err
	}
	return nil
}

func (r *permissionRepository) AssignRoleToUser(ctx context.Context, userName, roleName string) error {
	if err := r.dao.AssignRoleToUser(ctx, userName, roleName); err != nil {
		return err
	}
	return nil
}

// RemovePermission 移除指定用户的权限
func (r *permissionRepository) RemovePermission(ctx context.Context, userName string, path string, method string) error {
	if err := r.dao.RemovePermission(ctx, userName, path, method); err != nil {
		return err
	}
	return nil
}

func (r *permissionRepository) RemoveRoleFromUser(ctx context.Context, userName, roleName string) error {
	if err := r.dao.RemoveRoleFromUser(ctx, userName, roleName); err != nil {
		return err
	}
	return nil
}
