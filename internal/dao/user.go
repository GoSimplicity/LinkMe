package dao

import (
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

type User struct {
	ID           uint       `gorm:"primarykey"`
	CreateTime   int64      `gorm:"column:created_at;type:bigint"`
	UpdatedTime  int64      `gorm:"column:updated_at;type:bigint"`
	DeletedTime  int64      `gorm:"column:deleted_at;type:bigint;index"`
	Nickname     string     `gorm:"size:50"`
	PasswordHash string     `gorm:"not null"`
	Birthday     *time.Time `gorm:"column:birthday;type:datetime"`
	Email        string     `gorm:"type:varchar(100);uniqueIndex"`
	Phone        *string    `gorm:"type:varchar(15);uniqueIndex"`
}

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
