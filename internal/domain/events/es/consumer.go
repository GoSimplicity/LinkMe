package es

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

// EsConsumer 结构体用于消费Kafka消息并将数据同步到Elasticsearch
type EsConsumer struct {
	client sarama.Client
	rs     repository.SearchRepository
	tc     *elasticsearch.TypedClient
	l      *zap.Logger
}

// Event 结构体表示Kafka消息的事件数据
type Event struct {
	Type     string                   `json:"type"`
	Database string                   `json:"database"`
	Table    string                   `json:"table"`
	Data     []map[string]interface{} `json:"data"`
}

// Post 结构体表示文章的数据结构
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

// User 结构体表示用户的数据结构
type User struct {
	ID        int64   `mapstructure:"id"`
	Phone     *string `mapstructure:"phone"`
	Email     string  `mapstructure:"email"`
	Password  string  `mapstructure:"password"`
	CreatedAt int64   `mapstructure:"created_at"`
	UpdatedAt int64   `mapstructure:"updated_at"`
	Deleted   bool    `mapstructure:"deleted"`
}

// consumerGroupHandler 结构体实现了Kafka Consumer Group的接口
type consumerGroupHandler struct {
	r *EsConsumer
}

// NewEsConsumer 创建并返回一个新的EsConsumer实例
func NewEsConsumer(client sarama.Client, l *zap.Logger, rs repository.SearchRepository, tc *elasticsearch.TypedClient) *EsConsumer {
	return &EsConsumer{
		client: client,
		rs:     rs,
		tc:     tc,
		l:      l,
	}
}

// Start 启动Kafka消费者，监听消息并进行处理
func (r *EsConsumer) Start(_ context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("es_consumer_group", r.client)
	if err != nil {
		return err
	}

	r.l.Info("EsConsumer 开始消费")

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

// ConsumeClaim 消费Kafka的消息
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.r.Consume(sess, msg)
	}
	return nil
}

// Consume 处理Kafka消息，根据不同的表名执行不同的操作
func (r *EsConsumer) Consume(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var e Event
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		r.l.Error("消息反序列化失败", zap.Error(err))
		return
	}

	switch e.Table {
	case "posts":
		var posts []Post
		if err := decodeEventDataToPosts(e.Data, &posts); err != nil {
			r.l.Error("数据映射到结构体失败", zap.Error(err))
			return
		}
		for _, post := range posts {
			if err := r.handleEsPost(sess.Context(), post); err != nil {
				r.l.Error("处理ES失败", zap.Uint("id", post.ID), zap.Error(err))
				return
			}
		}
	case "users":
		var users []User
		if err := decodeEventDataToUsers(e.Data, &users); err != nil {
			r.l.Error("数据映射到结构体失败", zap.Error(err))
			return
		}
		for _, user := range users {
			if err := r.handleEsUser(sess.Context(), user); err != nil {
				r.l.Error("处理ES失败", zap.Int64("id", user.ID), zap.Error(err))
				return
			}
		}
	}

	sess.MarkMessage(msg, "")
}

// handleEsPost 处理文章的ES操作，发布或删除索引
func (r *EsConsumer) handleEsPost(ctx context.Context, post Post) error {
	if post.Status == domain.Published {
		return r.pushOrUpdatePostIndex(ctx, post)
	}
	return r.deletePostIndex(ctx, post)
}

// handleEsUser 处理用户的ES操作，创建或删除用户索引
func (r *EsConsumer) handleEsUser(ctx context.Context, user User) error {
	if user.Deleted {
		return r.deleteUserIndex(ctx, user)
	}
	return r.pushOrUpdateUserIndex(ctx, user)
}

