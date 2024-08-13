package publish

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const TopicPublishEvent = "linkme_publish_events"

type Producer interface {
	ProducePublishEvent(evt PublishEvent) error
}

type PublishEvent struct {
	PostId   uint   `json:"post_id"`
	AuthorID int64  `json:"author_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
	logger   *zap.Logger
}

func NewSaramaSyncProducer(producer sarama.SyncProducer, logger *zap.Logger) Producer {
	return &SaramaSyncProducer{
		producer: producer,
		logger:   logger,
	}
}

func (s *SaramaSyncProducer) ProducePublishEvent(evt PublishEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		s.logger.Error("Failed to marshal publish event", zap.Error(err))
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: TopicPublishEvent,
		Value: sarama.StringEncoder(val),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		s.logger.Error("Failed to send publish event message", zap.Error(err))
		return err
	}

	s.logger.Info("Publish event message sent", zap.Int32("partition", partition), zap.Int64("offset", offset))
	return nil
}
