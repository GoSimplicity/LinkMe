package comment

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

const TopicCommentEvent = "comment_events"

type Producer interface {
	ProduceCommentEvent(evt CommentEvent) error
}

type CommentEvent struct {
	BizId   int64
	PostId  uint
	Uid     int64
	Title   string
	Content string
	PlateID int64
	Status  uint8
}

type SaramaCommentProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaCommentProducer(producer sarama.SyncProducer) Producer {
	return &SaramaCommentProducer{
		producer: producer,
	}
}

func (s *SaramaCommentProducer) ProduceCommentEvent(evt CommentEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicCommentEvent,
		Value: sarama.StringEncoder(val),
	})

	return err
}
