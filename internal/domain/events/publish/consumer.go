package publish

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const (
	TopicDeadLetter = "publish_events_dlq" // 死信队列主题
	MaxRetries      = 3                    // 最大重试次数
)

type PublishPostEventConsumer struct {
	repo    repository.PostRepository
	client  sarama.Client
	l       *zap.Logger
	dlqProd sarama.SyncProducer // 死信队列生产者
}

type consumerGroupHandler struct {
	consumer *PublishPostEventConsumer
}

func NewPublishPostEventConsumer(repo repository.PostRepository, client sarama.Client, dlqProd sarama.SyncProducer, l *zap.Logger) *PublishPostEventConsumer {
	return &PublishPostEventConsumer{
		repo:    repo,
		client:  client,
		l:       l,
		dlqProd: dlqProd,
	}
}

// Start 启动消费者，并开始消费 Kafka 中的消息
func (p *PublishPostEventConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("publish_event", p.client)
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

func (c *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 处理每一条消息
		if err := c.consumer.processMessage(sess, msg); err != nil {
			c.consumer.l.Error("处理消息失败", zap.Error(err), zap.ByteString("message", msg.Value))
			// 发送到死信队列
			if err := c.consumer.sendToDLQ(msg); err != nil {
				c.consumer.l.Error("发送到死信队列失败", zap.Error(err))
			}
			continue // 确保继续处理下一条消息
		}
		sess.MarkMessage(msg, "") // 只有成功处理后才标记
	}

	return nil
}

// sendToDLQ 发送消息到死信队列
func (p *PublishPostEventConsumer) sendToDLQ(msg *sarama.ConsumerMessage) error {
	dlqMsg := &sarama.ProducerMessage{
		Topic: TopicDeadLetter,
		Value: sarama.ByteEncoder(msg.Value),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("original_topic"),
				Value: []byte(msg.Topic),
			},
			{
				Key:   []byte("error_time"),
				Value: []byte(time.Now().Format(time.RFC3339)),
			},
		},
	}

	_, _, err := p.dlqProd.SendMessage(dlqMsg)
	return err
}

// processMessage 处理从 Kafka 消费的消息
func (p *PublishPostEventConsumer) processMessage(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error {
	// 参数校验
	if err := p.validateMessage(msg); err != nil {
		return err
	}

	var event PublishEvent
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
func (p *PublishPostEventConsumer) validateMessage(msg *sarama.ConsumerMessage) error {
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
func (p *PublishPostEventConsumer) validateEvent(event *PublishEvent) error {
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
func (p *PublishPostEventConsumer) handleEvent(ctx context.Context, event *PublishEvent) error {
	if ctx == nil {
		return errors.New("context为空")
	}

	// 先检查帖子是否存在
	post, err := p.repo.GetPostById(ctx, event.PostId, event.Uid)
	if err != nil {
		p.l.Error("获取帖子失败", zap.Error(err), zap.Uint("post_id", event.PostId))
		return fmt.Errorf("获取帖子失败: %w", err)
	}

	// 如果是草稿状态,说明审核被拒绝
	if event.Status == domain.Draft {
		if err := p.repo.UpdateStatus(ctx, event.PostId, event.Uid, domain.Draft); err != nil {
			p.l.Error("更新帖子状态为草稿失败",
				zap.Error(err),
				zap.Uint("post_id", event.PostId),
				zap.Int64("uid", event.Uid))
			return fmt.Errorf("更新帖子状态为草稿失败: %w", err)
		}
		return nil
	}

	// 检查帖子状态
	if post.Status == domain.Published {
		p.l.Warn("帖子已发布", zap.Uint("post_id", event.PostId), zap.Int64("uid", event.Uid))
		return nil
	}

	// 更新帖子状态为已发布
	if err := p.repo.UpdateStatus(ctx, event.PostId, event.Uid, domain.Published); err != nil {
		p.l.Error("更新帖子状态失败",
			zap.Error(err),
			zap.Uint("post_id", event.PostId),
			zap.Int64("uid", event.Uid))
		return fmt.Errorf("更新帖子状态失败: %w", err)
	}

	return nil
}
