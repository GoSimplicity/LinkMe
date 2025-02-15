package repository_test

import (
	"context"
	"encoding/json"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"testing"
)

// 辅助函数，用于创建一个模拟的Elasticsearch客户端，便于测试
func createMockElasticsearchClient() *elasticsearch.TypedClient {
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "elastic",
		Password:  "07X-o2S2eD7TrGodvKJw",
	})
	if err != nil {
		panic(err)
	}
	return client
}

// 辅助函数，用于创建一个模拟的GORM数据库实例，这里简单返回nil，根据实际情况可进一步完善模拟逻辑
func createMockGormDB() *gorm.DB {
	return nil
}

// 辅助函数，用于创建一个模拟的Zap日志记录器，这里简单返回一个空的实现，根据实际情况可替换为真实日志记录器或更完善的模拟
func createMockLogger() *zap.Logger {
	return zap.NewNop()
}

// 测试CreatePostIndex函数
func TestCreatePostIndex(t *testing.T) {
	searchDAO := dao.NewSearchDAO(createMockElasticsearchClient(), createMockLogger())
	err := searchDAO.CreatePostIndex(context.Background())
	if err != nil {
		t.Errorf("CreatePostIndex failed: %v", err)
	}
}

// 测试SearchPosts函数
func TestSearchPosts(t *testing.T) {
	searchDAO := dao.NewSearchDAO(createMockElasticsearchClient(), createMockLogger())
	keywords := []string{"test", "post"}
	posts, err := searchDAO.SearchPosts(context.Background(), keywords)
	if err != nil {
		t.Errorf("SearchPosts failed: %v", err)
		return
	}
	if len(posts) == 0 {
		t.Log("No posts found, but test passed as long as there's no error")
		return
	}
	// 可以进一步验证返回的帖子数据结构是否符合预期
	for _, post := range posts {
		postJSON, _ := json.Marshal(post)
		t.Logf("Retrieved post: %s", postJSON)
	}
}

// 测试SearchUsers函数
func TestSearchUsers(t *testing.T) {
	searchDAO := dao.NewSearchDAO(createMockElasticsearchClient(), createMockLogger())
	keywords := []string{"test", "user"}
	users, err := searchDAO.SearchUsers(context.Background(), keywords)
	if err != nil {
		t.Errorf("SearchUsers failed: %v", err)
	}
	if len(users) == 0 {
		t.Log("No users found, but test passed as long as there's no error")
	}
	// 可以进一步验证返回的用户数据结构是否符合预期
	for _, user := range users {
		userJSON, _ := json.Marshal(user)
		t.Logf("Retrieved user: %s", userJSON)
	}
}
