package check

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/comment"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/publish"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	aiCheck "github.com/GoSimplicity/LinkMe/utils/AiCheck"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

const (
	TopicDeadLetter = "check_events_dlq" // 死信队列主题
	MaxRetries      = 3                  // 最大重试次数
)

type CheckEventConsumer struct {
	checkRepo       repository.CheckRepository
	client          sarama.Client
	l               *zap.Logger
	dlqProd         sarama.SyncProducer // 死信队列生产者
	postProducer    publish.Producer
	commentProducer comment.Producer
}

func NewCheckEventConsumer(
	checkRepo repository.CheckRepository,
	client sarama.Client,
	dlqProd sarama.SyncProducer,
	l *zap.Logger,
	postProducer publish.Producer,
	commentProducer comment.Producer,
) *CheckEventConsumer {
	return &CheckEventConsumer{
		checkRepo:       checkRepo,
		client:          client,
		l:               l,
		dlqProd:         dlqProd,
		postProducer:    postProducer,
		commentProducer: commentProducer,
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
	// h.consumer.l.Info("ConsumeClaim  function called")
	// 输出获取的消息数量，确认消息是否被正确获取
	h.consumer.l.Info("消息数", zap.Int("message_count", len(claim.Messages())))
	for msg := range claim.Messages() {
		h.consumer.l.Info("收到消息", zap.ByteString("message", msg.Value))
		// 处理每一条消息
		if err := h.consumer.processMessage(msg); err != nil {
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
func (c *CheckEventConsumer) processMessage(msg *sarama.ConsumerMessage) error {
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

	if (evt.PostId == 0 && evt.BizId == 1) || evt.Uid == 0 || (evt.Title == "" && evt.BizId == 1) || evt.Content == "" {
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
		BizId:   evt.BizId,
		PostID:  evt.PostId,
		Uid:     evt.Uid,
		Title:   evt.Title,
		Content: evt.Content,
		PlateID: evt.PlateID,
	}
	// 先使用AI进行审核，如果AI审核失败，再交由人工审核。

	checkResult, err2 := aiCheck.CheckPostContent(evt.Content)
	if err2 != nil || checkResult == false {
		// 交由人工审核
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
	}
	// AI审核通过，直接更新帖子状态|传给各个消费者我

	// 创建对应的审核记录，然后发送给各个消费者使用
	go func() {
		// 创建带超时的context
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		// 并发执行发布事件和记录活动
		done := make(chan error, 2)
		go func() {
			// 用于区分审核的业务类型[1：帖子 2：评论]
			if check.BizId == 1 {
				done <- c.postProducer.ProducePublishEvent(publish.PublishEvent{
					PostId: check.PostID,
					Uid:    check.Uid,
					Status: domain.Published,
					BizId:  check.BizId,
				})
			} else if check.BizId == 2 {
				done <- c.commentProducer.ProduceCommentEvent(comment.CommentEvent{
					PostId: check.PostID,
					Uid:    check.Uid,
					Status: domain.Published,
					BizId:  check.BizId,
				})
			}

		}()

		// 等待所有goroutine完成或超时
		for i := 0; i < 2; i++ {
			select {
			case err := <-done:
				if err != nil {
					c.l.Error("异步任务执行失败", zap.Error(err))
				}
			case <-ctx.Done():
				c.l.Error("异步任务执行超时")
				return
			}
		}
	}()

	return nil
}
