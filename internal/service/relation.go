package service

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
)

type RelationService interface {
	ListFollowerRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]domain.Relation, error)
	ListFolloweeRelations(ctx context.Context, followeeID int64, pagination domain.Pagination) ([]domain.Relation, error)
	FollowUser(ctx context.Context, followerID, followeeID int64) error
	CancelFollowUser(ctx context.Context, followerID, followeeID int64) error
	GetFolloweeCount(ctx context.Context, UserID int64) (int64, error)
	GetFollowerCount(ctx context.Context, UserID int64) (int64, error)
}

type relationService struct {
	repo repository.RelationRepository
}

func NewRelationService(repo repository.RelationRepository) RelationService {
	return &relationService{
		repo: repo,
	}
}

// ListFollowerRelations 列出所有关注关系
func (r *relationService) ListFollowerRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]domain.Relation, error) {
	// 计算偏移量
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	return r.repo.ListFollowerRelations(ctx, followerID, pagination)
}

// ListFolloweeRelations 获取特定的关注关系信息
func (r *relationService) ListFolloweeRelations(ctx context.Context, followeeID int64, pagination domain.Pagination) ([]domain.Relation, error) {
	// 计算偏移量
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	return r.repo.ListFolloweeRelations(ctx, followeeID, pagination)
}

// FollowUser 关注用户
func (r *relationService) FollowUser(ctx context.Context, followerID, followeeID int64) error {
	return r.repo.FollowUser(ctx, followerID, followeeID)
}

// CancelFollowUser 取消关注用户
func (r *relationService) CancelFollowUser(ctx context.Context, followerID, followeeID int64) error {
	return r.repo.CancelFollowUser(ctx, followerID, followeeID)
}

func (r *relationService) GetFolloweeCount(ctx context.Context, UserID int64) (int64, error) {
	return r.repo.GetFolloweeCount(ctx, UserID)
}

func (r *relationService) GetFollowerCount(ctx context.Context, UserID int64) (int64, error) {
	return r.repo.GetFollowerCount(ctx, UserID)
}
