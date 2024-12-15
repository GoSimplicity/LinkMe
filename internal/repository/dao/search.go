package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"

	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const (
	PostIndex = "post_index"
	UserIndex = "user_index"
	LogsIndex = "logs_index"
)

type SearchDAO interface {
	CreateIndex(ctx context.Context, indexName string, properties ...interface{}) error
	CreatePostIndex(ctx context.Context, properties ...interface{}) error
	CreateUserIndex(ctx context.Context, properties ...interface{}) error
	CreateLogsIndex(ctx context.Context) error
	SearchPosts(ctx context.Context, keywords []string) ([]PostSearch, error)
	SearchUsers(ctx context.Context, keywords []string) ([]UserSearch, error)
	ListAllPostsWithAuthorId(ctx context.Context, authorid string) ([]PostSearch, error)
	IsExistsPost(ctx context.Context, postid string) (bool, error)
	IsExistsUser(ctx context.Context, userid string) (bool, error)
	InputUser(ctx context.Context, user UserSearch) error
	InputPost(ctx context.Context, post PostSearch) error
	BulkInputLogs(ctx context.Context, event []ReadEvent) error
	DeleteUserIndex(ctx context.Context, userId int64) error
	DeletePostIndex(ctx context.Context, postId uint) error
}

type searchDAO struct {
	db     *gorm.DB
	client *elasticsearch.TypedClient
	l      *zap.Logger
}

