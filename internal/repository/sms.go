package repository

import (
	"LinkMe/internal/repository/dao"
	"context"
)

// SendVCodeRepository 接口定义了异步 SMS 记录操作的相关方法
type SendVCodeRepository interface {
	CheckCode(ctx context.Context, mobile, smsID, vCode string) error
	SendCode(ctx context.Context, mobile, smsID string) error
}

// sendCodeRepository 实现了 SendVCodeRepository 接口
type sendCodeRepository struct {
	dao dao.SendVCodeDAO
}

// NewSendCodeRepository 创建并返回一个新的 sendCodeRepository 实例
func NewSendCodeRepository(dao dao.SendVCodeDAO) SendVCodeRepository {
	return &sendCodeRepository{
		dao: dao,
	}
}

// CheckCode 检查验证码是否正确
func (s *sendCodeRepository) CheckCode(ctx context.Context, mobile, smsID, vCode string) error {
	// TODO: 实现从 DAO 中检查验证码的逻辑
	return nil
}

// SendCode 记录发送的验证码
func (s *sendCodeRepository) SendCode(ctx context.Context, mobile, smsID string) error {
	// TODO: 实现将验证码存储在 DAO 中的逻辑
	return nil
}
