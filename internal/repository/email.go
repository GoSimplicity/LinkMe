package repository

import (
	"context"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	qqEmail "github.com/GoSimplicity/LinkMe/pkg/email"
	"github.com/GoSimplicity/LinkMe/utils"
	"go.uber.org/zap"
)

type EmailRepository interface {
	SendCode(ctx context.Context, email string) error
	CheckCode(ctx context.Context, email, vCode string) (bool, error)
}

// emailRepository 实现了 EmailRepository 接口
type emailRepository struct {
	cache cache.EmailCache
	l     *zap.Logger
}

// NewEmailRepository 创建并返回一个新的 smsRepository 实例
func NewEmailRepository(cache cache.EmailCache, l *zap.Logger) EmailRepository {
	return &emailRepository{
		cache: cache,
		l:     l,
	}
}

func (e emailRepository) SendCode(ctx context.Context, email string) error {
	e.l.Info("[emailRepository.SendCode]", zap.String("email", email))
	vCode := utils.GenRandomCode(6)
	e.l.Info("[emailRepository.SendCode]", zap.String("vCode", vCode))
	if err := e.cache.StoreVCode(ctx, email, vCode); err != nil {
		e.l.Error("[emailRepository.SendCode] StoreVCode失败", zap.Error(err))
		return err
	}
	body := fmt.Sprintf("您的验证码是：%s", vCode)
	return qqEmail.SendEmail(email, "【LinkMe】密码重置", body)
}

func (e emailRepository) CheckCode(ctx context.Context, email, vCode string) (bool, error) {
	storedCode, err := e.cache.GetVCode(ctx, email)
	return storedCode == vCode, err
}
