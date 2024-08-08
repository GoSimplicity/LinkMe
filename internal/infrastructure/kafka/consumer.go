package sync

import (
	"context"
	"encoding/json"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/IBM/sarama"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

const (
	brokerList = "localhost:9092"
	groupID    = "linkme-group"
	topic      = "schemahistory.linkme"
)

type SyncConsumer struct {
	Ready      chan bool         // 用于同步消费者的就绪状态
	collection *mongo.Collection // MongoDB集合
	logger     *zap.Logger
}

// NewSyncConsumer 创建并初始化一个新的Kafka消费者
func NewSyncConsumer(mongoClient *mongo.Client, logger *zap.Logger) *SyncConsumer {
	collection := mongoClient.Database("linkme").Collection("posts")
	return &SyncConsumer{
		Ready:      make(chan bool),
		collection: collection,
		logger:     logger,
	}
}

// Setup 在新的会话开始时运行，设置消费者的初始状态
func (consumer *SyncConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

// Cleanup 在会话结束时运行，清理资源
func (consumer *SyncConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 消费者组中的每个消费者循环处理消息
func (consumer *SyncConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var post domain.Post
		err := json.Unmarshal(message.Value, &post)
		if err != nil {
			consumer.logger.Error("反序列化消息时出错", zap.Error(err))
			continue
		}

		// 仅处理发布状态的帖子
		if post.Status == domain.Published {
			// 根据消息的键处理不同类型的操作
			switch string(message.Key) {
			case "c": // 创建
				if _, err := consumer.collection.InsertOne(context.TODO(), post); err != nil {
					consumer.logger.Error("插入文档时出错", zap.Error(err))
				}
			case "u": // 更新
				if _, err := consumer.collection.UpdateOne(
					context.TODO(),
					bson.M{"_id": post.ID},
					bson.D{{"$set", post}},
				); err != nil {
					consumer.logger.Error("更新文档时出错", zap.Error(err))
				}
			case "d": // 删除
				if _, err := consumer.collection.DeleteOne(context.TODO(), bson.M{"_id": post.ID}); err != nil {
					consumer.logger.Error("删除文档时出错", zap.Error(err))
				}
			default:
				consumer.logger.Warn("未知操作类型", zap.String("key", string(message.Key)))
			}
		}

		consumer.logger.Info("处理消息",
			zap.Int64("offset", message.Offset),
			zap.String("key", string(message.Key)),
			zap.String("value", string(message.Value)),
		)
		session.MarkMessage(message, "")
	}
	return nil
}

// Start 启动Kafka消费者组并开始处理消息
func (consumer *SyncConsumer) Start(_ context.Context) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup([]string{brokerList}, groupID, config)
	if err != nil {
		consumer.logger.Fatal("创建消费者组时出错", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, consumer); err != nil {
				consumer.logger.Error("消费者消费消息时出错", zap.Error(err))
			}
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()
	<-consumer.Ready
	consumer.logger.Info("Sarama消费者已启动并运行")
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)
	<-sigterm

	consumer.logger.Info("接收到终止信号，正在关闭消费者组")
	cancel()
	if err := consumerGroup.Close(); err != nil {
		consumer.logger.Error("关闭消费者组时出错", zap.Error(err))
	}
	return err
}
