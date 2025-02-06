package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	PostIndex    = "post_index"
	UserIndex    = "user_index"
	LogsIndex    = "logs_index"
	CommentIndex = "comment_index"
)

type SearchDAO interface {
	CreateIndex(ctx context.Context, indexName string, properties ...interface{}) error
	CreatePostIndex(ctx context.Context, properties ...interface{}) error
	CreateCommentIndex(ctx context.Context, properties ...interface{}) error
	CreateUserIndex(ctx context.Context, properties ...interface{}) error
	CreateLogsIndex(ctx context.Context) error
	SearchPosts(ctx context.Context, keywords []string) ([]PostSearch, error)
	SearchUsers(ctx context.Context, keywords []string) ([]UserSearch, error)
	SearchComments(ctx context.Context, keywords []string) ([]CommentSearch, error)
	ListAllPostsWithAuthorId(ctx context.Context, authorid string) ([]PostSearch, error)
	IsExistsPost(ctx context.Context, postid string) (bool, error)
	IsExistsUser(ctx context.Context, userid string) (bool, error)
	IsExistsComment(ctx context.Context, commentid string) (bool, error)
	InputUser(ctx context.Context, user UserSearch) error
	InputPost(ctx context.Context, post PostSearch) error
	InputComment(ctx context.Context, comment CommentSearch) error
	BulkInputLogs(ctx context.Context, event []ReadEvent) error
	DeleteUserIndex(ctx context.Context, userId int64) error
	DeletePostIndex(ctx context.Context, postId uint) error
	DeleteCommentIndex(ctx context.Context, commentId uint) error
}

type searchDAO struct {
	db     *gorm.DB
	client *elasticsearch.TypedClient
	l      *zap.Logger
}

