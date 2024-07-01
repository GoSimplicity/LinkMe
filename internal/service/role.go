package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
	"go.uber.org/zap"
)

// PermissionService 定义了权限服务接口
type PermissionService interface {
	GetPermissions(ctx context.Context) ([]domain.Permission, error)                         // 获取权限列表
	AssignPermission(ctx context.Context, userName string, path string, method string) error // 分配权限
	AssignRoleToUser(ctx context.Context, userName, roleName string) error
	RemovePermission(ctx context.Context, userName string, path string, method string) error // 移除权限
	RemoveRoleFromUser(ctx context.Context, userName, roleName string) error
}

// permissionService 是 PermissionService 的实现
type permissionService struct {
	repo repository.PermissionRepository // 权限仓库
	l    *zap.Logger                     // 日志记录器
}

// NewPermissionService 创建一个新的 PermissionService
func NewPermissionService(repo repository.PermissionRepository, l *zap.Logger) PermissionService {
	return &permissionService{
		repo: repo,
		l:    l,
	}
}

// GetPermissions 获取指定用户的权限列表
func (s *permissionService) GetPermissions(ctx context.Context) ([]domain.Permission, error) {
	permissions, err := s.repo.GetPermissions(ctx)
	if err != nil {
		s.l.Error("get permissions failed", zap.Error(err))
		return nil, err
	}
	return permissions, nil
}

// AssignPermission 分配权限给指定用户
func (s *permissionService) AssignPermission(ctx context.Context, userName string, path string, method string) error {
	if err := s.repo.AssignPermission(ctx, userName, path, method); err != nil {
		s.l.Error("assign permissions failed", zap.Error(err))
		return err
	}
	return nil
}

func (s *permissionService) AssignRoleToUser(ctx context.Context, userName, roleName string) error {
	if err := s.repo.AssignRoleToUser(ctx, userName, roleName); err != nil {
		s.l.Error("assign role to user failed", zap.Error(err))
		return err
	}
	return nil
}

// RemovePermission 移除指定用户的权限
func (s *permissionService) RemovePermission(ctx context.Context, userName string, path string, method string) error {
	if err := s.repo.RemovePermission(ctx, userName, path, method); err != nil {
		s.l.Error("remove permissions failed", zap.Error(err))
		return err
	}
	return nil
}

func (s *permissionService) RemoveRoleFromUser(ctx context.Context, userName, roleName string) error {
	if err := s.repo.RemoveRoleFromUser(ctx, userName, roleName); err != nil {
		s.l.Error("remove role from user failed", zap.Error(err))
		return err
	}
	return nil
}
