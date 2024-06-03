package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"context"
	"go.uber.org/zap"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrUserNotFound
)

type UserRepository interface {
	CreateUser(ctx context.Context, u domain.User) error
	FindByID(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
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
	go func() {
		if setErr := ur.cache.Set(ctx, du); setErr != nil {
			ur.l.Error("set cache failed", zap.Error(setErr))
		}
	}()
	return du, nil
}

// FindByPhone 通过电话查询用户
func (ur *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, phone)
	return toDomainUser(u), err
}

// FindByEmail 通过Email查询用户
func (ur *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	return toDomainUser(u), err
}

// 将领域层对象转为dao层对象
func fromDomainUser(u domain.User) models.User {
	return models.User{
		ID:           u.ID,
		PasswordHash: u.Password,
		Nickname:     u.Nickname,
		Birthday:     u.Birthday,
		Email:        u.Email,
		Phone:        u.Phone,
		CreateTime:   u.CreateTime,
	}
}

// 将dao层对象转为领域层对象
func toDomainUser(u models.User) domain.User {
	return domain.User{
		ID:         u.ID,
		Password:   u.PasswordHash,
		Nickname:   u.Nickname,
		Birthday:   u.Birthday,
		Email:      u.Email,
		Phone:      u.Phone,
		CreateTime: u.CreateTime,
	}
}
