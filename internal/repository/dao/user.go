package dao

import (
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
	UpdateProfile(ctx context.Context, profile Profile) error
	GetProfileByUserID(ctx context.Context, UserID int64) (*Profile, error)
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
	err := ud.db.WithContext(ctx).Create(&u).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeDuplicateEmailNumber {
			return ErrDuplicateEmail
		}
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
func (ud *userDAO) UpdateProfile(ctx context.Context, profile Profile) error {
	return ud.db.WithContext(ctx).Model(&Profile{}).Where("user_id = ?", profile.UserID).Save(&profile).Error
}
func (ud *userDAO) GetProfileByUserID(ctx context.Context, UserID int64) (*Profile, error) {
	var profile Profile
	err := ud.db.WithContext(ctx).Where("user_id = ?", UserID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &profile, nil
}
