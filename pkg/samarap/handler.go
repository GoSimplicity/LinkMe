package samarap

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

// BatchHandler BatchHandler[T] 是一个泛型结构体，用于处理 Kafka 消息的批次
type BatchHandler[T any] struct {
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error // 处理消息的函数
	l  *zap.Logger                                        // 日志记录器
}

// NewBatchHandler NewBatchHandler[T] 是一个工厂函数，用于创建 BatchHandler[T] 的实例
func NewBatchHandler[T any](l *zap.Logger, fn func(msgs []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{fn: fn, l: l}
}

// Setup 方法是 Sarama ConsumerGroupHandler 接口的实现，它在消费者组开始时被调用
func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil // 当前实现不需要执行任何设置操作
}

// Cleanup 方法是 Sarama ConsumerGroupHandler 接口的实现，它在消费者组结束时被调用
func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil // 当前实现不需要执行任何清理操作
}

// ConsumeClaim 方法是 Sarama ConsumerGroupHandler 接口的实现，它在处理消息时被调用
func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages() // 获取分配给当前消费者的消息通道
	const batchSize = 100    // 定义批处理的大小
	const timeout = time.Second * 10
	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		var done bool
		for i := 0; i < batchSize && !done; i++ { // 循环直到批次大小达到或者上下文超时
			select {
			case <-ctx.Done(): // 是否超时
				done = true // 标记为完成
			case msg, ok := <-msgs: // 从消息通道中接收消息
				if !ok {
					cancel()
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列消息体失败", zap.Error(err))
					continue // 跳过当前消息，继续处理下一个
				}
				batch = append(batch, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		if len(batch) == 0 { // 如果批次为空，跳过处理
			continue
		}
		err := b.fn(batch, ts) // 调用处理函数处理批次中的消息
		if err != nil {        // 如果处理失败
			b.l.Error("处理消息失败", zap.Error(err))
			continue
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "") // 标记消息为已处理
		}
	}
}
