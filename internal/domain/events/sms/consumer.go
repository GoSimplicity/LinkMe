package sms

import (
	"context"
	"encoding/json"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"time"

	//"LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/pkg/samarap"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type SMSConsumer struct {
	repo   repository.SmsRepository
	client sarama.Client
	l      *zap.Logger
	rdb    cache.SMSCache
}

func NewSMSConsumer(repo repository.SmsRepository, client sarama.Client, l *zap.Logger, rdb cache.SMSCache) *SMSConsumer {
	return &SMSConsumer{repo: repo, client: client, l: l, rdb: rdb}
}

func (s *SMSConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("sms_consumer_group", s.client)
	s.l.Info("SMSConsumer 开始消费")
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
		s.l.Error("json.Unmarshal 失败", zap.Any("msg", msg))
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.SendCode(ctx, smsEvent.Number)
}
