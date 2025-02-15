package es

import (
	"context"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func initKafka() sarama.Client {
	scfg := sarama.NewConfig()
	scfg.Consumer.Offsets.Initial = sarama.OffsetOldest //从头消费
	client, _ := sarama.NewClient([]string{"192.168.84.130:9092"}, scfg)
	return client
}

func initLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	logger, _ := cfg.Build()
	return logger
}

func initEsClient() *elasticsearch.TypedClient {
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "elastic",
		Password:  "gNWFCTYAobqKRl09LiPj",
	})
	if err != nil {
		panic(err)
	}
	return client
}

// 通过修改或添加数据库中的数据,判断是否数据同步到es中
func TestEsConsumer(t *testing.T) {
	logger := initLogger()
	es := initEsClient()
	searchDao := dao.NewSearchDAO(es, logger)

	esConsumer := NewEsConsumer(initKafka(), logger, repository.NewSearchRepository(searchDao))
	go func() {
		err := esConsumer.Start(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(2 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	users, err := searchDao.SearchUsers(ctx, []string{"bob"})
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		fmt.Println("已成功找到用户:")
		fmt.Println(user)
	}

	posts, err := searchDao.SearchPosts(ctx, []string{"HisLife"})
	if err != nil {
		panic(err)
	}
	fmt.Println("已成功找到文章:")
	fmt.Println(posts)

	// 优雅退出
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigchan:
		fmt.Printf("consumer test terminated")
	}

}
