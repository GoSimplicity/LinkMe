package es

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"reflect"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

type EsConsumer struct {
	client sarama.Client
	rs     repository.SearchRepository
	tc     *elasticsearch.TypedClient
	l      *zap.Logger
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
	r *EsConsumer
}

func NewEsConsumer(client sarama.Client, l *zap.Logger, rs repository.SearchRepository, tc *elasticsearch.TypedClient) *EsConsumer {
	// 创建MongoDB客户端
	return &EsConsumer{
		client: client,
		rs:     rs,
		tc:     tc,
		l:      l,
	}
}

func (r *EsConsumer) Start(_ context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("es_consumer_group", r.client)

	r.l.Info("EsConsumer 开始消费")

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

func (r *EsConsumer) Consume(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
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
		if err := r.handleEs(sess.Context(), post); err != nil {
			r.l.Error("处理es失败", zap.Uint("id", post.ID), zap.Error(err))
			return
		}
	}

	// 标记消息为已处理
	sess.MarkMessage(msg, "")
}

func (r *EsConsumer) handleEs(ctx context.Context, post Post) error {
	switch post.Status {
	case domain.Published:
		return r.pushOrUpdatePostIndex(ctx, post)
	default:
		return r.deletePostIndex(ctx, post)
	}
}

func (r *EsConsumer) pushOrUpdatePostIndex(ctx context.Context, post Post) error {
	// 检查索引中是否已经存在该 Post
	exists, err := r.isPostIndexExists(ctx, post.ID)
	if err != nil {
		return err
	}

	// 如果索引中已存在该 Post，跳过处理
	if exists {
		r.l.Debug("Post 已存在，跳过处理", zap.Uint("id", post.ID))
		return nil
	}

	// 如果索引中不存在该 Post，创建索引
	err = r.rs.InputPost(ctx, domain.PostSearch{
		Id:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Status:  post.Status,
	})
	if err != nil {
		r.l.Error("创建索引失败", zap.Uint("id", post.ID), zap.Error(err))
		return err
	}

	r.l.Info("Post 索引创建成功", zap.Uint("id", post.ID))

	return nil
}

func (r *EsConsumer) deletePostIndex(ctx context.Context, post Post) error {
	// 检查索引中是否存在该 Post
	exists, err := r.isPostIndexExists(ctx, post.ID)
	if err != nil {
		return err
	}

	// 如果索引中不存在该 Post，跳过删除
	if !exists {
		r.l.Debug("Post 不存在于索引中，跳过删除", zap.Uint("id", post.ID))
		return nil
	}

	// 删除索引
	if err := r.rs.DeletePostIndex(ctx, post.ID); err != nil {
		r.l.Error("删除索引失败", zap.Uint("id", post.ID), zap.Error(err))
		return err
	}

	r.l.Info("Post 索引删除成功", zap.Uint("id", post.ID))

	return nil
}

// isPostIndexExists 查询Post索引是否存在
func (r *EsConsumer) isPostIndexExists(ctx context.Context, postID uint) (bool, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"id": postID, // 使用 term 查询精确匹配 postID
			},
		},
	}

	// 将查询转换为 JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		r.l.Error("查询构建失败", zap.Uint("id", postID), zap.Error(err))
		return false, err
	}

	// 构建 Elasticsearch 搜索请求
	req := esapi.SearchRequest{
		Index: []string{"post_index"}, // 替换为你的实际索引名
		Body:  bytes.NewReader(queryJSON),
	}

	// 执行查询请求
	res, err := req.Do(ctx, r.tc)
	if err != nil {
		r.l.Error("查询索引失败", zap.Uint("id", postID), zap.Error(err))
		return false, err
	}
	defer res.Body.Close()

	// 检查响应是否包含错误
	if res.IsError() {
		r.l.Error("Elasticsearch 查询返回错误", zap.String("status", res.Status()), zap.Uint("id", postID))
		return false, fmt.Errorf("elasticsearch returned an error: %s", res.Status())
	}

	// 解析查询结果并返回是否存在该 Post
	var searchResult struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		r.l.Error("解析查询结果失败", zap.Uint("id", postID), zap.Error(err))
		return false, err
	}

	// 如果命中数大于0，则表示Post存在
	exists := searchResult.Hits.Total.Value > 0

	r.l.Debug("Post 索引查询结果", zap.Uint("id", postID), zap.Bool("exists", exists))

	return exists, nil
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
