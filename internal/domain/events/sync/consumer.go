package sync

import (
	"context"
	"encoding/json"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/GoSimplicity/LinkMe/pkg/canalp"
	"github.com/IBM/sarama"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type SyncConsumer struct {
	client      sarama.Client
	l           *zap.Logger
	db          *gorm.DB
	mongoClient *mongo.Client
	postDAO     dao.PostDAO
}

type consumerGroupHandler struct {
	r *SyncConsumer
}

func NewSyncConsumer(client sarama.Client, l *zap.Logger, db *gorm.DB, mongoClient *mongo.Client, postDAO dao.PostDAO) *SyncConsumer {
	// 创建MongoDB客户端
	return &SyncConsumer{
		client:      client,
		l:           l,
		db:          db,
		mongoClient: mongoClient,
		postDAO:     postDAO,
	}
}

func (r *SyncConsumer) Start(_ context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("sync_consumer_group", r.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			if err := cg.Consume(context.Background(), []string{"linkme_binlog"}, &consumerGroupHandler{r: r}); err != nil {
				r.l.Error("退出了消费循环异常", zap.Error(err))
			}
		}
	}()
	return nil
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.r.Consume(sess, msg)
	}
	return nil
}

func (r *SyncConsumer) Consume(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var binlogMsg canalp.Message[map[string]interface{}]
	if err := json.Unmarshal(msg.Value, &binlogMsg); err != nil {
		r.l.Error("消息反序列化失败", zap.Error(err))
		return
	}
	var post dao.Post
	for _, data := range binlogMsg.Data {
		if err := mapToStruct(data, &post); err != nil {
			r.l.Error("数据映射到结构体失败", zap.Error(err))
			continue
		}
		sess.MarkMessage(msg, "")
	}
}
func mapToStruct(data map[string]interface{}, post *dao.Post) error {
	if idStr, ok := data["id"].(string); ok {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}
		post.ID = uint(id)
	}

	if authorIDStr, ok := data["author_id"].(string); ok {
		authorID, err := strconv.ParseInt(authorIDStr, 10, 64)
		if err != nil {
			return err
		}
		post.AuthorID = authorID
	}

	if statusStr, ok := data["status"].(string); ok {
		status, err := strconv.ParseUint(statusStr, 10, 8)
		if err != nil {
			return err
		}
		post.Status = uint8(status)
	}

	if plateIDStr, ok := data["plate_id"].(string); ok {
		plateID, err := strconv.ParseInt(plateIDStr, 10, 64)
		if err != nil {
			return err
		}
		post.PlateID = plateID
	}

	if categoryIDStr, ok := data["category_id"].(string); ok {
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err != nil {
			return err
		}
		post.CategoryID = categoryID
	}

	if commentCountStr, ok := data["comment_count"].(string); ok {
		commentCount, err := strconv.ParseInt(commentCountStr, 10, 64)
		if err != nil {
			return err
		}
		post.CommentCount = commentCount
	}

	// Convert other fields directly
	post.Title = data["title"].(string)
	post.Content = data["content"].(string)

	if createdAtStr, ok := data["created_at"].(string); ok {
		createdAt, err := time.Parse("2006-01-02 15:04:05.999", createdAtStr)
		if err != nil {
			return err
		}
		post.CreatedAt = createdAt
	}

	if updatedAtStr, ok := data["updated_at"].(string); ok {
		updatedAt, err := time.Parse("2006-01-02 15:04:05.999", updatedAtStr)
		if err != nil {
			return err
		}
		post.UpdatedAt = updatedAt
	}

	post.Slug = data["slug"].(string)
	post.Tags = data["tags"].(string)

	return nil
}
