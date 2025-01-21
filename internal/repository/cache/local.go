package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"go.uber.org/zap"
)

var (
	ErrCacheEmpty   = errors.New("本地缓存为空")
	ErrCacheExpired = errors.New("本地缓存已过期")
)

type RankingLocalCache interface {
	Set(ctx context.Context, posts []domain.Post) error
	Get(ctx context.Context) ([]domain.Post, error)
	ForceGet(ctx context.Context) ([]domain.Post, error)
}

type rankingLocalCache struct {
	posts    []domain.Post
	expireAt time.Time
	ttl      time.Duration
	mu       sync.RWMutex
	logger   *zap.Logger
}

func NewRankingLocalCache(logger *zap.Logger) RankingLocalCache {
	return &rankingLocalCache{
		ttl:    10 * time.Minute,
		logger: logger,
	}
}

func (r *rankingLocalCache) Set(ctx context.Context, posts []domain.Post) error {
	if len(posts) == 0 {
		r.logger.Warn("试图设置空缓存")
		return ErrCacheEmpty
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts = make([]domain.Post, len(posts))
	copy(r.posts, posts)
	r.expireAt = time.Now().Add(r.ttl)

	r.logger.Info("本地缓存已更新",
		zap.Int("post_count", len(posts)),
		zap.Time("expire_at", r.expireAt))
	return nil
}

func (r *rankingLocalCache) Get(ctx context.Context) ([]domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if err := r.validateCache(); err != nil {
		return nil, err
	}

	return r.copyPosts(), nil
}

func (r *rankingLocalCache) ForceGet(ctx context.Context) ([]domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.posts) == 0 {
		r.logger.Warn("强制获取:本地缓存为空")
		return nil, ErrCacheEmpty
	}

	result := r.copyPosts()
	r.logger.Debug("强制获取本地缓存成功",
		zap.Int("post_count", len(result)),
		zap.Time("expire_at", r.expireAt))
	return result, nil
}

func (r *rankingLocalCache) validateCache() error {
	if len(r.posts) == 0 {
		r.logger.Warn("本地缓存为空")
		return ErrCacheEmpty
	}

	if time.Now().After(r.expireAt) {
		r.logger.Warn("本地缓存已过期",
			zap.Time("expire_at", r.expireAt))
		return ErrCacheExpired
	}

	return nil
}

func (r *rankingLocalCache) copyPosts() []domain.Post {
	result := make([]domain.Post, len(r.posts))
	copy(result, r.posts)
	r.logger.Debug("本地缓存命中",
		zap.Int("post_count", len(result)))
	return result
}
