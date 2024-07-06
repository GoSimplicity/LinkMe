package dao

import (
	"LinkMe/internal/domain"
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

// PermissionDAO 定义了权限数据库访问接口
type PermissionDAO interface {
	GetPermissions(ctx context.Context) ([]domain.Permission, error)
	AssignPermission(ctx context.Context, userName string, path string, method string) error
	AssignRoleToUser(ctx context.Context, userName, roleName string) error
	RemovePermission(ctx context.Context, userName string, path string, method string) error
	RemoveRoleFromUser(ctx context.Context, userName, roleName string) error
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
func (d *permissionDAO) GetPermissions(ctx context.Context) ([]domain.Permission, error) {
	var policies []domain.Permission
	err := d.db.WithContext(ctx).Table("casbin_rule").Find(&policies).Error
	if err != nil {
		return nil, err
	}
	return policies, nil
}

// AssignPermission 分配权限给指定用户
func (d *permissionDAO) AssignPermission(ctx context.Context, userName, path, method string) error {
	userID, err := d.getUserIDByEmail(ctx, userName)
	if err != nil {
		return err
	}
	// 将 userID 转换为字符串
	userIDStr := strconv.FormatInt(userID, 10)
	ok, err := d.ce.AddPolicy(userIDStr, path, method)
	if err != nil {
		d.l.Error("failed to add policy", zap.Error(err))
		return err
	}
	if !ok {
		d.l.Error("policy already exists", zap.Error(err))
		return fmt.Errorf("policy already exists for user %d, path %s, method %s", userID, path, method)
	}
	return nil
}

// AssignRoleToUser 分配角色给指定用户
func (d *permissionDAO) AssignRoleToUser(ctx context.Context, userName, roleName string) error {
	userID, err := d.getUserIDByEmail(ctx, userName)
	if err != nil {
		d.l.Error("failed to get user ID", zap.Error(err))
		return err
	}
	roleID, err := d.getUserIDByEmail(ctx, roleName)
	if err != nil {
		d.l.Error("failed to get user ID", zap.Error(err))
		return err
	}
	// 将 userID 转换为字符串
	userIDStr := strconv.FormatInt(userID, 10)
	roleIDStr := strconv.FormatInt(roleID, 10)
	// 分配角色给用户
	ok, err := d.ce.AddGroupingPolicy(userIDStr, roleIDStr)
	if err != nil {
		d.l.Error("failed to add role to user", zap.Error(err))
		return err
	}
	if !ok {
		d.l.Error("role already assigned to user", zap.Error(err))
		return fmt.Errorf("role %s already assigned to user %s", roleName, userName)
	}
	return nil
}

// RemovePermission 移除指定用户的权限
func (d *permissionDAO) RemovePermission(ctx context.Context, userName, path, method string) error {
	userID, err := d.getUserIDByEmail(ctx, userName)
	if err != nil {
		return err
	}
	// 将 userID 转换为字符串
	userIDStr := strconv.FormatInt(userID, 10)
	ok, err := d.ce.RemovePolicy(userIDStr, path, method)
	if err != nil {
		d.l.Error("failed to remove policy", zap.Error(err))
		return err
	}
	if !ok {
		d.l.Error("policy does not exist", zap.Error(err))
		return fmt.Errorf("policy does not exist for user %d, path %s, method %s", userID, path, method)
	}
	return nil
}

// RemoveRoleFromUser 移除指定用户的角色
func (d *permissionDAO) RemoveRoleFromUser(ctx context.Context, userName, roleName string) error {
	userID, err := d.getUserIDByEmail(ctx, userName)
	if err != nil {
		d.l.Error("failed to get user ID", zap.Error(err))
		return err
	}
	roleID, err := d.getUserIDByEmail(ctx, roleName)
	if err != nil {
		d.l.Error("failed to get role ID", zap.Error(err))
		return err
	}
	// 将 userID 和 roleID 转换为字符串
	userIDStr := strconv.FormatInt(userID, 10)
	roleIDStr := strconv.FormatInt(roleID, 10)
	// 移除用户的角色
	ok, err := d.ce.RemoveGroupingPolicy(userIDStr, roleIDStr)
	if err != nil {
		d.l.Error("failed to remove role from user", zap.Error(err))
		return err
	}
	if !ok {
		d.l.Error("role not assigned to user", zap.Error(err))
		return fmt.Errorf("role %s not assigned to user %s", roleName, userName)
	}
	return nil
}

// getUserIDByEmail 根据用户的邮箱获取用户ID
func (d *permissionDAO) getUserIDByEmail(ctx context.Context, email string) (int64, error) {
	var user User
	if err := d.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Select("id").First(&user).Error; err != nil {
		d.l.Error("failed to get user by email", zap.Error(err))
		return 0, err
	}
	return user.ID, nil
}
