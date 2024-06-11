package repository

import (
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

// SmsRepository 接口定义了异步 SMS 记录操作的相关方法
type SmsRepository interface {
	CheckCode(ctx context.Context, smsID, number, vCode string) (bool, error)
	AddUserOperationLog(ctx context.Context, log models.VCodeSmsLog) error
	SetNX(ctx context.Context, number string, value interface{}, expiration time.Duration) (*redis.BoolCmd, error)
	StoreVCode(ctx context.Context, smsID, number string, vCode string) error
	Exist(ctx context.Context, number string) bool
	Count(ctx context.Context, number string) int
	IncrCnt(ctx context.Context, number string) error
}

// smsRepository 实现了 SmsRepository 接口
type smsRepository struct {
	dao   dao.SmsDAO
	cache cache.SMSCache
}

// NewSmsRepository 创建并返回一个新的 smsRepository 实例
func NewSmsRepository(dao dao.SmsDAO, cache cache.SMSCache) SmsRepository {
	return &smsRepository{
		dao:   dao,
		cache: cache,
	}
}

// CheckCode 检查验证码是否正确
func (s *smsRepository) CheckCode(ctx context.Context, smsID, number, vCode string) (bool, error) {
	storedVCode, err := s.cache.GetVCode(ctx, smsID, number)
	if err != nil {
		return false, err
	}
	if storedVCode != vCode || storedVCode == "" {
		return false, errors.New("该验证码无效或已过期")
	}
	return true, nil
}

// AddUserOperationLog 添加用户验证码行为日志
func (s *smsRepository) AddUserOperationLog(ctx context.Context, log models.VCodeSmsLog) error {
	return s.dao.Insert(ctx, log)
}

func (s *smsRepository) SetNX(ctx context.Context, number string, value interface{}, expiration time.Duration) (*redis.BoolCmd, error) {
	return s.cache.SetNX(ctx, number, value, expiration)
}

func (s *smsRepository) StoreVCode(ctx context.Context, smsID, number, vCode string) error {
	return s.cache.StoreVCode(ctx, smsID, number, vCode)
}

func (s *smsRepository) Exist(ctx context.Context, key string) bool {
	return s.cache.Exist(ctx, key)
}

func (s *smsRepository) Count(ctx context.Context, number string) int {
	return s.cache.Count(ctx, number)
}

func (s *smsRepository) IncrCnt(ctx context.Context, number string) error {
	return s.IncrCnt(ctx, number)
}
