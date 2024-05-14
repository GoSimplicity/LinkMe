package repository

import (
	"LinkMe/internal/dao"
	"LinkMe/internal/domain"
	"context"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrUserNotFound
)

type UserRepository interface {
	CreateUser(ctx context.Context, u domain.User) error
	FindByID(ctx context.Context, id int64) (domain.User, error)
	FindByUsername(ctx context.Context, username string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type userRepository struct {
	dao dao.UserDAO
}

func NewUserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{
		dao: dao,
	}
}

func (ur *userRepository) CreateUser(ctx context.Context, u domain.User) error {
	return ur.dao.CreateUser(ctx, fromDomainUser(u))
}

func (ur *userRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	u, err := ur.dao.FindByID(ctx, id)
	return toDomainUser(u), err
}

func (ur *userRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}
func (ur *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (ur *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	return toDomainUser(u), err
}

func fromDomainUser(u domain.User) dao.User {
	return dao.User{
		PasswordHash: u.Password,
		Nickname:     u.Nickname,
		Birthday:     u.Birthday,
		Email:        u.Email,
		Phone:        u.Phone,
	}
}

func toDomainUser(u dao.User) domain.User {
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
