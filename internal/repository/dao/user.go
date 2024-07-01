package dao

import (
	"LinkMe/internal/domain"
	. "LinkMe/internal/repository/models"
	"context"
	"errors"
	sf "github.com/bwmarrin/snowflake"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

var (
	ErrCodeDuplicateEmailNumber uint16 = 1062
	ErrDuplicateEmail                  = errors.New("duplicate email")
	ErrUserNotFound                    = errors.New("user not found")
)

type UserDAO interface {
	CreateUser(ctx context.Context, u User) error
	FindByID(ctx context.Context, id int64) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	UpdatePasswordByEmail(ctx context.Context, email string, newPassword string) error
	DeleteUser(ctx context.Context, email string, uid int64) error
	UpdateProfile(ctx context.Context, profile domain.Profile) error
	GetProfileByUserID(ctx context.Context, userId int64) (domain.Profile, error)
	GetAllUser(ctx context.Context) ([]domain.UserWithProfile, error)
}

type userDAO struct {
	db   *gorm.DB
	node *sf.Node
	l    *zap.Logger
}

func NewUserDAO(db *gorm.DB, node *sf.Node, l *zap.Logger) UserDAO {
	return &userDAO{
		db:   db,
		node: node,
		l:    l,
	}
}

// 获取当前时间的时间戳
func (ud *userDAO) currentTime() int64 {
	return time.Now().UnixMilli()
}

// CreateUser 创建用户
func (ud *userDAO) CreateUser(ctx context.Context, u User) error {
	u.CreateTime = ud.currentTime()
	u.UpdatedTime = u.CreateTime
	// 使用雪花算法生成id
	u.ID = ud.node.Generate().Int64()
	// 初始化用户资料
	profile := Profile{
		UserID:   u.ID,
		NickName: "",
		Avatar:   "",
		About:    "",
		Birthday: "",
	}
	// 开始事务
	tx := ud.db.WithContext(ctx).Begin()
	// 创建用户
	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeDuplicateEmailNumber {
			ud.l.Error("duplicate email error", zap.String("email", u.Email), zap.Error(err))
			return ErrDuplicateEmail
		}
		ud.l.Error("failed to create user", zap.Error(err))
		return err
	}
	// 创建用户资料
	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		ud.l.Error("failed to create profile", zap.Error(err))
		return err
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ud.l.Error("transaction commit failed", zap.Error(err))
		return err
	}
	return nil
}

// FindByID 根据ID查询用户数据
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

// FindByEmail 根据Email查询用户信息
func (ud *userDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("email = ? AND deleted = ?", email, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

// FindByPhone 根据phone查询用户信息
func (ud *userDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("phone = ? AND deleted = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

func (ud *userDAO) UpdatePasswordByEmail(ctx context.Context, email string, newPassword string) error {
	tx := ud.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		ud.l.Error("failed to begin transaction", zap.Error(tx.Error))
		return tx.Error
	}
	// 更新密码
	if err := tx.Model(&User{}).Where("email = ? AND deleted = ?", email, false).Update("password_hash", newPassword).Error; err != nil {
		ud.l.Error("update password failed", zap.String("email", email), zap.Error(err))
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			ud.l.Error("failed to rollback transaction", zap.Error(rollbackErr))
		}
		return err
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ud.l.Error("failed to commit transaction", zap.Error(err))
		return err
	}
	ud.l.Info("password updated successfully", zap.String("email", email))
	return nil
}

func (ud *userDAO) DeleteUser(ctx context.Context, email string, uid int64) error {
	tx := ud.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	if err := tx.Model(&User{}).Where("email = ? AND deleted = ? AND id = ?", email, false, uid).Update("deleted", true).Error; err != nil {
		tx.Rollback()
		ud.l.Error("failed to mark user as deleted", zap.String("email", email), zap.Error(err))
		return err
	}
	if err := tx.Commit().Error; err != nil {
		ud.l.Error("failed to commit transaction", zap.String("email", email), zap.Error(err))
		return err
	}
	ud.l.Info("user marked as deleted", zap.String("email", email))
	return nil
}

// UpdateProfile 更新用户资料
func (ud *userDAO) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	// 创建一个更新用的结构体
	updates := domain.Profile{
		NickName: profile.NickName,
		Avatar:   profile.Avatar,
		About:    profile.About,
		Birthday: profile.Birthday,
	}
	// 更新操作
	err := ud.db.WithContext(ctx).Model(&Profile{}).Where("user_id = ?", profile.UserID).Updates(updates).Error
	if err != nil {
		ud.l.Error("failed to update profile", zap.Error(err))
		return err
	}
	return nil
}

func (ud *userDAO) GetProfileByUserID(ctx context.Context, userId int64) (domain.Profile, error) {
	var profile domain.Profile
	if err := ud.db.WithContext(ctx).Where("user_id = ?", userId).First(&profile).Error; err != nil {
		ud.l.Error("failed to get profile by user id", zap.Error(err))
		return domain.Profile{}, err
	}
	return profile, nil
}

func (ud *userDAO) GetAllUser(ctx context.Context) ([]domain.UserWithProfile, error) {
	var usersWithProfiles []domain.UserWithProfile
	// 执行连接查询
	if err := ud.db.WithContext(ctx).
		Table("users").
		Select(`users.id, users.password_hash, users.deleted, users.email, users.phone,
				profiles.id as profile_id, profiles.user_id, profiles.nick_name, profiles.avatar, profiles.about, profiles.birthday`).
		Joins("left join profiles on profiles.user_id = users.id").
		Scan(&usersWithProfiles).Error; err != nil {
		ud.l.Error("failed to get all users with profiles", zap.Error(err))
		return nil, err
	}
	return usersWithProfiles, nil
}
