package dao

import (
	"context"
	"errors"
	"strconv"
	"strings"
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
	ErrDuplicateUsername = errors.New("duplicate username")
	// ErrUserNotFound 表示用户未找到错误
	ErrUserNotFound = errors.New("user not found")
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
	ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfileAndRule, error)
	GetUserCount(ctx context.Context) (int64, error)
}

type userDAO struct {
	db   *gorm.DB
	node *sf.Node
	l    *zap.Logger
	ce   *casbin.Enforcer
}

// User 用户结构体
type User struct {
	ID           int64   `gorm:"primarykey"`
	CreateTime   int64   `gorm:"column:created_at;type:bigint"`
	UpdatedTime  int64   `gorm:"column:updated_at;type:bigint"`
	DeletedTime  int64   `gorm:"column:deleted_at;type:bigint;index"`
	Username     string  `gorm:"column:username;type:varchar(100);uniqueIndex"`
	PasswordHash string  `gorm:"not null"`
	Deleted      bool    `gorm:"column:deleted;default:false"`
	Email        string  `gorm:"type:varchar(100)"`
	Phone        *string `gorm:"type:varchar(15);uniqueIndex"`
	Profile      Profile `gorm:"foreignKey:UserID;references:ID"`
}

// Profile 用户资料信息结构体
type Profile struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	UserID   int64  `gorm:"not null;index"`
	NickName string `gorm:"size:50"`
	Avatar   string `gorm:"type:text"`
	About    string `gorm:"type:text"`
	Birthday string `gorm:"column:birthday;type:varchar(10)"`
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

func (ud *userDAO) CreateUser(ctx context.Context, u User) error {
	now := ud.currentTime()
	u.CreateTime = now
	u.UpdatedTime = now
	u.ID = ud.node.Generate().Int64()

	profile := Profile{
		UserID:   u.ID,
		NickName: "",
		Avatar:   "",
		About:    "",
		Birthday: "",
	}

	tx := ud.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeDuplicateUsernameNumber {
			ud.l.Error("duplicate username error", zap.String("username", u.Username), zap.Error(err))
			return ErrDuplicateUsername
		}
		ud.l.Error("failed to create user", zap.Error(err))
		return err
	}

	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		ud.l.Error("failed to create profile", zap.Error(err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		ud.l.Error("transaction commit failed", zap.Error(err))
		return err
	}

	return nil
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
		ud.l.Error("update password failed", zap.String("username", username), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (ud *userDAO) DeleteUser(ctx context.Context, username string, uid int64) error {
	tx := ud.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	result := tx.Model(&User{}).
		Where("username = ? AND deleted = ? AND id = ?", username, false, uid).
		Update("deleted", true)

	if result.Error != nil {
		tx.Rollback()
		ud.l.Error("failed to mark user as deleted", zap.String("username", username), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return ErrUserNotFound
	}

	if err := tx.Commit().Error; err != nil {
		ud.l.Error("failed to commit transaction", zap.String("username", username), zap.Error(err))
		return err
	}

	return nil
}

func (ud *userDAO) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	updates := domain.Profile{
		NickName: profile.NickName,
		Avatar:   profile.Avatar,
		About:    profile.About,
		Birthday: profile.Birthday,
	}

	result := ud.db.WithContext(ctx).Model(&Profile{}).
		Where("user_id = ?", profile.UserID).
		Updates(updates)

	if result.Error != nil {
		ud.l.Error("failed to update profile", zap.Error(result.Error))
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
		ud.l.Error("failed to get profile by user id", zap.Error(err))
		return domain.Profile{}, err
	}
	return profile, nil
}

func (ud *userDAO) ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfileAndRule, error) {
	var usersWithProfiles []domain.UserWithProfileAndRule
	intSize := int(*pagination.Size)
	intOffset := int(*pagination.Offset)

	err := ud.db.WithContext(ctx).
		Table("users").
		Select(`users.id, users.password_hash, users.deleted, users.username, users.phone,
                profiles.id as profile_id, profiles.user_id, profiles.nick_name, profiles.avatar, profiles.about, profiles.birthday`).
		Joins("left join profiles on profiles.user_id = users.id").
		Where("users.deleted = ?", false).
		Limit(intSize).
		Offset(intOffset).
		Scan(&usersWithProfiles).Error

	if err != nil {
		ud.l.Error("failed to get all users with profiles", zap.Error(err))
		return nil, err
	}

	for i, user := range usersWithProfiles {
		roleUsernames, err := ud.getUserRoleUsernames(ctx, user.ID)
		if err != nil {
			ud.l.Error("failed to get role usernames for user", zap.Int64("userID", user.ID), zap.Error(err))
			return nil, err
		}
		if len(roleUsernames) > 0 {
			usersWithProfiles[i].Role = strings.Join(roleUsernames, ",")
		}
	}

	return usersWithProfiles, nil
}

func (ud *userDAO) GetUserCount(ctx context.Context) (int64, error) {
	var count int64
	err := ud.db.WithContext(ctx).Model(&User{}).Where("deleted = ?", false).Count(&count).Error
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (ud *userDAO) getUserRoleUsernames(ctx context.Context, userID int64) ([]string, error) {
	userIDStr := strconv.FormatInt(userID, 10)
	roles, err := ud.ce.GetRolesForUser(userIDStr)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return nil, nil
	}

	roleIDs := make([]int64, 0, len(roles))
	for _, roleIDStr := range roles {
		roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
		if err != nil {
			return nil, err
		}
		roleIDs = append(roleIDs, roleID)
	}

	var roleUsers []struct {
		Username string
	}
	err = ud.db.WithContext(ctx).
		Table("users").
		Select("username").
		Where("id IN (?) AND deleted = ?", roleIDs, false).
		Scan(&roleUsers).Error
	if err != nil {
		return nil, err
	}

	roleUsernames := make([]string, 0, len(roleUsers))
	for _, roleUser := range roleUsers {
		roleUsernames = append(roleUsernames, roleUser.Username)
	}

	return roleUsernames, nil
}
