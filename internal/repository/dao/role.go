package dao

import (
	. "LinkMe/internal/repository/models"
	"context"
	"gorm.io/gorm"
)

type RoleDAO interface {
	CreateRole(ctx context.Context, r Role) error
	CreatePermission(ctx context.Context, p Permission) error
	AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error
}

type roleDAO struct {
	db *gorm.DB
}

func NewRoleDAO(db *gorm.DB) RoleDAO {
	return &roleDAO{
		db: db,
	}
}

func (rd *roleDAO) CreateRole(ctx context.Context, r Role) error {
	return rd.db.WithContext(ctx).Create(&r).Error
}

func (rd *roleDAO) CreatePermission(ctx context.Context, p Permission) error {
	return rd.db.WithContext(ctx).Create(&p).Error
}

func (rd *roleDAO) AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error {
	rolePermission := RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return rd.db.WithContext(ctx).Create(&rolePermission).Error
}
