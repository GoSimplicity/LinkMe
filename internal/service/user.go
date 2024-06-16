package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"LinkMe/internal/repository/models"
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("username or password is incorrect")
)

type ProfileService interface {
	UpdateProfile(ctx context.Context, profile *models.Profile) (err error)
	GetProfileByUserID(ctx context.Context, UserID int64) (profile *models.Profile, err error)
}

type profileServiceImpl struct {
	profileRepo repository.ProfileRepository
}

func NewProfileService(profileRepo repository.ProfileRepository) ProfileService {
	return &profileServiceImpl{profileRepo: profileRepo}
}

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	ChangePassword(ctx context.Context, email string, password string, newPassword string, confirmPassword string) error
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
		us.l.Error("password conversion failed")
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
		us.l.Error("password conversion failed", zap.Error(err))
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
		us.l.Error("failed to find user", zap.Error(err))
		return err
	}
	if er := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); er != nil {
		us.l.Error("password verification failed", zap.Error(er))
		return ErrInvalidUserOrPassword
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		us.l.Error("failed to hash new password", zap.Error(err))
		return err
	}
	if er := us.repo.ChangePassword(ctx, email, string(newHash)); er != nil {
		us.l.Error("failed to change password", zap.Error(er))
		return er
	}
	return nil
}

func (s *profileServiceImpl) UpdateProfile(ctx context.Context, profile *models.Profile) (err error) {
	return s.profileRepo.UpdateProfile(ctx, profile)
}

func (s *profileServiceImpl) GetProfileByUserID(ctx context.Context, UserID int64) (profile *models.Profile, err error) {
	return s.profileRepo.GetProfile(ctx, UserID)

}
