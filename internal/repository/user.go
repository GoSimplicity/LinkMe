package repository

import (
	"context"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

var (
	// ErrDuplicateUsername 表示用户名重复错误
	ErrDuplicateUsername = dao.ErrDuplicateUsername
)

type UserRepository interface {
	CreateUser(ctx context.Context, u domain.User) error
	FindByID(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByUsername(ctx context.Context, username string) (domain.User, error)
	ChangePassword(ctx context.Context, username string, newPassword string) error
	DeleteUser(ctx context.Context, username string, uid int64) error
	UpdateProfile(ctx context.Context, profile domain.Profile) error
	GetProfile(ctx context.Context, UserID int64) (domain.Profile, error)
	ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfile, error)
}

type userRepository struct {
	l     *zap.Logger
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache, l *zap.Logger) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

// ChangePassword 修改密码
func (ur *userRepository) ChangePassword(ctx context.Context, username string, newPassword string) error {
	err := ur.dao.UpdatePasswordByUsername(ctx, username, newPassword)
	if err != nil {
		ur.l.Error("修改密码失败", zap.Error(err))
		return err
	}

	// 重新设置缓存,避免缓存不一致
	user, err := ur.dao.FindByUsername(ctx, username)
	if err != nil {
		ur.l.Error("密码修改后获取用户信息失败", zap.Error(err))
		return err
	}

	// 异步设置缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := ur.cache.Set(ctx, toDomainUser(user)); err != nil {
			ur.l.Error("密码修改后更新缓存失败", zap.Error(err))
		}
	}()

	return nil
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
		return du, nil
	}

	// 缓存中未找到数据，从数据库中查找
	u, err := ur.dao.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	du = toDomainUser(u)
	// 异步将用户信息写入缓存
	go func() {
		ctx := context.Background()
		if setErr := ur.cache.Set(ctx, du); setErr != nil {
			ur.l.Error("设置缓存失败", zap.Error(setErr))
		}
	}()

	return du, nil
}

// FindByPhone 通过电话查询用户
func (ur *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := ur.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}

	return toDomainUser(u), nil
}

// FindByUsername 通过用户名查询用户
func (ur *userRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	u, err := ur.dao.FindByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}

	return toDomainUser(u), nil
}

// DeleteUser 删除用户
func (ur *userRepository) DeleteUser(ctx context.Context, username string, uid int64) error {
	err := ur.dao.DeleteUser(ctx, username, uid)
	if err != nil {
		ur.l.Error("删除用户失败", zap.Error(err))
		return err
	}

	// 异步删除缓存
	go func() {
		ctx := context.Background()
		du, err := ur.cache.Get(ctx, uid)
		if err == nil {
			if err := ur.cache.Set(ctx, domain.User{ID: du.ID, Deleted: true}); err != nil {
				ur.l.Error("删除用户后更新缓存失败", zap.Error(err))
			}
		}
	}()

	return nil
}

// UpdateProfile 更新用户资料
func (ur *userRepository) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	err := ur.dao.UpdateProfile(ctx, profile)
	if err != nil {
		return err
	}

	// 异步更新缓存
	go func() {
		ctx := context.Background()
		du, err := ur.cache.Get(ctx, profile.UserID)
		if err == nil {
			du.Profile = profile
			if err := ur.cache.Set(ctx, du); err != nil {
				ur.l.Error("更新用户资料后更新缓存失败", zap.Error(err))
			}
		}
	}()

	return nil
}

// GetProfile 通过用户ID获取用户资料
func (ur *userRepository) GetProfile(ctx context.Context, UserID int64) (domain.Profile, error) {
	profile, err := ur.dao.GetProfileByUserID(ctx, UserID)
	if err != nil {
		ur.l.Error("获取用户资料失败", zap.Error(err))
		return domain.Profile{}, err
	}

	return profile, nil
}

// ListUser 获取用户列表
func (ur *userRepository) ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfile, error) {
	users, err := ur.dao.ListUser(ctx, pagination)
	if err != nil {
		ur.l.Error("获取用户列表失败", zap.Error(err))
		return nil, err
	}

	return users, nil
}

// fromDomainUser 将领域层对象转为dao层对象
func fromDomainUser(u domain.User) dao.User {
	return dao.User{
		ID:           u.ID,
		PasswordHash: u.Password,
		Username:     u.Username,
		Phone:        u.Phone,
		CreateTime:   u.CreateTime,
		UpdatedTime:  u.UpdatedTime,
		Deleted:      u.Deleted,
	}
}

// toDomainUser 将dao层对象转为领域层对象
func toDomainUser(u dao.User) domain.User {
	return domain.User{
		ID:          u.ID,
		Password:    u.PasswordHash,
		Username:    u.Username,
		Phone:       u.Phone,
		CreateTime:  u.CreateTime,
		UpdatedTime: u.UpdatedTime,
		Deleted:     u.Deleted,
	}
}
