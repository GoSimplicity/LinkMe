package post

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

const TopicReadEvent = "linkme_read_events"

type Producer interface {
	ProduceReadEvent(evt ReadEvent) error
}

type ReadEvent struct {
	PostId uint
	Uid    int64
}

type BatchReadEvent struct {
	PostId []uint
	Uid    []int64
}

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{producer: producer}
}

func (s *SaramaSyncProducer) ProduceReadEvent(evt ReadEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicReadEvent,
		Value: sarama.StringEncoder(val),
	})
	return err
}
