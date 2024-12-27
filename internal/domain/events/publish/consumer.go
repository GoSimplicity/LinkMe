package publish

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type PublishPostEventConsumer struct {
	repo   repository.PostRepository
	client sarama.Client
	l      *zap.Logger
}

type consumerGroupHandler struct {
	consumer *PublishPostEventConsumer
}

func NewPublishPostEventConsumer(repo repository.PostRepository, client sarama.Client, l *zap.Logger) *PublishPostEventConsumer {
	return &PublishPostEventConsumer{
		repo:   repo,
		client: client,
		l:      l,
	}
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
		} else {
			// 如果消息处理成功，标记消息为已消费
			sess.MarkMessage(msg, "")
		}
	}
	return nil
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

// processMessage 处理从 Kafka 消费的消息
func (p *PublishPostEventConsumer) processMessage(_ sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error {
	// 添加日志记录消息处理开始
	p.l.Debug("开始处理消息",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset))

	var event PublishEvent
	// 将消息内容反序列化为 PublishEvent 结构体
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		p.l.Error("反序列化消息失败",
			zap.Error(err),
			zap.String("message", string(msg.Value)))
		return fmt.Errorf("反序列化消息失败: %w", err)
	}

	// 添加基本的参数校验
	if event.PostId == 0 || event.Uid == 0 {
		p.l.Error("消息参数无效",
			zap.Uint("post_id", event.PostId),
			zap.Int64("uid", event.Uid))
		return fmt.Errorf("无效的消息参数: post_id=%d, uid=%d", event.PostId, event.Uid)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	post, err := p.repo.GetPublishPostById(ctx, event.PostId)
	if err != nil && err != dao.ErrPostNotFound {
		p.l.Error("获取帖子失败", zap.Error(err), zap.Uint("post_id", event.PostId))
		return fmt.Errorf("获取帖子失败: %w", err)
	}

	if post.ID != 0 {
		p.l.Warn("帖子已发布", zap.Uint("post_id", event.PostId), zap.Int64("uid", event.Uid))
		return nil
	}

	if err := p.repo.UpdateStatus(ctx, event.PostId, event.Uid, domain.Published); err != nil {
		p.l.Error("更新帖子状态失败",
			zap.Error(err),
			zap.Uint("post_id", event.PostId),
			zap.Int64("uid", event.Uid))
		return fmt.Errorf("更新帖子状态失败: %w", err)
	}

	p.l.Info("消息处理成功",
		zap.Uint("post_id", event.PostId),
		zap.Int64("uid", event.Uid),
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset))

	return nil
}
