package prometheus

import (
	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type KafkaMetricsHook struct {
	operationMetrics *prometheus.SummaryVec
}

// NewKafkaMetricsHook 初始化 KafkaMetricsHook 实例，并注册 Prometheus 指标
func NewKafkaMetricsHook(opts prometheus.SummaryOpts) *KafkaMetricsHook {
	operationMetrics := prometheus.NewSummaryVec(opts, []string{"operation", "topic"})
	prometheus.MustRegister(operationMetrics)
	return &KafkaMetricsHook{
		operationMetrics: operationMetrics,
	}
}

// WrapProducer 为 Sarama 同步生产者包装一个带有指标记录的生产者
func (h *KafkaMetricsHook) WrapProducer(producer sarama.SyncProducer) sarama.SyncProducer {
	return &instrumentedProducer{
		SyncProducer:     producer,
		operationMetrics: h.operationMetrics,
	}
}

// instrumentedProducer 是一个包装了 Sarama SyncProducer 的结构体
type instrumentedProducer struct {
	sarama.SyncProducer
	operationMetrics *prometheus.SummaryVec
}

// SendMessage 包装后的 SendMessage 方法，用于监控消息发送时间
func (p *instrumentedProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	startTime := time.Now()
	partition, offset, err = p.SyncProducer.SendMessage(msg)
	duration := time.Since(startTime).Seconds()
	p.operationMetrics.WithLabelValues("send", msg.Topic).Observe(duration)
	return partition, offset, err
}
