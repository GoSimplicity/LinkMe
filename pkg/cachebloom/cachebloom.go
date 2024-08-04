package cachebloom

import (
	"context"
	"errors"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/redis/go-redis/v9"
	"time"
)

// CacheBloom 是一个包含布隆过滤器和 Redis 客户端的结构体
type CacheBloom struct {
	bf     *bloom.BloomFilter
	client *redis.Client
}

// NewCacheBloom 创建并初始化一个 CacheBloom 实例
func NewCacheBloom(client *redis.Client, filterSize uint, falsePositiveRate float64) *CacheBloom {
	return &CacheBloom{
		bf:     bloom.NewWithEstimates(filterSize, falsePositiveRate),
		client: client,
	}
}

// QueryData 查询数据，如果缓存和布隆过滤器都没有命中，则查询数据库
func (cb *CacheBloom) QueryData(ctx context.Context, key string, queryDB func(string) string, ttl time.Duration) (string, error) {
	// 检查布隆过滤器
	if !cb.bf.TestString(key) {
		return "", nil // 数据不存在，直接返回空结果
	}
	// 检查缓存
	data, err := cb.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存中没有数据，查询数据库
			data = queryDB(key)
			if data != "" {
				if er := cb.client.Set(ctx, key, data, ttl).Err(); er != nil {
					return "", er // 返回设置缓存时的错误
				}
				cb.bf.AddString(key) // 更新布隆过滤器
			} else {
				if er := cb.client.Set(ctx, key, "", ttl).Err(); er != nil {
					return "", er // 返回设置缓存空对象时的错误
				}
				cb.bf.AddString(key) // 更新布隆过滤器
			}
		} else {
			return "", err // 返回Redis操作时发生的其他错误
		}
	}
	return data, nil
}
