package sync

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/IBM/sarama"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type SyncConsumer struct {
	client      sarama.Client
	l           *zap.Logger
	db          *gorm.DB
	mongoClient *mongo.Client
	postDAO     dao.PostDAO
}

type Event struct {
	Type     string              `json:"type"`
	Database string              `json:"database"`
	Table    string              `json:"table"`
	Data     []map[string]string `json:"data"`
}

type Post struct {
	ID           uint         `mapstructure:"id"`
	Title        string       `mapstructure:"title"`
	Content      string       `mapstructure:"content"`
	CreatedAt    time.Time    `mapstructure:"created_at"`
	UpdatedAt    time.Time    `mapstructure:"updated_at"`
	DeletedAt    sql.NullTime `mapstructure:"deleted_at"`
	AuthorID     int64        `mapstructure:"author_id"`
	Status       uint8        `mapstructure:"status"`
	PlateID      int64        `mapstructure:"plate_id"`
	Slug         string       `mapstructure:"slug"`
	CategoryID   int64        `mapstructure:"category_id"`
	Tags         string       `mapstructure:"tags"`
	CommentCount int64        `mapstructure:"comment_count"`
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
	var e Event

	if err := json.Unmarshal(msg.Value, &e); err != nil {
		panic(err)
	}

	var p []Post

	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &p,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToTimeHookFunc("2006-01-02 15:04:05.999"),
			stringToNullTimeHookFunc("2006-01-02 15:04:05.999"),
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		r.l.Error("解码器创建失败", zap.Error(err))
		return
	}

	if err := decoder.Decode(e.Data); err != nil {
		r.l.Error("数据映射到结构体失败", zap.Error(err))
		return
	}

	sess.MarkMessage(msg, "")
}

func stringToTimeHookFunc(layout string) mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
		if f != reflect.String || t != reflect.Struct {
			return data, nil
		}

		str := data.(string)
		if str == "" {
			return time.Time{}, nil
		}

		return time.Parse(layout, str)
	}
}

func stringToNullTimeHookFunc(layout string) mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
		if f != reflect.String || t != reflect.Struct {
			return data, nil
		}

		str := data.(string)
		if str == "" {
			return sql.NullTime{Valid: false}, nil
		}

		parsedTime, err := time.Parse(layout, str)
		if err != nil {
			return nil, err
		}
		return sql.NullTime{Time: parsedTime, Valid: true}, nil
	}
}
