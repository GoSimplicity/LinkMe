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

type PostDeadLetterConsumer struct {
	repo    repository.InteractiveRepository
	hisRepo repository.HistoryRepository
	client  sarama.Client
	l       *zap.Logger
}

func NewPostDeadLetterConsumer(
	repo repository.InteractiveRepository,
	hisRepo repository.HistoryRepository,
	client sarama.Client,
	l *zap.Logger,
) *PostDeadLetterConsumer {
	return &PostDeadLetterConsumer{
		repo:    repo,
		hisRepo: hisRepo,
		client:  client,
		l:       l,
	}
}

func (i *PostDeadLetterConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("post_dlq", i.client)
	if err != nil {
		i.l.Error("创建死信队列消费者组失败", zap.Error(err))
		return err
	}

	i.l.Info("DeadLetterConsumer 开始消费死信队列")

	// 启动死信队列消息消费
	go func() {
		defer cg.Close()
		for {
			select {
			case <-ctx.Done():
				i.l.Info("死信队列消费者停止")
				return
			default:
				if err := cg.Consume(ctx, []string{TopicDeadLetter}, &dlqConsumerGroupHandler{consumer: i}); err != nil {
					i.l.Error("死信队列消费循环出错", zap.Error(err))
					continue
				}
			}
		}
	}()

	return nil
}

type dlqConsumerGroupHandler struct {
	consumer *PostDeadLetterConsumer
}

func (h *dlqConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *dlqConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *dlqConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	const (
		maxRetries   = 5
		baseWaitTime = 5 * time.Second
	)

	for msg := range claim.Messages() {
		var err error
		for i := 0; i < maxRetries; i++ {
			if err = h.consumer.processDLQMessage(msg); err == nil {
				break
			}

			if i < maxRetries-1 { // 最后一次失败不需要记录重试日志
				h.consumer.l.Error("处理死信消息失败,准备重试",
					zap.Error(err),
					zap.Int("重试次数", i+1),
					zap.Int("剩余重试次数", maxRetries-i-1),
					zap.ByteString("message", msg.Value))

				// 指数退避策略,等待时间随重试次数指数增长
				waitTime := baseWaitTime * time.Duration(1<<uint(i))
				time.Sleep(waitTime)
			}
		}

		if err != nil {
			h.consumer.l.Error("处理死信消息最终失败",
				zap.Error(err),
				zap.Int("重试次数", maxRetries),
				zap.ByteString("message", msg.Value))
		} else {
			h.consumer.l.Info("死信消息处理成功",
				zap.ByteString("message", msg.Value))
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

// processDLQMessage 处理死信队列中的消息
func (i *PostDeadLetterConsumer) processDLQMessage(msg *sarama.ConsumerMessage) error {
	var evt ReadEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		i.l.Error("死信消息反序列化失败", zap.Error(err), zap.ByteString("message", msg.Value))
		return fmt.Errorf("死信消息反序列化失败: %w", err)
	}

	// 从死信队列获取原始主题、时间等信息
	originalTopic := ""
	for _, header := range msg.Headers {
		if string(header.Key) == "original_topic" {
			originalTopic = string(header.Value)
		}
	}

	i.l.Info("处理死信消息",
		zap.String("original_topic", originalTopic),
		zap.ByteString("message", msg.Value))

	// 验证事件参数
	if evt.PostId == 0 || evt.Uid == 0 {
		i.l.Error("死信消息参数无效",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid))
		return fmt.Errorf("无效的消息参数: post_id=%d, uid=%d", evt.PostId, evt.Uid)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 重新处理消息
	if err := i.handleDeadLetterMessage(ctx, &evt); err != nil {
		i.l.Error("处理死信消息失败", zap.Error(err))
		return err
	}

	return nil
}

// handleDeadLetterMessage 处理死信消息的具体业务逻辑
func (i *PostDeadLetterConsumer) handleDeadLetterMessage(ctx context.Context, evt *ReadEvent) error {
	post := domain.Post{
		ID:      evt.PostId,
		Content: evt.Content,
		Title:   evt.Title,
		Tags:    strconv.FormatInt(evt.PlateID, 10),
		Uid:     evt.Uid,
	}

	// 保存历史记录
	if err := i.hisRepo.SetHistory(ctx, post); err != nil {
		i.l.Error("保存历史记录失败",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid),
			zap.Error(err))
		return fmt.Errorf("保存历史记录失败: %w", err)
	}

	// 增加阅读计数
	if err := i.repo.IncrReadCnt(ctx, evt.PostId); err != nil {
		i.l.Error("增加阅读计数失败",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid),
			zap.Error(err))
		return fmt.Errorf("增加阅读计数失败: %w", err)
	}

	i.l.Info("成功处理死信消息",
		zap.Uint("post_id", evt.PostId),
		zap.Int64("uid", evt.Uid))

	return nil
}
