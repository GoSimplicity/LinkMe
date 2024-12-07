package cache

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
)

var (
	ErrCacheEmpty   = errors.New("local cache is empty")
	ErrCacheExpired = errors.New("local cache expired")
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
	if len(arts) == 0 {
		return ErrCacheEmpty
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.topN = make([]domain.Post, len(arts))
	copy(r.topN, arts) // 深拷贝避免外部修改
	r.ddl = time.Now().Add(r.expiration)

	log.Printf("已设置缓存，共 %d 篇文章，过期时间为 %s", len(arts), r.ddl)
	return nil
}

// Get 获取缓存内容，如果缓存已过期或为空则返回错误
func (r *rankingLocalCache) Get(ctx context.Context) ([]domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.topN) == 0 {
		log.Println("本地缓存为空")
		return nil, ErrCacheEmpty
	}

	if time.Now().After(r.ddl) {
		log.Println("本地缓存已过期")
		return nil, ErrCacheExpired
	}

	result := make([]domain.Post, len(r.topN))
	copy(result, r.topN) // 返回副本避免外部修改

	log.Printf("缓存命中，共 %d 篇文章", len(result))
	return result, nil
}

// ForceGet 强制获取缓存内容，即使缓存已过期也返回
func (r *rankingLocalCache) ForceGet(ctx context.Context) ([]domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.topN) == 0 {
		log.Println("强制获取：本地缓存为空")
		return nil, ErrCacheEmpty
	}

	result := make([]domain.Post, len(r.topN))
	copy(result, r.topN) // 返回副本避免外部修改

	log.Printf("强制获取：缓存命中，共 %d 篇文章", len(result))
	return result, nil
}
