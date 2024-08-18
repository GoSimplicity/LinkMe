package repository

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
	"time"
)

// RelationRepository 定义了关注关系的存储库接口
type RelationRepository interface {
	ListFollowerRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]domain.Relation, error)
	ListFolloweeRelations(ctx context.Context, followeeID int64, pagination domain.Pagination) ([]domain.Relation, error)
	FollowUser(ctx context.Context, followerID, followeeID int64) error
	CancelFollowUser(ctx context.Context, followerID, followeeID int64) error
	GetFolloweeCount(ctx context.Context, userID int64) (int64, error)
	GetFollowerCount(ctx context.Context, userID int64) (int64, error)
}

type relationRepository struct {
	dao    dao.RelationDAO
	cache  cache.RelationCache
	logger *zap.Logger
}

func NewRelationRepository(dao dao.RelationDAO, cache cache.RelationCache, logger *zap.Logger) RelationRepository {
	return &relationRepository{
		dao:    dao,
		cache:  cache,
		logger: logger,
	}
}

// FollowUser 关注用户
func (r *relationRepository) FollowUser(ctx context.Context, followerID, followeeID int64) error {
	// 关注成功后清除相关缓存
	if err := r.dao.FollowUser(ctx, followerID, followeeID); err != nil {
		return err
	}
	r.cache.ClearFollowCache(ctx, followerID, followeeID)

	return nil
}

func (r *relationRepository) CancelFollowUser(ctx context.Context, followerID, followeeID int64) error {
	// 取消关注后清除相关缓存
	if err := r.dao.CancelFollowUser(ctx, followerID, followeeID); err != nil {
		return err
	}
	r.cache.ClearFollowCache(ctx, followerID, followeeID)

	return nil
}

// ListFollowerRelations 列出粉丝列表
func (r *relationRepository) ListFollowerRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]domain.Relation, error) {
	cacheKey := r.cache.GenerateCacheKey(followerID, "followers", pagination)
	if cachedRelations, err := r.cache.GetCache(ctx, cacheKey); err == nil && cachedRelations != nil {
		r.logger.Info("Cache hit for follower relations", zap.String("key", cacheKey))
		return cachedRelations, nil
	}

	// 如果缓存未命中，则从数据库获取
	relations, err := r.dao.ListFollowerRelations(ctx, followerID, pagination)
	if err != nil {
		return nil, err
	}

	// 将结果存入缓存
	r.cache.SetCache(ctx, cacheKey, r.toDomainRelationSlice(relations), 5*time.Minute)

	return r.toDomainRelationSlice(relations), nil
}

// ListFolloweeRelations 列出关注列表
func (r *relationRepository) ListFolloweeRelations(ctx context.Context, followeeID int64, pagination domain.Pagination) ([]domain.Relation, error) {
	cacheKey := r.cache.GenerateCacheKey(followeeID, "followees", pagination)
	if cachedRelations, err := r.cache.GetCache(ctx, cacheKey); err == nil && cachedRelations != nil {
		r.logger.Info("Cache hit for followee relations", zap.String("key", cacheKey))
		return cachedRelations, nil
	}

	// 如果缓存未命中，则从数据库获取
	relations, err := r.dao.ListFolloweeRelations(ctx, followeeID, pagination)
	if err != nil {
		return nil, err
	}

	// 将结果存入缓存
	r.cache.SetCache(ctx, cacheKey, r.toDomainRelationSlice(relations), 5*time.Minute)

	return r.toDomainRelationSlice(relations), nil
}

func (r *relationRepository) GetFolloweeCount(ctx context.Context, userID int64) (int64, error) {
	cacheKey := r.cache.GenerateCountCacheKey(userID, "followees")
	if cachedCount, err := r.cache.GetCountCache(ctx, cacheKey); err == nil {
		r.logger.Info("Cache hit for followee count", zap.String("key", cacheKey))
		return cachedCount, nil
	}

	count, err := r.dao.FollowCount(ctx, userID)
	if err != nil {
		return 0, err
	}

	// 缓存关注数
	r.cache.SetCountCache(ctx, cacheKey, count.FolloweeCount, 5*time.Minute)

	return count.FolloweeCount, nil
}

func (r *relationRepository) GetFollowerCount(ctx context.Context, userID int64) (int64, error) {
	cacheKey := r.cache.GenerateCountCacheKey(userID, "followers")
	if cachedCount, err := r.cache.GetCountCache(ctx, cacheKey); err == nil {
		r.logger.Info("Cache hit for follower count", zap.String("key", cacheKey))
		return cachedCount, nil
	}

	count, err := r.dao.FollowCount(ctx, userID)
	if err != nil {
		return 0, err
	}

	// 缓存粉丝数
	err = r.cache.SetCountCache(ctx, cacheKey, count.FollowerCount, 5*time.Minute)
	if err != nil {
		return 0, err
	}

	return count.FollowerCount, nil
}

func (r *relationRepository) toDomainRelation(relation dao.Relation) domain.Relation {
	return domain.Relation{
		FolloweeId: relation.FolloweeID,
		FollowerId: relation.FollowerID,
	}
}

func (r *relationRepository) toDomainRelationSlice(relations []dao.Relation) []domain.Relation {
	relationSlice := make([]domain.Relation, len(relations))

	for i, relation := range relations {
		relationSlice[i] = r.toDomainRelation(relation)
	}

	return relationSlice
}
