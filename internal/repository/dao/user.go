package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	sf "github.com/bwmarrin/snowflake"
	"github.com/casbin/casbin/v2"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// ErrCodeDuplicateUsernameNumber 表示用户名重复的错误码
	ErrCodeDuplicateUsernameNumber uint16 = 1062
	// ErrDuplicateUsername 表示用户名重复错误
	ErrDuplicateUsername = errors.New("用户名已存在")
	// ErrUserNotFound 表示用户未找到错误
	ErrUserNotFound = errors.New("用户不存在")
)

type UserDAO interface {
	CreateUser(ctx context.Context, u User) error
	FindByID(ctx context.Context, id int64) (User, error)
	FindByUsername(ctx context.Context, username string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	UpdatePasswordByUsername(ctx context.Context, username string, newPassword string) error
	DeleteUser(ctx context.Context, username string, uid int64) error
	UpdateProfile(ctx context.Context, profile domain.Profile) error
	GetProfileByUserID(ctx context.Context, userId int64) (domain.Profile, error)
	ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfile, error)
	UpdateProfileAdmin(ctx context.Context, profile domain.Profile) error
}

type userDAO struct {
	db   *gorm.DB
	node *sf.Node
	l    *zap.Logger
	ce   *casbin.Enforcer
}

// User 用户模型
type User struct {
	ID           int64   `gorm:"primarykey"`
	CreateTime   int64   `gorm:"column:created_at;type:bigint;not null"`
	UpdatedTime  int64   `gorm:"column:updated_at;type:bigint;not null"`
	DeletedTime  int64   `gorm:"column:deleted_at;type:bigint;index"`
	Username     string  `gorm:"column:username;type:varchar(100);uniqueIndex;not null"`
	PasswordHash string  `gorm:"not null"`
	Deleted      bool    `gorm:"column:deleted;default:false;not null"`
	Profile      Profile `gorm:"foreignKey:UserID;references:ID"`
	Roles        string  `gorm:"column:roles;type:json;comment:用户角色ID列表"`
}

// Profile 用户资料信息模型
type Profile struct {
	ID       int64   `gorm:"primaryKey;autoIncrement"`
	UserID   int64   `gorm:"not null;index"`
	RealName string  `gorm:"size:50"`
	Avatar   string  `gorm:"type:text"`
	About    string  `gorm:"type:text"`
	Birthday string  `gorm:"column:birthday;type:varchar(10)"`
	Email    string  `gorm:"type:varchar(100)"`
	Phone    *string `gorm:"type:varchar(15);uniqueIndex"`
}

func NewUserDAO(db *gorm.DB, node *sf.Node, l *zap.Logger, ce *casbin.Enforcer) UserDAO {
	return &userDAO{
		db:   db,
		node: node,
		l:    l,
		ce:   ce,
	}
}

func (ud *userDAO) currentTime() int64 {
	return time.Now().UnixMilli()
}

// CreateUser 创建用户
func (ud *userDAO) CreateUser(ctx context.Context, u User) error {
	now := ud.currentTime()
	u.CreateTime = now
	u.UpdatedTime = now
	u.ID = ud.node.Generate().Int64()

	if u.Roles == "" {
		u.Roles = "[]"
	}

	profile := Profile{
		UserID:   u.ID,
		RealName: "",
		Avatar:   "",
		About:    "",
		Birthday: "",
	}

	err := ud.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&u).Error; err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeDuplicateUsernameNumber {
				ud.l.Error("用户名重复错误", zap.String("username", u.Username), zap.Error(err))
				return ErrDuplicateUsername
			}
			ud.l.Error("创建用户失败", zap.Error(err))
			return err
		}

		if err := tx.Create(&profile).Error; err != nil {
			ud.l.Error("创建用户资料失败", zap.Error(err))
			return err
		}

		return nil
	})

	return err
}

