package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
	"go.uber.org/zap"
)

type RoleService interface {
	CreateRole(ctx context.Context, r domain.Role) error
	CreatePermission(ctx context.Context, p domain.Permission) error
	AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error
}

type roleService struct {
	repo repository.RoleRepository
	l    *zap.Logger
}

func NewRoleService(repo repository.RoleRepository, l *zap.Logger) RoleService {
	return &roleService{
		repo: repo,
		l:    l,
	}
}

func (rs *roleService) CreateRole(ctx context.Context, r domain.Role) error {
	return rs.repo.CreateRole(ctx, r)
}

func (rs *roleService) CreatePermission(ctx context.Context, p domain.Permission) error {
	return rs.repo.CreatePermission(ctx, p)
}

func (rs *roleService) AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error {
	return rs.repo.AssignPermissionToRole(ctx, roleID, permissionID)
}
