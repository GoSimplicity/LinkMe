package sms

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const TopicSMS = "linkme_sms_events"

type Producer interface {
	ProduceSMSCode(ctx context.Context, evt SMSCodeEvent) error
}

// SMSCodeEvent 代表单个短信验证码事件
type SMSCodeEvent struct {
	Number string
}

// SaramaSyncProducer 实现Producer接口的结构体
type SaramaSyncProducer struct {
	producer sarama.SyncProducer
	logger   *zap.Logger
}

// NewSaramaSyncProducer 创建一个新的SaramaSyncProducer实例
func NewSaramaSyncProducer(producer sarama.SyncProducer, logger *zap.Logger) Producer {
	return &SaramaSyncProducer{
		producer: producer,
		logger:   logger,
	}
}

// ProduceSMSCode 发送短信验证码事件到Kafka
func (s *SaramaSyncProducer) ProduceSMSCode(ctx context.Context, evt SMSCodeEvent) error {
	// 序列化事件
	data, err := json.Marshal(evt)
	if err != nil {
		s.logger.Error("序列化事件失败", zap.Error(err))
		return err
	}
	// 发送消息到Kafka
	partition, offset, err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicSMS,
		Value: sarama.StringEncoder(data),
	})
	if err != nil {
		s.logger.Error("发送信息到Kafka失败", zap.Error(err))
		return err
	}
	s.logger.Info("成功发送消息到Kafka", zap.String("topic", TopicSMS), zap.Int32("partition", partition), zap.Int64("offset", offset))
	return nil
}
