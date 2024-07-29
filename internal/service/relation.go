package service

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
)

type RelationService interface {
	ListRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]domain.Relation, error)
	GetRelationInfo(ctx context.Context, followerID, followeeID int64) (domain.Relation, error)
	FollowUser(ctx context.Context, followerID, followeeID int64) error
	CancelFollowUser(ctx context.Context, followerID, followeeID int64) error
}

type relationService struct {
	repo repository.RelationRepository
}

func NewRelationService(repo repository.RelationRepository) RelationService {
	return &relationService{
		repo: repo,
	}
}

// ListRelations 列出所有关注关系
func (r *relationService) ListRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]domain.Relation, error) {
	// TODO 实现方法
	panic("implement me")
}

// GetRelationInfo 获取特定的关注关系信息
func (r *relationService) GetRelationInfo(ctx context.Context, followerID, followeeID int64) (domain.Relation, error) {
	// TODO 实现方法
	panic("implement me")
}

// FollowUser 关注用户
func (r *relationService) FollowUser(ctx context.Context, followerID, followeeID int64) error {
	// TODO 实现方法
	panic("implement me")
}

// CancelFollowUser 取消关注用户
func (r *relationService) CancelFollowUser(ctx context.Context, followerID, followeeID int64) error {
	// TODO 实现方法
	panic("implement me")
}
