package post

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/pkg/samarap"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type ReadEventConsumer struct {
	repo    repository.InteractiveRepository
	hisRepo repository.HistoryRepository
	client  sarama.Client
	l       *zap.Logger
}

func NewReadEventConsumer(repo repository.InteractiveRepository,
	client sarama.Client, l *zap.Logger, hisRepo repository.HistoryRepository) *ReadEventConsumer {
	return &ReadEventConsumer{
		repo:    repo,
		hisRepo: hisRepo,
		client:  client,
		l:       l,
	}
}

func (i *ReadEventConsumer) Start(_ context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}

	i.l.Info("PostConsumer 开始消费")

	go func() {
		attempts := 0   // 当前重试次数
		maxRetries := 3 // 最大重试次数

		for attempts < maxRetries {
			// 启动消费者组并开始消费指定主题 read_post 的消息
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
func (i *ReadEventConsumer) BatchConsume(_ []*sarama.ConsumerMessage, events []ReadEvent) error {
	// 提前计算需要的上下文超时时间
	//ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel() // 确保上下文在函数结束时被取消

	// 提前分配内存，避免循环内的重复分配
	posts := make([]domain.Post, len(events))
	bizs := make([]string, len(events))
	bizIds := make([]uint, len(events))

	for idx, evt := range events {
		posts[idx] = domain.Post{
			ID:       evt.PostId,
			Content:  evt.Content,
			Title:    evt.Title,
			Tags:     strconv.FormatInt(evt.PlateID, 10),
			AuthorID: evt.Uid,
		}
		bizs[idx] = "post"
		bizIds[idx] = evt.PostId
	}

	// 保存历史记录
	if err := i.hisRepo.SetHistory(ctx, posts); err != nil {
		i.l.Error("保存历史记录失败", zap.Error(err))
		return err
	}

	// 增加阅读计数
	if err := i.repo.BatchIncrReadCnt(ctx, bizs, bizIds); err != nil {
		i.l.Error("增加阅读计数失败", zap.Error(err))
		return err
	}

	return nil
}
