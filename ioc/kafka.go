package ioc

import (
	"github.com/GoSimplicity/LinkMe/internal/domain/events"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/email"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/post"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/sms"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/sync"
	"github.com/IBM/sarama"
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
	// 创建Sarama配置对象，并设置生产者需要的配置项
	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true // 配置生产者需要返回确认成功的消息
	// 使用配置好的scfg创建Sarama客户端连接到Kafka
	client, err := sarama.NewClient(cfg.Addr, scfg)
	if err != nil {
		panic(err)
	}
	return client
}

// InitSyncProducer 使用已有的Sarama客户端初始化同步生产者
func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
	// 根据现有的客户端实例创建同步生产者
	// 确保每个消息在被确认已经成功写入至少一个副本之前，生产者会阻塞等待响应
	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return p
}

// InitConsumers 初始化并返回一个事件消费者，当前仅包含InteractiveReadEventConsumer
func InitConsumers(postConsumer *post.InteractiveReadEventConsumer, smsConsumer *sms.SMSConsumer, emailConsumer *email.EmailConsumer, syncConsumer *sync.SyncConsumer) []events.Consumer {
	// 返回消费者切片
	return []events.Consumer{postConsumer, smsConsumer, emailConsumer, syncConsumer}
}
