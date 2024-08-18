package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RelationCache interface {
	GetCache(ctx context.Context, key string) ([]domain.Relation, error)
	SetCache(ctx context.Context, key string, relations []domain.Relation, ttl time.Duration) error
	GetCountCache(ctx context.Context, key string) (int64, error)
	SetCountCache(ctx context.Context, key string, count int64, ttl time.Duration) error
	ClearFollowCache(ctx context.Context, followerID, followeeID int64) error
	GenerateCacheKey(userID int64, relationType string, pagination domain.Pagination) string
	GenerateCountCacheKey(userID int64, relationType string) string
}

type relationCache struct {
	client redis.Cmdable
}

func NewRelationCache(client redis.Cmdable) RelationCache {
	return &relationCache{client: client}
}

func (c *relationCache) GetCache(ctx context.Context, key string) ([]domain.Relation, error) {
	// 从Redis获取数据
	data, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // 缓存未命中
	} else if err != nil {
		return nil, err // Redis调用失败
	}

	// 解码数据
	var relations []domain.Relation

	if err := json.Unmarshal([]byte(data), &relations); err != nil {
		return nil, err
	}

	return relations, nil
}

func (c *relationCache) SetCache(ctx context.Context, key string, relations []domain.Relation, ttl time.Duration) error {
	// 编码数据为JSON格式
	data, err := json.Marshal(relations)
	if err != nil {
		return err
	}

	// 设置缓存并设置过期时间
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *relationCache) GetCountCache(ctx context.Context, key string) (int64, error) {
	count, err := c.client.Get(ctx, key).Int64()
	if errors.Is(err, redis.Nil) {
		return 0, nil // 缓存未命中
	} else if err != nil {
		return 0, err // Redis调用失败
	}

	return count, nil
}

func (c *relationCache) SetCountCache(ctx context.Context, key string, count int64, ttl time.Duration) error {
	// 设置计数缓存并设置过期时间
	return c.client.Set(ctx, key, count, ttl).Err()
}

func (c *relationCache) ClearFollowCache(ctx context.Context, followerID, followeeID int64) error {
	// 批量删除相关缓存
	keys := []string{
		c.GenerateCacheKey(followerID, "followers", domain.Pagination{}),
		c.GenerateCacheKey(followeeID, "followees", domain.Pagination{}),
		c.GenerateCountCacheKey(followerID, "followers"),
		c.GenerateCountCacheKey(followeeID, "followees"),
	}

	// 使用Pipeline进行批量删除
	pipe := c.client.Pipeline()

	for _, key := range keys {
		pipe.Del(ctx, key)
	}

	_, err := pipe.Exec(ctx)

	return err
}

func (c *relationCache) GenerateCacheKey(userID int64, relationType string, pagination domain.Pagination) string {
	// 生成缓存键
	return fmt.Sprintf("relation:%s:%d:offset=%d:size=%d",
		relationType, userID, *pagination.Offset, *pagination.Size)
}

func (c *relationCache) GenerateCountCacheKey(userID int64, relationType string) string {
	// 生成计数缓存键
	return fmt.Sprintf("relation:count:%s:%d", relationType, userID)
}
