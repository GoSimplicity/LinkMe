package check

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
	TopicDeadLetter = "check_events_dlq" // 死信队列主题
	MaxRetries      = 3                  // 最大重试次数
)

type CheckEventConsumer struct {
	checkRepo repository.CheckRepository
	client    sarama.Client
	l         *zap.Logger
	dlqProd   sarama.SyncProducer // 死信队列生产者
}

func NewCheckEventConsumer(
	checkRepo repository.CheckRepository,
	client sarama.Client,
	dlqProd sarama.SyncProducer,
	l *zap.Logger,
) *CheckEventConsumer {
	return &CheckEventConsumer{
		checkRepo: checkRepo,
		client:    client,
		l:         l,
		dlqProd:   dlqProd,
	}
}

func (c *CheckEventConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("check", c.client)
	if err != nil {
		c.l.Error("创建消费者组失败", zap.Error(err))
		return err
	}

	c.l.Info("CheckConsumer 开始消费")

	go func() {
		defer cg.Close()
		for {
			select {
			case <-ctx.Done():
				c.l.Info("审核消息消费者停止")
				return
			default:
				if err := cg.Consume(ctx, []string{TopicCheckEvent}, &consumerGroupHandler{consumer: c}); err != nil {
					c.l.Error("审核消息消费循环出错", zap.Error(err))
					continue
				}
			}
		}
	}()

	return nil
}

type consumerGroupHandler struct {
	consumer *CheckEventConsumer
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 处理每一条消息
		if err := h.consumer.processMessage(sess, msg); err != nil {
			h.consumer.l.Error("处理消息失败", zap.Error(err), zap.ByteString("message", msg.Value))
			// 发送到死信队列
			if err := h.consumer.sendToDLQ(msg); err != nil {
				h.consumer.l.Error("发送到死信队列失败", zap.Error(err))
			}
			continue // 确保继续处理下一条消息
		}
		sess.MarkMessage(msg, "") // 只有成功处理后才标记
	}
	return nil
}

// sendToDLQ 发送消息到死信队列
func (c *CheckEventConsumer) sendToDLQ(msg *sarama.ConsumerMessage) error {
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

	_, _, err := c.dlqProd.SendMessage(dlqMsg)
	return err
}

// processMessage 处理从 Kafka 消费的消息
func (c *CheckEventConsumer) processMessage(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error {
	// 参数校验
	if err := c.validateMessage(msg); err != nil {
		return err
	}

	var evt CheckEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		c.l.Error("反序列化消息失败",
			zap.Error(err),
			zap.String("message", string(msg.Value)))
		return fmt.Errorf("反序列化消息失败: %w", err)
	}

	// 参数校验
	if err := c.validateEvent(&evt); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 处理消息
	if err := c.handleEvent(ctx, &evt); err != nil {
		return err
	}

	c.l.Info("消息处理成功",
		zap.Uint("post_id", evt.PostId),
		zap.Int64("uid", evt.Uid),
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset))

	return nil
}

// validateMessage 验证消息是否有效
func (c *CheckEventConsumer) validateMessage(msg *sarama.ConsumerMessage) error {
	if msg == nil || msg.Value == nil {
		c.l.Error("消息为空")
		return errors.New("消息为空")
	}

	c.l.Debug("开始处理消息",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset))

	return nil
}

// validateEvent 验证事件参数是否有效
func (c *CheckEventConsumer) validateEvent(evt *CheckEvent) error {
	if evt == nil {
		return errors.New("事件为空")
	}

	if evt.PostId == 0 || evt.Uid == 0 || evt.Title == "" || evt.Content == "" {
		c.l.Error("消息参数无效",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid))
		return fmt.Errorf("无效的消息参数: post_id=%d, uid=%d", evt.PostId, evt.Uid)
	}

	return nil
}

// handleEvent 处理审核事件
func (c *CheckEventConsumer) handleEvent(ctx context.Context, evt *CheckEvent) error {
	if ctx == nil {
		return errors.New("context为空")
	}

	check := domain.Check{
		PostID:  evt.PostId,
		Uid:     evt.Uid,
		Title:   evt.Title,
		Content: evt.Content,
		PlateID: evt.PlateID,
	}

	// 创建审核记录
	code, err := c.checkRepo.Create(ctx, check)
	if err != nil {
		c.l.Error("创建审核记录失败",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid),
			zap.Error(err))
		return fmt.Errorf("创建审核记录失败: %w", err)
	}

	// 记录已存在,标记消息处理完成
	if code == -1 {
		c.l.Info("审核记录已存在",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid))
		return nil
	}

	return nil
}
