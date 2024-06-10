package dao

import (
	"LinkMe/internal/domain"
	"context"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

// PermissionDAO 定义了权限数据库访问接口
type PermissionDAO interface {
	GetPermissions(ctx context.Context, userID int64) ([]domain.Permission, error)
	AssignPermission(ctx context.Context, userID int64, path string, method string) error
	RemovePermission(ctx context.Context, userID int64, path string, method string) error
}

// permissionDAO 是 PermissionDAO 的实现
type permissionDAO struct {
	db *gorm.DB
	ce *casbin.Enforcer
	l  *zap.Logger
}

// NewPermissionDAO 创建一个新的 PermissionDAO
func NewPermissionDAO(ce *casbin.Enforcer, l *zap.Logger, db *gorm.DB) PermissionDAO {
	return &permissionDAO{
		db: db,
		ce: ce,
		l:  l,
	}
}

// GetPermissions 获取指定用户的权限列表
func (d *permissionDAO) GetPermissions(ctx context.Context, userID int64) ([]domain.Permission, error) {
	uid := strconv.FormatInt(userID, 10)
	var policies []domain.Permission
	err := d.db.WithContext(ctx).Table("casbin_rule").Where("v0 = ?", uid).Find(&policies).Error
	if err != nil {
		return nil, err
	}
	return policies, nil
}

// AssignPermission 分配权限给指定用户
func (d *permissionDAO) AssignPermission(ctx context.Context, userID int64, path string, method string) error {
	id := strconv.FormatInt(userID, 10)
	ok, err := d.ce.AddPolicy(id, path, method)
	if err != nil {
		d.l.Error("failed to add policy", zap.Error(err))
		return err
	}
	if !ok {
		d.l.Error("policy already exists", zap.Error(err))
		return err
	}
	return nil
}

// RemovePermission 移除指定用户的权限
func (d *permissionDAO) RemovePermission(ctx context.Context, userID int64, path string, method string) error {
	user := strconv.FormatInt(userID, 10)
	_, err := d.ce.RemovePolicy(user, path, method)
	if err != nil {
		d.l.Error("failed to delete role for user", zap.Error(err))
		return err
	}
	return nil
}
