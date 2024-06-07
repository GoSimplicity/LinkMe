package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"context"
	"go.uber.org/zap"
)

type RoleRepository interface {
	CreateRole(ctx context.Context, r domain.Role) error
	CreatePermission(ctx context.Context, p domain.Permission) error
	AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error
}

type roleRepository struct {
	dao dao.RoleDAO
	l   *zap.Logger
}

func NewRoleRepository(dao dao.RoleDAO, l *zap.Logger) RoleRepository {
	return &roleRepository{
		dao: dao,
		l:   l,
	}
}

func (rr *roleRepository) CreateRole(ctx context.Context, r domain.Role) error {
	return rr.dao.CreateRole(ctx, r)
}

func (rr *roleRepository) CreatePermission(ctx context.Context, p domain.Permission) error {
	return rr.dao.CreatePermission(ctx, p)
}

func (rr *roleRepository) AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error {
	return rr.dao.AssignPermissionToRole(ctx, roleID, permissionID)
}
