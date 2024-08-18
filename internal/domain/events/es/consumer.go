package es

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type EsConsumer struct {
	client sarama.Client
	l      *zap.Logger
}

type Event struct {
	Type     string              `json:"type"`
	Database string              `json:"database"`
	Table    string              `json:"table"`
	Data     []map[string]string `json:"data"`
}

type consumerGroupHandler struct {
	r *EsConsumer
}

func NewEsConsumer(client sarama.Client, l *zap.Logger) *EsConsumer {
	// 创建MongoDB客户端
	return &EsConsumer{
		client: client,
		l:      l,
	}
}

func (r *EsConsumer) Start(_ context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("es_consumer_group", r.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			if err := cg.Consume(context.Background(), []string{"linkme_es"}, &consumerGroupHandler{r: r}); err != nil {
				r.l.Error("退出了消费循环异常", zap.Error(err))
			}
		}
	}()
	return nil
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.r.Consume(sess, msg)
	}

	return nil
}

func (r *EsConsumer) Consume(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var e Event

	if err := json.Unmarshal(msg.Value, &e); err != nil {
		panic(err)
	}

	sess.MarkMessage(msg, "")
}
