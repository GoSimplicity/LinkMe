package ioc

import (
	"github.com/GoSimplicity/LinkMe/internal/domain/events"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/check"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/email"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/es"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/post"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/publish"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/sms"
	"github.com/GoSimplicity/LinkMe/pkg/samarap/prometheus"
	"github.com/IBM/sarama"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

// InitSaramaClient 初始化Sarama客户端，用于连接到Kafka集群
func InitSaramaClient() sarama.Client {
	type Config struct {
		Addr []string `yaml:"addr"`
	}

	var cfg Config

	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}

	scfg := sarama.NewConfig()
	// 配置生产者需要返回确认成功的消息
	scfg.Producer.Return.Successes = true

	client, err := sarama.NewClient(cfg.Addr, scfg)
	if err != nil {
		panic(err)
	}

	return client
}

// InitSyncProducer 使用已有的Sarama客户端初始化同步生产者
func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
	// 根据现有的客户端实例创建同步生产者
	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}

	// 创建并注册自定义的 KafkaMetricsHook 插件
	kafkaMetricsHook := prometheus.NewKafkaMetricsHook(prometheus2.SummaryOpts{
		Namespace: "linkme",
		Subsystem: "kafka",
		Name:      "operation_duration_seconds",
		Help:      "Duration of Kafka operations in seconds",
		Objectives: map[float64]float64{
			0.5:  0.01,
			0.75: 0.01,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	// 包装生产者
	return kafkaMetricsHook.WrapProducer(p)
}

// InitConsumers 初始化并返回一个事件消费者
func InitConsumers(
	postConsumer *post.EventConsumer,
	smsConsumer *sms.SMSConsumer,
	emailConsumer *email.EmailConsumer,
	publishConsumer *publish.PublishPostEventConsumer,
	esConsumer *es.EsConsumer,
	checkConsumer *check.CheckEventConsumer,
	postDLQConsumer *post.PostDeadLetterConsumer,
	publishDLQConsumer *publish.PublishDeadLetterConsumer,
	checkDLQConsumer *check.CheckDeadLetterConsumer,
) []events.Consumer {
	// 返回消费者切片
	return []events.Consumer{
		postConsumer,
		smsConsumer,
		emailConsumer,
		publishConsumer,
		esConsumer,
		checkConsumer,
		postDLQConsumer,
		publishDLQConsumer,
		checkDLQConsumer,
	}
}
