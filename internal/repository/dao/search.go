package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const (
	PostIndex = "post_index"
	UserIndex = "user_index"
)

type SearchDAO interface {
	SearchPosts(ctx context.Context, keywords []string) ([]PostSearch, error)
	SearchUsers(ctx context.Context, keywords []string) ([]UserSearch, error)
	InputUser(ctx context.Context, user UserSearch) error
	InputPost(ctx context.Context, post PostSearch) error
	DeleteUserIndex(ctx context.Context, userId int64) error
	DeletePostIndex(ctx context.Context, postId uint) error
}

type searchDAO struct {
	db     *gorm.DB
	client *elasticsearch.TypedClient
	l      *zap.Logger
}

type PostSearch struct {
	Id      uint     `json:"id"`
	Title   string   `json:"title"`
	Status  uint8    `json:"status"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type UserSearch struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

// NewSearchDAO 创建并返回一个新的 SearchDAO 实例
func NewSearchDAO(db *gorm.DB, client *elasticsearch.TypedClient, l *zap.Logger) SearchDAO {
	return &searchDAO{
		db:     db,
		client: client,
		l:      l,
	}
}

// SearchPosts 根据关键词搜索帖子，返回匹配的结果
func (s *searchDAO) SearchPosts(ctx context.Context, keywords []string) ([]PostSearch, error) {
	queryString := strings.Join(keywords, " ")
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"status.keyword": "Published", // 仅匹配发布状态的帖子
						},
					},
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":  queryString,
							"fields": []string{"title", "content"}, // 在标题和内容中搜索
						},
					},
				},
			},
		},
	}

	// 序列化查询 JSON
	queryBytes, err := json.Marshal(query)
	if err != nil {
		s.l.Error("Failed to serialize search query", zap.Error(err))
		return nil, err
	}

	// 创建搜索请求
	req := esapi.SearchRequest{
		Index: []string{PostIndex},
		Body:  strings.NewReader(string(queryBytes)),
	}

	// 执行搜索请求
	resp, err := req.Do(ctx, s.client)
	if err != nil {
		s.l.Error("Search request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	// 检查并处理响应中的错误
	if resp.IsError() {
		return nil, s.handleElasticsearchError(resp)
	}

	// 解析响应结果
	var searchResult struct {
		Hits struct {
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		s.l.Error("Failed to decode search response", zap.Error(err))
		return nil, err
	}

	// 将查询结果反序列化为 PostSearch 对象
	var posts []PostSearch

	for _, hit := range searchResult.Hits.Hits {
		var post PostSearch
		if err := json.Unmarshal(hit.Source, &post); err != nil {
			s.l.Error("Failed to unmarshal search hit", zap.Error(err))
			return nil, err
		}
		posts = append(posts, post)
	}

	s.l.Info("Successfully completed SearchPosts", zap.Int("resultCount", len(posts)))
	return posts, nil
}

// SearchUsers 根据关键词搜索用户，返回匹配的结果
func (s *searchDAO) SearchUsers(ctx context.Context, keywords []string) ([]UserSearch, error) {
	queryString := strings.Join(keywords, " ")
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"email": queryString, // 邮箱匹配查询
						},
					},
					map[string]interface{}{
						"match": map[string]interface{}{
							"nickname": queryString, // 昵称匹配查询
						},
					},
					map[string]interface{}{
						"match": map[string]interface{}{
							"phone": queryString, // 电话匹配查询
						},
					},
				},
			},
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		s.l.Error("Failed to serialize search query", zap.Error(err))
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{UserIndex},
		Body:  strings.NewReader(string(queryBytes)),
	}

	resp, err := req.Do(ctx, s.client)
	if err != nil {
		s.l.Error("Search request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, s.handleElasticsearchError(resp)
	}

	var searchResult struct {
		Hits struct {
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		s.l.Error("Failed to decode search response", zap.Error(err))
		return nil, err
	}

	var users []UserSearch
	for _, hit := range searchResult.Hits.Hits {
		var user UserSearch
		if err := json.Unmarshal(hit.Source, &user); err != nil {
			s.l.Error("Failed to unmarshal search hit", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}

	s.l.Info("Successfully completed SearchUsers", zap.Int("resultCount", len(users)))
	return users, nil
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
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: docID,
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		s.l.Error(fmt.Sprintf("Failed to delete %s index", index), zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		s.l.Error(fmt.Sprintf("Delete %s index response error", index), zap.String("status", res.Status()))
		return fmt.Errorf("error deleting %s index: %s", index, res.Status())
	}

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
