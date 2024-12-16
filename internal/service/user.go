package service

import (
	"context"
	"errors"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrDuplicateUsername 表示用户名重复错误
	ErrDuplicateUsername = repository.ErrDuplicateUsername
	// ErrInvalidUserOrPassword 表示用户名或密码错误
	ErrInvalidUserOrPassword = errors.New("用户名或密码错误")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, username string, password string) (domain.User, error)
	ChangePassword(ctx context.Context, username string, password string, newPassword string, confirmPassword string) error
	DeleteUser(ctx context.Context, username string, password string, uid int64) error
	UpdateProfile(ctx context.Context, profile domain.Profile) error
	GetProfileByUserID(ctx context.Context, UserID int64) (domain.Profile, error)
	ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfile, error)
	UpdateProfileAdmin(ctx context.Context, profile domain.Profile) error
}

type userService struct {
	repo       repository.UserRepository
	l          *zap.Logger
	searchRepo repository.SearchRepository
}

func NewUserService(repo repository.UserRepository, l *zap.Logger, searchRepo repository.SearchRepository) UserService {
	return &userService{
		repo:       repo,
		searchRepo: searchRepo,
		l:          l,
	}
}

// SignUp 用户注册
func (us *userService) SignUp(ctx context.Context, u domain.User) error {
	if err := u.ValidateUsername(); err != nil {
		return errors.New("用户名需要至少六位的字母数字组合")
	}

	if err := u.ValidatePassword(); err != nil {
		return errors.New("密码需要至少8位的字母数字符号组合")
	}

	if err := u.HashPassword(); err != nil {
		return err
	}

	// 异步更新搜索索引
	go func() {
		ctx := context.Background()
		err := us.searchRepo.InputUser(ctx, domain.UserSearch{
			Username: u.Username,
			Id:       u.ID,
			RealName: u.Profile.RealName,
		})
		if err != nil {
			us.l.Error("更新搜索索引失败",
				zap.String("用户名", u.Username),
				zap.Error(err))
		}
	}()

	return us.repo.CreateUser(ctx, u)
}

// Login 用户登录
func (us *userService) Login(ctx context.Context, username string, password string) (domain.User, error) {
	if username == "" || password == "" {
		return domain.User{}, errors.New("用户名和密码不能为空")
	}

	u, err := us.repo.FindByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}

	if err := u.VerifyPassword(password); err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}

// ChangePassword 修改密码
func (us *userService) ChangePassword(ctx context.Context, username string, password string, newPassword string, confirmPassword string) error {
	if newPassword != confirmPassword {
		return errors.New("新密码与确认密码不匹配")
	}

	u, err := us.repo.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	if err := u.VerifyPassword(password); err != nil {
		return ErrInvalidUserOrPassword
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return us.repo.ChangePassword(ctx, username, string(newHash))
}

// DeleteUser 删除用户
func (us *userService) DeleteUser(ctx context.Context, username string, password string, uid int64) error {
	u, err := us.repo.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	if err := u.VerifyPassword(password); err != nil {
		return ErrInvalidUserOrPassword
	}

	u.MarkAsDeleted()

	if err := us.repo.DeleteUser(ctx, username, uid); err != nil {
		return err
	}

	// 异步删除搜索索引
	go func() {
		ctx := context.Background()
		if err := us.searchRepo.DeleteUserIndex(ctx, uid); err != nil {
			us.l.Error("从搜索引擎中删除用户失败",
				zap.Int64("用户ID", uid),
				zap.Error(err))
		}
	}()

	return nil
}

// UpdateProfile 更新用户资料
func (us *userService) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	user, err := us.repo.FindByID(ctx, profile.UserID)
	if err != nil {
		return err
	}

	if user.Profile.RealName != profile.RealName {
		// 异步更新搜索索引
		go func() {
			ctx := context.Background()
			err := us.searchRepo.InputUser(ctx, domain.UserSearch{
				Id:       profile.UserID,
				RealName: profile.RealName,
			})
			if err != nil {
				us.l.Error("更新用户搜索索引失败",
					zap.Int64("用户ID", profile.UserID),
					zap.Error(err))
			}
		}()
	}

	return us.repo.UpdateProfile(ctx, profile)
}

// GetProfileByUserID 获取用户资料
func (us *userService) GetProfileByUserID(ctx context.Context, UserID int64) (domain.Profile, error) {
	if UserID <= 0 {
		return domain.Profile{}, errors.New("无效的用户ID")
	}

	return us.repo.GetProfile(ctx, UserID)
}

// ListUser 获取用户列表
func (us *userService) ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfile, error) {
	if pagination.Page <= 0 || *pagination.Size <= 0 {
		return nil, errors.New("无效的分页参数")
	}

	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset

	return us.repo.ListUser(ctx, pagination)
}

// UpdateProfileAdmin 更新用户资料(管理员)
func (us *userService) UpdateProfileAdmin(ctx context.Context, profile domain.Profile) error {
	return us.repo.UpdateProfileAdmin(ctx, profile)
}
