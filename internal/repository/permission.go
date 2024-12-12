package repository

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

type PermissionRepository interface {
	AssignRole(ctx context.Context, roleId int, menuIds []int, apiIds []int) error
	AssignRoleToUser(ctx context.Context, userId int, roleIds []int, menuIds []int, apiIds []int) error
	AssignRoleToUsers(ctx context.Context, userIds []int, roleIds []int, menuIds []int, apiIds []int) error
	RemoveUserPermissions(ctx context.Context, userId int) error
	RemoveRolePermissions(ctx context.Context, roleId int) error
	RemoveUsersPermissions(ctx context.Context, userIds []int) error
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

// AssignRole 实现角色权限分配
func (p *permissionRepository) AssignRole(ctx context.Context, roleId int, menuIds []int, apiIds []int) error {
	return p.dao.AssignRole(ctx, roleId, menuIds, apiIds)
}

// AssignRoleToUser 实现用户角色权限分配
func (p *permissionRepository) AssignRoleToUser(ctx context.Context, userId int, roleIds []int, menuIds []int, apiIds []int) error {
	return p.dao.AssignRoleToUser(ctx, userId, roleIds, menuIds, apiIds)
}

// RemoveRolePermissions 实现角色权限移除
func (p *permissionRepository) RemoveRolePermissions(ctx context.Context, roleId int) error {
	return p.dao.RemoveRolePermissions(ctx, roleId)
}

// RemoveUserPermissions 实现用户权限移除
func (p *permissionRepository) RemoveUserPermissions(ctx context.Context, userId int) error {
	return p.dao.RemoveUserPermissions(ctx, userId)
}

// AssignRoleToUsers 实现批量用户角色权限分配
func (p *permissionRepository) AssignRoleToUsers(ctx context.Context, userIds []int, roleIds []int, menuIds []int, apiIds []int) error {
	return p.dao.AssignRoleToUsers(ctx, userIds, roleIds, menuIds, apiIds)
}

// RemoveUsersPermissions 实现批量用户权限移除
func (p *permissionRepository) RemoveUsersPermissions(ctx context.Context, userIds []int) error {
	return p.dao.RemoveUsersPermissions(ctx, userIds)
}
