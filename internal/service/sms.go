package service

import (
	"context"
	"fmt"
	"time"

	"LinkMe/internal/repository"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
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
	repo     repository.SendVCodeRepository
	l        *zap.Logger
	client   *sms.Client
	appId    string
	signName string
}

// NewSendCodeService 创建并返回一个新的 sendCodeService 实例
func NewSendCodeService(repo repository.SendVCodeRepository, l *zap.Logger, client *sms.Client, appId string, signName string) SendCodeService {
	return &sendCodeService{
		repo:     repo,
		l:        l,
		client:   client,
		appId:    appId,
		signName: signName,
	}
}

// SendCode 发送验证码到指定手机号
func (s *sendCodeService) SendCode(ctx context.Context, tplId string, args []string, numbers ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = &s.appId
	request.SignName = &s.signName
	request.TemplateId = &tplId
	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(numbers)

	response, err := s.client.SendSms(request)
	if err != nil {
		s.l.Error("发送验证码失败", zap.Error(err))
		return err
	}

	for _, status := range response.Response.SendStatusSet {
		if status == nil || status.Code == nil || *status.Code != "Ok" {
			// 发送失败
			errMsg := fmt.Errorf("发送短信失败 code: %s, msg: %s", *status.Code, *status.Message)
			s.l.Error(errMsg.Error())
			return errMsg
		}
	}

	s.l.Info("验证码发送成功", zap.Strings("numbers", numbers), zap.String("templateId", tplId))
	return nil
}

// CheckCode 检查验证码是否正确
func (s *sendCodeService) CheckCode(ctx context.Context, mobile, vCode string) (bool, error) {
	// 假设存储库有记录 smsID
	smsID := fmt.Sprintf("%s-%d", mobile, time.Now().UnixNano())

	err := s.repo.CheckCode(ctx, mobile, smsID, vCode)
	if err != nil {
		s.l.Error("验证验证码失败", zap.Error(err))
		return false, err
	}

	s.l.Info("验证码验证成功", zap.String("mobile", mobile), zap.String("code", vCode))
	return true, nil
}
