package samarap

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

// BatchHandler BatchHandler[T] 是一个泛型结构体，用于处理 Kafka 消息的批次。
type BatchHandler[T any] struct {
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error // 处理消息的函数
	l  *zap.Logger                                        // 日志记录器
}

// NewBatchHandler NewBatchHandler[T] 是一个工厂函数，用于创建 BatchHandler[T] 的实例。
func NewBatchHandler[T any](l *zap.Logger, fn func(msgs []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{fn: fn, l: l}
}

// Setup 方法是 Sarama ConsumerGroupHandler 接口的实现，它在消费者组开始时被调用。
func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil // 当前实现不需要执行任何设置操作
}

// Cleanup 方法是 Sarama ConsumerGroupHandler 接口的实现，它在消费者组结束时被调用。
func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil // 当前实现不需要执行任何清理操作
}

// ConsumeClaim 方法是 Sarama ConsumerGroupHandler 接口的实现，它在处理消息时被调用。
func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages() // 获取分配给当前消费者的消息通道
	const batchSize = 10     // 定义批处理的大小
	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)                   // 创建一个消息批次
		ts := make([]T, 0, batchSize)                                            // 创建一个泛型类型 T 的切片，用于存储反序列化的消息
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10) // 创建一个带有超时的上下文
		var done = false
		for i := 0; i < batchSize && !done; i++ { // 循环直到批次大小达到或者上下文超时
			select {
			case <-ctx.Done(): // 如果上下文超时
				done = true // 标记为完成
			case msg, ok := <-msgs: // 从消息通道中接收消息
				if !ok { // 如果通道关闭
					cancel()   // 取消上下文
					return nil // 返回 nil，表示消费完成
				}
				var t T                              // 创建一个类型 T 的变量
				err := json.Unmarshal(msg.Value, &t) // 反序列化消息体
				if err != nil {                      // 如果反序列化失败
					b.l.Error("反序列消息体失败", zap.Error(err))
					continue // 跳过当前消息，继续处理下一个
				}
				batch = append(batch, msg) // 将消息添加到批次中
				ts = append(ts, t)         // 将反序列化的消息添加到类型 T 的切片中
			}
		}
		cancel() // 取消上下文
		// 凑够一批就处理
		err := b.fn(batch, ts) // 调用处理函数处理批次中的消息
		if err != nil {        // 如果处理失败
			b.l.Error("处理消息失败", zap.Error(err))
		}
		for _, msg := range batch { // 遍历批次中的消息
			session.MarkMessage(msg, "") // 标记消息为已处理
		}
	}
}
