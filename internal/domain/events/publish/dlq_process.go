package publish

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

type PublishDeadLetterConsumer struct {
	repo   repository.PostRepository
	client sarama.Client
	l      *zap.Logger
}

func NewPublishDeadLetterConsumer(
	repo repository.PostRepository,
	client sarama.Client,
	l *zap.Logger,
) *PublishDeadLetterConsumer {
	return &PublishDeadLetterConsumer{
		repo:   repo,
		client: client,
		l:      l,
	}
}

func (p *PublishDeadLetterConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("publish_dlq", p.client)
	if err != nil {
		p.l.Error("创建死信队列消费者组失败", zap.Error(err))
		return err
	}

	p.l.Info("DeadLetterConsumer 开始消费死信队列")

	// 启动死信队列消息消费
	go func() {
		defer cg.Close()
		for {
			select {
			case <-ctx.Done():
				p.l.Info("死信队列消费者停止")
				return
			default:
				if err := cg.Consume(ctx, []string{TopicDeadLetter}, &dlqConsumerGroupHandler{consumer: p}); err != nil {
					p.l.Error("死信队列消费循环出错", zap.Error(err))
					continue
				}
			}
		}
	}()

	return nil
}

type dlqConsumerGroupHandler struct {
	consumer *PublishDeadLetterConsumer
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
func (p *PublishDeadLetterConsumer) processDLQMessage(msg *sarama.ConsumerMessage) error {
	var evt PublishEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		p.l.Error("死信消息反序列化失败", zap.Error(err), zap.ByteString("message", msg.Value))
		return fmt.Errorf("死信消息反序列化失败: %w", err)
	}

	// 从死信队列获取原始主题、时间等信息
	originalTopic := ""
	for _, header := range msg.Headers {
		if string(header.Key) == "original_topic" {
			originalTopic = string(header.Value)
		}
	}

	p.l.Info("处理死信消息",
		zap.String("original_topic", originalTopic),
		zap.ByteString("message", msg.Value))

	// 验证事件参数
	if evt.PostId == 0 || evt.Uid == 0 {
		p.l.Error("死信消息参数无效",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid))
		return fmt.Errorf("无效的消息参数: post_id=%d, uid=%d", evt.PostId, evt.Uid)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 重新处理消息
	if err := p.handleDeadLetterMessage(ctx, &evt); err != nil {
		p.l.Error("处理死信消息失败", zap.Error(err))
		return err
	}

	return nil
}

// handleDeadLetterMessage 处理死信消息的具体业务逻辑
func (p *PublishDeadLetterConsumer) handleDeadLetterMessage(ctx context.Context, evt *PublishEvent) error {
	// 先检查帖子是否存在
	post, err := p.repo.GetPostById(ctx, evt.PostId, evt.Uid)
	if err != nil {
		p.l.Error("获取帖子失败", zap.Error(err), zap.Uint("post_id", evt.PostId))
		return fmt.Errorf("获取帖子失败: %w", err)
	}

	// 如果是草稿状态,说明审核被拒绝
	if evt.Status == domain.Draft {
		if err := p.repo.UpdateStatus(ctx, evt.PostId, evt.Uid, domain.Draft); err != nil {
			p.l.Error("更新帖子状态为草稿失败",
				zap.Error(err),
				zap.Uint("post_id", evt.PostId),
				zap.Int64("uid", evt.Uid))
			return fmt.Errorf("更新帖子状态为草稿失败: %w", err)
		}
		return nil
	}

	// 检查帖子状态
	if post.Status == domain.Published {
		p.l.Warn("帖子已发布", zap.Uint("post_id", evt.PostId), zap.Int64("uid", evt.Uid))
		return nil
	}

	// 更新帖子状态为已发布
	if err := p.repo.UpdateStatus(ctx, evt.PostId, evt.Uid, domain.Published); err != nil {
		p.l.Error("更新帖子状态失败",
			zap.Error(err),
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid))
		return fmt.Errorf("更新帖子状态失败: %w", err)
	}

	p.l.Info("成功处理死信消息",
		zap.Uint("post_id", evt.PostId),
		zap.Int64("uid", evt.Uid))

	return nil
}
