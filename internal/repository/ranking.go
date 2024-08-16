package repository

import (
	"context"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"go.uber.org/zap"
	"sync"
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
	// 尝试从本地缓存获取数据
	res, err := rc.localCache.Get(ctx)
	if err == nil {
		rc.l.Info("本地缓存命中")
		return res, nil
	}

	rc.l.Warn("本地缓存未命中", zap.Error(err))

	// 尝试从 Redis 缓存获取数据
	res, err = rc.redisCache.Get(ctx)
	if err != nil {
		rc.l.Warn("Redis 缓存未命中", zap.Error(err))
		// 尝试强制从本地缓存获取数据
		return rc.localCache.ForceGet(ctx)
	}

	rc.l.Info("redis缓存命中")

	// 将数据设置到本地缓存
	if er := rc.localCache.Set(ctx, res); er != nil {
		rc.l.Error("设置本地缓存时出错", zap.Error(er))
	}

	return res, nil
}

// ReplaceTopN 替换缓存中的排名前 N 的帖子
func (rc *rankingRepository) ReplaceTopN(ctx context.Context, posts []domain.Post) error {
	var wg sync.WaitGroup
	var localErr, redisErr error

	wg.Add(2)

	// 并发设置本地缓存
	go func() {
		defer wg.Done()
		if err := rc.localCache.Set(ctx, posts); err != nil {
			localErr = err
			rc.l.Error("设置本地缓存时出错", zap.Error(err))
		}
	}()

	// 并发设置 Redis 缓存
	go func() {
		defer wg.Done()
		if err := rc.redisCache.Set(ctx, posts); err != nil {
			redisErr = err
			rc.l.Error("设置 Redis 缓存时出错", zap.Error(err))
		}
	}()

	wg.Wait()

	if localErr != nil || redisErr != nil {
		return fmt.Errorf("替换前 N 失败: localErr=%v, redisErr=%v", localErr, redisErr)
	}

	return nil
}