// pushOrUpdatePostIndex 创建或更新文章索引
func (r *EsConsumer) pushOrUpdatePostIndex(ctx context.Context, post Post) error {
	exists, err := r.isPostIndexExists(ctx, post.ID)
	if err != nil {
		return err
	}
	if exists {
		r.l.Debug("Post 已存在，跳过处理", zap.Uint("id", post.ID))
		return nil
	}

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

// pushOrUpdateUserIndex 创建或更新用户索引
func (r *EsConsumer) pushOrUpdateUserIndex(ctx context.Context, user User) error {
	exists, err := r.isUserIndexExists(ctx, user.ID)
	if err != nil {
		return err
	}
	if exists {
		r.l.Debug("User 已存在，跳过处理", zap.Int64("id", user.ID))
		return nil
	}

	err = r.rs.InputUser(ctx, domain.UserSearch{
		Id:    user.ID,
		Email: user.Email,
		// 可以继续填充其他字段
	})
	if err != nil {
		r.l.Error("创建索引失败", zap.Int64("id", user.ID), zap.Error(err))
		return err
	}

	r.l.Info("User 索引创建成功", zap.Int64("id", user.ID))
	return nil
}

// deletePostIndex 删除文章索引
func (r *EsConsumer) deletePostIndex(ctx context.Context, post Post) error {
	exists, err := r.isPostIndexExists(ctx, post.ID)
	if err != nil {
		return err
	}
	if !exists {
		r.l.Debug("Post 不存在于索引中，跳过删除", zap.Uint("id", post.ID))
		return nil
	}

	if err := r.rs.DeletePostIndex(ctx, post.ID); err != nil {
		r.l.Error("删除索引失败", zap.Uint("id", post.ID), zap.Error(err))
		return err
	}

	r.l.Info("Post 索引删除成功", zap.Uint("id", post.ID))
	return nil
}

// deleteUserIndex 删除用户索引
func (r *EsConsumer) deleteUserIndex(ctx context.Context, user User) error {
	exists, err := r.isUserIndexExists(ctx, user.ID)
	if err != nil {
		return err
	}
	if !exists {
		r.l.Debug("User 不存在于索引中，跳过删除", zap.Int64("id", user.ID))
		return nil
	}

	if err := r.rs.DeleteUserIndex(ctx, user.ID); err != nil {
		r.l.Error("删除索引失败", zap.Int64("id", user.ID), zap.Error(err))
		return err
	}

	r.l.Info("User 索引删除成功", zap.Int64("id", user.ID))
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

// isUserIndexExists 查询User索引是否存在
func (r *EsConsumer) isUserIndexExists(ctx context.Context, userID int64) (bool, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"id": userID, // 使用 term 查询精确匹配 userID
			},
		},
	}

	// 将查询转换为 JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		r.l.Error("查询构建失败", zap.Int64("id", userID), zap.Error(err))
		return false, err
	}

	// 构建 Elasticsearch 搜索请求
	req := esapi.SearchRequest{
		Index: []string{"user_index"}, // 替换为你的实际索引名
		Body:  bytes.NewReader(queryJSON),
	}

	// 执行查询请求
	res, err := req.Do(ctx, r.tc)
	if err != nil {
		r.l.Error("查询索引失败", zap.Int64("id", userID), zap.Error(err))
		return false, err
	}
	defer res.Body.Close()

	// 检查响应是否包含错误
	if res.IsError() {
		r.l.Error("Elasticsearch 查询返回错误", zap.String("status", res.Status()), zap.Int64("id", userID))
		return false, fmt.Errorf("elasticsearch returned an error: %s", res.Status())
	}

	// 解析查询结果并返回是否存在该 User
	var searchResult struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		r.l.Error("解析查询结果失败", zap.Int64("id", userID), zap.Error(err))
		return false, err
	}

	// 如果命中数大于0，则表示User存在
	exists := searchResult.Hits.Total.Value > 0

	r.l.Debug("User 索引查询结果", zap.Int64("id", userID), zap.Bool("exists", exists))

	return exists, nil
}

// decodeEventDataToPosts 解析事件数据为 Post 结构体
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

// decodeEventDataToUsers 解析事件数据为 User 结构体
func decodeEventDataToUsers(data interface{}, users *[]User) error {
	config := &mapstructure.DecoderConfig{
		Result:           users,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToNullTimeHookFunc("2006-01-02 15:04:05.999"),
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(data)
}

// stringToTimeHookFunc 转换字符串到时间类型
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

// stringToNullTimeHookFunc 转换字符串到 NullTime 类型
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
