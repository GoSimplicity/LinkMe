package es

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// EsConsumer 结构体用于消费Kafka消息并将数据同步到Elasticsearch
type EsConsumer struct {
	client sarama.Client
	rs     repository.SearchRepository
	l      *zap.Logger

	flushUser *FlushUser
	flushPost *FlushPost
}

// Event 结构体表示Kafka消息的事件数据
type ChangeEvent struct {
	Payload struct { //数据的主体部分
		Before map[string]interface{} `json:"before"`
		After  map[string]interface{} `json:"after"`
		Source struct {
			Connector string      `json:"connector"` //创建链接器的对象即被备份的数据库
			Name      string      `json:"name"`      //自定义topic前缀
			Snapshot  string      `json:"snapshot"`  //是否是快照
			Db        string      `json:"db"`        //数据库名
			Table     string      `json:"table"`     //表名
			Gtid      interface{} `json:"gtid"`
			File      string      `json:"file"` //从哪个binlog日志中取得数据的
			Pos       int         `json:"pos"`  //偏移量
		} `json:"source"`
		Op string `json:"op"` //指令类型(c:INSERT u:UPDATE d:DELETE
	} `json:"payload"`
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
	uid          int64        `mapstructure:"uid"`
	IsSubmit     int8         `mapstructure:"is_submit"`
}

// User 结构体表示用户的数据结构
type User struct {
	Id           int64     `mapstructure:"id"`
	CreatedAt    int64     `mapstructure:"created_at"`
	UpdatedAt    int64     `mapstructure:"updated_at"`
	DeletedAt    int64     `mapstructure:"deleted_at"`
	NickName     string    `mapstructure:"nick_name"`
	PasswordHash string    `mapstructure:"password_hash"`
	Birthday     time.Time `mapstructure:"birthday"`
	Deleted      int8      `mapstructure:"deleted"`
	Email        string    `mapstructure:"email"`
	Phone        string    `mapstructure:"phone"`
	About        string    `mapstructure:"about"`
}

// consumerGroupHandler 结构体实现了Kafka Consumer Group的接口
type consumerGroupHandler struct {
	r *EsConsumer
}

