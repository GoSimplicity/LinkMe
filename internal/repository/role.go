package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"context"
	"go.uber.org/zap"
)

// PermissionRepository 定义了权限仓库接口
type PermissionRepository interface {
	GetPermissions(ctx context.Context, userID int64) ([]domain.Permission, error)        // 获取权限列表
	AssignPermission(ctx context.Context, userID int64, path string, method string) error // 分配权限
	RemovePermission(ctx context.Context, userID int64, path string, method string) error // 移除权限
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
func (r *permissionRepository) GetPermissions(ctx context.Context, userID int64) ([]domain.Permission, error) {
	permissions, err := r.dao.GetPermissions(ctx, userID)
	if err != nil {
		r.l.Error("获取权限失败", zap.Error(err))
		return nil, err
	}
	return permissions, nil
}

// AssignPermission 分配权限给指定用户
func (r *permissionRepository) AssignPermission(ctx context.Context, userID int64, path string, method string) error {
	if err := r.dao.AssignPermission(ctx, userID, path, method); err != nil {
		r.l.Error("分配权限失败", zap.Error(err))
		return err
	}
	return nil
}

// RemovePermission 移除指定用户的权限
func (r *permissionRepository) RemovePermission(ctx context.Context, userID int64, path string, method string) error {
	if err := r.dao.RemovePermission(ctx, userID, path, method); err != nil {
		r.l.Error("移除权限失败", zap.Error(err))
		return err
	}
	return nil
}
