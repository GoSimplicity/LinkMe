package service

import (
	"LinkMe/internal/repository"
	"context"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.uber.org/zap"
)

// SendCodeService 定义了发送验证码的服务接口
type SendCodeService interface {
	SendCode(ctx context.Context, tplId string, args []string, numbers ...string) error
	CheckCode(ctx context.Context, mobile, vCode string) (bool, error)
}

// sendCodeService 实现了 SendCodeService 接口
type sendCodeService struct {
	repo     repository.SmsRepository
	l        *zap.Logger
	client   *sms.Client
	appId    string
	signName string
}

// NewSendCodeService 创建并返回一个新的 sendCodeService 实例
func NewSendCodeService(repo repository.SmsRepository, l *zap.Logger, client *sms.Client, appId string, signName string) SendCodeService {
	s := &sendCodeService{
		repo:     repo,
		l:        l,
		client:   client,
		appId:    appId,
		signName: signName,
	}
	return s
}

func (s sendCodeService) SendCode(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//TODO implement me

	//使用分布式锁，保证每个手机号一分钟内只能请求一次

	//随机生成长度为6的验证码

	//redis存储验证码

	//构造req

	//向第三方发送req

	//存储用户操作日志

	panic("implement me")
}

func (s sendCodeService) CheckCode(ctx context.Context, mobile, vCode string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
