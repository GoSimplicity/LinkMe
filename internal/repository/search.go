package repository

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
)

type searchRepository struct {
	dao dao.SearchDAO
}

type SearchRepository interface {
	SearchPosts(ctx context.Context, userID int64, expression string) ([]domain.UserSearch, error)
	SearchUsers(ctx context.Context, userID int64, expression string) ([]domain.PostSearch, error)
}

func NewSearchRepository(dao dao.SearchDAO) SearchRepository {
	return &searchRepository{
		dao: dao,
	}
}

func (s searchRepository) SearchPosts(ctx context.Context, userID int64, expression string) ([]domain.UserSearch, error) {
	//TODO implement me
	panic("implement me")
}

func (s searchRepository) SearchUsers(ctx context.Context, userID int64, expression string) ([]domain.PostSearch, error) {
	//TODO implement me
	panic("implement me")
}
