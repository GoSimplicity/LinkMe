package email

import (
	"LinkMe/internal/service"
	"context"
	"encoding/json"
	"time"

	"LinkMe/pkg/samarap"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type EmailConsumer struct {
	service service.EmailService
	client  sarama.Client
	l       *zap.Logger
}

func NewEmailConsumer(service service.EmailService, client sarama.Client, l *zap.Logger) *EmailConsumer {
	return &EmailConsumer{service: service, client: client, l: l}
}

func (e *EmailConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("email_consumer_group", e.client)
	if err != nil {
		return err
	}
	go func() {
		attempts := 0
		maxRetries := 3
		e.l.Info("emailConsumer 开始消费")
		for attempts < maxRetries {
			er := cg.Consume(ctx, []string{TopicEmail}, samarap.NewHandler(e.l, e.HandleMessage))
			if er != nil {
				e.l.Error("消费错误", zap.Error(er), zap.Int("重试次数", attempts+1))
				attempts++
				time.Sleep(time.Second * time.Duration(attempts))
				continue
			}
			break
		}
		if attempts >= maxRetries {
			e.l.Error("达到最大重试次数，退出消费")
		}
	}()
	return nil
}

func (e *EmailConsumer) HandleMessage(msg *sarama.ConsumerMessage, emailEvent EmailEvent) error {
	err := json.Unmarshal(msg.Value, &emailEvent)
	if err != nil {
		e.l.Error("json.Unmarshal 失败", zap.Any("msg", msg))
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return e.service.SendCode(ctx, emailEvent.Email)
}
