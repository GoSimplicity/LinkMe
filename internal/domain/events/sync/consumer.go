package sync

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/IBM/sarama"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	Type     string                   `json:"type"`
	Database string                   `json:"database"`
	Table    string                   `json:"table"`
	Data     []map[string]interface{} `json:"data"`
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
	var posts []Post

	// 反序列化消息
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		r.l.Error("消息反序列化失败", zap.Error(err))
		return
	}

	if e.Table != "posts" {
		r.l.Info("不是帖子表，跳过", zap.String("table", e.Table))
		return
	}

	// 数据映射到结构体
	if err := decodeEventDataToPosts(e.Data, &posts); err != nil {
		r.l.Error("数据映射到结构体失败", zap.Error(err))
		return
	}

	// 处理每个 Post
	for _, post := range posts {
		if err := r.handlePost(sess.Context(), post); err != nil {
			r.l.Error("处理帖子失败", zap.Uint("id", post.ID), zap.Error(err))
			return
		}
	}

	// 标记消息为已处理
	sess.MarkMessage(msg, "")
}

// handlePost 根据状态处理帖子
func (r *SyncConsumer) handlePost(ctx context.Context, post Post) error {
	switch post.Status {
	case domain.Published:
		return r.pushOrUpdateMongo(ctx, post)
	default:
		return r.deleteMongo(ctx, post)
	}
}

// 自定义解析数据配置
func decodeEventDataToPosts(data interface{}, posts *[]Post) error {
	config := &mapstructure.DecoderConfig{
		Result:           posts,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToTimeHookFunc("2006-01-02 15:04:05.999"),
			stringToNullTimeHookFunc("2006-01-02 15:04:05.999"),
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(data)
}

// 上传或更新数据
func (r *SyncConsumer) pushOrUpdateMongo(ctx context.Context, post Post) error {
	post.UpdatedAt = time.Now()

	collection := r.mongoClient.Database("linkme").Collection("posts")
	filter := bson.M{"id": post.ID}

	// 尝试查询并更新文档，如果文档不存在则插入
	update := bson.M{"$set": post}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.l.Error("更新或插入数据失败", zap.Error(err))
		return err
	}

	return nil
}

// 删除数据
func (r *SyncConsumer) deleteMongo(ctx context.Context, post Post) error {
	collection := r.mongoClient.Database("linkme").Collection("posts")
	filter := bson.M{"id": post.ID}

	// 尝试删除文档
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		r.l.Error("删除帖子时出现错误", zap.Error(err))
		return err
	}

	// 检查删除的文档数量，处理帖子不存在的情况
	if result.DeletedCount == 0 {
		r.l.Info("帖子不存在，可能已经被删除", zap.Uint("id", post.ID))
		// 认为这是正常情况，返回 nil 表示没有错误
		return nil
	}

	r.l.Info("帖子已成功删除", zap.Uint("id", post.ID))
	return nil
}

// 转换字符串到时间类型
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

// 转换字符串到 NullTime 类型
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
