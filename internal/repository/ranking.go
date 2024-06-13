package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/cache"
	"context"
	"fmt"
	"log"
	"sync"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, posts []domain.Post) error
	GetTopN(ctx context.Context) ([]domain.Post, error)
}

type rankingRepository struct {
	redisCache cache.RankingRedisCache
	localCache cache.RankingLocalCache
}

func NewRankingCache(redisCache cache.RankingRedisCache, localCache cache.RankingLocalCache) RankingRepository {
	return &rankingRepository{
		redisCache: redisCache,
		localCache: localCache,
	}
}

// GetTopN 从缓存中获取排名前 N 的帖子
func (rc *rankingRepository) GetTopN(ctx context.Context) ([]domain.Post, error) {
	// 尝试从本地缓存获取数据
	res, err := rc.localCache.Get(ctx)
	if err == nil {
		log.Println("本地缓存命中")
		return res, nil
	}
	log.Printf("本地缓存未命中: %v", err)
	// 尝试从 Redis 缓存获取数据
	res, err = rc.redisCache.Get(ctx)
	if err != nil {
		log.Printf("Redis 缓存未命中: %v", err)
		// 尝试强制从本地缓存获取数据
		return rc.localCache.ForceGet(ctx)
	}
	log.Println("Redis 缓存命中")
	// 将数据设置到本地缓存
	if er := rc.localCache.Set(ctx, res); er != nil {
		log.Printf("设置本地缓存时出错: %v", er)
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
			log.Printf("设置本地缓存时出错: %v", err)
		}
	}()
	// 并发设置 Redis 缓存
	go func() {
		defer wg.Done()
		if err := rc.redisCache.Set(ctx, posts); err != nil {
			redisErr = err
			log.Printf("设置 Redis 缓存时出错: %v", err)
		}
	}()
	// 等待所有并发操作完成
	wg.Wait()
	// 如果任何一个缓存操作失败，返回错误
	if localErr != nil || redisErr != nil {
		return fmt.Errorf("替换前 N 失败: localErr=%v, redisErr=%v", localErr, redisErr)
	}
	return nil
}
