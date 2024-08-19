package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

type CheckCache interface {
	GetCache(ctx context.Context, key string) (*domain.Check, error)
	SetCache(ctx context.Context, key string, check domain.Check, ttl time.Duration) error
	ClearCache(ctx context.Context, key string) error
	DeleteKeysWithPattern(ctx context.Context, pattern string) error
	GetCacheList(ctx context.Context, key string) ([]domain.Check, error)
	SetCacheList(ctx context.Context, key string, checks []domain.Check, ttl time.Duration) error
	GetCountCache(ctx context.Context, key string) (int64, error)
	SetCountCache(ctx context.Context, key string, count int64, ttl time.Duration) error
	GenerateCacheKey(id interface{}) string
	GeneratePaginationCacheKey(pagination domain.Pagination) string
	GenerateCountCacheKey() string
}

type checkCache struct {
	client redis.Cmdable
}

func NewCheckCache(client redis.Cmdable) CheckCache {
	return &checkCache{client: client}
}

// GetCache 从缓存中获取单个审核记录
func (c *checkCache) GetCache(ctx context.Context, key string) (*domain.Check, error) {
	// 从Redis获取数据
	data, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 缓存未命中
		return nil, nil
	} else if err != nil {
		// Redis调用出错
		return nil, err
	}

	var check domain.Check
	// 将JSON数据反序列化为domain.Check对象
	if err := json.Unmarshal([]byte(data), &check); err != nil {
		return nil, err
	}

	return &check, nil
}

// SetCache 将单个审核记录写入缓存
func (c *checkCache) SetCache(ctx context.Context, key string, check domain.Check, ttl time.Duration) error {
	// 将domain.Check对象序列化为JSON格式
	data, err := json.Marshal(check)
	if err != nil {
		return err
	}

	// 设置随机化的TTL，防止缓存雪崩
	randomTTL := ttl + time.Duration(randomInt(0, 300))*time.Second

	// 将数据写入Redis
	return c.client.Set(ctx, key, data, randomTTL).Err()
}

// ClearCache 清除缓存中的指定键
func (c *checkCache) ClearCache(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *checkCache) DeleteKeysWithPattern(ctx context.Context, pattern string) error {
	var cursor uint64

	for {
		// 执行 SCAN 命令
		keys, cursor, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan Redis keys: %w", err)
		}

		// 如果找到匹配的键，进行删除
		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("failed to delete Redis keys: %w", err)
			}
			fmt.Printf("Deleted keys: %v\n", keys)
		}

		// 如果 cursor 为 0，说明遍历结束
		if cursor == 0 {
			break
		}
	}

	return nil
}

// GetCacheList 从缓存中获取审核记录列表
func (c *checkCache) GetCacheList(ctx context.Context, key string) ([]domain.Check, error) {
	// 从Redis获取数据
	data, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 缓存未命中
		return nil, nil
	} else if err != nil {
		// Redis调用出错
		return nil, err
	}

	var checks []domain.Check
	// 将JSON数据反序列化为[]domain.Check对象
	if err := json.Unmarshal([]byte(data), &checks); err != nil {
		return nil, err
	}

	return checks, nil
}

// SetCacheList 将审核记录列表写入缓存
func (c *checkCache) SetCacheList(ctx context.Context, key string, checks []domain.Check, ttl time.Duration) error {
	// 将[]domain.Check对象序列化为JSON格式
	data, err := json.Marshal(checks)
	if err != nil {
		return err
	}

	// 设置随机化的TTL，防止缓存雪崩
	randomTTL := ttl + time.Duration(randomInt(0, 300))*time.Second

	// 将数据写入Redis
	return c.client.Set(ctx, key, data, randomTTL).Err()
}

// GetCountCache 从缓存中获取审核记录的数量
func (c *checkCache) GetCountCache(ctx context.Context, key string) (int64, error) {
	// 从Redis获取计数数据
	count, err := c.client.Get(ctx, key).Int64()
	if errors.Is(err, redis.Nil) {
		// 缓存未命中
		return 0, nil
	} else if err != nil {
		// Redis调用出错
		return 0, err
	}

	return count, nil
}

// SetCountCache 将审核记录的数量写入缓存
func (c *checkCache) SetCountCache(ctx context.Context, key string, count int64, ttl time.Duration) error {
	// 设置随机化的TTL，防止缓存雪崩
	randomTTL := ttl + time.Duration(randomInt(0, 300))*time.Second

	// 将计数数据写入Redis
	return c.client.Set(ctx, key, count, randomTTL).Err()
}

// GenerateCacheKey 生成单个审核记录的缓存键
func (c *checkCache) GenerateCacheKey(id interface{}) string {
	return fmt.Sprintf("linkme:check:%v", id)
}

// GeneratePaginationCacheKey 生成分页查询的缓存键
func (c *checkCache) GeneratePaginationCacheKey(pagination domain.Pagination) string {
	return fmt.Sprintf("linkme:check:list:%d", pagination.Page)
}

// GenerateCountCacheKey 生成审核记录数量的缓存键
func (c *checkCache) GenerateCountCacheKey() string {
	return "linkme:check:count"
}

// randomInt 生成 min 到 max 之间的随机整数，用于随机化 TTL
func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}
