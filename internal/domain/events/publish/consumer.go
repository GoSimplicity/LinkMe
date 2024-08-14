package publish

import (
	"context"
	"encoding/json"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type PublishPostEventConsumer struct {
	repo   repository.CheckRepository
	client sarama.Client
	l      *zap.Logger
}

type consumerGroupHandler struct {
	consumer *PublishPostEventConsumer
}

// NewPublishPostEventConsumer 创建一个新的 PublishPostEventConsumer 实例
func NewPublishPostEventConsumer(repo repository.CheckRepository, client sarama.Client, l *zap.Logger) *PublishPostEventConsumer {
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
			c.consumer.l.Error("Failed to process message", zap.Error(err), zap.ByteString("message", msg.Value))
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
	p.l.Info("PublishConsumer 开始消费")
	if err != nil {
		return err
	}

	go func() {
		for {
			// 开始消费指定的 Kafka 主题
			err := cg.Consume(ctx, []string{TopicPublishEvent}, &consumerGroupHandler{consumer: p})
			if err != nil {
				p.l.Error("Error occurred in consume loop", zap.Error(err))
				continue // 继续循环，即使出现错误，也不会退出
			}
		}
	}()

	return nil
}

// processMessage 处理从 Kafka 消费的消息
func (p *PublishPostEventConsumer) processMessage(_ sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error {
	var event PublishEvent
	// 将消息内容反序列化为 PublishEvent 结构体
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	// 创建检查记录
	check := domain.Check{
		Content: event.Content,
		PostID:  event.PostId,
		Title:   event.Title,
		UserID:  event.AuthorID,
	}

	// 将检查记录保存到数据库中
	checkId, err := p.repo.Create(context.Background(), check)
	if err != nil {
		p.l.Error("Failed to create check", zap.Error(err), zap.Int64("check_id", checkId))
		return err
	}

	p.l.Info("Successfully processed message", zap.Uint("post_id", event.PostId), zap.Int64("check_id", checkId))

	return nil
}
