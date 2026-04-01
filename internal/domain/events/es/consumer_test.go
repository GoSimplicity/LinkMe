//go:build integration
// +build integration

package es

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

func initKafka(t *testing.T) sarama.Client {
	t.Helper()
	scfg := sarama.NewConfig()
	scfg.Consumer.Offsets.Initial = sarama.OffsetOldest //从头消费
	addr := os.Getenv("LINKME_KAFKA_ADDR")
	if addr == "" {
		addr = "localhost:9094"
	}
	client, err := sarama.NewClient([]string{addr}, scfg)
	if err != nil {
		t.Fatalf("初始化 Kafka 客户端失败: %v", err)
	}
	return client
}

func initLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	logger, _ := cfg.Build()
	return logger
}

func initEsClient(t *testing.T) *elasticsearch.TypedClient {
	t.Helper()
	addr := os.Getenv("LINKME_ES_ADDR")
	if addr == "" {
		addr = "http://localhost:19200"
	}
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{addr},
	})
	if err != nil {
		t.Fatalf("初始化 ES 客户端失败: %v", err)
	}
	return client
}

func TestEsConsumer(t *testing.T) {
	logger := initLogger()
	es := initEsClient(t)
	searchDao := dao.NewSearchDAO(es, logger)

	esConsumer := NewEsConsumer(initKafka(t), logger, repository.NewSearchRepository(searchDao))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := esConsumer.Start(ctx); err != nil {
		t.Fatalf("启动 ES consumer 失败: %v", err)
	}
	time.Sleep(500 * time.Millisecond)
}
