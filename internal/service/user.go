package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("username or password is incorrect")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	l    *zap.Logger
}

func NewUserService(repo repository.UserRepository, l *zap.Logger) UserService {
	return &userService{
		repo: repo,
		l:    l,
	}
}

// SignUp 注册逻辑
func (us *userService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		us.l.Error("password conversion filed")
		return err
	}
	u.Password = string(hash)
	return us.repo.CreateUser(ctx, u)
}

// Login 登陆逻辑
func (us *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := us.repo.FindByEmail(ctx, email)
	// 如果用户没有找到(未注册)，则返回空对象
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, err
	} else if err != nil {
		us.l.Error("user not found", zap.Error(err))
		return domain.User{}, err
	}
	// 将密文密码转为明文
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		us.l.Error("password conversion filed", zap.Error(err))
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}
