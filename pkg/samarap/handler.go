package samarap

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// Handler 是一个通用的消息处理器
type Handler[T any] struct {
	logger *zap.Logger
	handle func(msg *sarama.ConsumerMessage, event T) error
}

// NewHandler 创建一个新的 Handler 实例
func NewHandler[T any](logger *zap.Logger, handle func(msg *sarama.ConsumerMessage, event T) error) *Handler[T] {
	return &Handler[T]{logger: logger, handle: handle}
}

// Setup 在消费组会话开始时调用
func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	h.logger.Info("Consumer group session setup")
	return nil
}

// Cleanup 在消费组会话结束时调用
func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	h.logger.Info("Consumer group session cleanup")
	return nil
}

// ConsumeClaim 处理消费组的消息
func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event T
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			h.logger.Error("Failed to unmarshal message", zap.Error(err), zap.ByteString("value", msg.Value))
			continue // 跳过无法反序列化的消息
		}
		if err := h.handle(msg, event); err != nil {
			h.logger.Error("Failed to process message", zap.Error(err), zap.ByteString("key", msg.Key), zap.ByteString("value", msg.Value))
			// 你可以在这里引入重试逻辑，根据具体需求
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
