package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

const (
	PostIndex = "post_index"
	UserIndex = "user_index"
)

type SearchDAO interface {
	//SearchPosts(ctx context.Context, PostIds []int64, keywords []string) ([]PostSearch, error)
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
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

func NewSearchDAO(db *gorm.DB, client *elasticsearch.TypedClient, l *zap.Logger) SearchDAO {
	return &searchDAO{
		db:     db,
		client: client,
		l:      l,
	}
}

func (s *searchDAO) SearchPosts(ctx context.Context, keywords []string) ([]PostSearch, error) {
	queryString := strings.Join(keywords, " ")
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"status.keyword": "Published", // 使用 .keyword 来匹配未分词的字段值
						},
					},
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":  queryString,
							"fields": []string{"title", "content"},
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
		Index: []string{PostIndex},
		Body:  strings.NewReader(string(queryBytes)),
	}
	resp, err := req.Do(ctx, s.client)
	if err != nil {
		s.l.Error("Search request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		var errMsg map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errMsg); err == nil {
			s.l.Error("Elasticsearch returned an error response",
				zap.String("status", resp.Status()),
				zap.Any("error", errMsg))
			return nil, fmt.Errorf("elasticsearch returned an error response: %s", resp.Status())
		}
		s.l.Error("Elasticsearch returned an error response",
			zap.String("status", resp.Status()))
		return nil, fmt.Errorf("elasticsearch returned an error response: %s", resp.Status())
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
	s.l.Info("Decoded search response", zap.Int("numHits", len(searchResult.Hits.Hits)))
	res := make([]PostSearch, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		var post PostSearch
		if err := json.Unmarshal(hit.Source, &post); err != nil {
			s.l.Error("Failed to unmarshal search hit", zap.Error(err))
			return nil, err
		}
		res = append(res, post)
	}
	s.l.Info("Successfully completed SearchPosts", zap.Int("resultCount", len(res)))
	return res, nil
}

func (s *searchDAO) SearchUsers(ctx context.Context, keywords []string) ([]UserSearch, error) {
	// 将关键词数组拼接成一个字符串
	queryString := strings.Join(keywords, " ")
	// 构建查询 JSON
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
	// 将查询 JSON 序列化为字节数组
	queryBytes, err := json.Marshal(query)
	if err != nil {
		s.l.Error("Failed to serialize search query", zap.Error(err))
		return nil, err
	}
	// 创建搜索请求
	req := esapi.SearchRequest{
		Index: []string{UserIndex}, // 索引名称
		Body:  strings.NewReader(string(queryBytes)),
	}
	// 执行搜索请求
	resp, err := req.Do(ctx, s.client)
	if err != nil {
		s.l.Error("Search request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	// 检查响应是否包含错误
	if resp.IsError() {
		var errMsg map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errMsg); err == nil {
			s.l.Error("Elasticsearch returned an error response",
				zap.String("status", resp.Status()),
				zap.Any("error", errMsg))
			return nil, fmt.Errorf("elasticsearch returned an error response: %s", resp.Status())
		}
		s.l.Error("Elasticsearch returned an error response",
			zap.String("status", resp.Status()))
		return nil, fmt.Errorf("elasticsearch returned an error response: %s", resp.Status())
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
	// 初始化结果切片
	res := make([]UserSearch, 0, len(searchResult.Hits.Hits))
	// 遍历查询结果，将每个命中的文档反序列化为 UserSearch 对象
	for _, hit := range searchResult.Hits.Hits {
		var user UserSearch
		if err := json.Unmarshal(hit.Source, &user); err != nil {
			s.l.Error("Failed to unmarshal search hit", zap.Error(err))
			return nil, err
		}
		res = append(res, user)
	}
	return res, nil
}

func (s *searchDAO) InputUser(ctx context.Context, user UserSearch) error {
	_, err := s.client.Index(UserIndex).Id(strconv.FormatInt(user.Id, 10)).Document(user).Do(ctx)
	if err != nil {
		s.l.Error("create user index failed", zap.Error(err))
		return err
	}
	return nil
}

func (s *searchDAO) InputPost(ctx context.Context, post PostSearch) error {
	_, err := s.client.Index(PostIndex).Id(strconv.FormatInt(int64(post.Id), 10)).Document(post).Do(ctx)
	if err != nil {
		s.l.Error("create post index failed", zap.Error(err))
		return err
	}
	return nil
}

func (s *searchDAO) DeleteUserIndex(ctx context.Context, userId int64) error {
	req := esapi.DeleteRequest{
		Index:      UserIndex,
		DocumentID: strconv.FormatInt(userId, 10),
	}
	res, err := req.Do(ctx, s.client)
	if err != nil {
		s.l.Error("delete user index failed", zap.Error(err))
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		s.l.Error("delete user index response error", zap.String("status", res.Status()))
		return fmt.Errorf("error deleting user index: %s", res.Status())
	}
	return nil
}

func (s *searchDAO) DeletePostIndex(ctx context.Context, postId uint) error {
	req := esapi.DeleteRequest{
		Index:      PostIndex,
		DocumentID: strconv.FormatInt(int64(postId), 10),
	}
	res, err := req.Do(ctx, s.client)
	if err != nil {
		s.l.Error("delete post index failed", zap.Error(err))
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		s.l.Error("delete post index response error", zap.String("status", res.Status()))
		return fmt.Errorf("error deleting post index: %s", res.Status())
	}
	return nil
}
