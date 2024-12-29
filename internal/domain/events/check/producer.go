package check

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

const TopicCheckEvent = "check_events"

type Producer interface {
	ProduceCheckEvent(evt CheckEvent) error
}

type CheckEvent struct {
	PostId  uint
	Uid     int64
	Title   string
	Content string
	PlateID int64
}

type SaramaCheckProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaCheckProducer(producer sarama.SyncProducer) Producer {
	return &SaramaCheckProducer{
		producer: producer,
	}
}

func (s *SaramaCheckProducer) ProduceCheckEvent(evt CheckEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicCheckEvent,
		Value: sarama.StringEncoder(val),
	})

	return err
}
