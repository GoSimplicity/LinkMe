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
	PostId  uint
	Uid     int64
	Title   string
	Content string
	PlateID int64
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

	// 创建消息
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicReadEvent,            // 订阅主题
		Value: sarama.StringEncoder(val), // 消息内容
	})

	return err
}