// PostSearch 定义帖子搜索模型
type PostSearch struct {
	Id       uint     `json:"id"`
	Title    string   `json:"title"`
	AuthorId int64    `json:"author_id"`
	Status   uint8    `json:"status"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
}

// UserSearch 定义用户搜索模型
type UserSearch struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
}
<<<<<<< HEAD
type CommentSearch struct {
	Id       int64  `json:"id"`
	AuthorId int64  `json:"author_id"`
	Content  string `json:"content"`
	Status   uint8  `json:"status"`
}
=======

// ReadEvent 定义日志事件模型
>>>>>>> db4d0af (update)
type ReadEvent struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// NewSearchDAO 创建并返回一个新的 SearchDAO 实例
func NewSearchDAO(db *gorm.DB, client *elasticsearch.TypedClient, l *zap.Logger) SearchDAO {
	return &searchDAO{
		db:     db,
		client: client,
		l:      l,
	}
}

// CreateIndex 创建一个新的索引,可指定mapping属性
func (s *searchDAO) CreateIndex(ctx context.Context, indexName string, properties ...interface{}) error {
	exists, err := s.client.Indices.Exists(indexName).IsSuccess(ctx)
	if err != nil {
		s.l.Error("检查索引是否存在失败", zap.Error(err))
		return err
	}
	if exists {
		return nil
	}

	prop := make(map[string]types.Property)
	if len(properties) > 0 {
		if p, ok := properties[0].(map[string]types.Property); ok {
			prop = p
		} else {
			s.l.Warn("无效的属性类型", zap.Any("properties", properties))
		}
	}

	_, err = s.client.Indices.Create(indexName).Request(&create.Request{
		Mappings: &types.TypeMapping{
			Properties: prop,
		},
	}).Do(ctx)
	if err != nil {
		s.l.Error("创建索引失败", zap.Error(err))
		return err
	}
	return nil
}

<<<<<<< HEAD
// CreateCommentIndex 创建comment的es索引

func (s *searchDAO) CreateCommentIndex(ctx context.Context, properties ...interface{}) error {
	var prop = map[string]types.Property{}
	if len(properties) != 0 {
		prop = properties[0].(map[string]types.Property)
	} else {
		prop = map[string]types.Property{
			"id":        types.NewUnsignedLongNumberProperty(),
			"author.id": types.NewLongNumberProperty(),
			"content":   types.NewTextProperty(),
			"status":    types.NewByteNumberProperty(),
		}
	}
	return s.CreateIndex(ctx, CommentIndex, prop)
}

// CreatePostIndex 创建post的es索引
=======
// CreatePostIndex 创建帖子索引
>>>>>>> db4d0af (update)
func (s *searchDAO) CreatePostIndex(ctx context.Context, properties ...interface{}) error {
	prop := map[string]types.Property{
		"id":        types.NewUnsignedLongNumberProperty(),
		"title":     types.NewTextProperty(),
		"author.id": types.NewLongNumberProperty(),
		"status":    types.NewByteNumberProperty(),
		"content":   types.NewTextProperty(),
		"tags":      types.NewKeywordProperty(),
	}

	if len(properties) > 0 {
		if p, ok := properties[0].(map[string]types.Property); ok {
			prop = p
		}
	}

	return s.CreateIndex(ctx, PostIndex, prop)
}

// CreateUserIndex 创建用户索引
func (s *searchDAO) CreateUserIndex(ctx context.Context, properties ...interface{}) error {
	prop := map[string]types.Property{
		"id":       types.NewUnsignedLongNumberProperty(),
		"email":    types.NewKeywordProperty(),
		"nickname": types.NewTextProperty(),
		"phone":    types.NewKeywordProperty(),
	}

	if len(properties) > 0 {
		if p, ok := properties[0].(map[string]types.Property); ok {
			prop = p
		}
	}
	return s.CreateIndex(ctx, UserIndex, prop)
}

// CreateLogsIndex 创建日志索引
func (s *searchDAO) CreateLogsIndex(ctx context.Context) error {
	prop := map[string]types.Property{
		"timestamp": types.NewDateProperty(),
		"level":     types.NewKeywordProperty(),
		"message":   types.NewTextProperty(),
	}
	return s.CreateIndex(ctx, LogsIndex, prop)
}

<<<<<<< HEAD
// SearchComment 根据关键词搜索评论，返回匹配的结果
func (s *searchDAO) SearchComments(ctx context.Context, keywords []string) ([]CommentSearch, error) {
	queryString := strings.Join(keywords, " ")
	query := types.NewQuery()
	query.Bool = &types.BoolQuery{
		Must: []types.Query{

			types.Query{
				Term: map[string]types.TermQuery{
					"status": {
						Value: 1,
					},
				},
			},
			types.Query{
				MultiMatch: &types.MultiMatchQuery{
					Query:  queryString,
					Fields: []string{"content"},
				},
			},
		},
	}
	// 创建并执行搜索请求
	resp, err := s.client.Search().Index(CommentIndex).Query(query).Do(ctx)
	if err != nil {
		s.l.Error("Search request failed", zap.Error(err))
		return nil, err
	}
	// 将查询结果反序列化为 CommentSearch 对象
	var comments []CommentSearch
	for _, hit := range resp.Hits.Hits {
		var comment CommentSearch
		if err := json.Unmarshal(hit.Source_, &comment); err != nil {
			s.l.Error("Failed to unmarshal search hit", zap.Error(err))
			return nil, err
		}
		comments = append(comments, comment)
	}
	s.l.Info("Successfully completed SearchComments", zap.Int("resultCount", len(comments)))
	return comments, nil
}

// SearchPosts 根据关键词搜索帖子，返回匹配的结果
=======
// SearchPosts 根据关键词搜索帖子
>>>>>>> db4d0af (update)
func (s *searchDAO) SearchPosts(ctx context.Context, keywords []string) ([]PostSearch, error) {
	queryString := strings.Join(keywords, " ")

	query := types.NewQuery()
	query.Bool = &types.BoolQuery{
		Must: []types.Query{
			{
				Term: map[string]types.TermQuery{
					"status.keyword": {Value: "Published"},
				},
			},
			{
				MultiMatch: &types.MultiMatchQuery{
					Query:  queryString,
					Fields: []string{"title", "content"},
				},
			},
		},
	}

	resp, err := s.client.Search().Index(PostIndex).Query(query).Do(ctx)
	if err != nil {
		s.l.Error("搜索请求失败", zap.Error(err))
		return nil, err
	}

	posts := make([]PostSearch, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var post PostSearch
		if err := json.Unmarshal(hit.Source_, &post); err != nil {
			s.l.Error("解析搜索结果失败", zap.Error(err))
			return nil, err
		}
		posts = append(posts, post)
	}

	s.l.Info("搜索帖子完成", zap.Int("结果数量", len(posts)))
	return posts, nil
}

// ListAllPostsWithAuthorId 根据作者ID查找所有帖子
func (s *searchDAO) ListAllPostsWithAuthorId(ctx context.Context, authorid string) ([]PostSearch, error) {
	query := types.NewQuery()
	query.Term = map[string]types.TermQuery{
		"author.id": {Value: authorid},
	}

	resp, err := s.client.Search().Index(PostIndex).Query(query).Do(ctx)
	if err != nil {
		s.l.Error("搜索请求失败", zap.Error(err))
		return nil, err
	}

	posts := make([]PostSearch, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var post PostSearch
		if err := json.Unmarshal(hit.Source_, &post); err != nil {
			s.l.Error("解析搜索结果失败", zap.Error(err))
			return nil, err
		}
		posts = append(posts, post)
	}

	s.l.Info("查找作者帖子完成", zap.Int("结果数量", len(posts)))
	return posts, nil
}

// SearchUsers 根据关键词搜索用户
func (s *searchDAO) SearchUsers(ctx context.Context, keywords []string) ([]UserSearch, error) {
	queryString := strings.Join(keywords, " ")

	query := types.NewQuery()
	query.Bool = &types.BoolQuery{
		Should: []types.Query{
			{
				Match: map[string]types.MatchQuery{
					"email": {Query: queryString},
				},
			},
			{
				Match: map[string]types.MatchQuery{
					"nickname": {Query: queryString},
				},
			},
			{
				Match: map[string]types.MatchQuery{
					"phone": {Query: queryString},
				},
			},
		},
	}

	resp, err := s.client.Search().Index(UserIndex).Query(query).Do(ctx)
	if err != nil {
		s.l.Error("搜索请求失败", zap.Error(err))
		return nil, err
	}

	users := make([]UserSearch, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var user UserSearch
		if err := json.Unmarshal(hit.Source_, &user); err != nil {
			s.l.Error("解析搜索结果失败", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}

	s.l.Info("搜索用户完成", zap.Int("结果数量", len(users)))
	return users, nil
}

// IsExistsPost 检查帖子是否存在
func (s *searchDAO) IsExistsPost(ctx context.Context, postid string) (bool, error) {
	return s.client.Exists(PostIndex, postid).Do(ctx)
}

// IsExistsUser 检查用户是否存在
func (s *searchDAO) IsExistsUser(ctx context.Context, userid string) (bool, error) {
	return s.client.Exists(UserIndex, userid).Do(ctx)
}

<<<<<<< HEAD
// IsExistsComment 查看指定commentId的comment是否存在
func (s *searchDAO) IsExistsComment(ctx context.Context, commentid string) (bool, error) {
	return s.client.Exists(CommentIndex, commentid).Do(ctx)
}

// InputUser 将用户信息输入到 Elasticsearch 索引中
=======
// InputUser 添加用户到搜索索引
>>>>>>> db4d0af (update)
func (s *searchDAO) InputUser(ctx context.Context, user UserSearch) error {
	_, err := s.client.Index(UserIndex).
		Id(strconv.FormatInt(user.Id, 10)).
		Document(user).
		Do(ctx)
	if err != nil {
		s.l.Error("创建用户索引失败", zap.Error(err))
		return err
	}
	return nil
}

// InputPost 添加帖子到搜索索引
func (s *searchDAO) InputPost(ctx context.Context, post PostSearch) error {
	_, err := s.client.Index(PostIndex).
		Id(strconv.FormatInt(int64(post.Id), 10)).
		Document(post).
		Do(ctx)
	if err != nil {
		s.l.Error("创建帖子索引失败", zap.Error(err))
		return err
	}
	return nil
}

<<<<<<< HEAD
// InputComment 将评论信息输入到 Elasticsearch 索引中
func (s *searchDAO) InputComment(ctx context.Context, comment CommentSearch) error {
	_, err := s.client.Index(CommentIndex).
		Id(strconv.FormatInt(int64(comment.Id), 10)).
		Document(comment).
		Do(ctx)
	if err != nil {
		s.l.Error("Failed to create comment index", zap.Error(err))
		return err
	}
	return nil
}

// BulkInputLogs 批量向es插入日志
func (s *searchDAO) BulkInputLogs(ctx context.Context, event []ReadEvent) error {
	var req bulk.Request
	for _, eve := range event {
		req = append(req, eve)
	}
	if _, err := s.client.Bulk().Index(LogsIndex).Request(&req).Do(ctx); err != nil {
		s.l.Error("bulk index failed", zap.Error(err))
=======
// BulkInputLogs 批量添加日志
func (s *searchDAO) BulkInputLogs(ctx context.Context, events []ReadEvent) error {
	req := make(bulk.Request, len(events))
	for i, event := range events {
		req[i] = event
>>>>>>> db4d0af (update)
	}

	if _, err := s.client.Bulk().Index(LogsIndex).Request(&req).Do(ctx); err != nil {
		s.l.Error("批量索引失败", zap.Error(err))
		return err
	}

	s.l.Info("批量索引成功")
	return nil
}

// DeleteUserIndex 删除用户索引
func (s *searchDAO) DeleteUserIndex(ctx context.Context, userId int64) error {
	return s.deleteIndex(ctx, UserIndex, strconv.FormatInt(userId, 10))
}

// DeletePostIndex 删除帖子索引
func (s *searchDAO) DeletePostIndex(ctx context.Context, postId uint) error {
	return s.deleteIndex(ctx, PostIndex, strconv.FormatInt(int64(postId), 10))
}

<<<<<<< HEAD
// DeleteCommentIndex 从 Elasticsearch 索引中删除指定评论
func (s *searchDAO) DeleteCommentIndex(ctx context.Context, commentId uint) error {
	return s.deleteIndex(ctx, CommentIndex, strconv.FormatInt(int64(commentId), 10))
}

// deleteIndex 根据索引名称和文档 ID 删除 Elasticsearch 中的文档
=======
// deleteIndex 删除指定索引
>>>>>>> db4d0af (update)
func (s *searchDAO) deleteIndex(ctx context.Context, index, docID string) error {
	resp, err := s.client.Delete(index, docID).Do(ctx)
	if err != nil {
		s.l.Error("删除索引失败", zap.String("index", index), zap.Error(err))
		return err
	}
	s.l.Info("成功删除索引", zap.String("index", resp.Index_), zap.String("docID", resp.Id_))
	return nil
}

// handleElasticsearchError 处理ES错误响应
func (s *searchDAO) handleElasticsearchError(resp *esapi.Response) error {
	var errMsg map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&errMsg); err == nil {
		s.l.Error("ES返回错误响应",
			zap.String("status", resp.Status()),
			zap.Any("error", errMsg))
		return fmt.Errorf("ES返回错误响应: %s", resp.Status())
	}

	s.l.Error("ES返回错误响应", zap.String("status", resp.Status()))
	return fmt.Errorf("ES返回错误响应: %s", resp.Status())
}
