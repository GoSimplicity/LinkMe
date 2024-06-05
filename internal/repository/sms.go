package repository

import (
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"LinkMe/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// SmsRepository 接口定义了异步 SMS 记录操作的相关方法
type SmsRepository interface {
	CheckCode(ctx context.Context, mobile, smsID, vCode string) error
	SendCode(ctx context.Context, mobile, smsID string) error
}

// smsRepository 实现了 SendVCodeRepository 接口
type smsRepository struct {
	dao   dao.SmsDAO
	cache cache.SMSCache
}

// NewSendCodeRepository 创建并返回一个新的 sendCodeRepository 实例
func NewSendCodeRepository(dao dao.SmsDAO) SmsRepository {
	return &smsRepository{
		dao: dao,
	}
}

// CheckCode 检查验证码是否正确
func (s *smsRepository) CheckCode(ctx context.Context, mobile, smsID, vCode string) error {
	storedVCode, err := s.cache.GetVCode(ctx, smsID, mobile)
	if err != nil {
		return err
	}
	if storedVCode != vCode || storedVCode == "" {
		return errors.New("该验证码无效或已过期")
	}
	return nil
}

// SendCode 记录发送的验证码
func (s *smsRepository) SendCode(ctx context.Context, mobile, smsID string) error {
	vCode := fmt.Sprintf("%06d", utils.GenRandomCode(100000)) //生成验证码

	err := s.cache.StoreVCode(ctx, smsID, mobile, vCode)
	if err != nil {
		return fmt.Errorf("存储验证码失败:%w", err)
	}
	//发送验证码
	message := fmt.Sprintf("你的验证码是: %s", vCode)
	_, err = s.cache.GetVCode(ctx, message, smsID)
	status := int64(1)
	if err != nil {
		status = int64(0)
		return fmt.Errorf("发送验证码失败:%w", err)
	}
	SmsId, _ := strconv.ParseInt(smsID, 10, 64)
	log := models.VCodeSmsLog{
		SmsId:       SmsId,
		Status:      status,
		Mobile:      mobile,
		VCode:       vCode,
		CreateTime:  time.Now().Unix(),
		UpdatedTime: time.Now().Unix(),
	}
	err = s.dao.Insert(ctx, log)
	if err != nil {
		return fmt.Errorf("记录发送日志失败:%w", err)
	}
	return nil
}
