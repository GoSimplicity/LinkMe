package repository

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

var (
	// ErrDuplicateEmail 表示邮箱重复错误
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	// ErrUserNotFound 表示用户未找到错误
	ErrUserNotFound = dao.ErrUserNotFound
)

// UserRepository 定义用户仓库接口
type UserRepository interface {
	// CreateUser 创建用户
	CreateUser(ctx context.Context, u domain.User) error
	// FindByID 根据用户ID查找用户
	FindByID(ctx context.Context, id int64) (domain.User, error)
	// FindByPhone 根据手机号查找用户
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, email string, newPassword string) error
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, email string, uid int64) error
	// UpdateProfile 更新用户资料
	UpdateProfile(ctx context.Context, profile domain.Profile) error
	// GetProfile 根据用户ID获取用户资料
	GetProfile(ctx context.Context, UserId int64) (domain.Profile, error)
	// ListUser 获取用户列表
	ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfileAndRule, error)
	// GetUserCount 获取用户总数
	GetUserCount(ctx context.Context) (int64, error)
}

type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
	l     *zap.Logger
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache, l *zap.Logger) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

// ChangePassword 修改密码
func (ur *userRepository) ChangePassword(ctx context.Context, email string, newPassword string) error {
	return ur.dao.UpdatePasswordByEmail(ctx, email, newPassword)
}

// CreateUser 创建用户
func (ur *userRepository) CreateUser(ctx context.Context, u domain.User) error {
	return ur.dao.CreateUser(ctx, fromDomainUser(u))
}

// FindByID 通过ID查询用户
func (ur *userRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	// 尝试从缓存中获取用户
	du, err := ur.cache.Get(ctx, id)
	if err == nil {
		// 如果在缓存中找到，直接返回
		return du, nil
	}
	// 如果缓存中未找到，从数据库中查找
	u, err := ur.dao.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	du = toDomainUser(u)
	// 异步将用户信息写入缓存
	go func() {
		if setErr := ur.cache.Set(ctx, du); setErr != nil {
			ur.l.Error("set cache failed", zap.Error(setErr))
		}
	}()
	return du, nil
}

// FindByPhone 通过电话查询用户
func (ur *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := ur.dao.FindByPhone(ctx, phone)
	return toDomainUser(u), err
}

// FindByEmail 通过Email查询用户
func (ur *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	return toDomainUser(u), err
}

// DeleteUser 删除用户
func (ur *userRepository) DeleteUser(ctx context.Context, email string, uid int64) error {
	return ur.dao.DeleteUser(ctx, email, uid)
}

// UpdateProfile 更新用户资料
func (ur *userRepository) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	return ur.dao.UpdateProfile(ctx, profile)
}

// GetProfile 通过用户ID获取用户资料
func (ur *userRepository) GetProfile(ctx context.Context, UserId int64) (domain.Profile, error) {
	return ur.dao.GetProfileByUserID(ctx, UserId)
}

// ListUser 获取用户列表
func (ur *userRepository) ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfileAndRule, error) {
	users, err := ur.dao.ListUser(ctx, pagination)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserCount 获取用户总数
func (ur *userRepository) GetUserCount(ctx context.Context) (int64, error) {
	count, err := ur.dao.GetUserCount(ctx)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// 将领域层对象转为dao层对象
func fromDomainUser(u domain.User) dao.User {
	return dao.User{
		ID:           u.ID,
		PasswordHash: u.Password,
		Email:        u.Email,
		Phone:        u.Phone,
		CreateTime:   u.CreateTime,
	}
}

// 将dao层对象转为领域层对象
func toDomainUser(u dao.User) domain.User {
	return domain.User{
		ID:         u.ID,
		Password:   u.PasswordHash,
		Email:      u.Email,
		Phone:      u.Phone,
		CreateTime: u.CreateTime,
	}
}
