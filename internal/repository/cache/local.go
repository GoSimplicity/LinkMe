package cache

import (
	"LinkMe/internal/domain"
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type RankingLocalCache interface {
	Set(ctx context.Context, arts []domain.Post) error   // 设置缓存内容并更新过期时间
	Get(ctx context.Context) ([]domain.Post, error)      // 获取缓存内容，如果缓存已过期或为空则返回错误
	ForceGet(ctx context.Context) ([]domain.Post, error) // 强制获取缓存内容，即使缓存已过期也返回
}

type rankingLocalCache struct {
	topN       []domain.Post // 缓存的排名帖子
	ddl        time.Time     // 缓存的过期时间
	expiration time.Duration // 缓存的持续时间
	mu         sync.RWMutex  // 读写锁，保证并发安全
}

func NewRankingLocalCache() RankingLocalCache {
	return &rankingLocalCache{
		expiration: 10 * time.Minute, // 固定为10分钟
	}
}

// Set 设置缓存内容并更新过期时间
func (r *rankingLocalCache) Set(ctx context.Context, arts []domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.topN = arts // 设置缓存内容
	r.ddl = time.Now().Add(r.expiration)
	log.Printf("Cache set with %d posts, expires at %s", len(arts), r.ddl)
	return nil
}

// Get 获取缓存内容，如果缓存已过期或为空则返回错误
func (r *rankingLocalCache) Get(ctx context.Context) ([]domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock() // 函数退出时释放读锁

	if len(r.topN) == 0 {
		log.Println("Local cache is empty.")
		return nil, errors.New("local cache is empty")
	}

	if r.ddl.Before(time.Now()) {
		log.Println("Local cache expired.")
		return nil, errors.New("local cache expired")
	}

	log.Printf("Cache hit with %d posts", len(r.topN))
	return r.topN, nil
}

// ForceGet 强制获取缓存内容，即使缓存已过期也返回
func (r *rankingLocalCache) ForceGet(ctx context.Context) ([]domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.topN) == 0 {
		log.Println("Force get: local cache is empty.")
		return nil, errors.New("local cache is empty")
	}

	log.Printf("Force get: cache hit with %d posts", len(r.topN))
	return r.topN, nil
}
