package repository

import (
	"LinkMe/internal/dao"
	"LinkMe/internal/domain"
	"context"
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
	//TODO implement me
	panic("implement me")
}

func (ur *userRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}
