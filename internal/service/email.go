package service

import (
	"LinkMe/internal/repository"
	qqEmail "LinkMe/pkg/email"
	"LinkMe/utils"
	"context"
	"fmt"
	"go.uber.org/zap"
)

// EmailService 定义了发送邮箱验证码的服务接口
type EmailService interface {
	SendCode(ctx context.Context, email string) error
	CheckCode(ctx context.Context, email, vCode string) (bool, error)
}

// emailService 实现了 EmailService 接口
type emailService struct {
	repo repository.EmailRepository
	l    *zap.Logger
}

// NewEmailService 创建并返回一个新的 sendCodeService 实例
func NewEmailService(r repository.EmailRepository, l *zap.Logger) EmailService {
	s := &emailService{
		repo: r,
		l:    l,
	}
	return s
}

func (e emailService) SendCode(ctx context.Context, email string) error {
	e.l.Info("[EmailService.SendCode]", zap.String("email", email))
	vCode := utils.GenRandomCode(6)
	e.l.Info("[EmailService.SendCode]", zap.String("vCode", vCode))
	//if err := e.repo.StoreVCode(ctx, email, vCode); err != nil {
	//	e.l.Error("[EmailService.SendCode] StoreVCode失败", zap.Error(err))
	//	return err
	//}
	body := fmt.Sprintf("您的验证码是：%s", vCode)
	return qqEmail.SendEmail(email, "【LinkMe】密码重置", body)
}

func (e emailService) CheckCode(ctx context.Context, email, vCode string) (bool, error) {
	return e.repo.CheckCode(ctx, email, vCode)
}
