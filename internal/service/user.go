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
	// ErrDuplicateEmail 表示邮箱重复错误
	ErrDuplicateEmail = repository.ErrDuplicateEmail
	// ErrInvalidUserOrPassword 表示用户名或密码错误
	ErrInvalidUserOrPassword = errors.New("username or password is incorrect")
)

// UserService 定义用户服务接口
type UserService interface {
	// SignUp 用户注册
	SignUp(ctx context.Context, u domain.User) error
	// Login 用户登录
	Login(ctx context.Context, email string, password string) (domain.User, error)
	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, email string, password string, newPassword string, confirmPassword string) error
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, email string, password string, uid int64) error
	// UpdateProfile 更新用户资料
	UpdateProfile(ctx context.Context, profile domain.Profile) error
	// GetProfileByUserID 通过用户ID获取用户资料
	GetProfileByUserID(ctx context.Context, UserID int64) (domain.Profile, error)
	// ListUser 获取用户列表
	ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfileAndRule, error)
	// GetUserCount 获取用户总数
	GetUserCount(ctx context.Context) (int64, error)
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

// SignUp 注册逻辑
func (us *userService) SignUp(ctx context.Context, u domain.User) error {
	// 生成密码哈希值
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	// 将用户信息输入到搜索库中
	go func() {
		err = us.searchRepo.InputUser(ctx, domain.UserSearch{
			Email:    u.Email,
			Id:       u.ID,
			Nickname: u.Profile.NickName,
		})
		if err != nil {
			us.l.Error("failed to input user to search repo", zap.Error(err))
		}
	}()
	// 创建用户
	return us.repo.CreateUser(ctx, u)
}

// Login 登陆逻辑
func (us *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	// 根据邮箱查找用户
	u, err := us.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, err
	} else if err != nil {
		return domain.User{}, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}

// ChangePassword 修改密码
func (us *userService) ChangePassword(ctx context.Context, email string, password string, newPassword string, confirmPassword string) error {
	// 根据邮箱查找用户
	u, err := us.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return err
		}
		return err
	}
	// 验证当前密码
	if er := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); er != nil {
		return ErrInvalidUserOrPassword
	}
	// 生成新密码哈希值
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// 更新密码
	if er := us.repo.ChangePassword(ctx, email, string(newHash)); er != nil {
		return er
	}
	return nil
}

// DeleteUser 删除用户
func (us *userService) DeleteUser(ctx context.Context, email string, password string, uid int64) error {
	// 根据邮箱查找用户
	u, err := us.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return err
		}
		return err
	}
	// 验证密码
	if er := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); er != nil {
		return ErrInvalidUserOrPassword
	}
	// 删除用户
	err = us.repo.DeleteUser(ctx, email, uid)
	if err != nil {
		return err
	}
	// 从搜索库中删除用户索引
	err = us.searchRepo.DeleteUserIndex(ctx, uid)
	if err != nil {
		us.l.Error("failed to delete user from search repo", zap.Error(err))
	}

	return nil
}

// UpdateProfile 更新用户资料
func (us *userService) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	// 根据用户ID查找用户
	user, _ := us.repo.FindByID(ctx, profile.UserID)
	// 如果昵称发生变化，更新搜索库中的用户信息
	if user.Profile.NickName != profile.NickName {
		err := us.searchRepo.InputUser(ctx, domain.UserSearch{
			Id:       profile.UserID,
			Nickname: profile.NickName,
		})
		if err != nil {
			us.l.Error("failed to input user to search repo", zap.Error(err))
		}
	}
	// 更新用户资料
	return us.repo.UpdateProfile(ctx, profile)
}

// GetProfileByUserID 根据用户ID获取用户资料
func (us *userService) GetProfileByUserID(ctx context.Context, UserID int64) (profile domain.Profile, err error) {
	// 从仓库中获取用户资料
	return us.repo.GetProfile(ctx, UserID)
}

// ListUser 获取用户列表
func (us *userService) ListUser(ctx context.Context, pagination domain.Pagination) ([]domain.UserWithProfileAndRule, error) {
	// 计算偏移量
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	// 从仓库中获取用户列表
	return us.repo.ListUser(ctx, pagination)
}

// GetUserCount 获取用户总数
func (us *userService) GetUserCount(ctx context.Context) (int64, error) {
	// 从仓库中获取用户总数
	count, err := us.repo.GetUserCount(ctx)
	if err != nil {
		return -1, err
	}
	return count, nil
}
