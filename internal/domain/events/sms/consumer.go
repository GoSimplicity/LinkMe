package post

import (
	"LinkMe/internal/service"
	"context"
	"encoding/json"
	"time"

	//"LinkMe/internal/repository"
	"LinkMe/internal/repository/cache"
	"LinkMe/pkg/samarap"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type SMSConsumer struct {
	service service.SendCodeService
	client  sarama.Client
	l       *zap.Logger
	rdb     cache.SMSCache
}

func NewSMSConsumer(service service.SendCodeService, client sarama.Client, l *zap.Logger, rdb cache.SMSCache) *SMSConsumer {
	return &SMSConsumer{service: service, client: client, l: l, rdb: rdb}
}

func (s *SMSConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("sms_consumer_group", s.client)
	if err != nil {
		return err
	}
	go func() {
		attempts := 0
		maxRetries := 3
		for attempts < maxRetries {
			er := cg.Consume(ctx, []string{TopicSMS}, samarap.NewHandler(s.l, s.HandleMessage))
			if er != nil {
				s.l.Error("消费错误", zap.Error(er), zap.Int("重试次数", attempts+1))
				attempts++
				time.Sleep(time.Second * time.Duration(attempts))
				continue
			}
			break
		}
		if attempts >= maxRetries {
			s.l.Error("达到最大重试次数，退出消费")
		}
	}()
	return nil
}

func (s *SMSConsumer) HandleMessage(msg *sarama.ConsumerMessage, smsEvent SMSCodeEvent) error {
	err := json.Unmarshal(msg.Value, &smsEvent)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.service.SendCode(ctx, smsEvent.TemplateId, smsEvent.Args, smsEvent.PhoneNumberSet...)
	if err != nil {
		s.l.Error("发送验证码失败", zap.Error(err))
		return err
	}
	s.l.Info("成功发送验证码", zap.String("phone", smsEvent.PhoneNumberSet[0]), zap.String("templateId", smsEvent.TemplateId))

	// TODO: 添加用户操作日志存储逻辑
	return nil
}

/*lockKey := "sms_lock_" + smsEvent.Phone
lock, err := s.rdb.SetNX(ctx, lockKey, "locked", time.Minute)
if err != nil {
	s.l.Error("获取锁失败", zap.Error(err))
	return err
}
if !lock.Val() {
	s.l.Warn("一分钟内只能发送一次验证码", zap.String("phone", smsEvent.Phone))
	return nil
}
code := utils.GenRandomCode(6) // 使用工具包生成随机验证码
smsEvent.Code = code
s.rdb.StoreVCode(ctx, smsEvent.Phone, smsEvent.Code, code)
// TODO: 调用第三方SMS服务发送验证码
/*if er := s.repo.SendCode(ctx, smsEvent.Phone, smsEvent.Code); er != nil {
	s.l.Error("发送验证码失败", zap.Error(er))
}修改前的逻辑
*/
