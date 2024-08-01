package service

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
)

type searchService struct {
	repo repository.SearchRepository
}

type SearchService interface {
	SearchPosts(ctx context.Context, userID int64, expression string) ([]domain.UserSearch, error)
	SearchUsers(ctx context.Context, userID int64, expression string) ([]domain.PostSearch, error)
}

func NewSearchService(repo repository.SearchRepository) SearchService {
	return &searchService{
		repo: repo,
	}
}

func (s searchService) SearchPosts(ctx context.Context, userID int64, expression string) ([]domain.UserSearch, error) {
	//TODO implement me
	panic("implement me")
}

func (s searchService) SearchUsers(ctx context.Context, userID int64, expression string) ([]domain.PostSearch, error) {
	//TODO implement me
	panic("implement me")
}