// NewEsConsumer 创建并返回一个新的EsConsumer实例
func NewEsConsumer(client sarama.Client, l *zap.Logger, rs repository.SearchRepository) *EsConsumer {
	esConsumer := &EsConsumer{
		client: client,
		rs:     rs,
		l:      l,
	}
	esConsumer.flushUser = NewFlushUser(esConsumer)
	esConsumer.flushPost = NewFlushPost(esConsumer)
	return esConsumer
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
			if err := cg.Consume(context.Background(), []string{"oracle.linkme.users", "oracle.linkme.posts"}, &consumerGroupHandler{r: r}); err != nil {
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
	var e ChangeEvent
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		r.l.Error("消息反序列化失败", zap.Error(err))
		return
	}

	r.l.Info("开始消费event")

	switch e.Payload.Source.Table {
	case "posts":
		var post Post
		var err error
		switch e.Payload.Op {
		case "r":
			if err := decodeEventDataToPost(e.Payload.After, &post); err != nil {
				r.l.Error("解析快照用户数据失败", zap.Error(err))
				return
			}
			// 加锁并写入缓冲区
			r.flushPost.bufferMutex.Lock()
			r.flushPost.postBuffer = append(r.flushPost.postBuffer, post)
			bufferLen := len(r.flushUser.userBuffer)
			isLast := e.Payload.Source.Snapshot == "last"
			r.flushUser.bufferMutex.Unlock()
			// 触发插入条件：达到批量大小或快照结束
			if bufferLen >= r.flushUser.bulkSize || isLast {
				r.flushUser.flushBuffer()
			}
			return
		case "c":
			err = decodeEventDataToPost(e.Payload.After, &post)
		case "u":
			err = decodeEventDataToPost(e.Payload.After, &post)
		case "d":
			err = decodeEventDataToPost(e.Payload.Before, &post)
		}
		if err != nil {
			r.l.Error("数据映射到结构体失败", zap.Error(err))
			return
		}

		if err := r.handleEsPost(sess.Context(), post); err != nil {
			r.l.Error("处理ES失败", zap.Uint("id", post.ID), zap.Error(err))
			return
		}

	case "users":
		var user User
		var err error
		switch e.Payload.Op {
		case "r":
			if err := decodeEventDataToUser(e.Payload.After, &user); err != nil {
				r.l.Error("解析快照用户数据失败", zap.Error(err))
				return
			}
			// 加锁并写入缓冲区
			r.flushUser.bufferMutex.Lock()
			r.flushUser.userBuffer = append(r.flushUser.userBuffer, user)
			bufferLen := len(r.flushUser.userBuffer)
			isLast := e.Payload.Source.Snapshot == "last"
			r.flushUser.bufferMutex.Unlock()
			// 触发插入条件：达到批量大小或快照结束
			if bufferLen >= r.flushUser.bulkSize || isLast {
				r.flushUser.flushBuffer()
			}
			return
		case "c":
			err = decodeEventDataToUser(e.Payload.After, &user)
		case "u":
			err = decodeEventDataToUser(e.Payload.After, &user)
		case "d":
			err = decodeEventDataToUser(e.Payload.Before, &user)
		}
		if err != nil {
			r.l.Error("数据映射到结构体失败", zap.Error(err))
			return
		}
		if err := r.handleEsUser(sess.Context(), user); err != nil {
			r.l.Error("处理ES失败", zap.Int64("id", user.Id), zap.Error(err))
			return
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
	if user.Deleted == 1 {
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
	exists, err := r.isUserExists(ctx, user.Id)
	if err != nil {
		return err
	}
	if exists {
		r.l.Debug("User 已存在，执行更新操作", zap.Int64("id", user.Id))
		return nil
	}

	err = r.rs.InputUser(ctx, domain.UserSearch{
		Id:       user.Id,
		Nickname: user.NickName,
		Birthday: user.Birthday,
		Email:    user.Email,
		Phone:    user.Phone,
		About:    user.About,
	})
	if err != nil {
		r.l.Error("创建索引失败", zap.Int64("id", user.Id), zap.Error(err))
		return err
	}

	r.l.Info("User 索引创建成功", zap.Int64("id", user.Id))
	return nil
}

func (r *EsConsumer) bulkInsertUser(ctx context.Context, users []User) error {
	var searchUsers []domain.UserSearch
	for _, user := range users {
		searchUsers = append(searchUsers, domain.UserSearch{
			Id:       user.Id,
			Nickname: user.NickName,
			Email:    user.Email,
			Phone:    user.Phone,
			Birthday: user.Birthday,
			About:    user.About,
		})
	}
	return r.rs.BulkInputUsers(ctx, searchUsers)
}

func (r *EsConsumer) bulkInsertPost(ctx context.Context, posts []Post) error {
	var searchPosts []domain.PostSearch
	for _, post := range posts {
		searchPosts = append(searchPosts, domain.PostSearch{
			Id:       post.ID,
			AuthorId: post.AuthorID,
			Title:    post.Title,
			Content:  post.Content,
			Status:   post.Status,
			Tags:     post.Tags,
		})
	}
	return r.rs.BulkInputPosts(ctx, searchPosts)
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

// deleteUserIndex 删除索引中的指定用户
func (r *EsConsumer) deleteUserIndex(ctx context.Context, user User) error {
	exists, err := r.isUserExists(ctx, user.Id)
	if err != nil {
		return err
	}
	if !exists {
		r.l.Debug("User 不存在于索引中，跳过删除", zap.Int64("id", user.Id))
		return nil
	}

	if err := r.rs.DeleteUserIndex(ctx, user.Id); err != nil {
		r.l.Error("删除索引失败", zap.Int64("id", user.Id), zap.Error(err))
		return err
	}

	r.l.Info("User 索引删除成功", zap.Int64("id", user.Id))
	return nil
}

// isPostIndexExists 查询Post索引是否存在
func (r *EsConsumer) isPostIndexExists(ctx context.Context, postID uint) (bool, error) {
	exist, err := r.rs.IsExistPost(ctx, postID)
	// 检查响应是否包含错误
	if err != nil {
		r.l.Error("Elasticsearch 查询返回错误", zap.Error(err))
		return false, fmt.Errorf("elasticsearch returned an error: %s", err)
	}

	r.l.Debug("Post 索引查询结果", zap.Uint("id", postID), zap.Bool("exists", exist))

	return exist, nil
}

// isUserExists 查询指定User是否存在
func (r *EsConsumer) isUserExists(ctx context.Context, userID int64) (bool, error) {
	exist, err := r.rs.IsExistUser(ctx, userID)
	if err != nil {
		r.l.Error("Elasticsearch 查询返回错误", zap.Error(err))
		return false, fmt.Errorf("elasticsearch returned an error: %s", err)
	}

	r.l.Debug("User 索引查询结果", zap.Int64("id", userID), zap.Bool("exists", exist))

	return exist, nil
}
