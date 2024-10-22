package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/pkg/cachep/local"
	"github.com/IBM/sarama"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"reflect"
	"time"
)

type CacheConsumer struct {
	client   sarama.Client
	l        *zap.Logger
	redis    redis.Cmdable
	hisCache cache.HistoryCache
	local    *local.CacheManager
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
	r *CacheConsumer
}

func NewCacheConsumer(client sarama.Client, l *zap.Logger, redis redis.Cmdable, local *local.CacheManager, hisCache cache.HistoryCache) *CacheConsumer {
	// 创建MongoDB客户端
	return &CacheConsumer{
		client:   client,
		hisCache: hisCache,
		l:        l,
		redis:    redis,
		local:    local,
	}
}

func (r *CacheConsumer) Start(_ context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("cache_consumer_group", r.client)
	r.l.Info("CacheConsumer 开始消费")
	if err != nil {
		return err
	}

	go func() {
		for {
			if err := cg.Consume(context.Background(), []string{"linkme_binlog"}, &consumerGroupHandler{r: r}); err != nil {
				r.l.Error("退出了消费循环异常", zap.Error(err))
				time.Sleep(time.Second * 5)
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

func (r *CacheConsumer) Consume(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var e Event
	var posts []Post

	// 反序列化消息
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		r.l.Error("消息反序列化失败", zap.Error(err))
		return
	}

	if e.Table != "posts" {
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
func (r *CacheConsumer) handlePost(ctx context.Context, post Post) error {
	pipe := r.redis.Pipeline() // 开启Redis管道 批量执行多个任务 提高性能

	if post.Status == domain.Published {
		// 删除公共详细缓存
		pipe.Del(ctx, fmt.Sprintf("post:pub:detail:%d", post.ID))
		// 使用 DeleteKeysWithPattern 删除匹配模式的公共列表缓存
		if err := r.DeleteKeysWithPattern(ctx, "post:pub:list:*"); err != nil {
			r.l.Error("Failed to delete public list keys", zap.Error(err))
			return err
		}
	} else {
		// 删除私有详细缓存
		pipe.Del(ctx, fmt.Sprintf("post:detail:%d:%d", post.ID, post.AuthorID))
		// 使用 DeleteKeysWithPattern 删除匹配模式的私有列表缓存
		if err := r.DeleteKeysWithPattern(ctx, fmt.Sprintf("post:pri:list:%d:*", post.AuthorID)); err != nil {
			r.l.Error("Failed to delete private list keys", zap.Error(err))
			return err
		}

		// 如果监测到 Post 不属于发布状态，删除相关历史记录缓存
		if err := r.hisCache.DeleteOneCache(ctx, post.ID, post.AuthorID); err != nil {
			r.l.Error("Failed to delete history cache", zap.Error(err))
			return err
		}
	}

	_, err := pipe.Exec(ctx) // 执行所有Redis删除操作
	if err != nil {
		return err
	}

	return nil
}

func (r *CacheConsumer) DeleteKeysWithPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string
	var err error

	for {
		// 执行 SCAN 命令
		// SCAN 命令的匹配规则是基于通配符
		// 这里会返回最多100个key 并返回一个游标cursor继续扫描
		keys, cursor, err = r.redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			r.l.Error("Failed to scan Redis keys", zap.Error(err))
			return err
		}

		if len(keys) > 0 {
			// 删除匹配的键
			if err := r.local.Delete(ctx, keys...); err != nil {
				r.l.Error("Failed to delete Redis keys", zap.Error(err))
				return err
			}
			r.l.Info("Deleted keys", zap.Strings("keys", keys))
		}

		// 如果 cursor 为 0，说明遍历结束
		if cursor == 0 {
			break
		}
	}

	return nil
}

// 自定义解析数据配置
func decodeEventDataToPosts(data interface{}, posts *[]Post) error {
	config := &mapstructure.DecoderConfig{
		Result:           posts,          // 结果保存在 posts 中
		TagName:          "mapstructure", // 使用 mapstructure 标签
		WeaklyTypedInput: true,           // 允许使用未导出字段
		DecodeHook: mapstructure.ComposeDecodeHookFunc( // 自定义解码钩子
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
