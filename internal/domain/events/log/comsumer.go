package log

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/pkg/samarap"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

const TopicZapLogs = "linkme_elk_events"

type EsLogsConsumer struct {
	client sarama.Client
	rs     repository.SearchRepository
	l      *zap.Logger
}

type Event struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

func NewEsLogsConsumer(client sarama.Client, rs repository.SearchRepository, l *zap.Logger) *EsLogsConsumer {
	return &EsLogsConsumer{
		client: client,
		rs:     rs,
		l:      l}
}

func (es *EsLogsConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("zap_logs_consumer_group", es.client)
	if err != nil {
		return err
	}

	go func() {
		es.l.Info("Es logs consumer 开始消费")
		for {
			err := cg.Consume(ctx, []string{TopicZapLogs}, samarap.NewBatchHandler[Event](es.l, es.BatchHandleLogs))
			if err != nil {
				es.l.Error("消费失败", zap.Error(err))
				time.Sleep(time.Second * 5) //重试间隔
				continue
			}
			break
		}
		es.l.Info("Es logs consumer 消费成功")
	}()
	return nil
}

func (es *EsLogsConsumer) BatchHandleLogs(_ []*sarama.ConsumerMessage, events []Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logs := make([]domain.ReadEvent, len(events))
	for i, event := range events {
		logs[i] = domain.ReadEvent{
			Timestamp: event.Timestamp,
			Level:     event.Level,
			Message:   event.Message,
		}
	}
	return es.rs.BulkInputLogs(ctx, logs)
}
