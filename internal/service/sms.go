package service

import (
	"LinkMe/internal/repository"
	"LinkMe/internal/repository/cache"
	"LinkMe/utils"
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.uber.org/zap"
	"time"
)

// SendCodeService 定义了发送验证码的服务接口
type SendCodeService interface {
	SendCode(ctx context.Context, tplId string, args []string, numbers ...string) error
	//CheckCode(ctx context.Context, mobile, vCode string) (bool, error)
}

// sendCodeService 实现了 SendCodeService 接口
type sendCodeService struct {
	repo     repository.SmsRepository
	l        *zap.Logger
	client   *sms.Client
	rdb      cache.SMSCache
	appId    string
	signName string
}

// NewSendCodeService 创建并返回一个新的 sendCodeService 实例
func NewSendCodeService(repo repository.SmsRepository, l *zap.Logger, client *sms.Client, appId string, signName string, rdb cache.SMSCache) SendCodeService {
	s := &sendCodeService{
		repo:     repo,
		l:        l,
		client:   client,
		appId:    appId,
		signName: signName,
		rdb:      rdb,
	}
	return s
}

func (s *sendCodeService) SendCode(ctx context.Context, tplId string, args []string, numbers ...string) error {

	//使用分布式锁，保证每个手机号一分钟内只能请求一次
	for _, number := range numbers {
		lockKey := "sms_lock" + number
		lock, err := s.rdb.SetNX(ctx, lockKey, "locked", time.Minute)
		if err != nil {
			s.l.Error("获取锁失败", zap.Error(err), zap.Error(err))
		}
		if !lock.Val() {
			s.l.Warn("一分钟只能发送一次验证码", zap.String("number", number))
			return fmt.Errorf("一分钟只能发送一次验证码")
		}
	}
	//随机生成长度为6的验证码
	code := fmt.Sprintf("%06d", utils.GenRandomCode(6))
	args = append(args, code)

	//构造req
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = &s.appId
	request.SignName = &s.signName
	request.TemplateId = &tplId
	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(numbers)
	//向第三方发送req
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

	//redis存储验证码
	for _, number := range numbers {
		err = s.rdb.StoreVCode(ctx, tplId, number, code)
		if err != nil {
			s.l.Error("存储验证码失败", zap.String("number", number), zap.Error(err))
			return fmt.Errorf("存储验证码失败:%w", err)
		}
		//记录发送日志
		err = s.repo.SendCode(ctx, number, code)
		if err != nil {
			s.l.Error("记录发送日志失败", zap.String("number", number), zap.Error(err))
			return err
		}
		err = s.repo.AddUserOperationLog(ctx, number, "发送验证码")
		if err != nil {
			s.l.Error("记录用户操作日志失败", zap.String("number", number), zap.Error(err))
			return err
		}
	}
	//存储用户操作日志

	return nil
}

/*// CheckCode 检查验证码是否正确
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
}*/
