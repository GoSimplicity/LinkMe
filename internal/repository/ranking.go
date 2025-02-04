package repository

import (
	"context"
	"fmt"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"go.uber.org/zap"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, posts []domain.Post) error
	GetTopN(ctx context.Context) ([]domain.Post, error)
}

type rankingRepository struct {
	redisCache cache.RankingRedisCache
	localCache cache.RankingLocalCache
	l          *zap.Logger
}

func NewRankingCache(redisCache cache.RankingRedisCache, localCache cache.RankingLocalCache, l *zap.Logger) RankingRepository {
	return &rankingRepository{
		redisCache: redisCache,
		localCache: localCache,
		l:          l,
	}
}

// GetTopN 从缓存中获取排名前 N 的帖子
func (rc *rankingRepository) GetTopN(ctx context.Context) ([]domain.Post, error) {
	// 优先从本地缓存获取
	if posts, err := rc.localCache.Get(ctx); err == nil {
		rc.l.Debug("本地缓存命中")
		return posts, nil
	}

	// 从 Redis 获取
	posts, err := rc.redisCache.Get(ctx)
	if err != nil {
		rc.l.Warn("Redis 缓存未命中", zap.Error(err))
		// Redis 未命中时强制从本地缓存获取
		return rc.localCache.ForceGet(ctx)
	}

	rc.l.Debug("Redis 缓存命中")

	// 异步更新本地缓存
	go func() {
		if err := rc.localCache.Set(context.Background(), posts); err != nil {
			rc.l.Error("更新本地缓存失败", zap.Error(err))
		}
	}()

	return posts, nil
}

// ReplaceTopN 替换缓存中的排名前 N 的帖子
func (rc *rankingRepository) ReplaceTopN(ctx context.Context, posts []domain.Post) error {
	errChan := make(chan error, 2)

	// 并发更新缓存
	go func() {
		errChan <- rc.localCache.Set(ctx, posts)
	}()

	go func() {
		errChan <- rc.redisCache.Set(ctx, posts)
	}()

	// 等待两个更新操作完成
	var errs []error
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("替换缓存失败: %v", errs)
	}

	return nil
}
