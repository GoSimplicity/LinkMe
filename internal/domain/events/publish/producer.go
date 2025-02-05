package publish

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const TopicPublishEvent = "publish_events"

type Producer interface {
	ProducePublishEvent(evt PublishEvent) error
}

type PublishEvent struct {
	PostId uint  `json:"post_id"`
	Uid    int64 `json:"uid"`
	Status uint8 `json:"status"`
	BizId  int64 `json:"biz_id"`
}

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
	l        *zap.Logger
}

func NewSaramaSyncProducer(producer sarama.SyncProducer, l *zap.Logger) Producer {
	return &SaramaSyncProducer{
		producer: producer,
		l:        l,
	}
}

func (s *SaramaSyncProducer) ProducePublishEvent(evt PublishEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		s.l.Error("Failed to marshal publish event", zap.Error(err))
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: TopicPublishEvent,
		Value: sarama.StringEncoder(val),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		s.l.Error("Failed to send publish event message", zap.Error(err))
		return err
	}

	s.l.Info("Publish event message sent", zap.Int32("partition", partition), zap.Int64("offset", offset))
	return nil
}
