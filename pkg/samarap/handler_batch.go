package samarap

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type BatchHandler[T any] struct {
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error
	l  *zap.Logger
}

func NewBatchHandler[T any](l *zap.Logger, fn func(msgs []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{fn: fn, l: l}
}

func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

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
