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

type CacheBloom struct {
	bf     *bloom.BloomFilter
	client redis.Cmdable
}

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
		//如果data数据不为空，则缓存数据并返回data
		if !isNil(data) {
			return CacheData(cb, ctx, key, data, ttl)
		}
		return zeroValue, nil
	}

	// 检查缓存
	cachedData, err := cb.client.Get(ctx, key).Result()
	if err != nil {
		// 判断错误是否为redis.Nil，即缓存不存在
		if errors.Is(err, redis.Nil) {
			// 如果data数据不为空，则缓存数据并返回data
			if !isNil(data) {
				return CacheData(cb, ctx, key, data, ttl)
			}
			return zeroValue, nil
		}
		return zeroValue, err
	}

	// 走到这里说明缓存中有数据，反序列化缓存中的数据
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

	// 设置缓存
	if err := cb.client.Set(ctx, key, serializedData, ttl).Err(); err != nil {
		return data, err
	}

	// 添加key到布隆过滤器
	cb.bf.AddString(key)

	return data, nil
}

// SetEmptyCache 缓存空对象，防止缓存穿透
func (cb *CacheBloom) SetEmptyCache(ctx context.Context, key string, ttl time.Duration) error {
	emptyValue, err := json.Marshal(struct{}{})

	if err != nil {
		return err
	}

	// 设置缓存空对象
	if err := cb.client.Set(ctx, key, emptyValue, ttl).Err(); err != nil {
		return err
	}

	// 添加key到布隆过滤器
	cb.bf.AddString(key)

	return nil
}

// rebuildBloomFilter 重建布隆过滤器
func (cb *CacheBloom) rebuildBloomFilter(ctx context.Context) {
	// 创建一个新的布隆过滤器
	newBF := bloom.NewWithEstimates(1000000, 0.01)

	// 从 Redis 获取所有键并将它们添加到新的布隆过滤器中
	keys, err := cb.client.Keys(ctx, "*").Result()
	if err != nil {
		return
	}

	for _, key := range keys {
		newBF.AddString(key)
	}

	// 用新的布隆过滤器替换旧的布隆过滤器
	cb.bf = newBF
}

// isNil 使用反射检查值是否为 nil
func isNil(value interface{}) bool {
	// 如果值本身就是 nil，直接返回 true
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)

	// 处理指针类型
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		return v.IsNil()
	}

	// 处理零值情况（例如空字符串、零值结构体等）
	return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())
}
