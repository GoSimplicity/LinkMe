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
	ChangePassword(ctx context.Context, email string, password string, newPassword string, confirmPassword string) error
	DeleteUser(ctx context.Context, email string, password string, uid int64) error
	UpdateProfile(ctx context.Context, profile domain.Profile) (err error)
	GetProfileByUserID(ctx context.Context, UserID int64) (profile domain.Profile, err error)
	ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfileAndRule, error)
	GetUserCount(ctx context.Context) (int64, error)
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
		return domain.User{}, err
	}
	// 将密文密码转为明文
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (us *userService) ChangePassword(ctx context.Context, email string, password string, newPassword string, confirmPassword string) error {
	u, err := us.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return err
		}
		return err
	}
	if er := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); er != nil {
		return ErrInvalidUserOrPassword
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if er := us.repo.ChangePassword(ctx, email, string(newHash)); er != nil {
		return er
	}
	return nil
}

func (us *userService) DeleteUser(ctx context.Context, email string, password string, uid int64) error {
	u, err := us.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return err
		}
		return err
	}
	if er := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); er != nil {
		return ErrInvalidUserOrPassword
	}
	err = us.repo.DeleteUser(ctx, email, uid)
	if err != nil {
		return err
	}
	return nil
}
func (us *userService) UpdateProfile(ctx context.Context, profile domain.Profile) (err error) {
	return us.repo.UpdateProfile(ctx, profile)
}

func (us *userService) GetProfileByUserID(ctx context.Context, UserID int64) (profile domain.Profile, err error) {
	return us.repo.GetProfile(ctx, UserID)
}

func (us *userService) ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfileAndRule, error) {
	// 计算偏移量
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	return us.repo.ListUser(ctx, pagination)
}

func (us *userService) GetUserCount(ctx context.Context) (int64, error) {
	count, err := us.repo.GetUserCount(ctx)
	if err != nil {
		return -1, err
	}
	return count, nil
}
