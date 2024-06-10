package service

import (
	"LinkMe/internal/repository"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/models"
	"LinkMe/pkg/sms"
	"LinkMe/utils"
	"context"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

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
	//todo 参数校验

	//一个系统中的每个用户一分钟内 只能发送一条vCode
	key := fmt.Sprintf("sms_locked:%s", number)
	nx, err := s.repo.SetNX(ctx, key, "sms_locked", time.Second*60)
	if err != nil || nx.Err() != nil {
		return fmt.Errorf("验证码发送过于频繁，请稍后再尝试")
	}
	//生成随机数
	vCode := utils.GenRandomCode(6)
	//发送sms req && 操作入库
	//Todo 记录用户今日sms发送次数，后续完成限流处理
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
	if err != nil {
		s.l.Warn("sms发送失败", zap.Error(err), zap.String("driver", driver), zap.String("smsID", smsID), zap.String("number", number))
		log.Status = 0
		log.VCode = "-1"
		s.repo.AddUserOperationLog(ctx, log)
		return err
	}
	if err = s.repo.StoreVCode(ctx, smsID, number, vCode); err != nil {
		log.Status = 0
		log.VCode = "-1"
		s.repo.AddUserOperationLog(ctx, log)
		return fmt.Errorf("存储随机数失败")
	}

	return s.repo.AddUserOperationLog(ctx, log)
}

// CheckCode 检查验证码是否正确
func (s smsService) CheckCode(ctx context.Context, smsID, number, vCode string) (bool, error) {
	return s.repo.CheckCode(ctx, smsID, number, vCode)
}
