package ioc

import (
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
