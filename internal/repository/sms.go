package repository

import (
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/pkg/sms"
	"LinkMe/utils"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	tencenterros "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const locked = "sms_locked"

// SmsRepository 接口定义了异步 SMS 记录操作的相关方法
type SmsRepository interface {
	CheckCode(ctx context.Context, smsID, number, vCode string) (bool, error)
	AddUserOperationLog(ctx context.Context, log dao.VCodeSmsLog) error
	SetNX(ctx context.Context, number string, value interface{}, expiration time.Duration) (*redis.BoolCmd, error)
	StoreVCode(ctx context.Context, smsID, number string, vCode string) error
	Exist(ctx context.Context, number string) bool
	Count(ctx context.Context, number string) int
	IncrCnt(ctx context.Context, number string) error
	ReleaseLock(ctx context.Context, number string) error
	SendCode(ctx context.Context, number string) error
}

// smsRepository 实现了 SmsRepository 接口
type smsRepository struct {
	dao    dao.SmsDAO
	cache  cache.SMSCache
	l      *zap.Logger
	client *sms.TencentSms //Todo 完成多个sms的集成
}

// NewSmsRepository 创建并返回一个新的 smsRepository 实例
func NewSmsRepository(dao dao.SmsDAO, cache cache.SMSCache, l *zap.Logger, client *sms.TencentSms) SmsRepository {
	return &smsRepository{
		dao:    dao,
		cache:  cache,
		l:      l,
		client: client,
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
func (s *smsRepository) AddUserOperationLog(ctx context.Context, log dao.VCodeSmsLog) error {
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
	return s.cache.IncrCnt(ctx, number)
}

func (s *smsRepository) ReleaseLock(ctx context.Context, number string) error {
	return s.cache.ReleaseLock(ctx, number)
}

func (s *smsRepository) SendCode(ctx context.Context, number string) error {
	// 限制：用户一分钟内只能发送一次sms请求 && 用户一天内只能发送5次SMS请求
	if s.cache.Count(ctx, number) > 5 {
		s.l.Warn("用户今日发送验证码次数过多", zap.String("number", number))
		return fmt.Errorf("用户今日发送验证码次数过多")
	}
	_, err := s.cache.SetNX(ctx, number, locked, time.Second*60)
	if err != nil {
		s.l.Warn("[smsService SendCode] s.repo.SetNX 报错: ", zap.Error(err))
		return fmt.Errorf("验证码发送过于频繁，请稍后再尝试")
	}
	vCode := utils.GenRandomCode(6)
	//todo: sms商无缝切换
	smsID, driver, err := s.client.Send(ctx, []string{vCode}, []string{number}...)
	id, _ := strconv.ParseInt(smsID, 10, 64)
	log := dao.VCodeSmsLog{
		SmsId:       id,
		SmsType:     "vCode", //todo
		Mobile:      number,
		VCode:       vCode,
		Driver:      driver,
		Status:      1,  //0为失败，1为成功；其中默认为 成功
		StatusCode:  "", //todo
		CreateTime:  time.Now().UnixNano(),
		UpdatedTime: time.Now().UnixNano(),
		DeletedTime: time.Now().UnixNano(),
	}

	// SDK异常
	if _, ok := err.(*tencenterros.TencentCloudSDKError); ok {
		s.l.Error("[smsService SendCode] an API error has returned: ", zap.Error(err))
		s.cache.ReleaseLock(ctx, number)
		return err
	}
	// 非SDK异常，直接失败
	if err != nil {
		s.l.Error("[smsService SendCode] an error has returned: ", zap.Error(err))
		return err
	}
	if err = s.StoreVCode(ctx, smsID, number, vCode); err != nil {
		log.Status = 0
		log.VCode = "-1"
		s.AddUserOperationLog(ctx, log)
		s.cache.ReleaseLock(ctx, number)
		s.l.Error("[smsService SendCode] s.repo.StoreVCode 报错: ", zap.Error(err))
		return fmt.Errorf("存储随机数失败")
	}

	if err = s.cache.IncrCnt(ctx, number); err != nil {
		s.l.Error("[smsService SendCode] s.rdb.IncrCnt 报错: ", zap.Error(err))
		return err
	}
	return s.AddUserOperationLog(ctx, log)
}
