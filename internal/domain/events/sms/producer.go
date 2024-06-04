package post

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const TopicSMS = "sms_events"

type Producer interface {
	ProduceSMSCode(ctx context.Context, evt SMSCodeEvent) error
}

// SMSCodeEvent 代表单个短信验证码事件
type SMSCodeEvent struct {
	Phone string
	Code  string
}

// SaramaSyncProducer 实现Producer接口的结构体
type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

// NewSaramaSyncProducer 创建一个新的SaramaSyncProducer实例
func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{producer: producer}
}

// ProduceSMSCode 发送短信验证码事件到Kafka
func (s *SaramaSyncProducer) ProduceSMSCode(ctx context.Context, evt SMSCodeEvent) error {
	// 序列化事件
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	// 发送消息到Kafka
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicSMS,
		Value: sarama.StringEncoder(data),
	})
	return err
}
