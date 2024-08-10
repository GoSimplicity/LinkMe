package local

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

// CacheManager 包含本地缓存和 Redis 客户端
type CacheManager struct {
	localCache  *cache.Cache
	redisClient redis.Cmdable
}

// NewLocalCacheManager 创建并初始化一个 CacheManager 实例
func NewLocalCacheManager(redisClient redis.Cmdable) *CacheManager {
	localCache := cache.New(5*time.Minute, 10*time.Minute)

	return &CacheManager{localCache: localCache, redisClient: redisClient}
}

// Set 缓存数据到本地缓存和 Redis
func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	cm.localCache.Set(key, value, cache.DefaultExpiration)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return cm.redisClient.Set(ctx, key, data, expiration).Err()
}

// Get 从缓存中获取数据，如果缓存未命中，则调用 loader 加载数据并缓存
func (cm *CacheManager) Get(ctx context.Context, key string, loader func() (interface{}, error), result interface{}) error {
	// 尝试从本地缓存获取
	if cachedValue, found := cm.localCache.Get(key); found {
		// 将缓存值反序列化为 result 类型
		data, err := json.Marshal(cachedValue)
		if err != nil {
			return err
		}

		return json.Unmarshal(data, result)
	}
	// 尝试从 Redis 获取
	data, err := cm.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存未命中，从 loader 加载数据
			value, err := loader()
			if err != nil {
				return err
			}

			// 缓存加载的数据
			err = cm.Set(ctx, key, value, cache.DefaultExpiration)
			if err != nil {
				return err
			}

			// 将加载的数据转换为字节数组
			data, err = json.Marshal(value)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return json.Unmarshal(data, result)
}

// Delete 从缓存中删除数据
func (cm *CacheManager) Delete(ctx context.Context, key string) error {
	cm.localCache.Delete(key)

	return cm.redisClient.Del(ctx, key).Err()
}

// SetEmptyCache 缓存空对象，防止缓存穿透
func (cm *CacheManager) SetEmptyCache(ctx context.Context, key string, ttl time.Duration) error {
	emptyValue, err := json.Marshal(struct{}{})
	if err != nil {
		return err
	}

	if err := cm.redisClient.Set(ctx, key, emptyValue, ttl).Err(); err != nil {
		return err
	}

	cm.localCache.Set(key, emptyValue, cache.DefaultExpiration)

	return nil
}
