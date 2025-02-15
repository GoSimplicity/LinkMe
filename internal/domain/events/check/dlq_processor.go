package check

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type CheckDeadLetterConsumer struct {
	checkRepo repository.CheckRepository
	client    sarama.Client
	l         *zap.Logger
}

func NewCheckDeadLetterConsumer(
	checkRepo repository.CheckRepository,
	client sarama.Client,
	l *zap.Logger,
) *CheckDeadLetterConsumer {
	return &CheckDeadLetterConsumer{
		checkRepo: checkRepo,
		client:    client,
		l:         l,
	}
}

func (c *CheckDeadLetterConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("check_dlq", c.client)
	if err != nil {
		c.l.Error("创建死信队列消费者组失败", zap.Error(err))
		return err
	}

	c.l.Info("DeadLetterConsumer 开始消费死信队列")

	go func() {
		defer cg.Close()
		for {
			select {
			case <-ctx.Done():
				c.l.Info("死信队列消费者停止")
				return
			default:
				if err := cg.Consume(ctx, []string{TopicDeadLetter}, &dlqConsumerGroupHandler{consumer: c}); err != nil {
					c.l.Error("死信队列消费循环出错", zap.Error(err))
					continue
				}
			}
		}
	}()

	return nil
}

type dlqConsumerGroupHandler struct {
	consumer *CheckDeadLetterConsumer
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

			if i < maxRetries-1 {
				h.consumer.l.Error("处理死信消息失败,准备重试",
					zap.Error(err),
					zap.Int("重试次数", i+1),
					zap.Int("剩余重试次数", maxRetries-i-1),
					zap.ByteString("message", msg.Value))

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

func (c *CheckDeadLetterConsumer) processDLQMessage(msg *sarama.ConsumerMessage) error {
	var evt CheckEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		c.l.Error("死信消息反序列化失败", zap.Error(err), zap.ByteString("message", msg.Value))
		return fmt.Errorf("死信消息反序列化失败: %w", err)
	}

	originalTopic := ""
	for _, header := range msg.Headers {
		if string(header.Key) == "original_topic" {
			originalTopic = string(header.Value)
		}
	}

	c.l.Info("处理死信消息",
		zap.String("original_topic", originalTopic),
		zap.ByteString("message", msg.Value))

	if (evt.PostId == 0 && evt.BizId == 1) || evt.Uid == 0 || (evt.Title == "" && evt.BizId == 1) || evt.Content == "" {
		c.l.Error("死信消息参数无效",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid))
		return fmt.Errorf("无效的消息参数: post_id=%d, uid=%d", evt.PostId, evt.Uid)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.handleDeadLetterMessage(ctx, &evt); err != nil {
		c.l.Error("处理死信消息失败", zap.Error(err))
		return err
	}

	return nil
}

func (c *CheckDeadLetterConsumer) handleDeadLetterMessage(ctx context.Context, evt *CheckEvent) error {
	check := domain.Check{
		BizId:   evt.BizId,
		PostID:  evt.PostId,
		Uid:     evt.Uid,
		Title:   evt.Title,
		Content: evt.Content,
		PlateID: evt.PlateID,
	}

	code, err := c.checkRepo.Create(ctx, check)
	if err != nil {
		c.l.Error("创建审核记录失败",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid),
			zap.Error(err))
		return fmt.Errorf("创建审核记录失败: %w", err)
	}

	if code == -1 {
		c.l.Info("审核记录已存在",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid))
		return nil
	}

	c.l.Info("成功处理死信消息",
		zap.Uint("post_id", evt.PostId),
		zap.Int64("uid", evt.Uid))

	return nil
}
