package repository

import (
	"LinkMe/internal/repository/dao"
	"context"
)

// SmsRepository 接口定义了异步 SMS 记录操作的相关方法
type SmsRepository interface {
	CheckCode(ctx context.Context, mobile, smsID, vCode string) error
	SendCode(ctx context.Context, mobile, smsID string) error
}

// smsRepository 实现了 SendVCodeRepository 接口
type smsRepository struct {
	dao dao.SmsDAO
}

// NewSendCodeRepository 创建并返回一个新的 sendCodeRepository 实例
func NewSendCodeRepository(dao dao.SmsDAO) SmsRepository {
	return &smsRepository{
		dao: dao,
	}
}

// CheckCode 检查验证码是否正确
func (s *smsRepository) CheckCode(ctx context.Context, mobile, smsID, vCode string) error {
	// TODO: 实现从 DAO 中检查验证码的逻辑
	return nil
}

// SendCode 记录发送的验证码
func (s *smsRepository) SendCode(ctx context.Context, mobile, smsID string) error {
	// TODO: 实现将验证码存储在 DAO 中的逻辑
	return nil
}