type PostSearch struct {
	Id       uint     `json:"id"`
	Title    string   `json:"title"`
	AuthorId int64    `json:"author_id"`
	Status   uint8    `json:"status"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
}

type UserSearch struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
}

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

// CreateIndex 创建一个新的index, 可指定mapping属性
func (s *searchDAO) CreateIndex(ctx context.Context, indexName string, properties ...interface{}) error {
	if success, err := s.client.Indices.Exists(indexName).IsSuccess(ctx); success || err != nil {
		if err != nil {
			s.l.Error("Failed to check if index exists", zap.Error(err))
		}
		return nil
	}

	prop := map[string]types.Property{}
	if len(properties) != 0 {
		if p, ok := properties[0].(map[string]types.Property); ok {
			prop = p
		} else {
			s.l.Info("invalid properties type", zap.Any("properties", properties))
		}
	}

	_, err := s.client.Indices.Create(indexName).Request(&create.Request{
		Mappings: &types.TypeMapping{
			Properties: prop,
		},
	}).Do(ctx)
	if err != nil {
		s.l.Error("create index failed", zap.Error(err))
	}
	return nil
}

// CreatePostIndex 创建post的es索引
func (s *searchDAO) CreatePostIndex(ctx context.Context, properties ...interface{}) error {
	var prop = map[string]types.Property{}
	if len(properties) != 0 {
		prop = properties[0].(map[string]types.Property)
	} else {
		prop = map[string]types.Property{
			"id":        types.NewUnsignedLongNumberProperty(),
			"title":     types.NewTextProperty(),
			"author.id": types.NewLongNumberProperty(),
			"status":    types.NewByteNumberProperty(),
			"content":   types.NewTextProperty(),
			"tags":      types.NewKeywordProperty(),
		}
	}

	return s.CreateIndex(ctx, PostIndex, prop)
}

// CreateUserIndex 创建uesr的es索引
func (s *searchDAO) CreateUserIndex(ctx context.Context, properties ...interface{}) error {
	var prop = map[string]types.Property{}
	if len(properties) != 0 {
		prop = properties[0].(map[string]types.Property)
	} else {
		prop = map[string]types.Property{
			"id":       types.NewUnsignedLongNumberProperty(),
			"email":    types.NewKeywordProperty(),
			"nickname": types.NewTextProperty(),
			"phone":    types.NewKeywordProperty(),
		}
	}
	return s.CreateIndex(ctx, UserIndex, prop)
}

// CreateLogsIndex 创建logs的es索引
func (s *searchDAO) CreateLogsIndex(ctx context.Context) error {
	prop := map[string]types.Property{
		"timestamp": types.NewDateProperty(),
		"level":     types.NewKeywordProperty(),
		"message":   types.NewTextProperty(),
	}
	return s.CreateIndex(ctx, LogsIndex, prop)
}

// SearchPosts 根据关键词搜索帖子，返回匹配的结果
func (s *searchDAO) SearchPosts(ctx context.Context, keywords []string) ([]PostSearch, error) {
	queryString := strings.Join(keywords, " ")

	query := types.NewQuery()
	query.Bool = &types.BoolQuery{
		Must: []types.Query{
			types.Query{
				Term: map[string]types.TermQuery{
					"status.keyword": {
						Value: "Published",
					},
				},
			},
			types.Query{
				MultiMatch: &types.MultiMatchQuery{
					Query:  queryString,
					Fields: []string{"title", "content"},
				},
			},
		},
	}

	// 创建并执行搜索请求
	resp, err := s.client.Search().Index(PostIndex).Query(query).Do(ctx)
	if err != nil {
		s.l.Error("Search request failed", zap.Error(err))
		return nil, err
	}

	// 将查询结果反序列化为 PostSearch 对象
	var posts []PostSearch

	for _, hit := range resp.Hits.Hits {
		var post PostSearch
		if err := json.Unmarshal(hit.Source_, &post); err != nil {
			s.l.Error("Failed to unmarshal search hit", zap.Error(err))
			return nil, err
		}
		posts = append(posts, post)
	}

	s.l.Info("Successfully completed SearchPosts", zap.Int("resultCount", len(posts)))
	return posts, nil
}

// ListAllPostsWithAuthorId 根据authorId 查找所有post
func (s *searchDAO) ListAllPostsWithAuthorId(ctx context.Context, authorid string) ([]PostSearch, error) {
	query := types.NewQuery()
	query.Term = map[string]types.TermQuery{
		"author.id": {
			Value: authorid,
		},
	}

	//创建并执行搜索请求
	resp, err := s.client.Search().Index(PostIndex).Query(query).Do(ctx)
	if err != nil {
		s.l.Error("Search request failed", zap.Error(err))
		return nil, err
	}

	// 将查询结果反序列化为 PostSearch 对象
	var posts []PostSearch

	for _, hit := range resp.Hits.Hits {
		var post PostSearch
		if err := json.Unmarshal(hit.Source_, &post); err != nil {
			s.l.Error("Failed to unmarshal search hit", zap.Error(err))
			return nil, err
		}
		posts = append(posts, post)
	}

	s.l.Info("Successfully completed ListAllPostsWithAuthor", zap.Int("resultCount", len(posts)))
	return posts, nil
}

// SearchUsers 根据关键词搜索用户，返回匹配的结果
func (s *searchDAO) SearchUsers(ctx context.Context, keywords []string) ([]UserSearch, error) {
	queryString := strings.Join(keywords, " ")

	query := types.NewQuery()
	query.Bool = &types.BoolQuery{
		Should: []types.Query{
			types.Query{
				Match: map[string]types.MatchQuery{
					"email": types.MatchQuery{
						Query: queryString,
					},
				},
			},
			types.Query{
				Match: map[string]types.MatchQuery{
					"nickname": types.MatchQuery{
						Query: queryString,
					},
				},
			},
			types.Query{
				Match: map[string]types.MatchQuery{
					"phone": types.MatchQuery{
						Query: queryString,
					},
				},
			},
		},
	}

	resp, err := s.client.Search().Index(UserIndex).Query(query).Do(ctx)
	if err != nil {
		s.l.Error("Search request failed", zap.Error(err))
		return nil, err
	}

	var users []UserSearch
	for _, hit := range resp.Hits.Hits {
		var user UserSearch
		if err := json.Unmarshal(hit.Source_, &user); err != nil {
			s.l.Error("Failed to unmarshal search hit", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}

	s.l.Info("Successfully completed SearchUsers", zap.Int("resultCount", len(users)))
	return users, nil
}

// IsExistsPost 查看指定postId的post是否存在
func (s *searchDAO) IsExistsPost(ctx context.Context, postid string) (bool, error) {
	return s.client.Exists(PostIndex, postid).Do(ctx)
}

// IsExistsUser 查看指定userId的user是否存在
func (s *searchDAO) IsExistsUser(ctx context.Context, userid string) (bool, error) {
	return s.client.Exists(UserIndex, userid).Do(ctx)
}

// InputUser 将用户信息输入到 Elasticsearch 索引中
func (s *searchDAO) InputUser(ctx context.Context, user UserSearch) error {
	_, err := s.client.Index(UserIndex).
		Id(strconv.FormatInt(user.Id, 10)).
		Document(user).
		Do(ctx)
	if err != nil {
		s.l.Error("Failed to create user index", zap.Error(err))
		return err
	}

	return nil
}

// InputPost 将帖子信息输入到 Elasticsearch 索引中
func (s *searchDAO) InputPost(ctx context.Context, post PostSearch) error {
	_, err := s.client.Index(PostIndex).
		Id(strconv.FormatInt(int64(post.Id), 10)).
		Document(post).
		Do(ctx)
	if err != nil {
		s.l.Error("Failed to create post index", zap.Error(err))
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
	}

	s.l.Info("bulk index successfully")
	return nil
}

// DeleteUserIndex 从 Elasticsearch 索引中删除指定用户
func (s *searchDAO) DeleteUserIndex(ctx context.Context, userId int64) error {
	return s.deleteIndex(ctx, UserIndex, strconv.FormatInt(userId, 10))
}

// DeletePostIndex 从 Elasticsearch 索引中删除指定帖子
func (s *searchDAO) DeletePostIndex(ctx context.Context, postId uint) error {
	return s.deleteIndex(ctx, PostIndex, strconv.FormatInt(int64(postId), 10))
}

// deleteIndex 根据索引名称和文档 ID 删除 Elasticsearch 中的文档
func (s *searchDAO) deleteIndex(ctx context.Context, index, docID string) error {

	resq, err := s.client.Delete(index, docID).Do(ctx)
	if err != nil {
		s.l.Error(fmt.Sprintf("Failed to delete %s index", index), zap.Error(err))
		return err
	}
	s.l.Info("Successfully deleted index", zap.String("index", resq.Index_), zap.String("docID", resq.Id_))

	return nil
}

// handleElasticsearchError 处理 Elasticsearch 返回的错误响应
func (s *searchDAO) handleElasticsearchError(resp *esapi.Response) error {
	var errMsg map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&errMsg); err == nil {
		s.l.Error("Elasticsearch returned an error response",
			zap.String("status", resp.Status()),
			zap.Any("error", errMsg))
		return fmt.Errorf("elasticsearch returned an error response: %s", resp.Status())
	}

	s.l.Error("Elasticsearch returned an error response", zap.String("status", resp.Status()))

	return fmt.Errorf("elasticsearch returned an error response: %s", resp.Status())
}
