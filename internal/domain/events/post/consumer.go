package post

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type EventConsumer struct {
	repo    repository.InteractiveRepository
	hisRepo repository.HistoryRepository
	client  sarama.Client
	l       *zap.Logger
}

func NewEventConsumer(
	repo repository.InteractiveRepository,
	hisRepo repository.HistoryRepository,
	client sarama.Client,
	l *zap.Logger,
) *EventConsumer {
	return &EventConsumer{
		repo:    repo,
		hisRepo: hisRepo,
		client:  client,
		l:       l,
	}
}

func (i *EventConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("post", i.client)
	if err != nil {
		i.l.Error("创建消费者组失败", zap.Error(err))
		return err
	}

	i.l.Info("PostConsumer 开始消费")

	// 启动阅读消息消费
	go func() {
		defer cg.Close()
		for {
			select {
			case <-ctx.Done():
				i.l.Info("阅读消息消费者停止")
				return
			default:
				if err := cg.Consume(ctx, []string{TopicReadEvent}, &consumerGroupHandler{r: i}); err != nil {
					i.l.Error("阅读消息消费循环出错", zap.Error(err))
					continue
				}
			}
		}
	}()

	return nil
}

type consumerGroupHandler struct {
	r *EventConsumer
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.r.ConsumeRead(msg); err != nil {
			h.r.l.Error("处理阅读消息失败", zap.Error(err))
		} else {
			sess.MarkMessage(msg, "")
		}
	}
	return nil
}

// ConsumeRead 处理单条阅读消息
func (i *EventConsumer) ConsumeRead(msg *sarama.ConsumerMessage) error {
	var evt ReadEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return fmt.Errorf("解析消息失败: %w", err)
	}

	// 参数校验
	if evt.PostId == 0 || evt.Uid == 0 {
		i.l.Warn("无效的阅读事件",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid))
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	post := domain.Post{
		ID:      evt.PostId,
		Content: evt.Content,
		Title:   evt.Title,
		Tags:    strconv.FormatInt(evt.PlateID, 10),
		Uid:     evt.Uid,
	}

	// 保存历史记录
	if err := i.hisRepo.SetHistory(ctx, []domain.Post{post}); err != nil {
		i.l.Error("保存历史记录失败", zap.Error(err))
		return fmt.Errorf("保存历史记录失败: %w", err)
	}

	// 增加阅读计数
	if err := i.repo.IncrReadCnt(ctx, evt.PostId); err != nil {
		i.l.Error("增加阅读计数失败", zap.Error(err))
		return fmt.Errorf("增加阅读计数失败: %w", err)
	}

	i.l.Info("处理阅读事件成功",
		zap.Uint("post_id", evt.PostId),
		zap.Int64("uid", evt.Uid))

	return nil
}
