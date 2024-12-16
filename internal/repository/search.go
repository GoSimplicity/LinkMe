package repository

import (
	"context"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
)

type searchRepository struct {
	dao dao.SearchDAO
}

type SearchRepository interface {
	SearchPosts(ctx context.Context, keywords []string) ([]domain.PostSearch, error) // 搜索文章
	SearchUsers(ctx context.Context, keywords []string) ([]domain.UserSearch, error) // 搜索用户
	IsExistPost(ctx context.Context, postId uint) (bool, error)
	IsExistUser(ctx context.Context, userId int64) (bool, error)
	InputUser(ctx context.Context, user domain.UserSearch) error // 处理输入用户
	InputPost(ctx context.Context, post domain.PostSearch) error
	BulkInputLogs(ctx context.Context, event []domain.ReadEvent) error // 处理输入文章
	DeleteUserIndex(ctx context.Context, userId int64) error           // 删除用户索引
	DeletePostIndex(ctx context.Context, postId uint) error            // 删除文章索引
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

func (s *searchRepository) IsExistPost(ctx context.Context, postId uint) (bool, error) {
	return s.dao.IsExistsPost(ctx, fmt.Sprint(postId))
}

func (s *searchRepository) IsExistUser(ctx context.Context, userId int64) (bool, error) {
	return s.dao.IsExistsUser(ctx, fmt.Sprint(userId))
}

func (s *searchRepository) InputUser(ctx context.Context, user domain.UserSearch) error {
	return s.dao.InputUser(ctx, s.toDaoUserSearch(user))
}

func (s *searchRepository) InputPost(ctx context.Context, post domain.PostSearch) error {
	return s.dao.InputPost(ctx, s.toDaoPostSearch(post))
}

func (s *searchRepository) BulkInputLogs(ctx context.Context, event []domain.ReadEvent) error {
	return s.dao.BulkInputLogs(ctx, s.toDaoReadEvent(event))
}

func (s *searchRepository) DeleteUserIndex(ctx context.Context, userId int64) error {
	return s.dao.DeleteUserIndex(ctx, userId)
}

func (s *searchRepository) DeletePostIndex(ctx context.Context, postId uint) error {
	return s.dao.DeletePostIndex(ctx, postId)
}

func (s *searchRepository) toDaoPostSearch(domainPosts domain.PostSearch) dao.PostSearch {
	return dao.PostSearch{
		Content: domainPosts.Content,
		Id:      domainPosts.Id,
		Status:  domainPosts.Status,
		Tags:    domainPosts.Tags,
		Title:   domainPosts.Title,
	}
}

func (s *searchRepository) toDaoUserSearch(domainUsers domain.UserSearch) dao.UserSearch {
	return dao.UserSearch{
		Username: domainUsers.Username,
		Id:       domainUsers.Id,
		RealName: domainUsers.RealName,
	}
}

func (s *searchRepository) toDaoReadEvent(events []domain.ReadEvent) []dao.ReadEvent {
	daoEvents := make([]dao.ReadEvent, len(events))
	for _, e := range events {
		daoEvent := dao.ReadEvent{
			Timestamp: e.Timestamp,
			Level:     e.Level,
			Message:   e.Message,
		}
		daoEvents = append(daoEvents, daoEvent)
	}
	return daoEvents
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
			Username: daoUser.Username,
			Id:       daoUser.Id,
			RealName: daoUser.RealName,
		}
	}
	return domainUsers
}
