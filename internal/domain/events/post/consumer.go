package post

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const (
	TopicDeadLetter = "post_events_dlq" // 死信队列主题
	MaxRetries      = 3                 // 最大重试次数
)

type EventConsumer struct {
	repo    repository.InteractiveRepository
	hisRepo repository.HistoryRepository
	client  sarama.Client
	l       *zap.Logger
	dlqProd sarama.SyncProducer // 死信队列生产者
}

func NewEventConsumer(
	repo repository.InteractiveRepository,
	hisRepo repository.HistoryRepository,
	client sarama.Client,
	dlqProd sarama.SyncProducer,
	l *zap.Logger,
) *EventConsumer {
	return &EventConsumer{
		repo:    repo,
		hisRepo: hisRepo,
		client:  client,
		l:       l,
		dlqProd: dlqProd,
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
				if err := cg.Consume(ctx, []string{TopicReadEvent}, &consumerGroupHandler{consumer: i}); err != nil {
					i.l.Error("阅读消息消费循环出错", zap.Error(err))
					continue
				}
			}
		}
	}()

	return nil
}

type consumerGroupHandler struct {
	consumer *EventConsumer
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
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
func (i *EventConsumer) sendToDLQ(msg *sarama.ConsumerMessage) error {
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

	_, _, err := i.dlqProd.SendMessage(dlqMsg)
	return err
}

// processMessage 处理从 Kafka 消费的消息
func (i *EventConsumer) processMessage(msg *sarama.ConsumerMessage) error {
	// 参数校验
	if err := i.validateMessage(msg); err != nil {
		return err
	}

	var evt ReadEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		i.l.Error("反序列化消息失败",
			zap.Error(err),
			zap.String("message", string(msg.Value)))
		return fmt.Errorf("反序列化消息失败: %w", err)
	}

	// 参数校验
	if err := i.validateEvent(&evt); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 处理消息
	if err := i.handleEvent(ctx, &evt); err != nil {
		return err
	}

	i.l.Info("消息处理成功",
		zap.Uint("post_id", evt.PostId),
		zap.Int64("uid", evt.Uid),
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset))

	return nil
}

// validateMessage 验证消息是否有效
func (i *EventConsumer) validateMessage(msg *sarama.ConsumerMessage) error {
	if msg == nil || msg.Value == nil {
		i.l.Error("消息为空")
		return errors.New("消息为空")
	}

	i.l.Debug("开始处理消息",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset))

	return nil
}

// validateEvent 验证事件参数是否有效
func (i *EventConsumer) validateEvent(evt *ReadEvent) error {
	if evt == nil {
		return errors.New("事件为空")
	}

	if evt.PostId == 0 || evt.Uid == 0 {
		i.l.Error("消息参数无效",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid))
		return fmt.Errorf("无效的消息参数: post_id=%d, uid=%d", evt.PostId, evt.Uid)
	}

	return nil
}

// handleEvent 处理阅读事件
func (i *EventConsumer) handleEvent(ctx context.Context, evt *ReadEvent) error {
	if ctx == nil {
		return errors.New("context为空")
	}

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

	return nil
}
