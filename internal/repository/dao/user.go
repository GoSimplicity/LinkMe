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
	ErrDuplicateEmail                  = errors.New("邮箱冲突")
	ErrUserNotFound                    = errors.New("用户未找到")
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

// CreateUser 创建用户
func (ud *userDAO) CreateUser(ctx context.Context, u User) error {
	var m *mysql.MySQLError
	u.CreateTime = time.Now().UnixMilli()
	u.UpdatedTime = time.Now().UnixMilli()
	err := ud.db.WithContext(ctx).Create(&u).Error
	if errors.As(err, &m) {
		if m.Number == ErrCodeDuplicateEmailNumber {
			return ErrDuplicateEmail
		}
		return err
	}
	return err
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
