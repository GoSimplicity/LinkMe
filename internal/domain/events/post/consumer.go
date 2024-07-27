package post

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/pkg/samarap"
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
	return &InteractiveReadEventConsumer{
		repo:   repo,
		client: client,
		l:      l,
	}
}

func (i *InteractiveReadEventConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		attempts := 0   // 当前重试次数
		maxRetries := 3 // 最大重试次数
		for attempts < maxRetries {
			// 启动消费者组并开始消费指定主题read_post的消息
			er := cg.Consume(context.Background(), []string{TopicReadEvent}, samarap.NewBatchHandler[ReadEvent](i.l, i.BatchConsume))
			if er != nil {
				i.l.Error("消费错误", zap.Error(er), zap.Int("重试次数", attempts+1))
				attempts++
				time.Sleep(time.Second * time.Duration(attempts)) // 退避策略，每次重试后等待的时间增加
				continue
			}
			break
		}
		if attempts >= maxRetries {
			i.l.Error("达到最大重试次数，退出消费")
		}
	}()
	return nil
}

// BatchConsume 处理函数，处理批次消息
func (i *InteractiveReadEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage, events []ReadEvent) error {
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
