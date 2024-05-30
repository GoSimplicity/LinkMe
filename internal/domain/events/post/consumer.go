package post

import (
	"LinkMe/internal/repository"
	"LinkMe/pkg/samarap"
	"context"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

type InteractiveReadEventConsumer struct {
	repo   repository.InteractiveRepository
	client sarama.Client
	l      *zap.Logger
}

func NewInteractiveReadEventConsumer(repo repository.InteractiveRepository,
	client sarama.Client, l *zap.Logger) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{repo: repo, client: client, l: l}
}

func (i *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicReadEvent},
			samarap.NewBatchHandler[ReadEvent](i.l, i.BatchConsume))
		if er != nil {
			i.l.Error("退出消费", zap.Error(er))
		}
	}()
	return err
}

func (i *InteractiveReadEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage,
	events []ReadEvent) error {
	bizs := make([]string, 0, len(events))
	bizIds := make([]int64, 0, len(events))
	for _, evt := range events {
		bizs = append(bizs, "post")
		bizIds = append(bizIds, evt.PostId)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.BatchIncrReadCnt(ctx, bizs, bizIds)
}
