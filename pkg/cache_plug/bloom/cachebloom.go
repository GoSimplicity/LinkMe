package bloom

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/redis/go-redis/v9"
)

// CacheBloom 包含布隆过滤器和 Redis 客户端
type CacheBloom struct {
	bf     *bloom.BloomFilter
	client redis.Cmdable
}

// NewCacheBloom 创建并初始化一个 CacheBloom 实例
func NewCacheBloom(client redis.Cmdable) *CacheBloom {
	return &CacheBloom{
		bf:     bloom.NewWithEstimates(1000000, 0.01),
		client: client,
	}
}

// QueryData 查询数据，如果缓存和布隆过滤器都没有命中，则存储并返回传递的数据
func QueryData[T any](cb *CacheBloom, ctx context.Context, key string, data T, ttl time.Duration) (T, error) {
	// 定义一个零值变量来处理错误返回
	var zeroValue T

	// 检查布隆过滤器
	if !cb.bf.TestString(key) {
		if !isNil(data) {
			return CacheData(cb, ctx, key, data, ttl)
		}
		return zeroValue, nil
	}

	// 检查缓存
	cachedData, err := cb.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			if !isNil(data) {
				return CacheData(cb, ctx, key, data, ttl)
			}
			return zeroValue, nil
		}
		return zeroValue, err
	}

	// 反序列化缓存中的数据
	var result T
	if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
		return zeroValue, err
	}
	return result, nil
}

// CacheData 序列化数据并存储到 Redis 和布隆过滤器中
func CacheData[T any](cb *CacheBloom, ctx context.Context, key string, data T, ttl time.Duration) (T, error) {
	serializedData, err := json.Marshal(data)
	if err != nil {
		return data, err
	}

	if err := cb.client.Set(ctx, key, serializedData, ttl).Err(); err != nil {
		return data, err
	}

	cb.bf.AddString(key)
	return data, nil
}

// SetEmptyCache 缓存空对象，防止缓存穿透
func (cb *CacheBloom) SetEmptyCache(ctx context.Context, key string, ttl time.Duration) error {
	emptyValue, err := json.Marshal(struct{}{})
	if err != nil {
		return err
	}
	if err := cb.client.Set(ctx, key, emptyValue, ttl).Err(); err != nil {
		return err
	}
	cb.bf.AddString(key)
	return nil
}

// isNil 使用反射检查值是否为 nil
func isNil(value interface{}) bool {
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Ptr && v.IsNil()
}
