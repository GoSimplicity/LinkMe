package repository

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
)

// RelationRepository 定义了关注关系的存储库接口
type RelationRepository interface {
	ListRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]domain.Relation, error)
	GetRelationInfo(ctx context.Context, followerID, followeeID int64) (domain.Relation, error)
	FollowUser(ctx context.Context, followerID, followeeID int64) error
	CancelFollowUser(ctx context.Context, followerID, followeeID int64) error
}

type relationRepository struct {
	dao dao.RelationDAO
}

func NewRelationRepository(dao dao.RelationDAO) RelationRepository {
	return &relationRepository{
		dao: dao,
	}
}

// ListRelations 列出所有关注关系
func (r *relationRepository) ListRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]domain.Relation, error) {
	// TODO 实现方法
	panic("implement me")
}

// GetRelationInfo 获取特定的关注关系信息
func (r *relationRepository) GetRelationInfo(ctx context.Context, followerID, followeeID int64) (domain.Relation, error) {
	// TODO 实现方法
	panic("implement me")
}

// FollowUser 关注用户
func (r *relationRepository) FollowUser(ctx context.Context, followerID, followeeID int64) error {
	// TODO 实现方法
	panic("implement me")
}

func (r *relationRepository) CancelFollowUser(ctx context.Context, followerID, followeeID int64) error {
	//TODO implement me
	panic("implement me")
}
