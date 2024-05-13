package dao

import (
	"LinkMe/internal/domain"
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO interface {
	CreateUser(ctx context.Context, u domain.User) error
	FindByID(ctx context.Context, id int64) (domain.User, error)
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
	gorm.Model
	Username     string    `gorm:"varchar(255);uniqueIndex(255)"` // 用户名，数据库中唯一索引
	PasswordHash string    `json:"-"`                             // 密码的哈希值
	Nickname     string    `gorm:"size:50"`                       // 昵称，限制长度为50字符
	Birthday     time.Time `gorm:"column:birthday;type:datetime"`
	Email        *string   `gorm:"type:varchar(100);uniqueIndex"` // 邮箱，可为空，唯一索引
	Phone        *string   `gorm:"type:varchar(15);uniqueIndex"`  // 手机号，可为空，唯一索引，最大长度为15位
}

func (u2 *userDAO) CreateUser(ctx context.Context, u domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (u2 *userDAO) FindByID(ctx context.Context, id int64) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}
