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
	SearchPosts(ctx context.Context, keywords []string) ([]domain.PostSearch, error)
	SearchUsers(ctx context.Context, keywords []string) ([]domain.UserSearch, error)
}

func NewSearchRepository(dao dao.SearchDAO) SearchRepository {
	return &searchRepository{
		dao: dao,
	}
}

func (s *searchRepository) SearchPosts(ctx context.Context, keywords []string) ([]domain.PostSearch, error) {
	posts, err := s.dao.SearchPosts(ctx, keywords)
	return s.toDomainPostSearch(posts), err
}

func (s *searchRepository) SearchUsers(ctx context.Context, keywords []string) ([]domain.UserSearch, error) {
	users, err := s.dao.SearchUsers(ctx, keywords)
	return s.toDomainUserSearch(users), err
}

func (s *searchRepository) toDomainPostSearch(daoPosts []dao.PostSearch) []domain.PostSearch {
	domainPosts := make([]domain.PostSearch, len(daoPosts))
	for i, daoPost := range daoPosts {
		domainPosts[i] = domain.PostSearch{
			Content: daoPost.Content,
			Id:      daoPost.Id,
			Status:  daoPost.Status,
			Tags:    daoPost.Tags,
			Title:   daoPost.Title,
		}
	}
	return domainPosts
}

func (s *searchRepository) toDomainUserSearch(daoUsers []dao.UserSearch) []domain.UserSearch {
	domainUsers := make([]domain.UserSearch, len(daoUsers))
	for i, daoUser := range daoUsers {
		domainUsers[i] = domain.UserSearch{
			Email:    daoUser.Email,
			Id:       daoUser.Id,
			Nickname: daoUser.Nickname,
			Phone:    daoUser.Phone,
		}
	}
	return domainUsers
}