func (ud *userDAO) FindByID(ctx context.Context, id int64) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("id = ? AND deleted = ?", id, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

func (ud *userDAO) FindByUsername(ctx context.Context, username string) (User, error) {
	var user User

	err := ud.db.WithContext(ctx).Where("username = ? AND deleted = ?", username, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, nil
}

func (ud *userDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("phone = ? AND deleted = ?", phone, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

func (ud *userDAO) UpdatePasswordByUsername(ctx context.Context, username string, newPassword string) error {
	result := ud.db.WithContext(ctx).Model(&User{}).
		Where("username = ? AND deleted = ?", username, false).
		Update("password_hash", newPassword)

	if result.Error != nil {
		ud.l.Error("更新密码失败", zap.String("username", username), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (ud *userDAO) DeleteUser(ctx context.Context, username string, uid int64) error {
	err := ud.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&User{}).
			Where("username = ? AND deleted = ? AND id = ?", username, false, uid).
			Update("deleted", true)

		if result.Error != nil {
			ud.l.Error("标记用户删除失败", zap.String("username", username), zap.Error(result.Error))
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ErrUserNotFound
		}

		return nil
	})

	return err
}

func (ud *userDAO) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	updates := domain.Profile{
		RealName: profile.RealName,
		Avatar:   profile.Avatar,
		About:    profile.About,
		Birthday: profile.Birthday,
		Phone:    profile.Phone,
	}

	result := ud.db.WithContext(ctx).Model(&Profile{}).
		Where("user_id = ?", profile.UserID).
		Updates(updates)

	if result.Error != nil {
		ud.l.Error("更新用户资料失败", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (ud *userDAO) GetProfileByUserID(ctx context.Context, userId int64) (domain.Profile, error) {
	var profile domain.Profile
	err := ud.db.WithContext(ctx).Where("user_id = ?", userId).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Profile{}, ErrUserNotFound
		}
		ud.l.Error("获取用户资料失败", zap.Error(err))
		return domain.Profile{}, err
	}
	return profile, nil
}

func (ud *userDAO) ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfile, error) {
	if pagination.Size == nil || pagination.Offset == nil {
		return nil, errors.New("分页参数不能为空")
	}

	intSize := int(*pagination.Size)
	intOffset := int(*pagination.Offset)

	// 定义查询结果结构
	type userResult struct {
		ID           int64   `json:"id"`
		PasswordHash string  `json:"password_hash"`
		Deleted      bool    `json:"deleted"`
		Username     string  `json:"username"`
		Phone        *string `json:"phone"`
		ProfileID    int64   `json:"profile_id"`
		UserID       int64   `json:"user_id"`
		RealName     string  `json:"real_name"`
		Avatar       string  `json:"avatar"`
		About        string  `json:"about"`
		Birthday     string  `json:"birthday"`
		Roles        string  `json:"roles"`
	}

	var results []userResult

	// 构建查询
	query := ud.db.WithContext(ctx).
		Table("users").
		Select(`users.id, users.password_hash, users.deleted, users.username, users.phone,
                profiles.id as profile_id, profiles.user_id, profiles.real_name, 
                profiles.avatar, profiles.about, profiles.birthday, users.roles`).
		Joins("LEFT JOIN profiles ON profiles.user_id = users.id").
		Where("users.deleted = ?", false).
		Limit(intSize).
		Offset(intOffset)

	// 执行查询
	if err := query.Scan(&results).Error; err != nil {
		ud.l.Error("获取用户列表失败", zap.Error(err))
		return nil, fmt.Errorf("获取用户列表失败: %v", err)
	}

	usersWithProfiles := make([]domain.UserWithProfile, 0, len(results))

	// 处理结果
	for _, r := range results {
		var roleIds []int
		if r.Roles != "" {
			if err := json.Unmarshal([]byte(r.Roles), &roleIds); err != nil {
				ud.l.Error("解析用户角色失败", zap.Error(err))
				return nil, fmt.Errorf("解析用户角色失败: %v", err)
			}
		}

		// 查询角色名称
		var roleNames []string
		if len(roleIds) > 0 {
			var roles []*Role
			if err := ud.db.WithContext(ctx).Where("id IN ? AND is_deleted = ?", roleIds, 0).Find(&roles).Error; err != nil {
				ud.l.Error("获取角色信息失败", zap.Error(err))
				return nil, fmt.Errorf("获取角色信息失败: %v", err)
			}

			for _, role := range roles {
				roleNames = append(roleNames, role.Name)
			}
		}

		usersWithProfiles = append(usersWithProfiles, domain.UserWithProfile{
			ID:           r.ID,
			PasswordHash: r.PasswordHash,
			Deleted:      r.Deleted,
			Username:     r.Username,
			Phone:        r.Phone,
			ProfileID:    r.ProfileID,
			UserID:       r.UserID,
			RealName:     r.RealName,
			Avatar:       r.Avatar,
			About:        r.About,
			Birthday:     r.Birthday,
			Roles:        roleNames,
		})
	}

	return usersWithProfiles, nil
}

// UpdateProfileAdmin implements UserDAO.
func (ud *userDAO) UpdateProfileAdmin(ctx context.Context, profile domain.Profile) error {
	// 更新用户资料
	if err := ud.db.WithContext(ctx).Model(&Profile{}).Where("user_id = ?", profile.UserID).Updates(map[string]interface{}{
		"real_name": profile.RealName,
		"avatar":    profile.Avatar,
		"about":     profile.About,
		"birthday":  profile.Birthday,
	}).Error; err != nil {
		ud.l.Error("更新用户资料失败", zap.Error(err))
		return fmt.Errorf("更新用户资料失败: %v", err)
	}

	return nil
}
