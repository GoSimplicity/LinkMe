package comment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/publish"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

const TopicPublishEvent = "publish_events"

type PublishCommentEventConsumer struct {
	repo   repository.CommentRepository
	client sarama.Client
	l      *zap.Logger
}

type consumerGroupHandler struct {
	consumer *PublishCommentEventConsumer
}

func (c consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 处理每一条消息
		if err := c.consumer.processMessage(msg); err != nil {
			c.consumer.l.Error("处理消息失败", zap.Error(err), zap.ByteString("message", msg.Value))
			// 发送到死信队列
			if err != nil {
				c.consumer.l.Error("处理消息失败", zap.Error(err))
			}
			continue // 确保继续处理下一条消息
		}
		sess.MarkMessage(msg, "") // 只有成功处理后才标记
	}

	return nil
}

func NewPublishCommentEventConsumer(repo repository.CommentRepository, client sarama.Client, l *zap.Logger) *PublishCommentEventConsumer {
	return &PublishCommentEventConsumer{
		repo:   repo,
		client: client,
		l:      l,
	}
}

// Start 启动消费者，并开始消费 Kafka 中的消息
func (p *PublishCommentEventConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("publish_comment_event", p.client)
	if err != nil {
		p.l.Error("创建消费者组失败", zap.Error(err))
		return err
	}

	p.l.Info("PublishConsumer 开始消费")

	go func() {
		defer cg.Close()
		for {
			select {
			case <-ctx.Done():
				p.l.Info("消费者停止")
				return
			default:
				// 开始消费指定的 Kafka 主题
				if err := cg.Consume(ctx, []string{TopicPublishEvent}, &consumerGroupHandler{consumer: p}); err != nil {
					p.l.Error("消费循环出错", zap.Error(err))
					continue
				}
			}
		}
	}()

	return nil
}

// processMessage 处理从 Kafka 消费的消息
func (p *PublishCommentEventConsumer) processMessage(msg *sarama.ConsumerMessage) error {
	// 参数校验
	if err := p.validateMessage(msg); err != nil {
		return err
	}

	var event publish.PublishEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		p.l.Error("反序列化消息失败",
			zap.Error(err),
			zap.String("message", string(msg.Value)))
		return fmt.Errorf("反序列化消息失败: %w", err)
	}

	// 参数校验
	if err := p.validateEvent(&event); err != nil {
		return err
	}
	// 判断这个审核业务类型是否是评论
	if event.BizId != 2 {
		p.l.Warn("无效的审核业务类型",
			zap.Int64("bizid", event.BizId),
			zap.Uint("post_id", event.PostId),
			zap.Int64("uid", event.Uid))
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 处理消息
	if err := p.handleEvent(ctx, &event); err != nil {
		return err
	}

	p.l.Info("消息处理成功",
		zap.Uint("post_id", event.PostId),
		zap.Int64("uid", event.Uid),
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset))

	return nil
}

// validateMessage 验证消息是否有效
func (p *PublishCommentEventConsumer) validateMessage(msg *sarama.ConsumerMessage) error {
	if msg == nil || msg.Value == nil {
		p.l.Error("消息为空")
		return errors.New("消息为空")
	}

	p.l.Debug("开始处理消息",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset))

	return nil
}

// validateEvent 验证事件参数是否有效
func (p *PublishCommentEventConsumer) validateEvent(event *publish.PublishEvent) error {
	if event == nil {
		return errors.New("事件为空")
	}

	if event.PostId == 0 || event.Uid == 0 {
		p.l.Error("消息参数无效",
			zap.Uint("post_id", event.PostId),
			zap.Int64("uid", event.Uid))
		return fmt.Errorf("无效的消息参数: post_id=%d, uid=%d", event.PostId, event.Uid)
	}

	if event.Status > domain.Published { // 检查状态值是否有效
		return fmt.Errorf("无效的状态值: %d", event.Status)
	}
	return nil
}

// handleEvent 处理发布事件
func (p *PublishCommentEventConsumer) handleEvent(ctx context.Context, event *publish.PublishEvent) error {
	if ctx == nil {
		return errors.New("context为空")
	}
	comment, err := p.repo.FindCommentByCommentId(ctx, int64(event.PostId))
	if err != nil {
		return err
	}
	// 更改评论状态
	comment.Status = domain.Published
	if err := p.repo.UpdateComment(ctx, comment); err != nil {
		return err
	}
	return nil
}
