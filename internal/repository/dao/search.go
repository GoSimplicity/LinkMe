package dao

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SearchDAO interface {
	SearchPosts(ctx context.Context, userID int64, expression string) ([]UserSearch, error)
	SearchUsers(ctx context.Context, userID int64, expression string) ([]PostSearch, error)
	InputUser(ctx context.Context, user User) error
	InputPost(ctx context.Context, post Post) error
}

type searchDAO struct {
	db     *gorm.DB
	client *elasticsearch.Client
	l      *zap.Logger
}

type PostSearch struct {
	Id      int64    `json:"id"`
	Title   string   `json:"title"`
	Status  int32    `json:"status"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type UserSearch struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

func NewSearchDAO(db *gorm.DB, client *elasticsearch.Client, l *zap.Logger) SearchDAO {
	return &searchDAO{
		db:     db,
		client: client,
		l:      l,
	}
}

func (s *searchDAO) test() {
	s.client.Ping.WithContext(context.Background())
}

func (s *searchDAO) SearchPosts(ctx context.Context, userID int64, expression string) ([]UserSearch, error) {
	//TODO implement me
	panic("implement me")
}

func (s *searchDAO) SearchUsers(ctx context.Context, userID int64, expression string) ([]PostSearch, error) {
	//TODO implement me
	panic("implement me")
}

func (s *searchDAO) InputUser(ctx context.Context, user User) error {
	//TODO implement me
	panic("implement me")
}

func (s *searchDAO) InputPost(ctx context.Context, post Post) error {
	//TODO implement me
	panic("implement me")
}
