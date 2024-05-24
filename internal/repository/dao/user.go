package dao

import (
	. "LinkMe/internal/repository/models"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
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
}

type userDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &userDAO{
		db: db,
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
	err := ud.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
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
	err := ud.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}
