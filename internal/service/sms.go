package service

import (
	"LinkMe/internal/repository"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/models"
	"LinkMe/pkg/sms"
	"LinkMe/utils"
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const locked = "sms_locked"

// SmsService 定义了发送验证码的服务接口
type SmsService interface {
	SendCode(ctx context.Context, number string) error
	CheckCode(ctx context.Context, smsID, mobile, vCode string) (bool, error)
}

// smsService 实现了 SmsService 接口
type smsService struct {
	repo   repository.SmsRepository
	l      *zap.Logger
	client *sms.TencentSms //Todo 完成多个sms的集成
	rdb    cache.SMSCache
}

// NewSmsService 创建并返回一个新的 sendCodeService 实例
func NewSmsService(r repository.SmsRepository, l *zap.Logger, client *sms.TencentSms, rdb cache.SMSCache) SmsService {
	s := &smsService{
		repo:   r,
		l:      l,
		client: client,
		rdb:    rdb,
	}
	return s
}

func (s smsService) SendCode(ctx context.Context, number string) error {
	// 限制：用户一分钟内只能发送一次sms请求 && 用户一天内只能发送5次SMS请求
	if s.rdb.Count(ctx, number) > 5 {
		s.l.Warn("用户今日发送验证码次数过多", zap.String("number", number))
		return fmt.Errorf("用户今日发送验证码次数过多")
	}
	_, err := s.repo.SetNX(ctx, number, locked, time.Second*60)
	if err != nil {
		s.l.Warn("[smsService SendCode] s.repo.SetNX 报错: ", zap.Error(err))
		return fmt.Errorf("验证码发送过于频繁，请稍后再尝试")
	}
	vCode := utils.GenRandomCode(6)
	//todo: sms商无缝切换
	smsID, driver, err := s.client.Send(ctx, []string{vCode}, []string{number}...)
	id, _ := strconv.ParseInt(smsID, 10, 64)
	log := models.VCodeSmsLog{
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
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		s.l.Error("[smsService SendCode] an API error has returned: ", zap.Error(err))
		s.rdb.ReleaseLock(ctx, number)
		return err
	}
	// 非SDK异常，直接失败
	if err != nil {
		s.l.Error("[smsService SendCode] an error has returned: ", zap.Error(err))
		return err
	}
	if err = s.repo.StoreVCode(ctx, smsID, number, vCode); err != nil {
		log.Status = 0
		log.VCode = "-1"
		s.repo.AddUserOperationLog(ctx, log)
		s.rdb.ReleaseLock(ctx, number)
		s.l.Error("[smsService SendCode] s.repo.StoreVCode 报错: ", zap.Error(err))
		return fmt.Errorf("存储随机数失败")
	}

	if err = s.rdb.IncrCnt(ctx, number); err != nil {
		s.l.Error("[smsService SendCode] s.rdb.IncrCnt 报错: ", zap.Error(err))
		return err
	}
	return s.repo.AddUserOperationLog(ctx, log)
}

// CheckCode 检查验证码是否正确
func (s smsService) CheckCode(ctx context.Context, smsID, number, vCode string) (bool, error) {
	return s.repo.CheckCode(ctx, smsID, number, vCode)
}
