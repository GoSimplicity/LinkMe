package service

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"golang.org/x/sync/errgroup"
	"strings"
)

type searchService struct {
	repo repository.SearchRepository
}

type SearchService interface {
	SearchPosts(ctx context.Context, expression string) ([]domain.PostSearch, error)
	SearchUsers(ctx context.Context, expression string) ([]domain.UserSearch, error)
	SearchComments(ctx context.Context, expression string) ([]domain.CommentSearch, error)
}

func NewSearchService(repo repository.SearchRepository) SearchService {
	return &searchService{
		repo: repo,
	}
}
func (s *searchService) SearchComments(ctx context.Context, expression string) ([]domain.CommentSearch, error) {
	// 将表达式拆分为关键字数组
	keywords := strings.Split(expression, " ")
	var eg errgroup.Group
	var comments []domain.CommentSearch
	eg.Go(func() error {
		// 搜索评论
		foundComments, err := s.repo.SearchComments(ctx, keywords)
		if err != nil {
			return err
		}
		comments = foundComments
		return nil
	})
	// 等待所有并发任务完成
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return comments, nil
}
func (s *searchService) SearchPosts(ctx context.Context, expression string) ([]domain.PostSearch, error) {
	// 将表达式拆分为关键字数组
	keywords := strings.Split(expression, " ")

	var eg errgroup.Group
	var posts []domain.PostSearch

	eg.Go(func() error {
		// 搜索帖子
		foundPosts, err := s.repo.SearchPosts(ctx, keywords)
		if err != nil {
			return err
		}

		posts = foundPosts

		return nil
	})

	// 等待所有并发任务完成
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *searchService) SearchUsers(ctx context.Context, expression string) ([]domain.UserSearch, error) {
	// 将表达式拆分为关键字数组
	keywords := strings.Split(expression, " ")

	var eg errgroup.Group
	var users []domain.UserSearch

	eg.Go(func() error {
		// 搜索用户
		foundUsers, err := s.repo.SearchUsers(ctx, keywords)
		if err != nil {
			return err
		}

		users = foundUsers

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return users, nil
}
