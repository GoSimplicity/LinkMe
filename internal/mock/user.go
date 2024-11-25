package mock

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type MockUserRepository interface {
	MockUser() error
}

type mockUserRepository struct {
	db *gorm.DB
	ce *casbin.Enforcer
	l  *zap.Logger
}

func NewMockUserRepository(db *gorm.DB, l *zap.Logger, ce *casbin.Enforcer) MockUserRepository {
	return &mockUserRepository{
		db: db,
		ce: ce,
		l:  l,
	}
}

// User 用户信息结构体
type User struct {
	ID           int64  `gorm:"primarykey"`                          // 用户ID，主键
	CreateTime   int64  `gorm:"column:created_at;type:bigint"`       // 创建时间，Unix时间戳
	UpdatedTime  int64  `gorm:"column:updated_at;type:bigint"`       // 更新时间，Unix时间戳
	DeletedTime  int64  `gorm:"column:deleted_at;type:bigint;index"` // 删除时间，Unix时间戳，用于软删除
	Username     string `gorm:"type:varchar(100);uniqueIndex"`
	PasswordHash string `gorm:"not null"` // 密码哈希值，不能为空
}

func (m *mockUserRepository) MockUser() error {
	var existingUser User

	if err := m.db.Where("username = ?", "admin").First(&existingUser).Error; err == nil {
		m.l.Info("user already exists, skipping creation", zap.String("email", "admin"))
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		m.l.Error("failed to query user", zap.Error(err))
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user := User{
		Username:     "admin",
		PasswordHash: string(hash),
	}
	user.ID = 1
	user.CreateTime = time.Now().Unix()
	user.UpdatedTime = time.Now().Unix()

	if err := m.db.Create(&user).Error; err != nil {
		m.l.Error("failed to create user", zap.Error(err))
		return err
	}

	// 将 userID 转换为字符串
	userIDStr := strconv.FormatInt(user.ID, 10)

	// 定义所有路径和方法的集合（可根据你的应用实际调整）
	paths := []string{"/*"} // 假设 /* 表示所有路径
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	// 为用户添加所有路径和方法的策略
	for _, path := range paths {
		for _, method := range methods {
			ok, err := m.ce.AddPolicy(userIDStr, path, method)
			if err != nil {
				m.l.Error("failed to add policy", zap.Error(err))
				return err
			}

			if !ok {
				m.l.Error("policy already exists", zap.Error(err))
				return fmt.Errorf("policy already exists for user %d, path %s, method %s", user.ID, path, method)
			}
		}
	}

	return nil
}
