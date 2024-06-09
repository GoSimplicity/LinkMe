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
	AddUserOperationLog(ctx context.Context, phone, action string) error
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
	_, err := s.cache.SetNX(ctx, fmt.Sprintf("sms_lock_:%s", mobile), "1", time.Second*60)
	status := int64(1)
	SmsId, _ := strconv.ParseInt(smsID, 10, 64)
	vCode := utils.GenRandomCode(6)
	if err != nil {
		status = int64(0)
		FailLog := models.VCodeSmsLog{
			SmsId:       SmsId,
			Status:      status,
			Mobile:      mobile,
			VCode:       vCode,
			CreateTime:  time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		}
		err = s.dao.Insert(ctx, FailLog)
		return fmt.Errorf("设置键值对失败")
	}

	err = s.cache.StoreVCode(ctx, smsID, mobile, vCode)
	if err != nil {
		status = int64(0)
		FailLog := models.VCodeSmsLog{
			SmsId:       SmsId,
			Status:      status,
			Mobile:      mobile,
			VCode:       vCode,
			CreateTime:  time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		}
		err = s.dao.Insert(ctx, FailLog)
		return fmt.Errorf("存储验证码失败:%w", err)
	}
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

func (s *smsRepository) AddUserOperationLog(ctx context.Context, phone, action string) error {
	log := models.UserOperationLog{
		Phone:      phone,
		Action:     action,
		CreateTime: time.Now().Unix(),
	}
	err := s.dao.InsertUserOperationLog(ctx, log)
	if err != nil {
		return fmt.Errorf("记录用户操作日志失败:%w", err)
	}
	return nil
}
