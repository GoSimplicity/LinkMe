package local

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

type CacheManager struct {
	localCache  *cache.Cache // 本地缓存
	redisClient redis.Cmdable
}

func NewLocalCacheManager(redisClient redis.Cmdable) *CacheManager {
	return &CacheManager{
		localCache:  cache.New(5*time.Minute, 10*time.Minute), // 默认缓存 5 分钟
		redisClient: redisClient,
	}
}

// Set 缓存数据到本地缓存和 Redis
func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// 并发设置本地缓存和Redis
	errChan := make(chan error, 2)

	go func() {
		cm.localCache.Set(key, value, expiration)
		errChan <- nil
	}()

	go func() {
		errChan <- cm.redisClient.Set(ctx, key, data, expiration).Err()
	}()

	// 等待两个操作完成
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

// Get 从缓存中获取数据，如果缓存未命中，则调用 loader 加载数据并缓存
func (cm *CacheManager) Get(ctx context.Context, key string, loader func() (interface{}, error), result interface{}) error {
	// 尝试从本地缓存获取
	if value, found := cm.localCache.Get(key); found {
		return cm.unmarshalValue(value, result)
	}

	// 尝试从 Redis 获取
	data, err := cm.redisClient.Get(ctx, key).Bytes()
	if err == nil {
		if err := json.Unmarshal(data, result); err != nil {
			return err
		}
		// 同步到本地缓存
		cm.localCache.Set(key, result, cache.DefaultExpiration)
		return nil
	}

	if !errors.Is(err, redis.Nil) {
		return err
	}

	// 缓存未命中，从 loader 加载数据
	value, err := loader()
	if err != nil {
		return err
	}

	// 缓存加载的数据
	if err := cm.Set(ctx, key, value, cache.DefaultExpiration); err != nil {
		return err
	}

	return cm.unmarshalValue(value, result)
}

// Delete 从本地缓存和 Redis 中删除一个或多个键
func (cm *CacheManager) Delete(ctx context.Context, keys ...string) error {
	// 并发删除本地缓存和Redis
	errChan := make(chan error, 2)

	go func() {
		for _, key := range keys {
			cm.localCache.Delete(key)
		}
		errChan <- nil
	}()

	go func() {
		errChan <- cm.redisClient.Del(ctx, keys...).Err()
	}()

	// 等待两个操作完成
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

// SetEmptyCache 缓存空对象，防止缓存穿透
func (cm *CacheManager) SetEmptyCache(ctx context.Context, key string, ttl time.Duration) error {
	emptyValue := struct{}{}
	return cm.Set(ctx, key, emptyValue, ttl)
}

// unmarshalValue 将缓存值反序列化为指定类型
func (cm *CacheManager) unmarshalValue(value interface{}, result interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}
