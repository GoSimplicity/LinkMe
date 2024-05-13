package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, username string, password string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) SignUp(ctx context.Context, u domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (us *userService) Login(ctx context.Context, username string, password string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}
